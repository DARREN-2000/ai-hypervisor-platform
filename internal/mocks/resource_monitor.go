package mocks

import (
	"context"

	"github.com/DARREN-2000/ai-hypervisor-platform/internal/models"
	"github.com/DARREN-2000/ai-hypervisor-platform/internal/orchestrator"
)

type MockResourceMonitor struct{}

func (m *MockResourceMonitor) GetNodeMetrics(ctx context.Context, nodeID string) (*models.ResourceMetrics, error) {
	return &models.ResourceMetrics{}, nil
}

func (m *MockResourceMonitor) GetVMMetrics(ctx context.Context, vmID string) (*models.ResourceMetrics, error) {
	return &models.ResourceMetrics{}, nil
}

func (m *MockResourceMonitor) GetClusterMetrics(ctx context.Context) (*orchestrator.ClusterMetrics, error) {
	return &orchestrator.ClusterMetrics{}, nil
}

func (m *MockResourceMonitor) CollectMetrics(ctx context.Context, nodeID string) error {
	return nil
}

func (m *MockResourceMonitor) PredictResourceUsage(ctx context.Context, vmID string) (*models.ResourceMetrics, error) {
	return &models.ResourceMetrics{}, nil
}
