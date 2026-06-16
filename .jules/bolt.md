## 2024-06-16 - Early Return Loop Breaks
**Learning:** Found a missing early return in the inner loop of `selectGPUsForVM` where we continued iterating over all available GPUs even after finding enough matching GPUs (`len(matched) == req.Count`).
**Action:** Always add early return `break` statements to stop iterating when sufficient matches are found to prevent unnecessary computation.
