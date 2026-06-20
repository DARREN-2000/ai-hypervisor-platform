package scheduler

import (
	"context"
	"sort"
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/DARREN-2000/ai-hypervisor-platform/internal/config"
	"github.com/DARREN-2000/ai-hypervisor-platform/internal/models"
	"github.com/DARREN-2000/ai-hypervisor-platform/internal/orchestrator"
	"github.com/DARREN-2000/ai-hypervisor-platform/internal/storage"
	apierrors "github.com/DARREN-2000/ai-hypervisor-platform/pkg/errors"
)

// Dependencies bundles scheduler dependencies.
type Dependencies struct {
	HostRepo storage.HostRepository
	GPURepo  storage.GPURepository
	VMRepo   storage.VMRepository
	ResourceMonitor orchestrator.ResourceMonitor
	Logger   *logrus.Logger
}

// Service implements GPU-aware scheduling policies.
type Service struct {
	cfg         config.SchedulerConfig
	hostRepo    storage.HostRepository
	gpuRepo     storage.GPURepository
	vmRepo      storage.VMRepository
	resMonitor  orchestrator.ResourceMonitor
	logger      *logrus.Logger
	policyMu    sync.RWMutex
	policies    map[string]Policy
	monitor     Monitor
	autoscaler  AutoscalerHook
	metrics     schedulerMetrics
}

// Option configures the scheduler service.
type Option func(*Service)

// WithMonitor sets a scheduling monitor.
func WithMonitor(m Monitor) Option {
	return func(s *Service) {
		if m != nil {
			s.monitor = m
		}
	}
}

// WithAutoscaler sets an autoscaler hook.
func WithAutoscaler(a AutoscalerHook) Option {
	return func(s *Service) {
		if a != nil {
			s.autoscaler = a
		}
	}
}

// WithPolicy registers a policy implementation.
func WithPolicy(p Policy) Option {
	return func(s *Service) {
		if p == nil {
			return
		}
		s.policyMu.Lock()
		s.policies[p.Name()] = p
		s.policyMu.Unlock()
	}
}

// NewService creates a scheduler service with default policies.
func NewService(cfg config.SchedulerConfig, deps Dependencies, opts ...Option) *Service {
	logger := deps.Logger
	if logger == nil {
		logger = logrus.New()
	}

	svc := &Service{
		cfg:        cfg,
		hostRepo:   deps.HostRepo,
		gpuRepo:    deps.GPURepo,
		vmRepo:     deps.VMRepo,
		resMonitor: deps.ResourceMonitor,
		logger:     logger,
		policies:   make(map[string]Policy),
		monitor:    NoopMonitor{},
		autoscaler: NoopAutoscaler{},
	}

	// Built-in policies
	svc.policies[BinPackPolicy{}.Name()] = BinPackPolicy{}
	svc.policies[SpreadPolicy{}.Name()] = SpreadPolicy{}
	svc.policies[NUMAAwarePolicy{}.Name()] = NUMAAwarePolicy{}

	for _, opt := range opts {
		if opt != nil {
			opt(svc)
		}
	}

	return svc
}

// ScheduleVM makes a placement decision for a VM.
func (s *Service) ScheduleVM(ctx context.Context, vm *models.VirtualMachine) (*models.SchedulingDecision, error) {
	start := time.Now()
	decision, err := s.schedule(ctx, vm, "")
	latency := time.Since(start)
	if err != nil {
		s.metrics.recordDecision(latency.Milliseconds(), false)
		s.monitor.RecordFailure(ctx, FailureRecord{
			VMID:      vmID(vm),
			Reason:    err.Error(),
			Latency:   latency,
			Timestamp: time.Now().UTC(),
		})
		return nil, err
	}

	s.metrics.recordDecision(latency.Milliseconds(), true)
	if decision != nil {
		s.monitor.RecordDecision(ctx, DecisionRecord{
			VMID:      decision.VMID,
			HostID:    decision.SelectedHostID,
			Policy:    decision.Policy,
			Score:     decision.Score,
			Reason:    decision.Reason,
			GPUIds:    decision.SelectedGPUIds,
			Latency:   latency,
			Timestamp: decision.DecisionTimestamp,
		})
	}

	return decision, nil
}

// RescheduleVM attempts to reschedule a VM to a different host.
func (s *Service) RescheduleVM(ctx context.Context, vmID string) (*models.SchedulingDecision, error) {
	if vmID == "" {
		return nil, apierrors.ValidationError("vmID is required")
	}
	if s.vmRepo == nil {
		return nil, apierrors.InternalError("vm repository is not configured")
	}

	vm, err := s.vmRepo.Get(ctx, vmID)
	if err != nil {
		return nil, apierrors.NotFound("vm").WithCause(err)
	}

	return s.schedule(ctx, vm, vm.HostNodeID)
}

// CheckNodeCapacity checks if a node can accommodate a VM flavor.
func (s *Service) CheckNodeCapacity(ctx context.Context, nodeID string, flavor models.VMFlavor) bool {
	if nodeID == "" || s.hostRepo == nil {
		return false
	}

	host, err := s.hostRepo.Get(ctx, nodeID)
	if err != nil || host == nil {
		return false
	}

	demand := ResourceDemand{CPU: flavor.CPU, MemoryGB: flavor.Memory, DiskGB: flavor.DiskSize}
	return s.fitsResources(host, demand)
}

// GetSchedulingMetrics returns current scheduling metrics.
func (s *Service) GetSchedulingMetrics(ctx context.Context) (*orchestrator.SchedulingMetrics, error) {
	return s.metrics.snapshot(), nil
}

func (s *Service) schedule(ctx context.Context, vm *models.VirtualMachine, excludeHostID string) (*models.SchedulingDecision, error) {
	if vm == nil {
		return nil, apierrors.ValidationError("vm is required")
	}
	if s.hostRepo == nil {
		return nil, apierrors.InternalError("host repository is not configured")
	}

	hosts, err := s.hostRepo.List(ctx, nil)
	if err != nil {
		return nil, apierrors.InternalError("failed to list hosts").WithCause(err)
	}

	policy := s.getPolicy(s.cfg.Algorithm)
	candidates := make([]candidate, 0, len(hosts))
	demand := calculateDemand(vm)

	for _, host := range hosts {
		if host == nil || host.Status != models.HostStatusReady {
			continue
		}
		if excludeHostID != "" && host.ID == excludeHostID {
			continue
		}

		// Bolt optimization: perform fast local resource check before expensive snapshot generation (which makes DB/API calls)
		if !s.fitsResources(host, demand) {
			continue
		}
		snapshot := s.buildSnapshot(ctx, host)

		selectedGPUs, err := selectGPUsForVM(vm.GPURequests, snapshot.GPUs, policy.GPUSelectionStrategy())
		if err != nil {
			continue
		}

		scoreResult, err := policy.Score(ScoreInput{
			VM:        vm,
			Snapshot:  snapshot,
			Demand:    demand,
			GPUChoice: selectedGPUs,
			Now:       time.Now().UTC(),
		})
		if err != nil {
			continue
		}

		candidates = append(candidates, candidate{
			host:     host,
			gpus:     selectedGPUs,
			score:    scoreResult.Score,
			reason:   scoreResult.Reason,
			metadata: scoreResult.Metadata,
		})
	}

	if len(candidates) == 0 {
		if s.autoscaler != nil {
			s.autoscaler.Notify(ctx, AutoscaleSignal{
				Reason:    "no_capacity",
				CPUShortage: demand.CPU > 0,
				MemoryShortage: demand.MemoryGB > 0,
				GPUShortage: demand.GPUSlots > 0,
				Demand:    demand,
				Timestamp: time.Now().UTC(),
			})
		}
		return nil, apierrors.InsufficientResourcesError("capacity")
	}

	sort.SliceStable(candidates, func(i, j int) bool {
		if candidates[i].score == candidates[j].score {
			return candidates[i].host.ID < candidates[j].host.ID
		}
		return candidates[i].score > candidates[j].score
	})

	winner := candidates[0]
	altHosts := make([]models.HostScore, 0, minInt(5, len(candidates)-1))
	for i := 1; i < len(candidates) && i <= 5; i++ {
		altHosts = append(altHosts, models.HostScore{
			HostID: candidates[i].host.ID,
			Score:  candidates[i].score,
			Reason: candidates[i].reason,
		})
	}

	selectedGPUIds := make([]string, 0, len(winner.gpus))
	gpuAssignments := make([]models.GPUAllocation, 0, len(winner.gpus))
	for _, gpu := range winner.gpus {
		selectedGPUIds = append(selectedGPUIds, gpu.ID)
		gpuAssignments = append(gpuAssignments, models.GPUAllocation{
			ID:         models.NewID(),
			GPUID:      gpu.ID,
			VMName:     vm.Name,
			HostNodeID: winner.host.ID,
			AllocatedAt: time.Now().UTC(),
			DeviceIndex: gpu.Index,
			PCIAddress:  gpu.PCIAddress,
		})
	}

	decision := &models.SchedulingDecision{
		VMID:              vm.ID,
		SelectedHostID:    winner.host.ID,
		SelectedGPUIds:    selectedGPUIds,
		Policy:            policy.Name(),
		Score:             winner.score,
		Reason:            winner.reason,
		Metadata:          winner.metadata,
		GPUAssignments:    gpuAssignments,
		AlternativeHosts:  altHosts,
		DecisionTimestamp: time.Now().UTC(),
	}

	return decision, nil
}

func (s *Service) buildSnapshot(ctx context.Context, host *models.HostNode) HostSnapshot {
	snapshot := HostSnapshot{Host: host}
	if s.gpuRepo != nil && host != nil {
		gpus, err := s.gpuRepo.List(ctx, map[string]string{"host_id": host.ID})
		if err == nil {
			snapshot.GPUs = gpus
		}
	}
	if s.resMonitor != nil && host != nil {
		metrics, err := s.resMonitor.GetNodeMetrics(ctx, host.ID)
		if err == nil {
			snapshot.Metrics = metrics
		}
	}
	return snapshot
}

func (s *Service) fitsResources(host *models.HostNode, demand ResourceDemand) bool {
	if host == nil {
		return false
	}

	cpuCapacity := applyOvercommit(host.Capacity.CPU, s.cfg.AllowOvercommit, s.cfg.OvercommitRatio)
	memCapacity := applyOvercommit(host.Capacity.Memory, s.cfg.AllowOvercommit, s.cfg.OvercommitRatio)

	cpuAvailable := cpuCapacity - float64(host.AllocatedResources.CPU)
	memAvailable := memCapacity - float64(host.AllocatedResources.Memory)
	diskAvailable := host.Capacity.DiskGB - host.AllocatedResources.DiskGB
	gpuAvailable := host.Capacity.GPUSlots - host.AllocatedResources.GPUSlots

	if float64(demand.CPU) > cpuAvailable {
		return false
	}
	if float64(demand.MemoryGB) > memAvailable {
		return false
	}
	if demand.DiskGB > diskAvailable {
		return false
	}
	if demand.GPUSlots > gpuAvailable {
		return false
	}

	return true
}

func (s *Service) getPolicy(name string) Policy {
	policyName := name
	if policyName == "" {
		policyName = BinPackPolicy{}.Name()
	}

	s.policyMu.RLock()
	policy, ok := s.policies[policyName]
	if !ok {
		policy = s.policies[BinPackPolicy{}.Name()]
	}
	s.policyMu.RUnlock()

	if policy == nil {
		policy = BinPackPolicy{}
	}
	return policy
}

func calculateDemand(vm *models.VirtualMachine) ResourceDemand {
	demand := ResourceDemand{}
	if vm == nil {
		return demand
	}
	if vm.Flavor.CPU > 0 {
		demand.CPU = vm.Flavor.CPU
	}
	if vm.Flavor.Memory > 0 {
		demand.MemoryGB = vm.Flavor.Memory
	}
	if vm.Flavor.DiskSize > 0 {
		demand.DiskGB = vm.Flavor.DiskSize
	}
	for _, req := range vm.GPURequests {
		if req.Count > 0 {
			demand.GPUSlots += req.Count
		}
		if req.MinMemoryGB > 0 {
			demand.GPUMemoryGB += req.MinMemoryGB * req.Count
		}
	}
	return demand
}

func applyOvercommit(capacity int, allow bool, ratio float64) float64 {
	if !allow {
		return float64(capacity)
	}
	if ratio <= 0 {
		ratio = 1.0
	}
	return float64(capacity) * ratio
}

func vmID(vm *models.VirtualMachine) string {
	if vm == nil {
		return ""
	}
	return vm.ID
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

type candidate struct {
	host     *models.HostNode
	gpus     []*models.GPU
	score    float64
	reason   string
	metadata map[string]string
}

