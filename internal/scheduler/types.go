package scheduler

import (
	"time"

	"github.com/DARREN-2000/ai-hypervisor-platform/internal/models"
)

// HostSnapshot captures placement-relevant host data.
type HostSnapshot struct {
	Host    *models.HostNode
	GPUs    []*models.GPU
	Metrics *models.ResourceMetrics
}

// ResourceDemand captures requested resources for placement.
type ResourceDemand struct {
	CPU         int
	MemoryGB    int
	DiskGB      int
	GPUSlots    int
	GPUMemoryGB int
}

// ScoreInput provides policy inputs.
type ScoreInput struct {
	VM        *models.VirtualMachine
	Snapshot  HostSnapshot
	Demand    ResourceDemand
	GPUChoice []*models.GPU
	Now       time.Time
}

// ScoreResult is the output of a policy evaluation.
type ScoreResult struct {
	Score    float64
	Reason   string
	Metadata map[string]string
}

// GPUSelectionStrategy controls how GPUs are chosen.
type GPUSelectionStrategy string

const (
	GPUSelectionPack   GPUSelectionStrategy = "pack"
	GPUSelectionSpread GPUSelectionStrategy = "spread"
)

// Policy scores candidate nodes.
type Policy interface {
	Name() string
	GPUSelectionStrategy() GPUSelectionStrategy
	Score(input ScoreInput) (ScoreResult, error)
}

// PlacementResult represents a scheduling outcome.
type PlacementResult struct {
	Decision      *models.SchedulingDecision
	SelectedHost  *models.HostNode
	SelectedGPUs  []*models.GPU
	Alternatives  []models.HostScore
}
