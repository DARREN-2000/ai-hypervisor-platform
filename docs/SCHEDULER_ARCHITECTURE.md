# GPU-Aware Scheduler Architecture

This document describes the GPU-aware scheduling subsystem for the AI Hypervisor Platform.

## Goals

- Place GPU-bound workloads with awareness of GPU memory and topology
- Balance cluster utilization while honoring policy and constraints
- Produce deterministic, explainable placement decisions
- Provide hooks for monitoring and future autoscaling

## High-Level Flow

```
API Request
   |
   v
VM Manager ----> Scheduler
                  |
                  v
     +----------------------------+
     | Snapshot Hosts + GPUs      |
     | - Host capacity/alloc      |
     | - GPU inventory + metrics  |
     | - Optional node metrics    |
     +----------------------------+
                  |
                  v
     +----------------------------+
     | Policy Scoring Engine      |
     | - bin-packing              |
     | - spread                   |
     | - numa-aware               |
     +----------------------------+
                  |
                  v
     +----------------------------+
     | Placement Decision         |
     | - Selected host            |
     | - Selected GPU IDs         |
     | - Alternatives + reasons   |
     +----------------------------+
                  |
                  v
        VM Manager + GPU Orchestrator
```

## Core Components

### Scheduler Service
- Accepts a `VirtualMachine` and produces a `SchedulingDecision`
- Evaluates CPU, memory, disk, and GPU capacity
- Computes GPU memory fit per request
- Returns a deterministic selection with alternatives

### GPU Resource Tracking
- GPU repository provides device inventory
- Memory awareness via `GPU.Metrics.MemoryFree` (MB) or `GPU.VRAM` (GB)
- Feature-aware matching (CUDA, tensor cores, NVLink, MIG)

### Policy Engine
- `bin-packing`: pack workloads to fewer nodes
- `spread`: distribute workloads evenly
- `numa-aware`: prefer NUMA-local placement when metadata is present
- Policies are pluggable via the `Policy` interface

### Monitoring Hooks
- `Monitor` interface records decisions and failures
- Supports emission to metrics, tracing, or audit logs

### Autoscaling Signals
- `AutoscalerHook` interface emits signals when capacity is exhausted
- Contains shortage indicators for CPU, memory, and GPU

## Extensibility

- Add new policies by implementing `Policy`
- Attach external scoring sources via `ResourceMonitor`
- Extend GPU selection strategy without changing API
- Emit autoscaling hints for future capacity automation

## Placement Contracts

- Decisions include:
  - Selected host ID
  - Selected GPU IDs
  - Policy name, score, and reason
  - Alternative hosts with scores
- Decisions are idempotent and traceable

## Scheduling Interfaces

- `orchestrator.Scheduler` for VM placement
- `Policy` for scoring implementations
- `Monitor` for observability
- `AutoscalerHook` for scaling signals

## Implementation Notes

- Concurrency-safe metrics and policy registry
- Deterministic sorting and tie-breaking
- Validation and error handling for missing dependencies
- Compatible with VM Manager and API handler integration
