package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	"github.com/DARREN-2000/ai-hypervisor-platform/internal/api/handlers"
	"github.com/DARREN-2000/ai-hypervisor-platform/internal/orchestrator"
	"github.com/DARREN-2000/ai-hypervisor-platform/pkg/telemetry"
)

// APIServer provides REST API endpoints for VM management
type APIServer struct {
	router        *mux.Router
	server        *http.Server
	logger        *logrus.Logger
	metrics       telemetry.Metrics
	baseHandler   *handlers.BaseHandler
	tlsCert       string
	tlsKey        string
	vmManager     orchestrator.VMManager
	scheduler     orchestrator.Scheduler
	gpuOrch       orchestrator.GPUOrchestrator
	taskExecutor  orchestrator.TaskExecutor
	resMonitor    orchestrator.ResourceMonitor
	eventBus      orchestrator.EventBus
	auditLogger   orchestrator.AuditLogger
	stateStore    orchestrator.StateStore
	healthChecks  map[string]HealthCheck
}

// HealthCheck validates the state of a dependency.
type HealthCheck func(context.Context) error

// Config represents API server configuration
type Config struct {
	Address         string
	Port            int
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
	MaxHeaderBytes  int
	TLSCert         string
	TLSKey          string
}

// NewAPIServer creates a new API server instance
func NewAPIServer(cfg *Config, logger *logrus.Logger) *APIServer {
	addr := fmt.Sprintf("%s:%d", cfg.Address, cfg.Port)
	return &APIServer{
		router: mux.NewRouter(),
		logger: logger,
		metrics: telemetry.NewNoopMetrics(),
		baseHandler: handlers.NewBaseHandler(logger),
		tlsCert: cfg.TLSCert,
		tlsKey: cfg.TLSKey,
		server: &http.Server{
			Addr:           addr,
			ReadTimeout:    cfg.ReadTimeout,
			WriteTimeout:   cfg.WriteTimeout,
			IdleTimeout:    cfg.IdleTimeout,
			MaxHeaderBytes: cfg.MaxHeaderBytes,
		},
	}
}

// SetDependencies injects service dependencies
func (s *APIServer) SetDependencies(
	vmMgr orchestrator.VMManager,
	scheduler orchestrator.Scheduler,
	gpuOrch orchestrator.GPUOrchestrator,
	taskExec orchestrator.TaskExecutor,
	resMon orchestrator.ResourceMonitor,
	eventBus orchestrator.EventBus,
	audit orchestrator.AuditLogger,
	store orchestrator.StateStore,
) {
	s.vmManager = vmMgr
	s.scheduler = scheduler
	s.gpuOrch = gpuOrch
	s.taskExecutor = taskExec
	s.resMonitor = resMon
	s.eventBus = eventBus
	s.auditLogger = audit
	s.stateStore = store
}

// SetMetrics injects a metrics recorder
func (s *APIServer) SetMetrics(metrics telemetry.Metrics) {
	if metrics == nil {
		metrics = telemetry.NewNoopMetrics()
	}
	s.metrics = metrics
}

// SetHealthChecks injects readiness checks for the health endpoints.
func (s *APIServer) SetHealthChecks(checks map[string]HealthCheck) {
	s.healthChecks = checks
}

// Handler returns the HTTP handler for the API server.
func (s *APIServer) Handler() http.Handler {
	return s.router
}

// RegisterRoutes registers all API routes
func (s *APIServer) RegisterRoutes() {
	// Middleware
	s.router.Use(s.observabilityMiddleware)
	s.router.Use(s.errorRecoveryMiddleware)

	// Health and status
	s.router.HandleFunc("/health", s.handleHealth).Methods(http.MethodGet)
	s.router.HandleFunc("/ready", s.handleReady).Methods(http.MethodGet)
	s.router.HandleFunc("/live", s.handleLive).Methods(http.MethodGet)
	s.router.HandleFunc("/api/v1/metrics", s.handleClusterMetrics).Methods(http.MethodGet)

	// VM endpoints
	vmRoutes := s.router.PathPrefix("/api/v1/vms").Subrouter()
	vmRoutes.HandleFunc("", s.handleListVMs).Methods(http.MethodGet)
	vmRoutes.HandleFunc("", s.handleCreateVM).Methods(http.MethodPost)
	vmRoutes.HandleFunc("/{vmId}", s.handleGetVM).Methods(http.MethodGet)
	vmRoutes.HandleFunc("/{vmId}", s.handleUpdateVM).Methods(http.MethodPatch)
	vmRoutes.HandleFunc("/{vmId}", s.handleDeleteVM).Methods(http.MethodDelete)
	vmRoutes.HandleFunc("/{vmId}/start", s.handleStartVM).Methods(http.MethodPost)
	vmRoutes.HandleFunc("/{vmId}/stop", s.handleStopVM).Methods(http.MethodPost)
	vmRoutes.HandleFunc("/{vmId}/reboot", s.handleRebootVM).Methods(http.MethodPost)

	// GPU endpoints
	gpuRoutes := s.router.PathPrefix("/api/v1/gpus").Subrouter()
	gpuRoutes.HandleFunc("", s.handleListGPUs).Methods(http.MethodGet)
	gpuRoutes.HandleFunc("/{gpuId}", s.handleGetGPU).Methods(http.MethodGet)

	// Host endpoints
	hostRoutes := s.router.PathPrefix("/api/v1/hosts").Subrouter()
	hostRoutes.HandleFunc("", s.handleListHosts).Methods(http.MethodGet)
	hostRoutes.HandleFunc("/{nodeId}", s.handleGetHost).Methods(http.MethodGet)
	hostRoutes.HandleFunc("/{nodeId}/metrics", s.handleGetHostMetrics).Methods(http.MethodGet)

	// Task endpoints
	taskRoutes := s.router.PathPrefix("/api/v1/tasks").Subrouter()
	taskRoutes.HandleFunc("/{taskId}", s.handleGetTask).Methods(http.MethodGet)

	// WebSocket endpoints
	s.router.HandleFunc("/ws/cluster/events", s.handleClusterEvents).Methods(http.MethodGet)
	s.router.HandleFunc("/ws/vm/{vmId}/metrics", s.handleVMMetricsStream).Methods(http.MethodGet)

	s.server.Handler = s.router
}

// Start starts the API server
func (s *APIServer) Start() error {
	s.logger.Infof("Starting API server on %s", s.server.Addr)
	var err error
	if s.tlsCert != "" && s.tlsKey != "" {
		err = s.server.ListenAndServeTLS(s.tlsCert, s.tlsKey)
	} else {
		err = s.server.ListenAndServe()
	}
	if err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

// Stop gracefully stops the API server
func (s *APIServer) Stop(ctx context.Context) error {
	s.logger.Info("Stopping API server")
	return s.server.Shutdown(ctx)
}

// Handlers

func (s *APIServer) handleHealth(w http.ResponseWriter, r *http.Request) {
	if len(s.healthChecks) == 0 {
		health := map[string]interface{}{
			"status":  "healthy",
			"version": "1.0.0",
			"components": map[string]interface{}{
				"api":      "operational",
				"database": "operational",
				"libvirt":  "operational",
			},
		}
		s.writeJSON(w, http.StatusOK, health)
		return
	}

	report, statusCode := s.buildHealthReport(r.Context())
	s.writeJSON(w, statusCode, report)
}

func (s *APIServer) handleReady(w http.ResponseWriter, r *http.Request) {
	report, statusCode := s.buildHealthReport(r.Context())
	s.writeJSON(w, statusCode, report)
}

func (s *APIServer) handleLive(w http.ResponseWriter, r *http.Request) {
	s.writeJSON(w, http.StatusOK, map[string]interface{}{
		"status":    "alive",
		"version":   "1.0.0",
		"timestamp": time.Now().UTC(),
	})
}

func (s *APIServer) buildHealthReport(ctx context.Context) (map[string]interface{}, int) {
	components := map[string]interface{}{}
	status := "healthy"
	httpStatus := http.StatusOK

	for name, check := range s.healthChecks {
		component := map[string]interface{}{
			"status":     "operational",
			"checked_at": time.Now().UTC(),
		}
		if check != nil {
			if err := check(ctx); err != nil {
				component["status"] = "degraded"
				// Do not expose raw error details on public endpoints
				component["error"] = "dependency check failed"
				status = "degraded"
				httpStatus = http.StatusServiceUnavailable
			}
		}
		components[name] = component
	}

	return map[string]interface{}{
		"status":     status,
		"version":    "1.0.0",
		"components": components,
		"timestamp":  time.Now().UTC(),
	}, httpStatus
}

func (s *APIServer) handleListVMs(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	namespace := r.URL.Query().Get("namespace")
	state := r.URL.Query().Get("state")

	filters := make(map[string]string)
	if namespace != "" {
		filters["namespace"] = namespace
	}
	if state != "" {
		filters["state"] = state
	}

	vms, err := s.vmManager.ListVMs(ctx, filters)
	if err != nil {
		s.writeError(w, err)
		return
	}
	for _, vm := range vms {
		if vm == nil || vm.ResourceUsage == nil {
			continue
		}
		s.metrics.ObserveVMResource(
			vm.ID,
			vm.HostNodeID,
			vm.ResourceUsage.CPU,
			int64(vm.ResourceUsage.Memory),
			int64(vm.ResourceUsage.DiskIORead),
			int64(vm.ResourceUsage.DiskIOWrite),
			int64(vm.ResourceUsage.NetworkIn),
			int64(vm.ResourceUsage.NetworkOut),
		)
	}

	s.writeJSON(w, http.StatusOK, vms)
}

func (s *APIServer) handleCreateVM(w http.ResponseWriter, r *http.Request) {
	// Implementation stub - full implementation in main codebase
	s.writeJSON(w, http.StatusCreated, map[string]string{"status": "Creating VM..."})
}

func (s *APIServer) handleGetVM(w http.ResponseWriter, r *http.Request) {
	// Implementation stub
	s.writeJSON(w, http.StatusOK, map[string]string{"status": "Getting VM..."})
}

func (s *APIServer) handleUpdateVM(w http.ResponseWriter, r *http.Request) {
	// Implementation stub
	s.writeJSON(w, http.StatusOK, map[string]string{"status": "Updating VM..."})
}

func (s *APIServer) handleDeleteVM(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

func (s *APIServer) handleStartVM(w http.ResponseWriter, r *http.Request) {
	s.writeJSON(w, http.StatusAccepted, map[string]string{"status": "Starting VM..."})
}

func (s *APIServer) handleStopVM(w http.ResponseWriter, r *http.Request) {
	s.writeJSON(w, http.StatusAccepted, map[string]string{"status": "Stopping VM..."})
}

func (s *APIServer) handleRebootVM(w http.ResponseWriter, r *http.Request) {
	s.writeJSON(w, http.StatusAccepted, map[string]string{"status": "Rebooting VM..."})
}

func (s *APIServer) handleListGPUs(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	gpus, err := s.gpuOrch.GetGPUAvailability(ctx)
	if err != nil {
		s.writeError(w, err)
		return
	}
	for _, gpu := range gpus {
		if gpu == nil || gpu.Metrics == nil {
			continue
		}
		s.metrics.ObserveGPUUsage(
			gpu.ID,
			gpu.HostNodeID,
			gpu.Model,
			float64(gpu.Metrics.Utilization),
			int64(gpu.Metrics.MemoryUsed)*1024*1024,
			int64(gpu.Metrics.MemoryFree)*1024*1024,
			float64(gpu.Metrics.TemperatureC),
			float64(gpu.Metrics.PowerDraw),
		)
	}
	s.writeJSON(w, http.StatusOK, gpus)
}

func (s *APIServer) handleGetGPU(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	gpuID := vars["gpuId"]
	ctx := r.Context()
	if s.gpuOrch == nil {
		s.writeJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "gpu orchestrator unavailable"})
		return
	}

	gpu, err := s.gpuOrch.GetGPUByID(ctx, gpuID)
	if err != nil {
		s.writeError(w, err)
		return
	}
	if gpu != nil && gpu.Metrics != nil {
		s.metrics.ObserveGPUUsage(
			gpu.ID,
			gpu.HostNodeID,
			gpu.Model,
			float64(gpu.Metrics.Utilization),
			int64(gpu.Metrics.MemoryUsed)*1024*1024,
			int64(gpu.Metrics.MemoryFree)*1024*1024,
			float64(gpu.Metrics.TemperatureC),
			float64(gpu.Metrics.PowerDraw),
		)
	}
	s.writeJSON(w, http.StatusOK, gpu)
}

func (s *APIServer) handleListHosts(w http.ResponseWriter, r *http.Request) {
	// Implementation stub
	s.writeJSON(w, http.StatusOK, []map[string]string{})
}

func (s *APIServer) handleGetHost(w http.ResponseWriter, r *http.Request) {
	// Implementation stub
	s.writeJSON(w, http.StatusOK, map[string]string{"status": "Getting host..."})
}

func (s *APIServer) handleGetHostMetrics(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	nodeID := vars["nodeId"]
	ctx := r.Context()
	if s.resMonitor == nil {
		s.writeJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "resource monitor unavailable"})
		return
	}

	metrics, err := s.resMonitor.GetNodeMetrics(ctx, nodeID)
	if err != nil {
		s.writeError(w, err)
		return
	}
	s.writeJSON(w, http.StatusOK, metrics)
}

func (s *APIServer) handleGetTask(w http.ResponseWriter, r *http.Request) {
	// Implementation stub
	s.writeJSON(w, http.StatusOK, map[string]string{"status": "Getting task..."})
}

func (s *APIServer) handleClusterMetrics(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	metrics, err := s.resMonitor.GetClusterMetrics(ctx)
	if err != nil {
		s.writeError(w, err)
		return
	}
	s.writeJSON(w, http.StatusOK, metrics)
}

func (s *APIServer) handleClusterEvents(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.logger.WithError(err).Error("WebSocket upgrade failed")
		return
	}
	defer conn.Close()

	ctx := r.Context()
	s.logger.Debug("Client connected to event stream")
	<-ctx.Done()
}

func (s *APIServer) handleVMMetricsStream(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.logger.WithError(err).Error("WebSocket upgrade failed")
		return
	}
	defer conn.Close()

	ctx := r.Context()
	<-ctx.Done()
}

// Middleware

func (s *APIServer) observabilityMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		requestID := r.Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = uuid.NewString()
		}

		routePath := routeTemplate(r)
		tracer := otel.Tracer("ai-hypervisor-platform/api")
		ctx, span := tracer.Start(
			r.Context(),
			routePath,
			trace.WithSpanKind(trace.SpanKindServer),
			trace.WithAttributes(
				attribute.String("http.request_id", requestID),
				attribute.String("http.method", r.Method),
				attribute.String("http.route", routePath),
				attribute.String("http.target", r.URL.Path),
			),
		)
		defer span.End()

		recorder := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		recorder.Header().Set("X-Request-ID", requestID)
		next.ServeHTTP(recorder, r.WithContext(ctx))

		duration := time.Since(start)
		status := recorder.status
		span.SetAttributes(
			attribute.Int("http.status_code", status),
			attribute.Float64("http.duration_ms", float64(duration.Milliseconds())),
		)
		if status >= http.StatusInternalServerError {
			span.SetStatus(codes.Error, http.StatusText(status))
		} else {
			span.SetStatus(codes.Ok, http.StatusText(status))
		}
		if s.metrics != nil {
			s.metrics.ObserveRequest(r.Method, routePath, status, duration)
		}

		fields := logrus.Fields{
			"request_id":  requestID,
			"trace_id":    span.SpanContext().TraceID().String(),
			"method":      r.Method,
			"path":        routePath,
			"status":      status,
			"duration_ms": duration.Milliseconds(),
			"remote_addr": r.RemoteAddr,
		}
		if status >= http.StatusInternalServerError {
			s.logger.WithFields(fields).Error("request completed")
			return
		}
		s.logger.WithFields(fields).Info("request completed")
	})
}

func (s *APIServer) errorRecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				s.logger.Errorf("Panic: %v", err)
				s.writeJSON(w, http.StatusInternalServerError, map[string]string{
					"error": "Internal server error",
				})
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// Utilities

func (s *APIServer) writeJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	if s.baseHandler == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		return
	}
	s.baseHandler.WriteJSON(w, statusCode, data)
}

func (s *APIServer) writeError(w http.ResponseWriter, err error) {
	if s.baseHandler == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	s.baseHandler.WriteError(w, err)
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

func routeTemplate(r *http.Request) string {
	if route := mux.CurrentRoute(r); route != nil {
		if template, err := route.GetPathTemplate(); err == nil && template != "" {
			return template
		}
	}
	return r.URL.Path
}
