## 2024-06-16 - Early Return Loop Breaks
**Learning:** Found a missing early return in the inner loop of `selectGPUsForVM` where we continued iterating over all available GPUs even after finding enough matching GPUs (`len(matched) == req.Count`).
**Action:** Always add early return `break` statements to stop iterating when sufficient matches are found to prevent unnecessary computation.
## 2026-06-19 - Pre-allocate slice capacity
**Learning:** Unnecessary array reallocation occurs when appending to slices without pre-allocated capacity inside a loop.
**Action:** Always specify the required capacity when initializing slices with `make` if the size is known, especially within loops, to avoid the overhead of Go's runtime array resizing.
## 2026-06-21 - Cheap checks before expensive calls
**Learning:** Performing expensive operations (like DB queries or external API calls) before verifying if a candidate meets baseline criteria causes unnecessary performance overhead, especially within loops.
**Action:** Always perform cheap local checks (like basic resource capacity comparison using in-memory structs) before making expensive external calls in loops to fail fast and avoid unnecessary computation.
