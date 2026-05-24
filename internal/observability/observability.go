package observability

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	otelsdk "go.opentelemetry.io/otel/sdk/trace"
	otlptracegrpc "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"

	"github.com/DARREN-2000/ai-hypervisor-platform/pkg/telemetry"
)

// Config controls metrics and tracing wiring for a service.
type Config struct {
	ServiceName    string
	ServiceVersion string
	Environment    string

	MetricsEnabled bool
	MetricsAddr    string
	MetricsPort    int

	TracingEnabled bool
	OTLPEndpoint   string
	OTLPInsecure   bool
}

// Manager owns the Prometheus registry, metrics server, and tracing provider.
type Manager struct {
	metricsRecorder telemetry.Metrics
	registry        *prometheus.Registry
	metricsServer   *http.Server
	tracerProvider  *otelsdk.TracerProvider
	logger          *logrus.Logger
	config          Config
}

// NewManager builds the observability stack.
func NewManager(ctx context.Context, cfg Config, logger *logrus.Logger) (*Manager, error) {
	if logger == nil {
		logger = logrus.New()
	}

	m := &Manager{
		logger: logger,
		config: cfg,
	}

	registry := prometheus.NewRegistry()
	if err := registry.Register(collectors.NewGoCollector()); err != nil {
		return nil, fmt.Errorf("register go collector: %w", err)
	}
	if err := registry.Register(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{})); err != nil {
		return nil, fmt.Errorf("register process collector: %w", err)
	}

	metrics := newPrometheusMetrics(registry, cfg.ServiceName, cfg.ServiceVersion, cfg.Environment)
	m.registry = registry
	m.metricsRecorder = metrics

	if cfg.TracingEnabled {
		tracerProvider, err := buildTracerProvider(ctx, cfg)
		if err != nil {
			return nil, err
		}
		m.tracerProvider = tracerProvider
		otel.SetTracerProvider(tracerProvider)
		otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		))
	}

	if cfg.MetricsEnabled {
		addr := cfg.MetricsAddr
		if addr == "" {
			addr = "0.0.0.0"
		}
		port := cfg.MetricsPort
		if port <= 0 {
			port = 8081
		}
		m.metricsServer = &http.Server{
			Addr: fmt.Sprintf("%s:%d", addr, port),
			Handler: promhttp.HandlerFor(registry, promhttp.HandlerOpts{
				EnableOpenMetrics: true,
			}),
		}
	}

	return m, nil
}

// Metrics returns the Prometheus-backed recorder.
func (m *Manager) Metrics() telemetry.Metrics {
	if m == nil || m.metricsRecorder == nil {
		return telemetry.NewNoopMetrics()
	}
	return m.metricsRecorder
}

// ServeMetrics starts the Prometheus HTTP endpoint.
func (m *Manager) ServeMetrics() error {
	if m == nil || m.metricsServer == nil {
		return nil
	}
	m.logger.Infof("Starting metrics server on %s", m.metricsServer.Addr)
	err := m.metricsServer.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

// Shutdown gracefully stops the metrics server and tracing provider.
func (m *Manager) Shutdown(ctx context.Context) error {
	if m == nil {
		return nil
	}

	var firstErr error
	if m.metricsServer != nil {
		if err := m.metricsServer.Shutdown(ctx); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	if m.tracerProvider != nil {
		shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		if err := m.tracerProvider.Shutdown(shutdownCtx); err != nil && firstErr == nil {
			firstErr = err
		}
	}

	return firstErr
}

func buildTracerProvider(ctx context.Context, cfg Config) (*otelsdk.TracerProvider, error) {
	attrs := []attribute.KeyValue{
		attribute.String("service.name", cfg.ServiceName),
		attribute.String("service.version", cfg.ServiceVersion),
		attribute.String("deployment.environment", cfg.Environment),
	}

	res, err := resource.New(ctx,
		resource.WithAttributes(attrs...),
	)
	if err != nil {
		return nil, fmt.Errorf("build resource: %w", err)
	}

	providerOptions := []otelsdk.TracerProviderOption{
		otelsdk.WithResource(res),
	}

	if cfg.OTLPEndpoint != "" {
		client := []otlptracegrpc.Option{
			otlptracegrpc.WithEndpoint(cfg.OTLPEndpoint),
		}
		if cfg.OTLPInsecure {
			client = append(client, otlptracegrpc.WithInsecure())
		}
		exporter, err := otlptracegrpc.New(ctx, client...)
		if err != nil {
			return nil, fmt.Errorf("build otlp trace exporter: %w", err)
		}
		providerOptions = append(providerOptions, otelsdk.WithBatcher(exporter))
	}

	return otelsdk.NewTracerProvider(providerOptions...), nil
}
