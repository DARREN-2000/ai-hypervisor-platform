package orchestrator

import (
	"context"

	"github.com/DARREN-2000/ai-hypervisor-platform/internal/models"
)

// VMManager handles VM lifecycle management
type VMManager interface {
	// CreateVM creates a new virtual machine
	CreateVM(ctx context.Context, vm *models.VirtualMachine) error
	
	// DeleteVM removes a virtual machine
	DeleteVM(ctx context.Context, vmID string) error
	
	// GetVM retrieves VM details
	GetVM(ctx context.Context, vmID string) (*models.VirtualMachine, error)
	
	// ListVMs lists all VMs, optionally filtered
	ListVMs(ctx context.Context, filters map[string]string) ([]*models.VirtualMachine, error)
	
	// UpdateVMState transitions a VM to a new state
	UpdateVMState(ctx context.Context, vmID string, state models.VMState) error
	
	// StartVM starts a stopped VM
	StartVM(ctx context.Context, vmID string) error
	
	// StopVM stops a running VM
	StopVM(ctx context.Context, vmID string) error

	// PauseVM pauses a running VM
	PauseVM(ctx context.Context, vmID string) error
	
	// RebootVM reboots a running VM
	RebootVM(ctx context.Context, vmID string) error
}

// Scheduler handles VM placement decisions
type Scheduler interface {
	// ScheduleVM makes a placement decision for a VM
	ScheduleVM(ctx context.Context, vm *models.VirtualMachine) (*models.SchedulingDecision, error)
	
	// RescheduleVM attempts to reschedule a VM to a different host
	RescheduleVM(ctx context.Context, vmID string) (*models.SchedulingDecision, error)
	
	// CheckNodeCapacity checks if a node can accommodate a VM
	CheckNodeCapacity(ctx context.Context, nodeID string, resources models.VMFlavor) bool
	
	// GetSchedulingMetrics returns current scheduling metrics
	GetSchedulingMetrics(ctx context.Context) (*SchedulingMetrics, error)
}

// SchedulingMetrics represents scheduler performance metrics
type SchedulingMetrics struct {
	AverageDecisionLatency int     // Milliseconds
	FailedSchedulingAttempts int
	SuccessfulSchedules    int
	OptimalPlacementRate   float64 // 0-1
}

// GPUOrchestrator handles GPU allocation and management
type GPUOrchestrator interface {
	// AllocateGPU allocates GPUs to a VM
	AllocateGPU(ctx context.Context, vmID string, requests []models.GPURequest) ([]models.GPUAllocation, error)
	
	// DeallocateGPU releases allocated GPUs
	DeallocateGPU(ctx context.Context, vmID string) error
	
	// GetGPUAvailability returns available GPUs
	GetGPUAvailability(ctx context.Context) ([]*models.GPU, error)
	
	// GetGPUByID retrieves GPU details
	GetGPUByID(ctx context.Context, gpuID string) (*models.GPU, error)
	
	// UpdateGPUMetrics updates GPU telemetry data
	UpdateGPUMetrics(ctx context.Context, gpuID string, metrics *models.GPUMetrics) error
	
	// CheckGPUHealth performs health checks on GPUs
	CheckGPUHealth(ctx context.Context, nodeID string) error
}

// TaskExecutor handles async task execution
type TaskExecutor interface {
	// EnqueueTask adds a task to the execution queue
	EnqueueTask(ctx context.Context, task *models.Task) error
	
	// DequeueTask retrieves the next task to execute
	DequeueTask(ctx context.Context) (*models.Task, error)
	
	// UpdateTaskStatus updates task status
	UpdateTaskStatus(ctx context.Context, taskID string, status models.TaskStatus, result map[string]interface{}) error
	
	// GetTask retrieves task details
	GetTask(ctx context.Context, taskID string) (*models.Task, error)
	
	// ListTasks lists tasks, optionally filtered
	ListTasks(ctx context.Context, filters map[string]string) ([]*models.Task, error)
	
	// RetryTask retries a failed task
	RetryTask(ctx context.Context, taskID string) error
}

// ResourceMonitor tracks resource utilization
type ResourceMonitor interface {
	// GetNodeMetrics returns current metrics for a node
	GetNodeMetrics(ctx context.Context, nodeID string) (*models.ResourceMetrics, error)
	
	// GetVMMetrics returns current metrics for a VM
	GetVMMetrics(ctx context.Context, vmID string) (*models.ResourceMetrics, error)
	
	// GetClusterMetrics returns aggregated cluster metrics
	GetClusterMetrics(ctx context.Context) (*ClusterMetrics, error)
	
	// CollectMetrics triggers metric collection
	CollectMetrics(ctx context.Context, nodeID string) error
	
	// PredictResourceUsage predicts future resource usage
	PredictResourceUsage(ctx context.Context, vmID string) (*models.ResourceMetrics, error)
}

// ClusterMetrics represents aggregated cluster metrics
type ClusterMetrics struct {
	TotalCPU       int
	TotalMemoryGB  int
	AllocatedCPU   int
	AllocatedMemGB int
	UtilizationCPU float64     // 0-1
	UtilizationMem float64     // 0-1
	NodeCount      int
	VMCount        int
	GPUCount       int
	AllocatedGPUs  int
	Timestamp      int64
}

// HostAgent manages node-level operations
type HostAgent interface {
	// GetNodeStatus returns node operational status
	GetNodeStatus(ctx context.Context, nodeID string) (*models.HostNode, error)
	
	// CreateVMOnNode creates a VM on a specific node
	CreateVMOnNode(ctx context.Context, nodeID string, vm *models.VirtualMachine) error
	
	// DeleteVMOnNode removes a VM from a node
	DeleteVMOnNode(ctx context.Context, nodeID string, vmID string) error
	
	// GetVMStatus returns status of a VM on a node
	GetVMStatus(ctx context.Context, nodeID string, vmID string) (*models.VirtualMachine, error)
	
	// ListVMsOnNode returns all VMs on a node
	ListVMsOnNode(ctx context.Context, nodeID string) ([]*models.VirtualMachine, error)
	
	// GetGPUsOnNode returns GPUs on a node
	GetGPUsOnNode(ctx context.Context, nodeID string) ([]*models.GPU, error)
}

// LibvirtClient provides low-level VM operations via libvirt
type LibvirtClient interface {
	// CreateDomain creates a libvirt domain
	CreateDomain(ctx context.Context, domainXML string) error
	
	// DefineDomain defines (but doesn't start) a libvirt domain
	DefineDomain(ctx context.Context, domainXML string) error
	
	// DestroyDomain destroys a libvirt domain
	DestroyDomain(ctx context.Context, domainName string) error
	
	// StartDomain starts a libvirt domain
	StartDomain(ctx context.Context, domainName string) error
	
	// StopDomain stops a libvirt domain
	StopDomain(ctx context.Context, domainName string) error

	// SuspendDomain pauses a libvirt domain
	SuspendDomain(ctx context.Context, domainName string) error

	// ResumeDomain resumes a suspended domain
	ResumeDomain(ctx context.Context, domainName string) error
	
	// RebootDomain reboots a libvirt domain
	RebootDomain(ctx context.Context, domainName string) error
	
	// GetDomainInfo returns domain information
	GetDomainInfo(ctx context.Context, domainName string) (*LibvirtDomainInfo, error)
	
	// ListDomains returns all domains
	ListDomains(ctx context.Context) ([]string, error)
	
	// AttachDevice attaches a device to a domain
	AttachDevice(ctx context.Context, domainName string, deviceXML string) error
	
	// DetachDevice detaches a device from a domain
	DetachDevice(ctx context.Context, domainName string, deviceXML string) error

	// UndefineDomain removes a libvirt domain definition
	UndefineDomain(ctx context.Context, domainName string) error
}

// LibvirtDomainInfo represents libvirt domain information
type LibvirtDomainInfo struct {
	Name       string
	UUID       string
	MaxMemory  int64 // KB
	Memory     int64 // KB
	MaxCPU     int
	CPU        int
	State      string
	CPUTime    int64 // Nanoseconds
}

// ConfigManager handles application configuration
type ConfigManager interface {
	// GetVMFlavor retrieves a VM flavor by name
	GetVMFlavor(ctx context.Context, flavorName string) (*models.VMFlavor, error)
	
	// ListVMFlavors lists all available VM flavors
	ListVMFlavors(ctx context.Context) ([]*models.VMFlavor, error)
	
	// GetVMImage retrieves a VM image
	GetVMImage(ctx context.Context, imageID string) (*models.VMImage, error)
	
	// ListVMImages lists all available VM images
	ListVMImages(ctx context.Context) ([]*models.VMImage, error)
	
	// GetGPUPolicy retrieves GPU allocation policy
	GetGPUPolicy(ctx context.Context) (*GPUAllocationPolicy, error)
	
	// GetSchedulingPolicy retrieves scheduler policy
	GetSchedulingPolicy(ctx context.Context) (*SchedulingPolicy, error)
}

// GPUAllocationPolicy defines GPU allocation rules
type GPUAllocationPolicy struct {
	AllowSharing          bool
	MaxVMsPerGPU          int
	PreferDedicatedGPUs   bool
	AllowNUMARemoteGPU    bool
	EnableMIG             bool
	DefaultMIGProfile     string
}

// SchedulingPolicy defines scheduling rules
type SchedulingPolicy struct {
	Algorithm             string // bin-packing, spread
	AntiAffinityRules     []string
	PreferredNodeLabels   map[string]string
	AllowOvercommit       bool
	OvercommitRatio       float64
}

// EventBus handles event publishing and subscription
type EventBus interface {
	// Publish publishes an event
	Publish(ctx context.Context, event *models.Event) error
	
	// Subscribe subscribes to events of a type
	Subscribe(ctx context.Context, eventType models.EventType) (chan *models.Event, error)
	
	// Unsubscribe unsubscribes from event channel
	Unsubscribe(ctx context.Context, ch chan *models.Event) error
}

// AuditLogger handles audit trail recording
type AuditLogger interface {
	// LogAction logs an action for audit purposes
	LogAction(ctx context.Context, log *models.AuditLog) error
	
	// GetLogs retrieves audit logs
	GetLogs(ctx context.Context, filters map[string]string) ([]*models.AuditLog, error)
	
	// GetLogsForResource retrieves logs for a specific resource
	GetLogsForResource(ctx context.Context, resourceType, resourceID string) ([]*models.AuditLog, error)
}

// StateStore provides persistence for system state
type StateStore interface {
	// VM operations
	SaveVM(ctx context.Context, vm *models.VirtualMachine) error
	GetVM(ctx context.Context, vmID string) (*models.VirtualMachine, error)
	DeleteVM(ctx context.Context, vmID string) error
	ListVMs(ctx context.Context) ([]*models.VirtualMachine, error)
	
	// GPU operations
	SaveGPU(ctx context.Context, gpu *models.GPU) error
	GetGPU(ctx context.Context, gpuID string) (*models.GPU, error)
	ListGPUs(ctx context.Context) ([]*models.GPU, error)
	
	// Node operations
	SaveNode(ctx context.Context, node *models.HostNode) error
	GetNode(ctx context.Context, nodeID string) (*models.HostNode, error)
	ListNodes(ctx context.Context) ([]*models.HostNode, error)
	
	// Task operations
	SaveTask(ctx context.Context, task *models.Task) error
	GetTask(ctx context.Context, taskID string) (*models.Task, error)
	ListTasks(ctx context.Context) ([]*models.Task, error)
}
