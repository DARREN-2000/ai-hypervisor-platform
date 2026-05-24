package telemetry

import "time"

// Metrics defines a small surface for HTTP metrics recording.
type Metrics interface {
	ObserveRequest(method, path string, status int, duration time.Duration)
	ObserveVMResource(vmID, hostID string, cpuCores float64, memoryBytes int64, diskReadBytes int64, diskWriteBytes int64, netInBytes int64, netOutBytes int64)
	ObserveGPUUsage(gpuID, hostID, model string, utilization float64, memoryUsedBytes int64, memoryFreeBytes int64, temperatureC float64, powerW float64)
}

// NoopMetrics is a placeholder metrics recorder.
type NoopMetrics struct{}

// ObserveRequest implements Metrics with no-op behavior.
func (NoopMetrics) ObserveRequest(method, path string, status int, duration time.Duration) {}
func (NoopMetrics) ObserveVMResource(vmID, hostID string, cpuCores float64, memoryBytes int64, diskReadBytes int64, diskWriteBytes int64, netInBytes int64, netOutBytes int64) {
}
func (NoopMetrics) ObserveGPUUsage(gpuID, hostID, model string, utilization float64, memoryUsedBytes int64, memoryFreeBytes int64, temperatureC float64, powerW float64) {
}

// NewNoopMetrics returns a no-op metrics recorder.
func NewNoopMetrics() Metrics {
	return NoopMetrics{}
}
