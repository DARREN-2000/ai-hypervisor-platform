package observability

import (
	"time"
	"strings"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/DARREN-2000/ai-hypervisor-platform/pkg/telemetry"
)

type prometheusMetrics struct {
	requestDuration *prometheus.HistogramVec
	requestTotal    *prometheus.CounterVec
	requestInFlight  prometheus.Gauge
	vmCPU           *prometheus.GaugeVec
	vmMemory        *prometheus.GaugeVec
	vmDiskRead      *prometheus.GaugeVec
	vmDiskWrite     *prometheus.GaugeVec
	vmNetworkIn     *prometheus.GaugeVec
	vmNetworkOut    *prometheus.GaugeVec
	gpuUtilization  *prometheus.GaugeVec
	gpuMemoryUsed   *prometheus.GaugeVec
	gpuMemoryFree   *prometheus.GaugeVec
	gpuTemperature  *prometheus.GaugeVec
	gpuPower        *prometheus.GaugeVec
	healthStatus    *prometheus.GaugeVec
}

func newPrometheusMetrics(registry *prometheus.Registry, serviceName, serviceVersion, environment string) telemetry.Metrics {
	namespace := sanitizeMetricName(serviceName)
	metrics := &prometheusMetrics{
		requestDuration: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: namespace,
			Subsystem: "http",
			Name:      "request_duration_seconds",
			Help:      "HTTP request duration in seconds.",
			Buckets:   prometheus.DefBuckets,
		}, []string{"method", "path", "status"}),
		requestTotal: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: "http",
			Name:      "requests_total",
			Help:      "Total number of HTTP requests.",
		}, []string{"method", "path", "status"}),
		requestInFlight: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "http",
			Name:      "requests_in_flight",
			Help:      "Current number of in-flight HTTP requests.",
		}),
		vmCPU: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "vm",
			Name:      "cpu_cores",
			Help:      "Observed VM CPU usage in cores.",
		}, []string{"vm_id", "host_id"}),
		vmMemory: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "vm",
			Name:      "memory_bytes",
			Help:      "Observed VM memory usage in bytes.",
		}, []string{"vm_id", "host_id"}),
		vmDiskRead: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "vm",
			Name:      "disk_read_bytes",
			Help:      "Observed VM disk read throughput in bytes.",
		}, []string{"vm_id", "host_id"}),
		vmDiskWrite: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "vm",
			Name:      "disk_write_bytes",
			Help:      "Observed VM disk write throughput in bytes.",
		}, []string{"vm_id", "host_id"}),
		vmNetworkIn: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "vm",
			Name:      "network_in_bytes",
			Help:      "Observed VM inbound network throughput in bytes.",
		}, []string{"vm_id", "host_id"}),
		vmNetworkOut: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "vm",
			Name:      "network_out_bytes",
			Help:      "Observed VM outbound network throughput in bytes.",
		}, []string{"vm_id", "host_id"}),
		gpuUtilization: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "gpu",
			Name:      "utilization_percent",
			Help:      "Observed GPU utilization in percent.",
		}, []string{"gpu_id", "host_id", "model"}),
		gpuMemoryUsed: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "gpu",
			Name:      "memory_used_bytes",
			Help:      "Observed GPU memory used in bytes.",
		}, []string{"gpu_id", "host_id", "model"}),
		gpuMemoryFree: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "gpu",
			Name:      "memory_free_bytes",
			Help:      "Observed GPU memory free in bytes.",
		}, []string{"gpu_id", "host_id", "model"}),
		gpuTemperature: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "gpu",
			Name:      "temperature_celsius",
			Help:      "Observed GPU temperature in Celsius.",
		}, []string{"gpu_id", "host_id", "model"}),
		gpuPower: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "gpu",
			Name:      "power_watts",
			Help:      "Observed GPU power draw in watts.",
		}, []string{"gpu_id", "host_id", "model"}),
		healthStatus: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "health",
			Name:      "component_status",
			Help:      "Health status by component where 1 means healthy.",
		}, []string{"component", "version", "environment"}),
	}

	collectors := []prometheus.Collector{
		metrics.requestDuration,
		metrics.requestTotal,
		metrics.requestInFlight,
		metrics.vmCPU,
		metrics.vmMemory,
		metrics.vmDiskRead,
		metrics.vmDiskWrite,
		metrics.vmNetworkIn,
		metrics.vmNetworkOut,
		metrics.gpuUtilization,
		metrics.gpuMemoryUsed,
		metrics.gpuMemoryFree,
		metrics.gpuTemperature,
		metrics.gpuPower,
		metrics.healthStatus,
		prometheus.NewGaugeFunc(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "build",
			Name:      "info",
			Help:      "Build and deployment metadata.",
		}, func() float64 { return 1 }),
	}

	for _, collector := range collectors {
		_ = registry.Register(collector)
	}

	metrics.healthStatus.WithLabelValues("service", serviceVersion, environment).Set(1)

	return metrics
}

func sanitizeMetricName(name string) string {
	name = strings.ToLower(name)
	name = strings.ReplaceAll(name, "-", "_")
	name = strings.ReplaceAll(name, ".", "_")
	return name
}

func (m *prometheusMetrics) ObserveRequest(method, path string, status int, duration time.Duration) {
	if m == nil {
		return
	}
	statusLabel := statusCategory(status)
	m.requestTotal.WithLabelValues(method, path, statusLabel).Inc()
	m.requestDuration.WithLabelValues(method, path, statusLabel).Observe(duration.Seconds())
}

func (m *prometheusMetrics) ObserveVMResource(vmID, hostID string, cpuCores float64, memoryBytes int64, diskReadBytes int64, diskWriteBytes int64, netInBytes int64, netOutBytes int64) {
	if m == nil {
		return
	}
	m.vmCPU.WithLabelValues(vmID, hostID).Set(cpuCores)
	m.vmMemory.WithLabelValues(vmID, hostID).Set(float64(memoryBytes))
	m.vmDiskRead.WithLabelValues(vmID, hostID).Set(float64(diskReadBytes))
	m.vmDiskWrite.WithLabelValues(vmID, hostID).Set(float64(diskWriteBytes))
	m.vmNetworkIn.WithLabelValues(vmID, hostID).Set(float64(netInBytes))
	m.vmNetworkOut.WithLabelValues(vmID, hostID).Set(float64(netOutBytes))
}

func (m *prometheusMetrics) ObserveGPUUsage(gpuID, hostID, model string, utilization float64, memoryUsedBytes int64, memoryFreeBytes int64, temperatureC float64, powerW float64) {
	if m == nil {
		return
	}
	m.gpuUtilization.WithLabelValues(gpuID, hostID, model).Set(utilization)
	m.gpuMemoryUsed.WithLabelValues(gpuID, hostID, model).Set(float64(memoryUsedBytes))
	m.gpuMemoryFree.WithLabelValues(gpuID, hostID, model).Set(float64(memoryFreeBytes))
	m.gpuTemperature.WithLabelValues(gpuID, hostID, model).Set(temperatureC)
	m.gpuPower.WithLabelValues(gpuID, hostID, model).Set(powerW)
}

func statusCategory(status int) string {
	switch {
	case status >= 500:
		return "5xx"
	case status >= 400:
		return "4xx"
	case status >= 300:
		return "3xx"
	default:
		return "2xx"
	}
}
