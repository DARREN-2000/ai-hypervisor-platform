## 2024-06-17 - GPU Selection Allocation Overhead
**Learning:** The initial implementation of `selectGPUsForVM` heavily utilized slice re-allocations and a `subtractGPUs` helper mapping, resulting in an O(N^2) complexity with high garbage collection overhead in hot scheduling paths.
**Action:** Replaced iterative slice modifications with a simple boolean `used` map and a short-circuiting counter variable to achieve O(N) selection performance. This prevents large allocation spikes when processing large pools of GPUs. Also, guard against nil pointer dereferences when scanning slices of structs.
