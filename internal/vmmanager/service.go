package vmmanager

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/DARREN-2000/ai-hypervisor-platform/internal/config"
	libvirtutil "github.com/DARREN-2000/ai-hypervisor-platform/internal/libvirt"
	"github.com/DARREN-2000/ai-hypervisor-platform/internal/models"
	"github.com/DARREN-2000/ai-hypervisor-platform/internal/orchestrator"
	"github.com/DARREN-2000/ai-hypervisor-platform/internal/storage"
	apierrors "github.com/DARREN-2000/ai-hypervisor-platform/pkg/errors"
)

// Dependencies bundles dependencies for the VM manager service.
type Dependencies struct {
	VMRepo        storage.VMRepository
	HostRepo      storage.HostRepository
	TemplateRepo  storage.TemplateRepository
	Libvirt       orchestrator.LibvirtClient
	Scheduler     orchestrator.Scheduler
	GPUOrch       orchestrator.GPUOrchestrator
	TaskExecutor  orchestrator.TaskExecutor
	Logger        *logrus.Logger
}

// Service implements VM lifecycle orchestration.
type Service struct {
	cfg          config.VMManagerConfig
	vmRepo       storage.VMRepository
	hostRepo     storage.HostRepository
	templateRepo storage.TemplateRepository
	libvirt      orchestrator.LibvirtClient
	scheduler    orchestrator.Scheduler
	gpuOrch      orchestrator.GPUOrchestrator
	taskExecutor orchestrator.TaskExecutor
	logger       *logrus.Logger
	locks        *vmLockManager
	async        *asyncRunner

	monitorMu     sync.Mutex
	monitorCancel context.CancelFunc
	monitorWG     sync.WaitGroup
}

// NewService creates a new VM manager service.
func NewService(cfg config.VMManagerConfig, deps Dependencies) *Service {
	logger := deps.Logger
	if logger == nil {
		logger = logrus.New()
	}

	workerCount := cfg.MaxConcurrentProvisioning
	if workerCount <= 0 {
		workerCount = 1
	}

	return &Service{
		cfg:          cfg,
		vmRepo:       deps.VMRepo,
		hostRepo:     deps.HostRepo,
		templateRepo: deps.TemplateRepo,
		libvirt:      deps.Libvirt,
		scheduler:    deps.Scheduler,
		gpuOrch:      deps.GPUOrch,
		taskExecutor: deps.TaskExecutor,
		logger:       logger,
		locks:        newVMLockManager(),
		async:        newAsyncRunner(workerCount, logger),
	}
}

// Close stops background workers.
func (s *Service) Close() {
	s.StopMonitoring()
	if s.async != nil {
		s.async.Close()
	}
}

// CreateVM provisions a new virtual machine.
func (s *Service) CreateVM(ctx context.Context, vm *models.VirtualMachine) error {
	if err := s.ensureDeps(); err != nil {
		return err
	}
	if vm == nil {
		return apierrors.ValidationError("vm is required")
	}
	if vm.ID == "" {
		vm.ID = models.NewID()
	}

	unlock := s.locks.Lock(vm.ID)
	defer unlock()

	if err := validateVM(vm); err != nil {
		return err
	}

	now := time.Now().UTC()
	if vm.CreatedAt.IsZero() {
		vm.CreatedAt = now
	}
	vm.UpdatedAt = now
	vm.LastStatusCheck = now
	vm.State = models.VMStateProvisioning
	if vm.Metadata == nil {
		vm.Metadata = make(map[string]string)
	}

	if err := s.allocateResources(ctx, vm); err != nil {
		return err
	}

	if err := s.vmRepo.Create(ctx, vm); err != nil {
		return apierrors.InternalError("failed to persist vm").WithCause(err)
	}

	if err := s.reserveHostResources(ctx, vm); err != nil {
		return err
	}

	task := s.newTask(models.TaskTypeProvisionVM, vm)
	if task != nil {
		vm.Metadata["task_id"] = task.ID
		_ = s.vmRepo.Update(ctx, vm)
	}

	return s.enqueueAsync(ctx, task, "provision_vm", func(runCtx context.Context) error {
		return s.provisionVM(runCtx, vm.ID, task)
	})
}

// DeleteVM destroys a virtual machine and releases resources.
func (s *Service) DeleteVM(ctx context.Context, vmID string) error {
	if err := s.ensureDeps(); err != nil {
		return err
	}
	if vmID == "" {
		return apierrors.ValidationError("vmID is required")
	}

	unlock := s.locks.Lock(vmID)
	defer unlock()

	vm, err := s.getVM(ctx, vmID)
	if err != nil {
		return err
	}

	vm.State = models.VMStateDeleting
	vm.UpdatedAt = time.Now().UTC()
	if err := s.vmRepo.Update(ctx, vm); err != nil {
		return apierrors.InternalError("failed to mark vm deleting").WithCause(err)
	}

	task := s.newTask(models.TaskTypeDeleteVM, vm)
	if task != nil {
		vm.Metadata["task_id"] = task.ID
		_ = s.vmRepo.Update(ctx, vm)
	}

	return s.enqueueAsync(ctx, task, "delete_vm", func(runCtx context.Context) error {
		return s.destroyVM(runCtx, vmID, task)
	})
}

// GetVM retrieves a VM by ID.
func (s *Service) GetVM(ctx context.Context, vmID string) (*models.VirtualMachine, error) {
	if vmID == "" {
		return nil, apierrors.ValidationError("vmID is required")
	}
	return s.getVM(ctx, vmID)
}

// ListVMs lists VMs optionally filtered by namespace or state.
func (s *Service) ListVMs(ctx context.Context, filters map[string]string) ([]*models.VirtualMachine, error) {
	if s.vmRepo == nil {
		return nil, apierrors.InternalError("vm repository is not configured")
	}

	vms, err := s.vmRepo.List(ctx, filters)
	if err != nil {
		return nil, apierrors.InternalError("failed to list vms").WithCause(err)
	}
	return vms, nil
}

// UpdateVMState updates a VM state.
func (s *Service) UpdateVMState(ctx context.Context, vmID string, state models.VMState) error {
	if vmID == "" {
		return apierrors.ValidationError("vmID is required")
	}

	unlock := s.locks.Lock(vmID)
	defer unlock()

	vm, err := s.getVM(ctx, vmID)
	if err != nil {
		return err
	}

	vm.State = state
	vm.UpdatedAt = time.Now().UTC()
	if err := s.vmRepo.Update(ctx, vm); err != nil {
		return apierrors.InternalError("failed to update vm state").WithCause(err)
	}

	return nil
}

// StartVM powers on a VM.
func (s *Service) StartVM(ctx context.Context, vmID string) error {
	if err := s.ensureDeps(); err != nil {
		return err
	}

	unlock := s.locks.Lock(vmID)
	defer unlock()

	vm, err := s.getVM(ctx, vmID)
	if err != nil {
		return err
	}

	switch vm.State {
	case models.VMStateRunning:
		return apierrors.InvalidStateError("vm", string(vm.State), string(models.VMStateRunning))
	case models.VMStatePaused:
		err = s.libvirt.ResumeDomain(ctx, vm.Name)
	default:
		err = s.libvirt.StartDomain(ctx, vm.Name)
	}

	if err != nil {
		return apierrors.InternalError("failed to start vm").WithCause(err)
	}

	now := time.Now().UTC()
	vm.State = models.VMStateRunning
	vm.UpdatedAt = now
	if vm.StartedAt == nil {
		vm.StartedAt = &now
	}

	if err := s.vmRepo.Update(ctx, vm); err != nil {
		return apierrors.InternalError("failed to update vm status").WithCause(err)
	}

	return nil
}

// StopVM powers off a VM.
func (s *Service) StopVM(ctx context.Context, vmID string) error {
	if err := s.ensureDeps(); err != nil {
		return err
	}

	unlock := s.locks.Lock(vmID)
	defer unlock()

	vm, err := s.getVM(ctx, vmID)
	if err != nil {
		return err
	}

	if vm.State != models.VMStateRunning && vm.State != models.VMStatePaused {
		return apierrors.InvalidStateError("vm", string(vm.State), string(models.VMStateCreated))
	}

	if err := s.libvirt.StopDomain(ctx, vm.Name); err != nil {
		return apierrors.InternalError("failed to stop vm").WithCause(err)
	}

	now := time.Now().UTC()
	vm.State = models.VMStateCreated
	vm.UpdatedAt = now
	vm.TerminatedAt = &now

	if err := s.vmRepo.Update(ctx, vm); err != nil {
		return apierrors.InternalError("failed to update vm status").WithCause(err)
	}

	return nil
}

// PauseVM pauses a running VM.
func (s *Service) PauseVM(ctx context.Context, vmID string) error {
	if err := s.ensureDeps(); err != nil {
		return err
	}

	unlock := s.locks.Lock(vmID)
	defer unlock()

	vm, err := s.getVM(ctx, vmID)
	if err != nil {
		return err
	}

	if vm.State != models.VMStateRunning {
		return apierrors.InvalidStateError("vm", string(vm.State), string(models.VMStatePaused))
	}

	if err := s.libvirt.SuspendDomain(ctx, vm.Name); err != nil {
		return apierrors.InternalError("failed to pause vm").WithCause(err)
	}

	vm.State = models.VMStatePaused
	vm.UpdatedAt = time.Now().UTC()

	if err := s.vmRepo.Update(ctx, vm); err != nil {
		return apierrors.InternalError("failed to update vm status").WithCause(err)
	}

	return nil
}

// RebootVM reboots a running VM.
func (s *Service) RebootVM(ctx context.Context, vmID string) error {
	if err := s.ensureDeps(); err != nil {
		return err
	}

	unlock := s.locks.Lock(vmID)
	defer unlock()

	vm, err := s.getVM(ctx, vmID)
	if err != nil {
		return err
	}

	if vm.State != models.VMStateRunning {
		return apierrors.InvalidStateError("vm", string(vm.State), string(models.VMStateRunning))
	}

	if err := s.libvirt.RebootDomain(ctx, vm.Name); err != nil {
		return apierrors.InternalError("failed to reboot vm").WithCause(err)
	}

	vm.UpdatedAt = time.Now().UTC()
	if err := s.vmRepo.Update(ctx, vm); err != nil {
		return apierrors.InternalError("failed to update vm status").WithCause(err)
	}

	return nil
}

// StartMonitoring starts VM status monitoring.
func (s *Service) StartMonitoring(ctx context.Context) {
	interval := s.cfg.StateCheckInterval
	if interval <= 0 {
		interval = 30 * time.Second
	}

	s.monitorMu.Lock()
	defer s.monitorMu.Unlock()
	if s.monitorCancel != nil {
		return
	}

	monitorCtx, cancel := context.WithCancel(ctx)
	s.monitorCancel = cancel
	s.monitorWG.Add(1)

	go func() {
		defer s.monitorWG.Done()
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-monitorCtx.Done():
				return
			case <-ticker.C:
				s.syncVMStates(monitorCtx)
			}
		}
	}()
}

// StopMonitoring stops VM status monitoring.
func (s *Service) StopMonitoring() {
	s.monitorMu.Lock()
	cancel := s.monitorCancel
	s.monitorCancel = nil
	s.monitorMu.Unlock()

	if cancel != nil {
		cancel()
		s.monitorWG.Wait()
	}
}

func (s *Service) ensureDeps() error {
	if s.vmRepo == nil {
		return apierrors.InternalError("vm repository is not configured")
	}
	if s.libvirt == nil {
		return apierrors.InternalError("libvirt client is not configured")
	}
	return nil
}

func (s *Service) getVM(ctx context.Context, vmID string) (*models.VirtualMachine, error) {
	if s.vmRepo == nil {
		return nil, apierrors.InternalError("vm repository is not configured")
	}
	vm, err := s.vmRepo.Get(ctx, vmID)
	if err != nil {
		return nil, apierrors.NotFound("vm").WithCause(err)
	}
	return vm, nil
}

func (s *Service) allocateResources(ctx context.Context, vm *models.VirtualMachine) error {
	if vm.HostNodeID == "" {
		if s.scheduler == nil {
			return apierrors.InternalError("scheduler is not configured")
		}
		decision, err := s.scheduler.ScheduleVM(ctx, vm)
		if err != nil {
			return apierrors.InternalError("failed to schedule vm").WithCause(err)
		}
		vm.HostNodeID = decision.SelectedHostID
	}

	if len(vm.GPURequests) > 0 {
		if s.gpuOrch == nil {
			return apierrors.InternalError("gpu orchestrator is not configured")
		}
		allocations, err := s.gpuOrch.AllocateGPU(ctx, vm.ID, vm.GPURequests)
		if err != nil {
			return apierrors.InternalError("failed to allocate gpus").WithCause(err)
		}
		vm.GPUAllocations = allocations
	}

	return nil
}

func (s *Service) provisionVM(ctx context.Context, vmID string, task *models.Task) error {
	ctx, cancel := withTimeout(ctx, s.cfg.ProvisioningTimeout)
	defer cancel()
	if err := s.updateTaskStatus(ctx, task, models.TaskStatusRunning, nil); err != nil {
		return err
	}

	vm, err := s.getVM(ctx, vmID)
	if err != nil {
		_ = s.updateTaskStatus(ctx, task, models.TaskStatusFailed, map[string]interface{}{"error": err.Error()})
		return err
	}

	domainXML, err := libvirtutil.BuildDomainXML(vm)
	if err != nil {
		_ = s.markVMFailed(ctx, vm, err)
		_ = s.updateTaskStatus(ctx, task, models.TaskStatusFailed, map[string]interface{}{"error": err.Error()})
		return err
	}

	if err := s.libvirt.DefineDomain(ctx, domainXML); err != nil {
		_ = s.markVMFailed(ctx, vm, err)
		_ = s.updateTaskStatus(ctx, task, models.TaskStatusFailed, map[string]interface{}{"error": err.Error()})
		return apierrors.InternalError("failed to define domain").WithCause(err)
	}

	now := time.Now().UTC()
	vm.State = models.VMStateCreated
	vm.UpdatedAt = now
	vm.LastStatusCheck = now
	vm.ErrorMessage = ""

	if err := s.vmRepo.Update(ctx, vm); err != nil {
		_ = s.updateTaskStatus(ctx, task, models.TaskStatusFailed, map[string]interface{}{"error": err.Error()})
		return apierrors.InternalError("failed to update vm state").WithCause(err)
	}

	_ = s.updateTaskStatus(ctx, task, models.TaskStatusCompleted, map[string]interface{}{"vm_id": vm.ID})
	return nil
}

func (s *Service) destroyVM(ctx context.Context, vmID string, task *models.Task) error {
	ctx, cancel := withTimeout(ctx, s.cfg.ProvisioningTimeout)
	defer cancel()
	_ = s.updateTaskStatus(ctx, task, models.TaskStatusRunning, nil)

	vm, err := s.getVM(ctx, vmID)
	if err != nil {
		_ = s.updateTaskStatus(ctx, task, models.TaskStatusFailed, map[string]interface{}{"error": err.Error()})
		return err
	}

	if err := s.libvirt.DestroyDomain(ctx, vm.Name); err != nil {
		s.logger.WithError(err).WithField("vm", vmID).Warn("Failed to destroy domain")
	}
	if err := s.libvirt.UndefineDomain(ctx, vm.Name); err != nil {
		s.logger.WithError(err).WithField("vm", vmID).Warn("Failed to undefine domain")
	}
	if s.gpuOrch != nil {
		if err := s.gpuOrch.DeallocateGPU(ctx, vm.ID); err != nil {
			s.logger.WithError(err).WithField("vm", vmID).Warn("Failed to deallocate gpus")
		}
	}
	if err := s.releaseHostResources(ctx, vm); err != nil {
		s.logger.WithError(err).WithField("vm", vmID).Warn("Failed to release host resources")
	}

	if err := s.vmRepo.Delete(ctx, vmID); err != nil {
		_ = s.updateTaskStatus(ctx, task, models.TaskStatusFailed, map[string]interface{}{"error": err.Error()})
		return apierrors.InternalError("failed to delete vm record").WithCause(err)
	}

	_ = s.updateTaskStatus(ctx, task, models.TaskStatusCompleted, map[string]interface{}{"vm_id": vmID})
	return nil
}

func (s *Service) markVMFailed(ctx context.Context, vm *models.VirtualMachine, cause error) error {
	if vm == nil {
		return nil
	}
	vm.State = models.VMStateFailed
	vm.UpdatedAt = time.Now().UTC()
	if cause != nil {
		vm.ErrorMessage = cause.Error()
	}
	return s.vmRepo.Update(ctx, vm)
}

func (s *Service) enqueueAsync(ctx context.Context, task *models.Task, name string, fn func(context.Context) error) error {
	if task != nil && s.taskExecutor != nil {
		if err := s.taskExecutor.EnqueueTask(ctx, task); err != nil {
			return apierrors.InternalError("failed to enqueue task").WithCause(err)
		}
	}

	if s.async == nil {
		return fn(ctx)
	}

	return s.async.Submit(ctx, name, fn)
}

func (s *Service) updateTaskStatus(ctx context.Context, task *models.Task, status models.TaskStatus, result map[string]interface{}) error {
	if task == nil || s.taskExecutor == nil {
		return nil
	}
	task.Status = status
	task.UpdatedAt = time.Now().UTC()
	return s.taskExecutor.UpdateTaskStatus(ctx, task.ID, status, result)
}

func (s *Service) syncVMStates(ctx context.Context) {
	if s.vmRepo == nil || s.libvirt == nil {
		return
	}

	vms, err := s.vmRepo.List(ctx, nil)
	if err != nil {
		s.logger.WithError(err).Warn("Failed to list VMs for monitoring")
		return
	}

	for _, vm := range vms {
		if vm == nil {
			continue
		}
		unlock := s.locks.Lock(vm.ID)
		state, err := s.fetchDomainState(ctx, vm.Name)
		if err != nil {
			unlock()
			s.logger.WithError(err).WithField("vm", vm.ID).Debug("Failed to fetch domain state")
			continue
		}
		if state != "" {
			vmState := mapDomainState(state)
			if vm.State != vmState {
				vm.State = vmState
				vm.UpdatedAt = time.Now().UTC()
				vm.LastStatusCheck = time.Now().UTC()
				_ = s.vmRepo.Update(ctx, vm)
			}
		}
		unlock()
	}
}

func (s *Service) fetchDomainState(ctx context.Context, domainName string) (string, error) {
	info, err := s.libvirt.GetDomainInfo(ctx, domainName)
	if err != nil {
		return "", err
	}
	return info.State, nil
}

func (s *Service) reserveHostResources(ctx context.Context, vm *models.VirtualMachine) error {
	if s.hostRepo == nil || vm == nil || vm.HostNodeID == "" {
		return nil
	}

	host, err := s.hostRepo.Get(ctx, vm.HostNodeID)
	if err != nil {
		return apierrors.InternalError("failed to load host for allocation").WithCause(err)
	}

	host.AllocatedResources.CPU += vm.Flavor.CPU
	host.AllocatedResources.Memory += vm.Flavor.Memory
	host.AllocatedResources.DiskGB += vm.Flavor.DiskSize
	host.AllocatedResources.GPUSlots += gpuSlotCount(vm)
	host.AllocatedResources.VMs++

	if err := s.hostRepo.Update(ctx, host); err != nil {
		return apierrors.InternalError("failed to update host allocation").WithCause(err)
	}

	return nil
}

func (s *Service) releaseHostResources(ctx context.Context, vm *models.VirtualMachine) error {
	if s.hostRepo == nil || vm == nil || vm.HostNodeID == "" {
		return nil
	}

	host, err := s.hostRepo.Get(ctx, vm.HostNodeID)
	if err != nil {
		return apierrors.InternalError("failed to load host for release").WithCause(err)
	}

	host.AllocatedResources.CPU = maxInt(0, host.AllocatedResources.CPU-vm.Flavor.CPU)
	host.AllocatedResources.Memory = maxInt(0, host.AllocatedResources.Memory-vm.Flavor.Memory)
	host.AllocatedResources.DiskGB = maxInt(0, host.AllocatedResources.DiskGB-vm.Flavor.DiskSize)
	host.AllocatedResources.GPUSlots = maxInt(0, host.AllocatedResources.GPUSlots-gpuSlotCount(vm))
	host.AllocatedResources.VMs = maxInt(0, host.AllocatedResources.VMs-1)

	if err := s.hostRepo.Update(ctx, host); err != nil {
		return apierrors.InternalError("failed to update host allocation").WithCause(err)
	}

	return nil
}

func validateVM(vm *models.VirtualMachine) error {
	if vm.Name == "" {
		return apierrors.ValidationError("vm name is required")
	}
	if vm.Namespace == "" {
		return apierrors.ValidationError("vm namespace is required")
	}
	if vm.Flavor.CPU <= 0 || vm.Flavor.Memory <= 0 {
		return apierrors.ValidationError("vm flavor is invalid")
	}
	if vm.Image.ID == "" {
		return apierrors.ValidationError("vm image is required")
	}
	return nil
}

func (s *Service) newTask(taskType models.TaskType, vm *models.VirtualMachine) *models.Task {
	if vm == nil {
		return nil
	}
	now := time.Now().UTC()
	return &models.Task{
		ID:         models.NewID(),
		Type:       taskType,
		Status:     models.TaskStatusQueued,
		VMRef:      &models.VMRef{ID: vm.ID, Name: vm.Name, Namespace: vm.Namespace},
		CreatedAt:  now,
		UpdatedAt:  now,
		MaxRetries: 3,
	}
}

func gpuSlotCount(vm *models.VirtualMachine) int {
	if vm == nil {
		return 0
	}
	if len(vm.GPUAllocations) > 0 {
		return len(vm.GPUAllocations)
	}
	count := 0
	for _, req := range vm.GPURequests {
		count += req.Count
	}
	return count
}

func withTimeout(ctx context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	if timeout <= 0 {
		return ctx, func() {}
	}
	return context.WithTimeout(ctx, timeout)
}

func mapDomainState(state string) models.VMState {
	switch strings.ToLower(state) {
	case "running":
		return models.VMStateRunning
	case "paused", "pmsuspended":
		return models.VMStatePaused
	case "shutdown", "shutoff":
		return models.VMStateCreated
	case "crashed":
		return models.VMStateFailed
	default:
		return models.VMStateFailed
	}
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// vmLockManager provides per-VM locks.
type vmLockManager struct {
	mu    sync.Mutex
	locks map[string]*vmLockEntry
}

type vmLockEntry struct {
	mu    sync.Mutex
	refs  int
}

func newVMLockManager() *vmLockManager {
	return &vmLockManager{locks: make(map[string]*vmLockEntry)}
}

func (m *vmLockManager) Lock(vmID string) func() {
	m.mu.Lock()
	entry := m.locks[vmID]
	if entry == nil {
		entry = &vmLockEntry{refs: 1}
		m.locks[vmID] = entry
	} else {
		entry.refs++
	}
	m.mu.Unlock()

	entry.mu.Lock()
	return func() {
		entry.mu.Unlock()
		m.mu.Lock()
		entry.refs--
		if entry.refs == 0 {
			delete(m.locks, vmID)
		}
		m.mu.Unlock()
	}
}
