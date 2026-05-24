package monitoring

import "time"

// Snapshot captures a point-in-time summary for a resource.
type Snapshot struct {
	Name      string
	Healthy   bool
	UpdatedAt time.Time
	Details   map[string]float64
}

// Series represents a lightweight time-series sample.
type Series struct {
	Metric string
	Value  float64
	At     time.Time
}

// Target describes a monitored service or host.
type Target struct {
	Name    string
	Address string
	Labels  map[string]string
}

// Summary combines the current state of a target set.
type Summary struct {
	HealthyTargets   int
	UnhealthyTargets int
	TotalTargets     int
	GeneratedAt      time.Time
}
