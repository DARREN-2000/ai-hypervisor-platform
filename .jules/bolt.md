## 2024-06-16 - Early Return Loop Breaks
**Learning:** Found a missing early return in the inner loop of `selectGPUsForVM` where we continued iterating over all available GPUs even after finding enough matching GPUs (`len(matched) == req.Count`).
**Action:** Always add early return `break` statements to stop iterating when sufficient matches are found to prevent unnecessary computation.
## 2026-06-18 - Prevent Expensive Host Snapshots Before Capacity Checks
**Learning:** In `internal/scheduler/service.go`, generating a host snapshot requires expensive API/database queries to get metrics and GPU statuses. Doing this before checking if the host even has the local capacity to accommodate a VM is wasteful, especially when iterating over many hosts.
**Action:** Always check local conditions (like simple capacity math via `fitsResources`) before proceeding to make expensive external queries (like `buildSnapshot`). This follows the early return / fast failure pattern.
