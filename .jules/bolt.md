## 2024-06-16 - Early Return Loop Breaks
**Learning:** Found a missing early return in the inner loop of `selectGPUsForVM` where we continued iterating over all available GPUs even after finding enough matching GPUs (`len(matched) == req.Count`).
**Action:** Always add early return `break` statements to stop iterating when sufficient matches are found to prevent unnecessary computation.
## 2026-06-19 - Pre-allocate slice capacity
**Learning:** Unnecessary array reallocation occurs when appending to slices without pre-allocated capacity inside a loop.
**Action:** Always specify the required capacity when initializing slices with `make` if the size is known, especially within loops, to avoid the overhead of Go's runtime array resizing.
## 2024-06-25 - Defer Expensive Operations
**Learning:** Found an unnecessary database query and metric retrieval () happening inside a loop *before* checking if the host even met the basic resource requirements.
**Action:** Always verify fast, in-memory basic requirements first before executing expensive operations like database queries or network calls within iterative loops.
## 2024-06-25 - Defer Expensive Operations
**Learning:** Found an unnecessary database query and metric retrieval happening inside a loop before checking if the host even met the basic resource requirements.
**Action:** Always verify fast, in-memory basic requirements first before executing expensive operations like database queries or network calls within iterative loops.
