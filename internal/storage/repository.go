package storage

import (
	"context"

	"github.com/DARREN-2000/ai-hypervisor-platform/internal/models"
)

// VMRepository provides persistence operations for virtual machines.
type VMRepository interface {
	Create(ctx context.Context, vm *models.VirtualMachine) error
	Update(ctx context.Context, vm *models.VirtualMachine) error
	Get(ctx context.Context, vmID string) (*models.VirtualMachine, error)
	List(ctx context.Context, filters map[string]string) ([]*models.VirtualMachine, error)
	Delete(ctx context.Context, vmID string) error
	UpdateState(ctx context.Context, vmID string, state models.VMState) error
}

// HostRepository provides persistence operations for hosts.
type HostRepository interface {
	Create(ctx context.Context, host *models.HostNode) error
	Update(ctx context.Context, host *models.HostNode) error
	Get(ctx context.Context, hostID string) (*models.HostNode, error)
	List(ctx context.Context, filters map[string]string) ([]*models.HostNode, error)
	Delete(ctx context.Context, hostID string) error
}

// GPURepository provides persistence operations for GPUs.
type GPURepository interface {
	Create(ctx context.Context, gpu *models.GPU) error
	Update(ctx context.Context, gpu *models.GPU) error
	Get(ctx context.Context, gpuID string) (*models.GPU, error)
	List(ctx context.Context, filters map[string]string) ([]*models.GPU, error)
	Delete(ctx context.Context, gpuID string) error
}

// TaskRepository provides persistence operations for tasks.
type TaskRepository interface {
	Create(ctx context.Context, task *models.Task) error
	Update(ctx context.Context, task *models.Task) error
	Get(ctx context.Context, taskID string) (*models.Task, error)
	List(ctx context.Context, filters map[string]string) ([]*models.Task, error)
	Delete(ctx context.Context, taskID string) error
}

// MetricsRepository provides access to resource metrics data.
type MetricsRepository interface {
	Create(ctx context.Context, metrics *models.ResourceMetrics) error
	GetForVM(ctx context.Context, vmID string, limit int) ([]*models.ResourceMetrics, error)
	GetForHost(ctx context.Context, hostID string, limit int) ([]*models.ResourceMetrics, error)
}

// EventRepository provides access to platform events.
type EventRepository interface {
	Create(ctx context.Context, event *models.Event) error
	List(ctx context.Context, filters map[string]string) ([]*models.Event, error)
}

// TemplateRepository provides persistence operations for VM templates.
type TemplateRepository interface {
	Create(ctx context.Context, template *models.VMTemplate) error
	Update(ctx context.Context, template *models.VMTemplate) error
	Get(ctx context.Context, templateID string) (*models.VMTemplate, error)
	List(ctx context.Context, filters map[string]string) ([]*models.VMTemplate, error)
	Delete(ctx context.Context, templateID string) error
}
