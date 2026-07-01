package scheduler

import (
	"github.com/DARREN-2000/ai-hypervisor-platform/internal/orchestrator"

	"context"
	"testing"

	"github.com/DARREN-2000/ai-hypervisor-platform/internal/config"
	"github.com/DARREN-2000/ai-hypervisor-platform/internal/models"
)

type mockHostRepo struct {
	hosts []*models.HostNode
	err   error
}

func (m *mockHostRepo) List(ctx context.Context, filter map[string]string) ([]*models.HostNode, error) {
	return m.hosts, m.err
}
func (m *mockHostRepo) Get(ctx context.Context, id string) (*models.HostNode, error) {
	if len(m.hosts) > 0 {
		return m.hosts[0], m.err
	}
	return nil, m.err
}
func (m *mockHostRepo) Create(ctx context.Context, host *models.HostNode) error { return nil }
func (m *mockHostRepo) Update(ctx context.Context, host *models.HostNode) error { return nil }
func (m *mockHostRepo) Delete(ctx context.Context, id string) error { return nil }

type mockVMRepo struct {
	vm  *models.VirtualMachine
	err error
}

func (m *mockVMRepo) List(ctx context.Context, filter map[string]string) ([]*models.VirtualMachine, error) {
	return nil, nil
}
func (m *mockVMRepo) Get(ctx context.Context, id string) (*models.VirtualMachine, error) {
	return m.vm, m.err
}
func (m *mockVMRepo) Create(ctx context.Context, vm *models.VirtualMachine) error { return nil }
func (m *mockVMRepo) Update(ctx context.Context, vm *models.VirtualMachine) error { return nil }
func (m *mockVMRepo) Delete(ctx context.Context, id string) error { return nil }
func (m *mockVMRepo) UpdateState(ctx context.Context, vmID string, state models.VMState) error { return nil }

func TestScheduleVM_Success(t *testing.T) {
	hostRepo := &mockHostRepo{
		hosts: []*models.HostNode{
			{
				ID:     "host-1",
				Status: models.HostStatusReady,
				Capacity: models.HostCapacity{
					CPU:      32,
					Memory:   128,
					DiskGB:   1000,
					GPUSlots: 4,
				},
				AllocatedResources: models.HostAllocated{
					CPU:      4,
					Memory:   16,
					DiskGB:   100,
					GPUSlots: 0,
				},
			},
		},
	}

	svc := NewService(config.SchedulerConfig{Algorithm: "bin-packing"}, Dependencies{
		HostRepo: hostRepo,
	})

	vm := &models.VirtualMachine{
		ID: "vm-1",
		Flavor: models.VMFlavor{
			CPU:      4,
			Memory:   16,
			DiskSize: 100,
		},
	}

	decision, err := svc.ScheduleVM(context.Background(), vm)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if decision == nil {
		t.Fatalf("expected decision, got nil")
	}
	if decision.SelectedHostID != "host-1" {
		t.Errorf("expected host-1, got %s", decision.SelectedHostID)
	}
}

func TestScheduleVM_NoCapacity(t *testing.T) {
	hostRepo := &mockHostRepo{
		hosts: []*models.HostNode{
			{
				ID:     "host-1",
				Status: models.HostStatusReady,
				Capacity: models.HostCapacity{
					CPU:      2,
					Memory:   4,
					DiskGB:   10,
					GPUSlots: 0,
				},
				AllocatedResources: models.HostAllocated{
					CPU:      2,
					Memory:   4,
					DiskGB:   10,
					GPUSlots: 0,
				},
			},
		},
	}

	svc := NewService(config.SchedulerConfig{Algorithm: "bin-packing"}, Dependencies{
		HostRepo: hostRepo,
	})

	vm := &models.VirtualMachine{
		ID: "vm-1",
		Flavor: models.VMFlavor{
			CPU:      4,
			Memory:   16,
			DiskSize: 100,
		},
	}

	decision, err := svc.ScheduleVM(context.Background(), vm)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if decision != nil {
		t.Fatalf("expected nil decision, got %v", decision)
	}
}

func TestRescheduleVM(t *testing.T) {
	vm := &models.VirtualMachine{
		ID: "vm-1",
		Flavor: models.VMFlavor{
			CPU:      4,
			Memory:   16,
			DiskSize: 100,
		},
		HostNodeID: "host-bad",
	}

	hostRepo := &mockHostRepo{
		hosts: []*models.HostNode{
			{
				ID:     "host-bad",
				Status: models.HostStatusReady,
				Capacity: models.HostCapacity{
					CPU:      32,
					Memory:   128,
					DiskGB:   1000,
					GPUSlots: 4,
				},
			},
			{
				ID:     "host-good",
				Status: models.HostStatusReady,
				Capacity: models.HostCapacity{
					CPU:      32,
					Memory:   128,
					DiskGB:   1000,
					GPUSlots: 4,
				},
			},
		},
	}

	vmRepo := &mockVMRepo{vm: vm}

	svc := NewService(config.SchedulerConfig{Algorithm: "bin-packing"}, Dependencies{
		HostRepo: hostRepo,
		VMRepo:   vmRepo,
	})

	decision, err := svc.RescheduleVM(context.Background(), "vm-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if decision.SelectedHostID != "host-good" {
		t.Errorf("expected host-good, got %s", decision.SelectedHostID)
	}
}

func TestCalculateDemand(t *testing.T) {
	d := calculateDemand(nil)
	if d.CPU != 0 {
		t.Error("expected 0 demand for nil vm")
	}

	vm := &models.VirtualMachine{
		Flavor: models.VMFlavor{
			CPU: 2, Memory: 4, DiskSize: 20,
		},
		GPURequests: []models.GPURequest{
			{Count: 2, MinMemoryGB: 16},
			{Count: 0},
		},
	}
	d = calculateDemand(vm)
	if d.CPU != 2 || d.MemoryGB != 4 || d.DiskGB != 20 || d.GPUSlots != 2 || d.GPUMemoryGB != 32 {
		t.Errorf("unexpected demand: %+v", d)
	}
}

func TestApplyOvercommit(t *testing.T) {
	if applyOvercommit(100, false, 2.0) != 100.0 {
		t.Error("expected no overcommit")
	}
	if applyOvercommit(100, true, 2.0) != 200.0 {
		t.Error("expected overcommit")
	}
	if applyOvercommit(100, true, -1.0) != 100.0 {
		t.Error("expected fallback to 1.0 ratio")
	}
}

func TestVMID(t *testing.T) {
	if vmID(nil) != "" {
		t.Error("expected empty string")
	}
	vm := &models.VirtualMachine{ID: "test"}
	if vmID(vm) != "test" {
		t.Error("expected test")
	}
}

func TestMinInt(t *testing.T) {
	if minInt(1, 2) != 1 {
		t.Error("expected 1")
	}
	if minInt(2, 1) != 1 {
		t.Error("expected 1")
	}
}

func TestGetPolicy(t *testing.T) {
	svc := NewService(config.SchedulerConfig{}, Dependencies{})
	p := svc.getPolicy("")
	if p.Name() != "bin-packing" {
		t.Error("expected bin-packing as default")
	}
	p = svc.getPolicy("unknown")
	if p.Name() != "bin-packing" {
		t.Error("expected bin-packing as fallback")
	}
	p = svc.getPolicy("spread")
	if p.Name() != "spread" {
		t.Error("expected spread")
	}
}

func TestCheckNodeCapacity_Success(t *testing.T) {
	hostRepo := &mockHostRepo{
		hosts: []*models.HostNode{
			{
				ID: "host-1",
				Capacity: models.HostCapacity{
					CPU: 10, Memory: 20, DiskGB: 100, GPUSlots: 1,
				},
			},
		},
	}
	svc := NewService(config.SchedulerConfig{}, Dependencies{HostRepo: hostRepo})
	if !svc.CheckNodeCapacity(context.Background(), "host-1", models.VMFlavor{CPU: 2, Memory: 4, DiskSize: 10}) {
		t.Error("expected true")
	}
}

func TestCheckNodeCapacity_MissingHost(t *testing.T) {
	hostRepo := &mockHostRepo{hosts: nil} // empty
	svc := NewService(config.SchedulerConfig{}, Dependencies{HostRepo: hostRepo})
	if svc.CheckNodeCapacity(context.Background(), "host-1", models.VMFlavor{}) {
		t.Error("expected false")
	}
}

func TestFitsResources(t *testing.T) {
	svc := NewService(config.SchedulerConfig{}, Dependencies{})
	if svc.fitsResources(nil, ResourceDemand{}) {
		t.Error("expected false for nil host")
	}

	host := &models.HostNode{
		Capacity: models.HostCapacity{CPU: 10, Memory: 20, DiskGB: 100, GPUSlots: 1},
		AllocatedResources: models.HostAllocated{CPU: 5, Memory: 10, DiskGB: 50, GPUSlots: 1},
	}
	if !svc.fitsResources(host, ResourceDemand{CPU: 5, MemoryGB: 10, DiskGB: 50, GPUSlots: 0}) {
		t.Error("expected true")
	}
	if svc.fitsResources(host, ResourceDemand{CPU: 6, MemoryGB: 10, DiskGB: 50, GPUSlots: 0}) {
		t.Error("expected false due to CPU")
	}
	if svc.fitsResources(host, ResourceDemand{CPU: 5, MemoryGB: 11, DiskGB: 50, GPUSlots: 0}) {
		t.Error("expected false due to Memory")
	}
	if svc.fitsResources(host, ResourceDemand{CPU: 5, MemoryGB: 10, DiskGB: 51, GPUSlots: 0}) {
		t.Error("expected false due to Disk")
	}
	if svc.fitsResources(host, ResourceDemand{CPU: 5, MemoryGB: 10, DiskGB: 50, GPUSlots: 1}) {
		t.Error("expected false due to GPU")
	}
}

type mockGPURepo struct {
	gpus []*models.GPU
}
func (m *mockGPURepo) List(ctx context.Context, filter map[string]string) ([]*models.GPU, error) {
	return m.gpus, nil
}
func (m *mockGPURepo) Get(ctx context.Context, id string) (*models.GPU, error) { return nil, nil }
func (m *mockGPURepo) Create(ctx context.Context, gpu *models.GPU) error { return nil }
func (m *mockGPURepo) Update(ctx context.Context, gpu *models.GPU) error { return nil }
func (m *mockGPURepo) Delete(ctx context.Context, id string) error { return nil }

type mockResourceMonitor struct {
	metrics *models.ResourceMetrics
}
func (m *mockResourceMonitor) GetNodeMetrics(ctx context.Context, id string) (*models.ResourceMetrics, error) {
	return m.metrics, nil
}
func (m *mockResourceMonitor) GetVMMetrics(ctx context.Context, id string) (*models.ResourceMetrics, error) {
	return nil, nil
}
func (m *mockResourceMonitor) GetClusterMetrics(ctx context.Context) (*orchestrator.ClusterMetrics, error) {
	return nil, nil
}
func (m *mockResourceMonitor) CollectMetrics(ctx context.Context, id string) error { return nil }
func (m *mockResourceMonitor) PredictResourceUsage(ctx context.Context, id string) (*models.ResourceMetrics, error) { return nil, nil }

func TestBuildSnapshot(t *testing.T) {
	gpuRepo := &mockGPURepo{gpus: []*models.GPU{{ID: "g1"}}}
	resMon := &mockResourceMonitor{metrics: &models.ResourceMetrics{CPU: 50.0}}

	svc := NewService(config.SchedulerConfig{}, Dependencies{
		GPURepo: gpuRepo,
		ResourceMonitor: resMon,
	})

	host := &models.HostNode{ID: "host-1"}
	snap := svc.buildSnapshot(context.Background(), host)

	if snap.Host != host {
		t.Error("host mismatch")
	}
	if len(snap.GPUs) != 1 || snap.GPUs[0].ID != "g1" {
		t.Error("gpus mismatch")
	}
	if snap.Metrics == nil || snap.Metrics.CPU != 50.0 {
		t.Error("metrics mismatch")
	}
}

func TestScheduleVM_InvalidInputs(t *testing.T) {
	svc := NewService(config.SchedulerConfig{}, Dependencies{})
	_, err := svc.ScheduleVM(context.Background(), nil)
	if err == nil {
		t.Error("expected error for nil vm")
	}
	vm := &models.VirtualMachine{ID: "vm-1"}
	_, err = svc.ScheduleVM(context.Background(), vm)
	if err == nil {
		t.Error("expected error for no host repo")
	}
}

func TestRescheduleVM_InvalidInputs(t *testing.T) {
	svc := NewService(config.SchedulerConfig{}, Dependencies{})
	_, err := svc.RescheduleVM(context.Background(), "")
	if err == nil {
		t.Error("expected error for empty vmID")
	}
	_, err = svc.RescheduleVM(context.Background(), "vm-1")
	if err == nil {
		t.Error("expected error for no vm repo")
	}
}
