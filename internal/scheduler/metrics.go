package scheduler

import (
	"sync/atomic"

	"github.com/DARREN-2000/ai-hypervisor-platform/internal/orchestrator"
)

type schedulerMetrics struct {
	decisions      int64
	failures       int64
	totalLatencyMs int64
}

func (m *schedulerMetrics) recordDecision(latencyMs int64, success bool) {
	if success {
		atomic.AddInt64(&m.decisions, 1)
		atomic.AddInt64(&m.totalLatencyMs, latencyMs)
	} else {
		atomic.AddInt64(&m.failures, 1)
	}
}

func (m *schedulerMetrics) snapshot() *orchestrator.SchedulingMetrics {
	decisions := atomic.LoadInt64(&m.decisions)
	failures := atomic.LoadInt64(&m.failures)
	totalLatency := atomic.LoadInt64(&m.totalLatencyMs)

	avg := 0
	if decisions > 0 {
		avg = int(totalLatency / decisions)
	}

	successful := int(decisions)
	failed := int(failures)
	total := successful + failed

	optimalRate := 0.0
	if total > 0 {
		optimalRate = float64(successful) / float64(total)
	}

	return &orchestrator.SchedulingMetrics{
		AverageDecisionLatency: avg,
		FailedSchedulingAttempts: failed,
		SuccessfulSchedules:    successful,
		OptimalPlacementRate:   optimalRate,
	}
}
