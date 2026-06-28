package mocks

import (
	"context"

	"github.com/DARREN-2000/ai-hypervisor-platform/internal/models"
)

type MockEventBus struct{}

func (m *MockEventBus) Publish(ctx context.Context, event *models.Event) error {
	return nil
}

func (m *MockEventBus) Subscribe(ctx context.Context, eventType models.EventType) (chan *models.Event, error) {
	return make(chan *models.Event), nil
}

func (m *MockEventBus) Unsubscribe(ctx context.Context, ch chan *models.Event) error {
	return nil
}
