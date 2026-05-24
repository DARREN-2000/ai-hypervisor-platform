package scheduler

import (
	"testing"

	"github.com/DARREN-2000/ai-hypervisor-platform/internal/models"
)

func TestSelectGPUsForVMPrefersLowerMemoryForPackAndHigherMemoryForSpread(t *testing.T) {
	requests := []models.GPURequest{{Type: "nvidia", Model: "A100", Count: 2, MinMemoryGB: 16, RequiredFeatures: []string{"cuda"}}}
	gpus := []*models.GPU{
		{ID: "gpu-32", Type: "nvidia", Model: "A100", Status: models.GPUStatusAvailable, VRAM: 32, Capabilities: models.GPUCapabilities{CUDA: true}},
		{ID: "gpu-24", Type: "nvidia", Model: "A100", Status: models.GPUStatusAvailable, VRAM: 24, Capabilities: models.GPUCapabilities{CUDA: true}},
		{ID: "gpu-offline", Type: "nvidia", Model: "A100", Status: models.GPUStatusAllocated, VRAM: 80, Capabilities: models.GPUCapabilities{CUDA: true}},
	}

	packed, err := selectGPUsForVM(requests, gpus, GPUSelectionPack)
	if err != nil {
		t.Fatalf("pack selection failed: %v", err)
	}
	if len(packed) != 2 || packed[0].ID != "gpu-24" || packed[1].ID != "gpu-32" {
		t.Fatalf("unexpected pack selection order: %#v", gpuIDs(packed))
	}

	spread, err := selectGPUsForVM(requests, gpus, GPUSelectionSpread)
	if err != nil {
		t.Fatalf("spread selection failed: %v", err)
	}
	if len(spread) != 2 || spread[0].ID != "gpu-32" || spread[1].ID != "gpu-24" {
		t.Fatalf("unexpected spread selection order: %#v", gpuIDs(spread))
	}
}

func TestSelectGPUsForVMReturnsCapacityErrorWhenMatchIsInsufficient(t *testing.T) {
	requests := []models.GPURequest{{Type: "nvidia", Count: 2, RequiredFeatures: []string{"cuda"}}}
	gpus := []*models.GPU{
		{ID: "gpu-1", Type: "nvidia", Status: models.GPUStatusAvailable, VRAM: 16, Capabilities: models.GPUCapabilities{CUDA: true}},
	}

	selected, err := selectGPUsForVM(requests, gpus, GPUSelectionPack)
	if err == nil {
		t.Fatal("expected capacity error")
	}
	if selected != nil {
		t.Fatalf("expected no selection, got %#v", selected)
	}
}

func gpuIDs(gpus []*models.GPU) []string {
	ids := make([]string, 0, len(gpus))
	for _, gpu := range gpus {
		if gpu != nil {
			ids = append(ids, gpu.ID)
		}
	}
	return ids
}
