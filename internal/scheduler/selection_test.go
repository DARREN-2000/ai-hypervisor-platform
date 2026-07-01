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

func TestGPUMatchesRequest(t *testing.T) {
	if gpuMatchesRequest(models.GPURequest{}, nil) {
		t.Error("expected false for nil gpu")
	}

	gpu := &models.GPU{
		Type: "nvidia", Model: "A100", VRAM: 40,
		Capabilities: models.GPUCapabilities{
			CUDA: true, TensorCores: true,
		},
	}
	if !gpuMatchesRequest(models.GPURequest{}, gpu) {
		t.Error("expected true for empty request")
	}
	if !gpuMatchesRequest(models.GPURequest{Type: "nvidia"}, gpu) {
		t.Error("expected true for matching type")
	}
	if gpuMatchesRequest(models.GPURequest{Type: "amd"}, gpu) {
		t.Error("expected false for mismatched type")
	}
	if !gpuMatchesRequest(models.GPURequest{Model: "A100"}, gpu) {
		t.Error("expected true for matching model")
	}
	if gpuMatchesRequest(models.GPURequest{Model: "V100"}, gpu) {
		t.Error("expected false for mismatched model")
	}
	if !gpuMatchesRequest(models.GPURequest{MinMemoryGB: 30}, gpu) {
		t.Error("expected true for sufficient memory")
	}
	if gpuMatchesRequest(models.GPURequest{MinMemoryGB: 50}, gpu) {
		t.Error("expected false for insufficient memory")
	}
	if !gpuMatchesRequest(models.GPURequest{RequiredFeatures: []string{"cuda", "tensor-cores"}}, gpu) {
		t.Error("expected true for supported features")
	}
	if gpuMatchesRequest(models.GPURequest{RequiredFeatures: []string{"rt-cores"}}, gpu) {
		t.Error("expected false for unsupported feature")
	}
}

func TestGPUAvailableMemoryGB(t *testing.T) {
	if gpuAvailableMemoryGB(nil) != 0 {
		t.Error("expected 0 for nil gpu")
	}

	gpu := &models.GPU{VRAM: 40}
	if gpuAvailableMemoryGB(gpu) != 40 {
		t.Error("expected 40 from VRAM")
	}

	gpu.Metrics = &models.GPUMetrics{MemoryFree: 16384} // 16 GB
	if gpuAvailableMemoryGB(gpu) != 16 {
		t.Error("expected 16 from Metrics")
	}
}

func TestGPUSupportsFeature(t *testing.T) {
	gpu := &models.GPU{
		Capabilities: models.GPUCapabilities{
			CUDA: true, TensorCores: true, RTCores: false,
			NVLinkSupported: true, MIGSupported: false,
		},
	}
	tests := []struct{
		feature string
		expected bool
	}{
		{"cuda", true},
		{"CUDA", true},
		{" tensor-cores ", true},
		{"tensor", true},
		{"rt-cores", false},
		{"rt", false},
		{"nvlink", true},
		{"mig", false},
		{"unknown", true},
	}
	for _, tt := range tests {
		if gpuSupportsFeature(gpu, tt.feature) != tt.expected {
			t.Errorf("feature %s: expected %v", tt.feature, tt.expected)
		}
	}
}

func TestFilterAvailableGPUs(t *testing.T) {
	gpus := []*models.GPU{
		nil,
		{Status: models.GPUStatusAllocated},
		{Status: models.GPUStatusAvailable, ID: "g1"},
	}
	filtered := filterAvailableGPUs(gpus)
	if len(filtered) != 1 || filtered[0].ID != "g1" {
		t.Error("unexpected filtered result")
	}
}

func TestSubtractGPUs(t *testing.T) {
	all := []*models.GPU{
		nil,
		{ID: "g1"}, {ID: "g2"}, {ID: "g3"},
	}
	remove := []*models.GPU{
		nil,
		{ID: "g2"},
	}
	res := subtractGPUs(all, remove)
	if len(res) != 2 || res[0].ID != "g1" || res[1].ID != "g3" {
		t.Error("unexpected subtract result")
	}
}

func TestSelectGPUsForVM_EmptyRequests(t *testing.T) {
	res, err := selectGPUsForVM(nil, nil, GPUSelectionPack)
	if err != nil || res != nil {
		t.Error("expected nil, nil")
	}
}
