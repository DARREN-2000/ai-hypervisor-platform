package scheduler

import (
	"testing"
	"github.com/DARREN-2000/ai-hypervisor-platform/internal/models"
)

func TestBinPackPolicy(t *testing.T) {
	policy := BinPackPolicy{}
	if policy.Name() != "bin-packing" {
		t.Errorf("expected bin-packing, got %s", policy.Name())
	}
	if policy.GPUSelectionStrategy() != GPUSelectionPack {
		t.Errorf("expected pack, got %v", policy.GPUSelectionStrategy())
	}

	input := ScoreInput{
		Snapshot: HostSnapshot{
			Host: &models.HostNode{
				Capacity:           models.HostCapacity{CPU: 100, Memory: 100, DiskGB: 1000, GPUSlots: 4},
				AllocatedResources: models.HostAllocated{CPU: 50, Memory: 50, DiskGB: 500, GPUSlots: 2},
			},
		},
		Demand: ResourceDemand{CPU: 10, MemoryGB: 10, DiskGB: 100, GPUSlots: 1},
	}

	res, err := policy.Score(input)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if res.Score <= 0 || res.Score > 1 {
		t.Errorf("invalid score: %f", res.Score)
	}
}

func TestSpreadPolicy(t *testing.T) {
	policy := SpreadPolicy{}
	if policy.Name() != "spread" {
		t.Errorf("expected spread, got %s", policy.Name())
	}

	input := ScoreInput{
		Snapshot: HostSnapshot{
			Host: &models.HostNode{
				Capacity:           models.HostCapacity{CPU: 100, Memory: 100, DiskGB: 1000, GPUSlots: 4},
				AllocatedResources: models.HostAllocated{CPU: 50, Memory: 50, DiskGB: 500, GPUSlots: 2},
			},
		},
		Demand: ResourceDemand{CPU: 10, MemoryGB: 10, DiskGB: 100, GPUSlots: 1},
	}

	res, err := policy.Score(input)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if res.Score <= 0 || res.Score > 1 {
		t.Errorf("invalid score: %f", res.Score)
	}
}

func TestNUMAAwarePolicy(t *testing.T) {
	policy := NUMAAwarePolicy{}
	if policy.Name() != "numa-aware" {
		t.Errorf("expected numa-aware, got %s", policy.Name())
	}

	input := ScoreInput{
		VM: &models.VirtualMachine{
			Metadata: map[string]string{"numa_node": "0"},
		},
		Snapshot: HostSnapshot{
			Host: &models.HostNode{
				Metadata: map[string]string{"numa_nodes": "0,1"},
				Capacity:           models.HostCapacity{CPU: 100, Memory: 100, DiskGB: 1000, GPUSlots: 4},
				AllocatedResources: models.HostAllocated{CPU: 50, Memory: 50, DiskGB: 500, GPUSlots: 2},
			},
		},
		Demand: ResourceDemand{CPU: 10, MemoryGB: 10, DiskGB: 100, GPUSlots: 1},
	}

	res, err := policy.Score(input)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if res.Score <= 0 || res.Score > 1 {
		t.Errorf("invalid score: %f", res.Score)
	}
}

func TestUtilizationSet(t *testing.T) {
	set := utilizationSet(ScoreInput{})
	if len(set) != 4 || set[0] != 0 || set[1] != 0 || set[2] != 0 || set[3] != 0 {
		t.Error("expected zeros for nil host")
	}
}

func TestWeightedAverage(t *testing.T) {
	if weightedAverage([]float64{1, 1}, defaultWeights) != 0 {
		t.Error("expected 0 for len < 4")
	}
	zeroWeights := weights{}
	if weightedAverage([]float64{1, 1, 1, 1}, zeroWeights) != 0 {
		t.Error("expected 0 for 0 weight sum")
	}
}

func TestUtilization(t *testing.T) {
	if utilization(0, 5, 5) != 1.0 {
		t.Error("expected 1.0 for <= 0 capacity")
	}
}

func TestBalanceScore(t *testing.T) {
	if balanceScore(nil) != 0 {
		t.Error("expected 0 for empty values")
	}
}

func TestClamp01(t *testing.T) {
	if clamp01(-1) != 0 {
		t.Error("expected 0")
	}
	if clamp01(2) != 1 {
		t.Error("expected 1")
	}
	if clamp01(0.5) != 0.5 {
		t.Error("expected 0.5")
	}
}

func TestNumaBonus(t *testing.T) {
	if numaBonus(ScoreInput{}) != 0 {
		t.Error("expected 0 for nil VM/Host")
	}

	vm := &models.VirtualMachine{}
	host := &models.HostNode{}
	if numaBonus(ScoreInput{VM: vm, Snapshot: HostSnapshot{Host: host}}) != 0 {
		t.Error("expected 0 for no metadata")
	}

	vm.Metadata = map[string]string{"other": "1"}
	if numaBonus(ScoreInput{VM: vm, Snapshot: HostSnapshot{Host: host}}) != 0 {
		t.Error("expected 0 for no numa_node on vm")
	}

	vm.Metadata["numa_node"] = "0"
	if numaBonus(ScoreInput{VM: vm, Snapshot: HostSnapshot{Host: host}}) != 0 {
		t.Error("expected 0 for no metadata on host")
	}

	host.Metadata = map[string]string{"other": "1"}
	if numaBonus(ScoreInput{VM: vm, Snapshot: HostSnapshot{Host: host}}) != 0 {
		t.Error("expected 0 for no numa_nodes on host")
	}

	host.Metadata["numa_nodes"] = "1,2"
	if numaBonus(ScoreInput{VM: vm, Snapshot: HostSnapshot{Host: host}}) != 0 {
		t.Error("expected 0 for mismatch")
	}
}
