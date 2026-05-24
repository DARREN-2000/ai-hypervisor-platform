package gpu

import "strings"

// AllocationPolicy describes the GPU allocation strategy.
type AllocationPolicy string

const (
	AllocationPolicyBinPacking AllocationPolicy = "bin-packing"
	AllocationPolicySpread     AllocationPolicy = "spread"
	AllocationPolicyNUMAAware  AllocationPolicy = "numa-aware"
)

// Request describes a GPU reservation request.
type Request struct {
	Count           int
	MinimumMemoryGB int
	RequireMIG      bool
	RequireShared   bool
	PreferredPolicy AllocationPolicy
}

// Decision captures the outcome of a GPU placement attempt.
type Decision struct {
	Policy   AllocationPolicy
	GPUCount int
	Accepted bool
	Reason   string
}

// NormalizePolicy returns a canonical policy value.
func NormalizePolicy(value string) AllocationPolicy {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case string(AllocationPolicySpread):
		return AllocationPolicySpread
	case string(AllocationPolicyNUMAAware):
		return AllocationPolicyNUMAAware
	default:
		return AllocationPolicyBinPacking
	}
}

// ValidateRequest checks whether the request contains a usable GPU count.
func ValidateRequest(request Request) bool {
	return request.Count > 0
}
