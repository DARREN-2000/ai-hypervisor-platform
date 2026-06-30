package mocks

import (
	"context"

	"github.com/DARREN-2000/ai-hypervisor-platform/internal/models"
)

type MockAuditLogger struct{}

func (m *MockAuditLogger) LogAction(ctx context.Context, log *models.AuditLog) error {
	return nil
}

func (m *MockAuditLogger) GetLogs(ctx context.Context, filters map[string]string) ([]*models.AuditLog, error) {
	return []*models.AuditLog{}, nil
}

func (m *MockAuditLogger) GetLogsForResource(ctx context.Context, resourceType, resourceID string) ([]*models.AuditLog, error) {
	return []*models.AuditLog{}, nil
}
