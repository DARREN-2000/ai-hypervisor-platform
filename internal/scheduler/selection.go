package scheduler

import (
	"fmt"
	"sort"
	"strings"

	"github.com/DARREN-2000/ai-hypervisor-platform/internal/models"
)

func selectGPUsForVM(requests []models.GPURequest, gpus []*models.GPU, strategy GPUSelectionStrategy) ([]*models.GPU, error) {
	if len(requests) == 0 {
		return nil, nil
	}

	available := filterAvailableGPUs(gpus)
	if len(available) == 0 {
		return nil, fmt.Errorf("no available gpus")
	}

	sort.SliceStable(available, func(i, j int) bool {
		memI := gpuAvailableMemoryGB(available[i])
		memJ := gpuAvailableMemoryGB(available[j])
		if strategy == GPUSelectionSpread {
			return memI > memJ
		}
		return memI < memJ
	})

	selected := make([]*models.GPU, 0)
	// ⚡ Bolt: Use a map to track used GPUs instead of repeated slice allocations
	used := make(map[string]bool)

	for _, req := range requests {
		if req.Count <= 0 {
			continue
		}
		pickedCount := 0
		for _, gpu := range available {
			if gpu == nil || used[gpu.ID] {
				continue
			}
			if gpuMatchesRequest(req, gpu) {
				selected = append(selected, gpu)
				used[gpu.ID] = true
				pickedCount++
				if pickedCount == req.Count {
					break
				}
			}
		}

		if pickedCount < req.Count {
			return nil, fmt.Errorf("insufficient gpu capacity")
		}
	}

	return selected, nil
}

func filterAvailableGPUs(gpus []*models.GPU) []*models.GPU {
	filtered := make([]*models.GPU, 0, len(gpus))
	for _, gpu := range gpus {
		if gpu == nil {
			continue
		}
		if gpu.Status != models.GPUStatusAvailable {
			continue
		}
		filtered = append(filtered, gpu)
	}
	return filtered
}

func gpuMatchesRequest(req models.GPURequest, gpu *models.GPU) bool {
	if gpu == nil {
		return false
	}
	if req.Type != "" && !strings.EqualFold(req.Type, gpu.Type) {
		return false
	}
	if req.Model != "" && !strings.EqualFold(req.Model, gpu.Model) {
		return false
	}
	if req.MinMemoryGB > 0 {
		available := gpuAvailableMemoryGB(gpu)
		if available < req.MinMemoryGB {
			return false
		}
	}
	for _, feature := range req.RequiredFeatures {
		if !gpuSupportsFeature(gpu, feature) {
			return false
		}
	}
	return true
}

func gpuAvailableMemoryGB(gpu *models.GPU) int {
	if gpu == nil {
		return 0
	}
	if gpu.Metrics != nil && gpu.Metrics.MemoryFree > 0 {
		return gpu.Metrics.MemoryFree / 1024
	}
	if gpu.VRAM > 0 {
		return gpu.VRAM
	}
	return 0
}

func gpuSupportsFeature(gpu *models.GPU, feature string) bool {
	feature = strings.ToLower(strings.TrimSpace(feature))
	switch feature {
	case "cuda":
		return gpu.Capabilities.CUDA
	case "tensor-cores", "tensor":
		return gpu.Capabilities.TensorCores
	case "rt-cores", "rt":
		return gpu.Capabilities.RTCores
	case "nvlink":
		return gpu.Capabilities.NVLinkSupported
	case "mig":
		return gpu.Capabilities.MIGSupported
	default:
		return true
	}
}

