package mocks

import (
	"context"

	"github.com/DARREN-2000/ai-hypervisor-platform/internal/models"
)

type MockGPUOrchestrator struct{}

func (m *MockGPUOrchestrator) AllocateGPU(ctx context.Context, vmID string, requests []models.GPURequest) ([]models.GPUAllocation, error) {
	return []models.GPUAllocation{}, nil
}

func (m *MockGPUOrchestrator) DeallocateGPU(ctx context.Context, vmID string) error {
	return nil
}

func (m *MockGPUOrchestrator) GetGPUAvailability(ctx context.Context) ([]*models.GPU, error) {
	return []*models.GPU{}, nil
}

func (m *MockGPUOrchestrator) GetGPUByID(ctx context.Context, gpuID string) (*models.GPU, error) {
	return &models.GPU{}, nil
}

func (m *MockGPUOrchestrator) UpdateGPUMetrics(ctx context.Context, gpuID string, metrics *models.GPUMetrics) error {
	return nil
}

func (m *MockGPUOrchestrator) CheckGPUHealth(ctx context.Context, nodeID string) error {
	return nil
}
