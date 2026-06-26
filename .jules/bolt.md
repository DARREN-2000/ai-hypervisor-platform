## 2024-06-16 - Early Return Loop Breaks
**Learning:** Found a missing early return in the inner loop of `selectGPUsForVM` where we continued iterating over all available GPUs even after finding enough matching GPUs (`len(matched) == req.Count`).
**Action:** Always add early return `break` statements to stop iterating when sufficient matches are found to prevent unnecessary computation.
## 2026-06-19 - Pre-allocate slice capacity
**Learning:** Unnecessary array reallocation occurs when appending to slices without pre-allocated capacity inside a loop.
**Action:** Always specify the required capacity when initializing slices with `make` if the size is known, especially within loops, to avoid the overhead of Go's runtime array resizing.

## 2026-06-26 - Fast-path Filtering Before Expensive Operations
**Learning:** In the scheduler loop, calling `buildSnapshot` (which performs DB queries and network calls) before a simple, in-memory capacity check (`fitsResources`) caused massive overhead when rejecting invalid hosts.
**Action:** Always place cheap, fast-path filter conditions before expensive operations like network requests or database queries when filtering candidates in a loop.
