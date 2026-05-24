package scheduler

import (
	"fmt"
	"math"
	"strings"
)

type weights struct {
	cpu    float64
	memory float64
	disk   float64
	gpu    float64
}

var defaultWeights = weights{cpu: 0.4, memory: 0.4, disk: 0.1, gpu: 0.1}

// BinPackPolicy packs workloads onto fewer hosts.
type BinPackPolicy struct{}

func (BinPackPolicy) Name() string { return "bin-packing" }
func (BinPackPolicy) GPUSelectionStrategy() GPUSelectionStrategy { return GPUSelectionPack }

func (BinPackPolicy) Score(input ScoreInput) (ScoreResult, error) {
	util := utilizationSet(input)
	avgUtil := weightedAverage(util, defaultWeights)
	balance := balanceScore(util)
	score := clamp01(avgUtil*0.9 + balance*0.1)
	return ScoreResult{
		Score:  score,
		Reason: fmt.Sprintf("bin-pack util=%.2f", avgUtil),
	}, nil
}

// SpreadPolicy distributes workloads across hosts.
type SpreadPolicy struct{}

func (SpreadPolicy) Name() string { return "spread" }
func (SpreadPolicy) GPUSelectionStrategy() GPUSelectionStrategy { return GPUSelectionSpread }

func (SpreadPolicy) Score(input ScoreInput) (ScoreResult, error) {
	util := utilizationSet(input)
	avgUtil := weightedAverage(util, defaultWeights)
	balance := balanceScore(util)
	score := clamp01((1-avgUtil)*0.9 + balance*0.1)
	return ScoreResult{
		Score:  score,
		Reason: fmt.Sprintf("spread util=%.2f", avgUtil),
	}, nil
}

// NUMAAwarePolicy favors host-local NUMA placements when metadata is available.
type NUMAAwarePolicy struct{}

func (NUMAAwarePolicy) Name() string { return "numa-aware" }
func (NUMAAwarePolicy) GPUSelectionStrategy() GPUSelectionStrategy { return GPUSelectionPack }

func (NUMAAwarePolicy) Score(input ScoreInput) (ScoreResult, error) {
	util := utilizationSet(input)
	avgUtil := weightedAverage(util, defaultWeights)
	base := avgUtil
	bonus := numaBonus(input)
	score := clamp01(base + bonus)
	return ScoreResult{
		Score:  score,
		Reason: fmt.Sprintf("numa util=%.2f bonus=%.2f", avgUtil, bonus),
	}, nil
}

func utilizationSet(input ScoreInput) []float64 {
	host := input.Snapshot.Host
	if host == nil {
		return []float64{0, 0, 0, 0}
	}

	cpu := utilization(host.Capacity.CPU, host.AllocatedResources.CPU, input.Demand.CPU)
	mem := utilization(host.Capacity.Memory, host.AllocatedResources.Memory, input.Demand.MemoryGB)
	disk := utilization(host.Capacity.DiskGB, host.AllocatedResources.DiskGB, input.Demand.DiskGB)
	gpu := utilization(host.Capacity.GPUSlots, host.AllocatedResources.GPUSlots, input.Demand.GPUSlots)
	return []float64{cpu, mem, disk, gpu}
}

func weightedAverage(values []float64, w weights) float64 {
	if len(values) < 4 {
		return 0
	}
	weightSum := w.cpu + w.memory + w.disk + w.gpu
	if weightSum == 0 {
		return 0
	}
	return (values[0]*w.cpu + values[1]*w.memory + values[2]*w.disk + values[3]*w.gpu) / weightSum
}

func utilization(capacity, allocated, demand int) float64 {
	if capacity <= 0 {
		return 1
	}
	value := float64(allocated+demand) / float64(capacity)
	return clamp01(value)
}

func balanceScore(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	mean := 0.0
	for _, v := range values {
		mean += v
	}
	mean = mean / float64(len(values))

	variance := 0.0
	for _, v := range values {
		delta := v - mean
		variance += delta * delta
	}
	variance = variance / float64(len(values))
	return clamp01(1 - math.Sqrt(variance))
}

func clamp01(value float64) float64 {
	if value < 0 {
		return 0
	}
	if value > 1 {
		return 1
	}
	return value
}

func numaBonus(input ScoreInput) float64 {
	if input.VM == nil || input.Snapshot.Host == nil {
		return 0
	}
	if input.VM.Metadata == nil {
		return 0
	}
	vmNode := input.VM.Metadata["numa_node"]
	if vmNode == "" {
		return 0
	}
	if input.Snapshot.Host.Metadata == nil {
		return 0
	}
	hostNodes := input.Snapshot.Host.Metadata["numa_nodes"]
	if hostNodes == "" {
		return 0
	}
	for _, node := range strings.Split(hostNodes, ",") {
		if strings.TrimSpace(node) == vmNode {
			return 0.1
		}
	}
	return 0
}
