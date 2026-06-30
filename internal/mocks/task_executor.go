package mocks

import (
	"context"

	"github.com/DARREN-2000/ai-hypervisor-platform/internal/models"
)

type MockTaskExecutor struct{}

func (m *MockTaskExecutor) EnqueueTask(ctx context.Context, task *models.Task) error {
	return nil
}

func (m *MockTaskExecutor) DequeueTask(ctx context.Context) (*models.Task, error) {
	return nil, nil
}

func (m *MockTaskExecutor) UpdateTaskStatus(ctx context.Context, taskID string, status models.TaskStatus, result map[string]interface{}) error {
	return nil
}

func (m *MockTaskExecutor) GetTask(ctx context.Context, taskID string) (*models.Task, error) {
	return &models.Task{}, nil
}

func (m *MockTaskExecutor) ListTasks(ctx context.Context, filters map[string]string) ([]*models.Task, error) {
	return []*models.Task{}, nil
}

func (m *MockTaskExecutor) RetryTask(ctx context.Context, taskID string) error {
	return nil
}
