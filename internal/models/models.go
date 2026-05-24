package models

import (
	"time"

	"github.com/google/uuid"
)

// VMState represents the lifecycle state of a virtual machine
type VMState string

const (
	VMStateCreated      VMState = "created"
	VMStateProvisioning VMState = "provisioning"
	VMStateRunning      VMState = "running"
	VMStatePaused       VMState = "paused"
	VMStateScaling      VMState = "scaling"
	VMStateFailed       VMState = "failed"
	VMStateDeleting     VMState = "deleting"
	VMStateDeleted      VMState = "deleted"
)

// VirtualMachine represents a VM definition and its current state
type VirtualMachine struct {
	ID              string            `json:"id"`
	Name            string            `json:"name"`
	Namespace       string            `json:"namespace"`
	Flavor          VMFlavor          `json:"flavor"`
	Image           VMImage           `json:"image"`
	State           VMState           `json:"state"`
	HostNodeID      string            `json:"host_node_id,omitempty"`
	NetworkConfig   NetworkConfig     `json:"network_config"`
	StorageConfig   StorageConfig     `json:"storage_config"`
	GPURequests     []GPURequest      `json:"gpu_requests"`
	GPUAllocations  []GPUAllocation   `json:"gpu_allocations,omitempty"`
	Metadata        map[string]string `json:"metadata,omitempty"`
	CreatedAt       time.Time         `json:"created_at"`
	UpdatedAt       time.Time         `json:"updated_at"`
	StartedAt       *time.Time        `json:"started_at,omitempty"`
	TerminatedAt    *time.Time        `json:"terminated_at,omitempty"`
	LastStatusCheck time.Time         `json:"last_status_check"`
	ResourceUsage   *ResourceMetrics  `json:"resource_usage,omitempty"`
	ErrorMessage    string            `json:"error_message,omitempty"`
}

// VMFlavor defines CPU, memory, and disk resources
type VMFlavor struct {
	Name      string `json:"name"`
	CPU       int    `json:"cpu"`       // Number of vCPUs
	Memory    int    `json:"memory"`    // Memory in GB
	DiskSize  int    `json:"disk_size"` // Disk size in GB
	CPUModel  string `json:"cpu_model,omitempty"`
	Features  []string `json:"features,omitempty"` // CPU features
}

// VMImage represents the disk image for a VM
type VMImage struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Format       string `json:"format"`       // qcow2, raw, etc
	Source       string `json:"source"`       // Local path or remote URL
	SizeGB       int    `json:"size_gb"`
	Checksum     string `json:"checksum,omitempty"`
	ChecksumType string `json:"checksum_type,omitempty"` // sha256, md5
}

// VMTemplate defines reusable VM defaults for provisioning.
type VMTemplate struct {
	ID            string            `json:"id"`
	Name          string            `json:"name"`
	Description   string            `json:"description,omitempty"`
	Flavor        VMFlavor          `json:"flavor"`
	Image         VMImage           `json:"image"`
	NetworkConfig NetworkConfig     `json:"network_config"`
	StorageConfig StorageConfig     `json:"storage_config"`
	GPURequests   []GPURequest      `json:"gpu_requests,omitempty"`
	Metadata      map[string]string `json:"metadata,omitempty"`
	CreatedAt     time.Time         `json:"created_at"`
	UpdatedAt     time.Time         `json:"updated_at"`
}

// NetworkConfig defines VM network configuration
type NetworkConfig struct {
	Type         string              `json:"type"` // bridge, nat, vlan
	Interfaces   []NetworkInterface  `json:"interfaces"`
	DNSServers   []string            `json:"dns_servers,omitempty"`
	DNSSearch    []string            `json:"dns_search,omitempty"`
	MTU          int                 `json:"mtu,omitempty"`
}

// NetworkInterface represents a single network interface
type NetworkInterface struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	MAC         string `json:"mac,omitempty"`
	IP          string `json:"ip,omitempty"`
	Netmask     string `json:"netmask,omitempty"`
	Gateway     string `json:"gateway,omitempty"`
	Bridge      string `json:"bridge,omitempty"`
	VLAN        int    `json:"vlan,omitempty"`
	MTU         int    `json:"mtu,omitempty"`
}

// StorageConfig defines VM storage configuration
type StorageConfig struct {
	Volumes []Volume `json:"volumes"`
}

// Volume represents a storage volume attached to a VM
type Volume struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Type       string `json:"type"`       // disk, cdrom
	Source     string `json:"source"`     // Path to volume
	TargetDev  string `json:"target_dev"` // vda, vdb, etc
	Format     string `json:"format,omitempty"`
	Readonly   bool   `json:"readonly,omitempty"`
	SizeGB     int    `json:"size_gb,omitempty"`
}

// GPURequest represents a GPU resource request
type GPURequest struct {
	ID              string `json:"id"`
	Type            string `json:"type"` // nvidia, amd, intel
	Model           string `json:"model"`
	Count           int    `json:"count"`
	MinMemoryGB     int    `json:"min_memory_gb,omitempty"`
	RequiredFeatures []string `json:"required_features,omitempty"` // cuda, tensor-cores
}

// GPUAllocation represents an allocated GPU resource
type GPUAllocation struct {
	ID              string    `json:"id"`
	GPUID           string    `json:"gpu_id"`
	VMName          string    `json:"vm_name"`
	HostNodeID      string    `json:"host_node_id"`
	AllocatedAt     time.Time `json:"allocated_at"`
	DeviceIndex     int       `json:"device_index"` // 0-7
	PCIAddress      string    `json:"pci_address"`
	MIGMode         bool      `json:"mig_mode,omitempty"`
	MIGProfile      string    `json:"mig_profile,omitempty"` // e.g., "1g.5gb"
	UtilizationData *GPUUtil  `json:"utilization_data,omitempty"`
}

// GPU represents a physical GPU device
type GPU struct {
	ID          string             `json:"id"`
	HostNodeID  string             `json:"host_node_id"`
	Type        string             `json:"type"` // nvidia, amd, intel
	Model       string             `json:"model"`
	VRAM        int                `json:"vram"` // VRAM in GB
	Index       int                `json:"index"`
	PCIAddress  string             `json:"pci_address"`
	Status      GPUStatus          `json:"status"`
	Capabilities GPUCapabilities   `json:"capabilities"`
	Metrics     *GPUMetrics        `json:"metrics,omitempty"`
	LastUpdated time.Time          `json:"last_updated"`
}

// GPUStatus represents the operational status of a GPU
type GPUStatus string

const (
	GPUStatusAvailable GPUStatus = "available"
	GPUStatusAllocated GPUStatus = "allocated"
	GPUStatusFaulty    GPUStatus = "faulty"
	GPUStatusDisabled  GPUStatus = "disabled"
)

// GPUCapabilities describes GPU capabilities
type GPUCapabilities struct {
	CUDA              bool   `json:"cuda,omitempty"`
	TensorCores       bool   `json:"tensor_cores,omitempty"`
	RTCores           bool   `json:"rt_cores,omitempty"`
	PCIeGen           int    `json:"pcie_gen,omitempty"`
	MaxPower          int    `json:"max_power,omitempty"` // Watts
	MIGSupported      bool   `json:"mig_supported,omitempty"`
	NVLinkSupported   bool   `json:"nvlink_supported,omitempty"`
	ComputeCapability string `json:"compute_capability,omitempty"` // e.g., "8.0"
}

// GPUMetrics represents current GPU metrics
type GPUMetrics struct {
	Utilization      int `json:"utilization"`      // 0-100%
	MemoryUsed       int `json:"memory_used"`      // MB
	MemoryFree       int `json:"memory_free"`      // MB
	TemperatureC     int `json:"temperature_c"`
	PowerDraw        int `json:"power_draw"`       // Watts
	ClockCore        int `json:"clock_core"`       // MHz
	ClockMemory      int `json:"clock_memory"`     // MHz
}

// GPUUtil represents GPU utilization data
type GPUUtil struct {
	Utilization  int `json:"utilization"`
	MemoryUsed   int `json:"memory_used"`
	TemperatureC int `json:"temperature_c"`
	PowerDraw    int `json:"power_draw"`
}

// HostNode represents a physical node in the cluster
type HostNode struct {
	ID              string              `json:"id"`
	Name            string              `json:"name"`
	Status          HostStatus          `json:"status"`
	Capacity        HostCapacity        `json:"capacity"`
	AllocatedResources HostAllocated    `json:"allocated_resources"`
	GPUs            []GPU               `json:"gpus,omitempty"`
	Metadata        map[string]string   `json:"metadata,omitempty"`
	LastHeartbeat   time.Time           `json:"last_heartbeat"`
	CreatedAt       time.Time           `json:"created_at"`
}

// HostStatus represents the operational status of a host
type HostStatus string

const (
	HostStatusReady      HostStatus = "ready"
	HostStatusNotReady   HostStatus = "not_ready"
	HostStatusDraining   HostStatus = "draining"
	HostStatusFailed     HostStatus = "failed"
)

// HostCapacity represents available resources on a host
type HostCapacity struct {
	CPU     int `json:"cpu"`          // Total CPU cores
	Memory  int `json:"memory"`       // Total memory in GB
	DiskGB  int `json:"disk_gb"`      // Total disk in GB
	GPUSlots int `json:"gpu_slots"`   // Total GPU slots
}

// HostAllocated represents allocated resources on a host
type HostAllocated struct {
	CPU     int `json:"cpu"`
	Memory  int `json:"memory"`
	DiskGB  int `json:"disk_gb"`
	GPUSlots int `json:"gpu_slots"`
	VMs     int `json:"vms"` // Number of running VMs
}

// ResourceMetrics represents resource usage metrics
type ResourceMetrics struct {
	CPU        float64 `json:"cpu"`         // CPU usage in cores
	Memory     int     `json:"memory"`      // Memory usage in MB
	DiskIORead int     `json:"disk_io_read"`     // KB/s
	DiskIOWrite int    `json:"disk_io_write"`    // KB/s
	NetworkIn  int     `json:"network_in"`       // Mbps
	NetworkOut int     `json:"network_out"`      // Mbps
	Timestamp  time.Time `json:"timestamp"`
}

// Task represents an async background task
type Task struct {
	ID              string            `json:"id"`
	Type            TaskType          `json:"type"`
	Status          TaskStatus        `json:"status"`
	VMRef           *VMRef            `json:"vm_ref,omitempty"`
	Payload         map[string]interface{} `json:"payload,omitempty"`
	Result          map[string]interface{} `json:"result,omitempty"`
	ErrorMessage    string            `json:"error_message,omitempty"`
	RetryCount      int               `json:"retry_count"`
	MaxRetries      int               `json:"max_retries"`
	CreatedAt       time.Time         `json:"created_at"`
	UpdatedAt       time.Time         `json:"updated_at"`
	CompletedAt     *time.Time        `json:"completed_at,omitempty"`
}

// TaskType represents the type of task
type TaskType string

const (
	TaskTypeProvisionVM      TaskType = "provision_vm"
	TaskTypeDeleteVM         TaskType = "delete_vm"
	TaskTypeStartVM          TaskType = "start_vm"
	TaskTypeStopVM           TaskType = "stop_vm"
	TaskTypeRebootVM         TaskType = "reboot_vm"
	TaskTypeResizeVM         TaskType = "resize_vm"
	TaskTypeAllocateGPU      TaskType = "allocate_gpu"
	TaskTypeDeallocateGPU    TaskType = "deallocate_gpu"
)

// TaskStatus represents the execution status of a task
type TaskStatus string

const (
	TaskStatusQueued     TaskStatus = "queued"
	TaskStatusRunning    TaskStatus = "running"
	TaskStatusCompleted  TaskStatus = "completed"
	TaskStatusFailed     TaskStatus = "failed"
	TaskStatusRetrying   TaskStatus = "retrying"
)

// VMRef is a lightweight reference to a VM
type VMRef struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

// AuditLog represents an audit trail entry
type AuditLog struct {
	ID        string                 `json:"id"`
	Timestamp time.Time              `json:"timestamp"`
	Action    string                 `json:"action"`
	Actor     string                 `json:"actor"`
	Resource  string                 `json:"resource"`
	ResourceID string                `json:"resource_id"`
	Changes   map[string]interface{} `json:"changes,omitempty"`
	Status    string                 `json:"status"` // success, failure
	Message   string                 `json:"message,omitempty"`
}

// EventType represents event types in the system
type EventType string

const (
	EventTypeVMCreated       EventType = "vm.created"
	EventTypeVMRunning       EventType = "vm.running"
	EventTypeVMFailed        EventType = "vm.failed"
	EventTypeVMDeleted       EventType = "vm.deleted"
	EventTypeGPUAllocated    EventType = "gpu.allocated"
	EventTypeGPUDeallocated  EventType = "gpu.deallocated"
	EventTypeGPUFaulty       EventType = "gpu.faulty"
	EventTypeHostDegraded    EventType = "host.degraded"
)

// Event represents a system event
type Event struct {
	ID        string                 `json:"id"`
	Type      EventType              `json:"type"`
	Timestamp time.Time              `json:"timestamp"`
	Source    string                 `json:"source"`
	Subject   string                 `json:"subject"`
	Data      map[string]interface{} `json:"data,omitempty"`
	Severity  string                 `json:"severity"` // info, warning, error
}

// SchedulingDecision represents a scheduler decision
type SchedulingDecision struct {
	VMID              string            `json:"vm_id"`
	SelectedHostID    string            `json:"selected_host_id"`
	SelectedGPUIds    []string          `json:"selected_gpu_ids,omitempty"`
	Policy            string            `json:"policy,omitempty"`
	Score             float64           `json:"score"`
	Reason            string            `json:"reason"`
	Metadata          map[string]string `json:"metadata,omitempty"`
	GPUAssignments    []GPUAllocation   `json:"gpu_assignments,omitempty"`
	AlternativeHosts  []HostScore       `json:"alternative_hosts,omitempty"`
	DecisionTimestamp time.Time         `json:"decision_timestamp"`
}

// HostScore represents a host scoring result
type HostScore struct {
	HostID string  `json:"host_id"`
	Score  float64 `json:"score"`
	Reason string  `json:"reason"`
}

// Generates a new UUID
func NewID() string {
	return uuid.New().String()
}
