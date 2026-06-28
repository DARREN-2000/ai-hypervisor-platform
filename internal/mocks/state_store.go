package mocks

import (
	"context"

	"github.com/DARREN-2000/ai-hypervisor-platform/internal/models"
)

type MockStateStore struct{}

func (m *MockStateStore) SaveVM(ctx context.Context, vm *models.VirtualMachine) error {
	return nil
}

func (m *MockStateStore) GetVM(ctx context.Context, vmID string) (*models.VirtualMachine, error) {
	return &models.VirtualMachine{}, nil
}

func (m *MockStateStore) DeleteVM(ctx context.Context, vmID string) error {
	return nil
}

func (m *MockStateStore) ListVMs(ctx context.Context) ([]*models.VirtualMachine, error) {
	return []*models.VirtualMachine{}, nil
}

func (m *MockStateStore) SaveGPU(ctx context.Context, gpu *models.GPU) error {
	return nil
}

func (m *MockStateStore) GetGPU(ctx context.Context, gpuID string) (*models.GPU, error) {
	return &models.GPU{}, nil
}

func (m *MockStateStore) ListGPUs(ctx context.Context) ([]*models.GPU, error) {
	return []*models.GPU{}, nil
}

func (m *MockStateStore) SaveNode(ctx context.Context, node *models.HostNode) error {
	return nil
}

func (m *MockStateStore) GetNode(ctx context.Context, nodeID string) (*models.HostNode, error) {
	return &models.HostNode{}, nil
}

func (m *MockStateStore) ListNodes(ctx context.Context) ([]*models.HostNode, error) {
	return []*models.HostNode{}, nil
}

func (m *MockStateStore) SaveTask(ctx context.Context, task *models.Task) error {
	return nil
}

func (m *MockStateStore) GetTask(ctx context.Context, taskID string) (*models.Task, error) {
	return &models.Task{}, nil
}

func (m *MockStateStore) ListTasks(ctx context.Context) ([]*models.Task, error) {
	return []*models.Task{}, nil
}
