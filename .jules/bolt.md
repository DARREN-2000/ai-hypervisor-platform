## 2024-06-16 - Early Return Loop Breaks
**Learning:** Found a missing early return in the inner loop of `selectGPUsForVM` where we continued iterating over all available GPUs even after finding enough matching GPUs (`len(matched) == req.Count`).
**Action:** Always add early return `break` statements to stop iterating when sufficient matches are found to prevent unnecessary computation.
## 2026-06-19 - Pre-allocate slice capacity
**Learning:** Unnecessary array reallocation occurs when appending to slices without pre-allocated capacity inside a loop.
**Action:** Always specify the required capacity when initializing slices with `make` if the size is known, especially within loops, to avoid the overhead of Go's runtime array resizing.
## 2024-06-24 - Avoid strings.ToLower allocation in feature matching
**Learning:** `strings.ToLower` allocates memory when converting characters to lowercase, which caused a minor bottleneck during frequent capability checks in `gpuSupportsFeature`.
**Action:** Use `strings.EqualFold` instead for case-insensitive matching in high-frequency checks to skip memory allocation.
