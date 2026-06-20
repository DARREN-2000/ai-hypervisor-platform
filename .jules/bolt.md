## 2024-06-16 - Early Return Loop Breaks
**Learning:** Found a missing early return in the inner loop of `selectGPUsForVM` where we continued iterating over all available GPUs even after finding enough matching GPUs (`len(matched) == req.Count`).
**Action:** Always add early return `break` statements to stop iterating when sufficient matches are found to prevent unnecessary computation.
## 2026-06-19 - Pre-allocate slice capacity
**Learning:** Unnecessary array reallocation occurs when appending to slices without pre-allocated capacity inside a loop.
**Action:** Always specify the required capacity when initializing slices with `make` if the size is known, especially within loops, to avoid the overhead of Go's runtime array resizing.
## 2024-05-19 - Fast Capacity Checks Before Expensive Queries
**Learning:** Found that we were generating a complete snapshot of host metrics and GPUs (which involves multiple DB/API calls) before doing a basic local capacity check (`fitsResources`). When capacity didn't fit, the expensive snapshot was thrown away.
**Action:** Always perform fast, local validation checks inside loops before initiating any expensive operations (DB queries, API calls) for the current iteration.
