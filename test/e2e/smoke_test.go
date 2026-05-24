//go:build e2e

package e2e

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DARREN-2000/ai-hypervisor-platform/internal/api"
	"github.com/DARREN-2000/ai-hypervisor-platform/internal/orchestrator"
	"github.com/DARREN-2000/ai-hypervisor-platform/internal/models"
	"github.com/sirupsen/logrus"
)

type smokeResourceMonitor struct{}

func (smokeResourceMonitor) GetNodeMetrics(context.Context, string) (*models.ResourceMetrics, error) { return &models.ResourceMetrics{}, nil }
func (smokeResourceMonitor) GetVMMetrics(context.Context, string) (*models.ResourceMetrics, error) { return &models.ResourceMetrics{}, nil }
func (smokeResourceMonitor) GetClusterMetrics(context.Context) (*orchestrator.ClusterMetrics, error) { return &orchestrator.ClusterMetrics{}, nil }
func (smokeResourceMonitor) CollectMetrics(context.Context, string) error { return nil }
func (smokeResourceMonitor) PredictResourceUsage(context.Context, string) (*models.ResourceMetrics, error) { return &models.ResourceMetrics{}, nil }

func TestAPISmokeEndpoints(t *testing.T) {
	server := api.NewAPIServer(&api.Config{Address: "127.0.0.1", Port: 0}, logrus.New())
	server.SetMetrics(nil)
	server.SetHealthChecks(map[string]api.HealthCheck{
		"database": func(context.Context) error { return nil },
		"nats":     func(context.Context) error { return nil },
	})
	server.SetDependencies(nil, nil, nil, nil, smokeResourceMonitor{}, nil, nil, nil)
	server.RegisterRoutes()

	ts := httptest.NewServer(server.Handler())
	defer ts.Close()

	assertStatusAndBodyContains(t, ts.URL+"/health", http.StatusOK, []string{"healthy", "database", "nats"})
	assertStatusAndBodyContains(t, ts.URL+"/ready", http.StatusOK, []string{"healthy", "database", "nats"})
	assertStatusAndBodyContains(t, ts.URL+"/live", http.StatusOK, []string{"alive"})
}

func assertStatusAndBodyContains(t *testing.T, url string, expectedStatus int, expectedSnippets []string) {
	t.Helper()

	resp, err := http.Get(url)
	if err != nil {
		t.Fatalf("GET %s: %v", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != expectedStatus {
		t.Fatalf("GET %s: expected status %d, got %d", url, expectedStatus, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read %s: %v", url, err)
	}

	var payload any
	if err := json.Unmarshal(body, &payload); err != nil {
		t.Fatalf("decode %s JSON: %v", url, err)
	}

	text := string(body)
	for _, snippet := range expectedSnippets {
		if !strings.Contains(text, snippet) {
			t.Fatalf("expected %s response to contain %q, got %s", url, snippet, text)
		}
	}
}
