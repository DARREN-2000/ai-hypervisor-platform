## 2024-06-16 - Early Return Loop Breaks
**Learning:** Found a missing early return in the inner loop of `selectGPUsForVM` where we continued iterating over all available GPUs even after finding enough matching GPUs (`len(matched) == req.Count`).
**Action:** Always add early return `break` statements to stop iterating when sufficient matches are found to prevent unnecessary computation.
## 2026-06-19 - Pre-allocate slice capacity
**Learning:** Unnecessary array reallocation occurs when appending to slices without pre-allocated capacity inside a loop.
**Action:** Always specify the required capacity when initializing slices with `make` if the size is known, especially within loops, to avoid the overhead of Go's runtime array resizing.
## 2026-06-23 - N+1 Query Avoidance in Scheduler
**Learning:** Checking in-memory conditions before making expensive data snapshot queries avoids N+1 queries. We were building a snapshot (which fetches GPUs and metrics from DB/monitor) for *all* hosts before checking if they even had the basic resources (CPU/Memory/Disk) to fit the VM demand.
**Action:** Always filter candidates using simple, in-memory checks (`fitsResources`) before executing expensive secondary queries or API calls (`buildSnapshot`).
