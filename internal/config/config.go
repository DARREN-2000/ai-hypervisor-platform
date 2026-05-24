package config

import "time"

// APIServerConfig represents API server configuration
type APIServerConfig struct {
	Address        string        `yaml:"address"`
	Port           int           `yaml:"port"`
	ReadTimeout    time.Duration `yaml:"read_timeout"`
	WriteTimeout   time.Duration `yaml:"write_timeout"`
	IdleTimeout    time.Duration `yaml:"idle_timeout"`
	MaxHeaderBytes int           `yaml:"max_header_bytes"`
	TLSCert        string        `yaml:"tls_cert"`
	TLSKey         string        `yaml:"tls_key"`
	Metrics        MetricsConfig `yaml:"metrics"`
	Database       DatabaseConfig `yaml:"database"`
	Redis          RedisConfig   `yaml:"redis"`
	NATS           NATSConfig    `yaml:"nats"`
}

// VMManagerConfig represents VM manager configuration
type VMManagerConfig struct {
	MaxConcurrentProvisioning int           `yaml:"max_concurrent_provisioning"`
	ProvisioningTimeout       time.Duration `yaml:"provisioning_timeout"`
	StateCheckInterval        time.Duration `yaml:"state_check_interval"`
	HeartbeatInterval         time.Duration `yaml:"heartbeat_interval"`
}

// SchedulerConfig represents scheduler configuration
type SchedulerConfig struct {
	Algorithm              string  `yaml:"algorithm"` // bin-packing, spread, numa-aware
	SchedulingTimeout      time.Duration `yaml:"scheduling_timeout"`
	AllowOvercommit        bool    `yaml:"allow_overcommit"`
	OvercommitRatio        float64 `yaml:"overcommit_ratio"`
	MaxVMsPerNode          int     `yaml:"max_vms_per_node"`
	PreferredNodeLabels    map[string]string `yaml:"preferred_node_labels"`
}

// GPUOrchestratorConfig represents GPU orchestrator configuration
type GPUOrchestratorConfig struct {
	AllocationPolicy      string `yaml:"allocation_policy"` // bin-packing, spread, numa-aware
	AllowSharing          bool   `yaml:"allow_sharing"`
	MaxVMsPerGPU          int    `yaml:"max_vms_per_gpu"`
	PreferDedicatedGPUs   bool   `yaml:"prefer_dedicated_gpus"`
	AllowNUMARemoteGPU    bool   `yaml:"allow_numa_remote_gpu"`
	EnableMIG             bool   `yaml:"enable_mig"`
	DefaultMIGProfile     string `yaml:"default_mig_profile"`
	HealthCheckInterval   time.Duration `yaml:"health_check_interval"`
	GPUMetricsInterval    time.Duration `yaml:"gpu_metrics_interval"`
}

// TaskExecutorConfig represents task executor configuration
type TaskExecutorConfig struct {
	MaxConcurrentTasks     int           `yaml:"max_concurrent_tasks"`
	MaxRetries             int           `yaml:"max_retries"`
	InitialBackoff         time.Duration `yaml:"initial_backoff"`
	MaxBackoff             time.Duration `yaml:"max_backoff"`
	TaskTimeout            time.Duration `yaml:"task_timeout"`
	FailedTaskRetentionTTL time.Duration `yaml:"failed_task_retention_ttl"`
}

// ResourceMonitorConfig represents resource monitor configuration
type ResourceMonitorConfig struct {
	MetricsInterval       time.Duration `yaml:"metrics_interval"`
	MetricsRetentionDays  int           `yaml:"metrics_retention_days"`
	PredictionWindow      time.Duration `yaml:"prediction_window"`
	AlertThresholds       AlertThresholds `yaml:"alert_thresholds"`
}

// AlertThresholds represents alert thresholds
type AlertThresholds struct {
	CPUUtilization    float64 `yaml:"cpu_utilization"`
	MemoryUtilization float64 `yaml:"memory_utilization"`
	DiskUtilization   float64 `yaml:"disk_utilization"`
	GPUUtilization    float64 `yaml:"gpu_utilization"`
	GPUTemperature    int     `yaml:"gpu_temperature"`
}

// DatabaseConfig represents database configuration
type DatabaseConfig struct {
	Host            string        `yaml:"host"`
	Port            int           `yaml:"port"`
	Username        string        `yaml:"username"`
	Password        string        `yaml:"password"`
	Database        string        `yaml:"database"`
	SSLMode         string        `yaml:"ssl_mode"`
	MaxConnections  int           `yaml:"max_connections"`
	ConnectTimeout  time.Duration `yaml:"connect_timeout"`
	IdleTimeout     time.Duration `yaml:"idle_timeout"`
	MaxIdleConns    int           `yaml:"max_idle_conns"`
	BackupEnabled   bool          `yaml:"backup_enabled"`
	BackupInterval  time.Duration `yaml:"backup_interval"`
}

// RedisConfig represents Redis configuration
type RedisConfig struct {
	Host            string        `yaml:"host"`
	Port            int           `yaml:"port"`
	Password        string        `yaml:"password"`
	Database        int           `yaml:"database"`
	PoolSize        int           `yaml:"pool_size"`
	ConnectTimeout  time.Duration `yaml:"connect_timeout"`
	ReadTimeout     time.Duration `yaml:"read_timeout"`
	WriteTimeout    time.Duration `yaml:"write_timeout"`
}

// NATSConfig represents NATS configuration
type NATSConfig struct {
	URLs           []string      `yaml:"urls"`
	Username       string        `yaml:"username"`
	Password       string        `yaml:"password"`
	ConnectTimeout time.Duration `yaml:"connect_timeout"`
	MaxReconnect   int           `yaml:"max_reconnect"`
	ReconnectWait  time.Duration `yaml:"reconnect_wait"`
}

// MetricsConfig represents metrics configuration
type MetricsConfig struct {
	Enabled         bool   `yaml:"enabled"`
	PrometheusAddr  string `yaml:"prometheus_addr"`
	PrometheusPort  int    `yaml:"prometheus_port"`
	ScrapeInterval  time.Duration `yaml:"scrape_interval"`
	RetentionDays   int    `yaml:"retention_days"`
}

// LoggingConfig represents logging configuration
type LoggingConfig struct {
	Level          string `yaml:"level"` // debug, info, warn, error
	Format         string `yaml:"format"` // json, text
	OutputPath     string `yaml:"output_path"`
	ErrorOutputPath string `yaml:"error_output_path"`
	Encoding       string `yaml:"encoding"`
}

// LibvirtConfig represents Libvirt configuration
type LibvirtConfig struct {
	URI                string        `yaml:"uri"` // qemu:///system
	ConnectionTimeout  time.Duration `yaml:"connection_timeout"`
	MaxConnections    int           `yaml:"max_connections"`
}

// KubernetesConfig represents Kubernetes integration configuration
type KubernetesConfig struct {
	Enabled        bool   `yaml:"enabled"`
	ConfigPath     string `yaml:"config_path"`
	Namespace      string `yaml:"namespace"`
	InCluster      bool   `yaml:"in_cluster"`
	WebhookPort    int    `yaml:"webhook_port"`
	CRDEnabled     bool   `yaml:"crd_enabled"`
}

// SecurityConfig represents security configuration
type SecurityConfig struct {
	TLSEnabled       bool   `yaml:"tls_enabled"`
	TLSCertPath      string `yaml:"tls_cert_path"`
	TLSKeyPath       string `yaml:"tls_key_path"`
	CAPath           string `yaml:"ca_path"`
	MutualTLSEnabled bool   `yaml:"mutual_tls_enabled"`
	APIKeyRequired   bool   `yaml:"api_key_required"`
	JWTEnabled       bool   `yaml:"jwt_enabled"`
	JWTSecret        string `yaml:"jwt_secret"`
	RBACEnabled      bool   `yaml:"rbac_enabled"`
	AuditEnabled     bool   `yaml:"audit_enabled"`
}

// PlatformConfig represents the entire platform configuration
type PlatformConfig struct {
	Environment     string              `yaml:"environment"`
	APIServer       APIServerConfig       `yaml:"api_server"`
	VMManager       VMManagerConfig       `yaml:"vm_manager"`
	Scheduler       SchedulerConfig       `yaml:"scheduler"`
	GPUOrchestrator GPUOrchestratorConfig `yaml:"gpu_orchestrator"`
	TaskExecutor    TaskExecutorConfig    `yaml:"task_executor"`
	ResourceMonitor ResourceMonitorConfig `yaml:"resource_monitor"`
	Database        DatabaseConfig        `yaml:"database"`
	Redis           RedisConfig           `yaml:"redis"`
	NATS            NATSConfig            `yaml:"nats"`
	Metrics         MetricsConfig         `yaml:"metrics"`
	Logging         LoggingConfig         `yaml:"logging"`
	Libvirt         LibvirtConfig         `yaml:"libvirt"`
	Kubernetes      KubernetesConfig      `yaml:"kubernetes"`
	Security        SecurityConfig        `yaml:"security"`
}

// DefaultConfig returns default configuration
func DefaultConfig() *PlatformConfig {
	return &PlatformConfig{
		Environment: "production",
		APIServer: APIServerConfig{
			Address:        "0.0.0.0",
			Port:           8080,
			ReadTimeout:    15 * time.Second,
			WriteTimeout:   15 * time.Second,
			IdleTimeout:    60 * time.Second,
			MaxHeaderBytes: 1 << 20,
		},
		VMManager: VMManagerConfig{
			MaxConcurrentProvisioning: 10,
			ProvisioningTimeout:       5 * time.Minute,
			StateCheckInterval:        30 * time.Second,
			HeartbeatInterval:         30 * time.Second,
		},
		Scheduler: SchedulerConfig{
			Algorithm:       "bin-packing",
			SchedulingTimeout: 10 * time.Second,
			AllowOvercommit: false,
			OvercommitRatio: 1.0,
			MaxVMsPerNode:   100,
		},
		GPUOrchestrator: GPUOrchestratorConfig{
			AllocationPolicy:    "bin-packing",
			AllowSharing:        false,
			MaxVMsPerGPU:        1,
			PreferDedicatedGPUs: true,
			EnableMIG:           false,
			HealthCheckInterval: 30 * time.Second,
			GPUMetricsInterval:  10 * time.Second,
		},
		TaskExecutor: TaskExecutorConfig{
			MaxConcurrentTasks:     20,
			MaxRetries:             3,
			InitialBackoff:         1 * time.Second,
			MaxBackoff:             5 * time.Minute,
			TaskTimeout:            10 * time.Minute,
			FailedTaskRetentionTTL: 24 * time.Hour,
		},
		ResourceMonitor: ResourceMonitorConfig{
			MetricsInterval:      10 * time.Second,
			MetricsRetentionDays: 30,
			PredictionWindow:     5 * time.Minute,
			AlertThresholds: AlertThresholds{
				CPUUtilization:    80,
				MemoryUtilization: 85,
				DiskUtilization:   90,
				GPUUtilization:    90,
				GPUTemperature:    80,
			},
		},
		Database: DatabaseConfig{
			Host:           "postgres.infra.svc.cluster.local",
			Port:           5432,
			Database:       "aihypervisor",
			SSLMode:        "require",
			MaxConnections: 100,
			ConnectTimeout: 10 * time.Second,
			IdleTimeout:    5 * time.Minute,
			MaxIdleConns:   25,
		},
		Redis: RedisConfig{
			Host:           "redis.infra.svc.cluster.local",
			Port:           6379,
			PoolSize:       50,
			ConnectTimeout: 5 * time.Second,
			ReadTimeout:    3 * time.Second,
			WriteTimeout:   3 * time.Second,
		},
		NATS: NATSConfig{
			URLs:           []string{"nats://nats.infra.svc.cluster.local:4222"},
			ConnectTimeout: 5 * time.Second,
			MaxReconnect:   10,
			ReconnectWait:  2 * time.Second,
		},
		Metrics: MetricsConfig{
			Enabled:        true,
			PrometheusAddr: "0.0.0.0",
			PrometheusPort: 9090,
			ScrapeInterval: 15 * time.Second,
			RetentionDays:  15,
		},
		Logging: LoggingConfig{
			Level:      "info",
			Format:     "json",
			OutputPath: "stdout",
		},
		Libvirt: LibvirtConfig{
			URI:               "qemu:///system",
			ConnectionTimeout: 10 * time.Second,
			MaxConnections:    10,
		},
		Kubernetes: KubernetesConfig{
			Enabled:   true,
			Namespace: "aihypervisor",
			InCluster: true,
			CRDEnabled: true,
		},
		Security: SecurityConfig{
			TLSEnabled:      true,
			MutualTLSEnabled: true,
			RBACEnabled:     true,
			AuditEnabled:    true,
		},
	}
}
