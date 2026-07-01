package scheduler

import (
	"context"
	"testing"
	"github.com/DARREN-2000/ai-hypervisor-platform/internal/config"
	"github.com/DARREN-2000/ai-hypervisor-platform/internal/models"
)

func TestMetrics(t *testing.T) {
	m := &schedulerMetrics{}
	m.recordDecision(10, true)
	m.recordDecision(20, false)

	snap := m.snapshot()
	if snap == nil {
		t.Fatal("expected snapshot")
	}
	if snap.SuccessfulSchedules != 1 {
		t.Errorf("expected 1 successful, got %d", snap.SuccessfulSchedules)
	}
	if snap.FailedSchedulingAttempts != 1 {
		t.Errorf("expected 1 failed, got %d", snap.FailedSchedulingAttempts)
	}
	if snap.AverageDecisionLatency != 10 {
		t.Errorf("expected 10 ms, got %d", snap.AverageDecisionLatency)
	}
	if snap.OptimalPlacementRate != 0.5 {
		t.Errorf("expected 0.5 optimal rate, got %f", snap.OptimalPlacementRate)
	}
}

type mockMonitor struct {
	decision *DecisionRecord
	failure  *FailureRecord
}

func (m *mockMonitor) RecordDecision(ctx context.Context, d DecisionRecord) {
	m.decision = &d
}

func (m *mockMonitor) RecordFailure(ctx context.Context, f FailureRecord) {
	m.failure = &f
}

func TestWithMonitor(t *testing.T) {
	m := &mockMonitor{}
	svc := NewService(config.SchedulerConfig{}, Dependencies{}, WithMonitor(m))
	if svc.monitor != m {
		t.Error("monitor not set")
	}
}

type mockAutoscaler struct {
	notified bool
}

func (m *mockAutoscaler) Notify(ctx context.Context, s AutoscaleSignal) {
	m.notified = true
}

func TestWithAutoscaler(t *testing.T) {
	m := &mockAutoscaler{}
	svc := NewService(config.SchedulerConfig{}, Dependencies{}, WithAutoscaler(m))
	if svc.autoscaler != m {
		t.Error("autoscaler not set")
	}
}

func TestWithPolicy(t *testing.T) {
	p := SpreadPolicy{}
	svc := NewService(config.SchedulerConfig{}, Dependencies{}, WithPolicy(p))
	if svc.policies["spread"] == nil {
		t.Error("policy not set")
	}
}

func TestCheckNodeCapacity(t *testing.T) {
	svc := NewService(config.SchedulerConfig{}, Dependencies{})
	// No repo
	if svc.CheckNodeCapacity(context.Background(), "node", models.VMFlavor{}) {
		t.Error("expected false with no repo")
	}
}

func TestGetSchedulingMetrics(t *testing.T) {
	svc := NewService(config.SchedulerConfig{}, Dependencies{})
	metrics, err := svc.GetSchedulingMetrics(context.Background())
	if err != nil {
		t.Errorf("unexpected err %v", err)
	}
	if metrics == nil {
		t.Error("expected metrics")
	}
}

func TestNoopMonitor(t *testing.T) {
	m := NoopMonitor{}
	m.RecordDecision(context.Background(), DecisionRecord{})
	m.RecordFailure(context.Background(), FailureRecord{})
}

func TestNoopAutoscaler(t *testing.T) {
	a := NoopAutoscaler{}
	a.Notify(context.Background(), AutoscaleSignal{})
}

func TestDefaultNewService(t *testing.T) {
	svc := NewService(config.SchedulerConfig{}, Dependencies{}, nil)
	if svc.monitor == nil {
		t.Error("monitor is nil")
	}
	if svc.autoscaler == nil {
		t.Error("autoscaler is nil")
	}
	if svc.logger == nil {
		t.Error("logger is nil")
	}
}

func TestRecordDecision(t *testing.T) {
	m := &mockMonitor{}
	m.RecordDecision(context.Background(), DecisionRecord{})
	if m.decision == nil {
		t.Error("decision not set")
	}
}

func TestRecordFailure(t *testing.T) {
	m := &mockMonitor{}
	m.RecordFailure(context.Background(), FailureRecord{})
	if m.failure == nil {
		t.Error("failure not set")
	}
}

func TestNotify(t *testing.T) {
	a := &mockAutoscaler{}
	a.Notify(context.Background(), AutoscaleSignal{})
	if !a.notified {
		t.Error("not notified")
	}
}
