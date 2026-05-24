package scheduler

import (
	"context"
	"time"
)

// DecisionRecord captures scheduling outcomes.
type DecisionRecord struct {
	VMID      string
	HostID    string
	Policy    string
	Score     float64
	Reason    string
	GPUIds    []string
	Latency   time.Duration
	Timestamp time.Time
}

// FailureRecord captures scheduling failures.
type FailureRecord struct {
	VMID      string
	Reason    string
	Latency   time.Duration
	Timestamp time.Time
}

// Monitor records scheduling events for observability.
type Monitor interface {
	RecordDecision(ctx context.Context, record DecisionRecord)
	RecordFailure(ctx context.Context, record FailureRecord)
}

// NoopMonitor provides a no-op implementation.
type NoopMonitor struct{}

func (NoopMonitor) RecordDecision(ctx context.Context, record DecisionRecord) {}
func (NoopMonitor) RecordFailure(ctx context.Context, record FailureRecord) {}

// AutoscaleSignal carries signals for future autoscaling.
type AutoscaleSignal struct {
	Reason         string
	CPUShortage    bool
	MemoryShortage bool
	GPUShortage    bool
	Demand         ResourceDemand
	Timestamp      time.Time
}

// AutoscalerHook allows scheduler to emit autoscaling hints.
type AutoscalerHook interface {
	Notify(ctx context.Context, signal AutoscaleSignal)
}

// NoopAutoscaler provides a no-op autoscaler hook.
type NoopAutoscaler struct{}

func (NoopAutoscaler) Notify(ctx context.Context, signal AutoscaleSignal) {}
