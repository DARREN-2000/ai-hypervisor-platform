# VM Lifecycle and GPU Allocation Strategy

## VM Lifecycle State Machine

### State Transitions

The VM lifecycle follows a strict state machine to ensure consistency and prevent invalid operations.

```
┌─────────────────────────────────────────────────────────────────┐
│                                                                   │
│                    CREATED                                      │
│                      │                                           │
│                      ↓                                           │
│              PROVISIONING ←──────── SCALING                     │
│                      │                ↑                          │
│                      ├────────────────┘                          │
│                      ↓                                           │
│                   RUNNING                                        │
│                      │                                           │
│                      ├──────────────────→ FAILED                │
│                      ↓                        ↓                  │
│                   DELETING ←──────────────────┘                 │
│                      ↓                                           │
│                   DELETED                                        │
│                                                                   │
└─────────────────────────────────────────────────────────────────┘
```

### State Descriptions

**Created**
- Initial state when VM is first created
- VM definition persisted to database
- No infrastructure provisioned
- No resources allocated
- Can transition to: PROVISIONING

**Provisioning**
- Infrastructure provisioning in progress
- Libvirt domain being created
- GPUs being allocated
- Network and storage being configured
- If error occurs: transitions to FAILED
- On success: transitions to RUNNING
- Can transition to: RUNNING, FAILED

**Running**
- VM operational and accepting workloads
- All infrastructure provisioned
- GPUs allocated and available
- Resources actively used
- Can transition to: SCALING, DELETING, FAILED

**Scaling**
- VM resources being modified (CPU, memory, GPU)
- Brief interruption expected
- Network still active
- After completion: transitions to RUNNING
- Can transition to: RUNNING, FAILED

**Failed**
- Error occurred during provisioning, operation, or execution
- Error details stored for debugging
- Resources may be partially allocated
- Manual intervention may be required
- Can transition to: DELETING, PROVISIONING (retry)
- Can transition to: DELETING

**Deleting**
- VM teardown in progress
- Resources being released
- GPUs deallocated
- Libvirt domain being destroyed
- Storage being cleaned up
- On success: transitions to DELETED
- If error: stays in DELETING and retries

**Deleted**
- VM completely removed
- All resources released
- No further transitions possible
- Terminal state

### State Transition Rules

1. **Only allowed transitions**: Only transitions listed in state descriptions are permitted
2. **Idempotency**: Transitioning to the same state is a no-op
3. **Error handling**: Failed operations transition to FAILED state
4. **Timeout handling**: Long-running transitions (> 30 min) may force FAILED state
5. **Audit logging**: All state transitions logged with actor, timestamp, and reason

## VM Provisioning Workflow

```
1. User API Request (Create VM)
   ↓
2. Validate Request
   - Check flavor exists
   - Check image available
   - Validate network config
   - Validate GPU requests
   ↓
3. Create VM Record
   - Persist to database with state=CREATED
   - Generate unique VM ID
   - Assign namespace/tenant
   ↓
4. Audit Log (VM_CREATED)
   ↓
5. Enqueue Provisioning Task
   - Create Task record with type=PROVISION_VM
   - Add to task queue
   - Emit event: vm.provisioning_queued
   ↓
6. Schedule Placement Decision
   - Call Scheduler service
   - Check node capacity
   - Check GPU availability
   - Apply anti-affinity rules
   - Return selected host
   ↓
7. Allocate GPU Resources
   - Call GPU Orchestrator
   - Reserve GPUs for VM
   - Generate GPU allocation records
   - Return device mappings
   ↓
8. Update VM State → PROVISIONING
   ↓
9. Execute on Host Agent
   - Generate libvirt XML from VM spec + GPU allocation
   - Generate network config
   - Prepare storage volumes
   - Call host agent to create domain
   ↓
10. Create Libvirt Domain
    - Define domain (not started)
    - Attach volumes
    - Attach network interfaces
    - Attach GPU devices
    - Configure memory and CPU
    ↓
11. Start Libvirt Domain
    - Power on VM
    - Boot OS
    - Wait for network readiness
    ↓
12. Verify VM Online
    - Ping VM network interface
    - Check SSH connectivity (if configured)
    - Verify GPU devices present (from guest OS)
    - Verify filesystem mounted
    ↓
13. Update VM State → RUNNING
    - Set host_node_id
    - Set last_status_check
    - Persist to database
    ↓
14. Update Task Status → COMPLETED
    - Set task result with VM details
    - Mark provisioning complete
    ↓
15. Emit Event: vm.running
    ↓
16. Return Success to API Client
```

## VM Scaling Workflow

```
1. User requests resize (CPU/memory)
   ↓
2. Validate resize request
   - New flavor must be valid
   - New size must be larger or smaller than current
   - Node must have capacity for new size
   ↓
3. Check GPU implications
   - GPU allocations typically don't change
   - Verify GPU resources still available
   ↓
4. Update VM State → SCALING
   ↓
5. For Online Resize (if supported):
   - Pause VM (optional)
   - Hot-resize CPU/memory via Libvirt
   - Resume VM
   ↓
6. For Offline Resize:
   - Stop VM
   - Modify domain configuration
   - Restart VM
   ↓
7. Verify new resources available
   - Check VM sees new CPU count
   - Check VM sees new memory
   ↓
8. Update VM flavor in database
   - Update flavor field
   - Set updated_at timestamp
   ↓
9. Update VM State → RUNNING
   ↓
10. Emit Event: vm.scaled
```

## VM Stop/Reboot Workflow

### Stop Workflow
```
1. Validate VM in RUNNING state
   ↓
2. Enqueue Stop Task
   ↓
3. Call host agent to stop domain
   - Graceful shutdown (SIGTERM to VM)
   - Wait up to 60 seconds
   - Force shutdown if needed (SIGKILL)
   ↓
4. Verify domain stopped
   ↓
5. Update VM metadata
   - Clear host assignments
   - Update last_status_check
   ↓
6. Emit Event: vm.stopped
```

### Reboot Workflow
```
1. Validate VM in RUNNING state
   ↓
2. Call host agent to reboot
   - Send reboot command to VM
   - Wait for VM to come back online
   ↓
3. Verify VM online
   ↓
4. Emit Event: vm.rebooted
```

## VM Deletion Workflow

```
1. Validate VM not already deleting
   ↓
2. Update VM State → DELETING
   ↓
3. Stop VM if running
   - Graceful shutdown
   - Force shutdown if needed
   ↓
4. Deallocate GPUs
   - Release GPU allocations
   - Clear device mappings
   ↓
5. Call host agent to destroy domain
   - Stop domain if running
   - Undefine domain
   - Remove transient volumes
   ↓
6. Clean up storage
   - Delete VM-specific volumes
   - Clean up snapshots
   ↓
7. Clean up networking
   - Release IP addresses
   - Remove interfaces
   ↓
8. Update VM State → DELETED
   - Clear host_node_id
   - Set terminated_at timestamp
   ↓
9. Emit Event: vm.deleted
   ↓
10. Audit Log (VM_DELETED)
```

## Error Handling in Workflows

### Retry Strategy

- **Maximum Retries**: 3 attempts for most operations
- **Backoff Strategy**: Exponential backoff (1s, 4s, 16s)
- **Retriable Errors**: Network timeouts, temporary resource unavailability
- **Non-Retriable Errors**: Invalid configuration, authentication failures

### Error Recovery

1. **Partial Provisioning Failure**
   - Deallocate any allocated GPUs
   - Clean up any created domains
   - Clean up partial storage
   - Set VM state to FAILED with error details

2. **Resource Exhaustion During Provisioning**
   - Immediately return error
   - Do not retry
   - Provide detailed resource breakdown in error message
   - Suggest alternative placements

3. **Timeout During Provisioning**
   - Set state to FAILED
   - Log error details
   - Do not automatic retry
   - Alert operator for manual intervention

4. **Host Failure During VM Runtime**
   - Detect via heartbeat failure
   - Mark host as DEGRADED
   - Attempt live migration (if supported)
   - If migration fails, mark VM as FAILED and alert operator

---

## GPU Allocation Strategy

### Architecture Overview

GPU allocation is a critical path operation that affects both provisioning latency and cluster efficiency. The strategy balances three competing concerns:

1. **Performance**: Fast allocation decisions
2. **Efficiency**: Optimal resource utilization
3. **Fairness**: Equitable distribution across tenants

### GPU Allocation Policies

#### Policy 1: Bin-Packing (Default)

Optimizes for resource consolidation and power efficiency.

```
Algorithm:
1. For each GPU request in VM spec:
   a. Get list of available GPUs matching type and capabilities
   b. Sort by: [Node Utilization] ascending, [GPU Utilization] ascending
   c. Select GPU with lowest combined utilization
   d. Mark GPU as reserved
   e. Record allocation
   
2. Verify:
   a. All GPU requests satisfied
   b. No GPU over-subscription
   c. All allocations on same or adjacent nodes (optional NUMA affinity)
```

**Pros**:
- Minimizes number of nodes used
- Reduces inter-node GPU communication overhead
- Improves power efficiency

**Cons**:
- Can lead to fragmentation on heavily loaded nodes
- May increase latency for large batches of small VMs

#### Policy 2: Spread

Optimizes for resilience and parallelism.

```
Algorithm:
1. For each GPU request in VM spec:
   a. Get list of available GPUs
   b. Sort by: [Node Load] ascending, [GPU Utilization] ascending
   c. Select GPU with lowest node load
   d. Mark GPU as reserved
   
2. Distribute GPU allocation across maximum number of nodes
```

**Pros**:
- Spreads failure domain
- Reduces single node dependency
- Better for I/O intensive workloads

**Cons**:
- Higher inter-node communication latency
- Less power efficient (more nodes powered on)

#### Policy 3: NUMA-Aware

Optimizes for latency-sensitive applications.

```
Algorithm:
1. Detect NUMA topology of requesting VM's target node
2. Prefer GPUs on same NUMA node as vCPUs
3. If insufficient on same node:
   a. Prefer adjacent NUMA nodes
   b. If still insufficient, spread across NUMA nodes
```

**Pros**:
- Minimized NUMA cross-node memory traffic
- Lower latency for GPU operations

**Cons**:
- More restrictive allocation
- May fail to allocate when other policies would succeed

### GPU Sharing Strategies

#### Strategy 1: No Sharing (Default)

Each GPU dedicated to a single VM.

```
- GPU reserved exclusively for one VM
- No multi-VM sharing
- Highest performance isolation
- Maximum latency predictability
```

#### Strategy 2: MIG (Multi-Instance GPU)

NVIDIA MIG profiles for fine-grained GPU partitioning.

```
Supported Profiles:
- 1g.5gb: 1/7 GPU with 5GB memory
- 2g.10gb: 2/7 GPU with 10GB memory
- 3g.20gb: 3/7 GPU with 20GB memory
- 4g.20gb: 4/7 GPU with 20GB memory
- 7g.40gb: Full GPU (no sharing)

Allocation:
1. Check if MIG supported and enabled
2. Select profile matching VM GPU requirements
3. Allocate MIG instance
4. Verify isolation
```

#### Strategy 3: Time-Sliced Sharing (Future)

Virtualization layer for GPU time-multiplexing.

```
- Schedule multiple workloads on single GPU
- Context switching between VMs
- Suitable for bursty, non-interactive workloads
- Requires hypervisor support
```

### Allocation Algorithm Pseudocode

```go
func AllocateGPUs(vm *VM, policy *GPUAllocationPolicy) ([]GPUAllocation, error) {
    allocations := []GPUAllocation{}
    
    for _, request := range vm.GPURequests {
        // Find candidate GPUs
        candidates := findAvailableGPUs(request, policy)
        
        if len(candidates) < request.Count {
            return nil, ErrInsufficientGPUs
        }
        
        // Apply allocation policy
        selected := applyAllocationPolicy(candidates, request, policy)
        
        // Verify NUMA affinity if required
        if policy.CheckNUMAAffinity {
            if !verifyNUMAOptimal(selected, vm.TargetNode) {
                // Try best-effort if affinity required but not strict
                if policy.StrictNUMAAfinity {
                    return nil, ErrNUMAAffinityViolation
                }
            }
        }
        
        // Check inter-GPU communication
        if policy.OptimizeInterGPUComm {
            if !hasOptimalP2PTopology(selected) {
                log.Warn("Selected GPUs have suboptimal P2P topology")
            }
        }
        
        // Reserve allocations
        for _, gpu := range selected {
            allocation := &GPUAllocation{
                GPUID: gpu.ID,
                VMName: vm.Name,
                DeviceIndex: len(allocations),
                PCIAddress: gpu.PCIAddress,
            }
            allocations = append(allocations, *allocation)
        }
    }
    
    // Persist allocations atomically
    if err := persistAllocations(allocations); err != nil {
        rollbackAllocations(allocations)
        return nil, err
    }
    
    return allocations, nil
}

func findAvailableGPUs(request *GPURequest, policy *GPUAllocationPolicy) []*GPU {
    candidates := []*GPU{}
    
    // Get all GPUs matching type
    gpus := getGPUsByType(request.Type)
    
    for _, gpu := range gpus {
        // Check status
        if gpu.Status != StatusAvailable {
            continue
        }
        
        // Check VRAM requirement
        if gpu.VRAM < request.MinMemoryGB {
            continue
        }
        
        // Check capabilities
        if !hasRequiredCapabilities(gpu, request.Features) {
            continue
        }
        
        // Check health
        if gpu.LastHealthCheck < (now - 5*minute) {
            continue // Stale health data
        }
        
        candidates = append(candidates, gpu)
    }
    
    return candidates
}

func applyAllocationPolicy(candidates []*GPU, request *GPURequest, 
    policy *GPUAllocationPolicy) []*GPU {
    
    switch policy.Algorithm {
    case "bin-packing":
        return binPackGPUs(candidates, request, policy)
    case "spread":
        return spreadGPUs(candidates, request, policy)
    case "numa-aware":
        return numaAwareAllocation(candidates, request, policy)
    default:
        return candidates[:request.Count]
    }
}
```

### GPU Health Monitoring

```
Health Check Interval: Every 30 seconds
Telemetry Collected:
- Temperature (threshold: > 80°C WARNING, > 85°C CRITICAL)
- Power draw (threshold: > max_power * 1.1 WARNING)
- Clock speed (throttling detection)
- ECC errors (single-bit errors → log, multi-bit errors → faulty)
- Memory bandwidth utilization
- Thermal throttling events

Actions on Anomalies:
- Temperature too high: Slow down, escalate to faulty if persists
- ECC errors: Log event, alert operator
- Clock throttling: Investigate node cooling
- Power draw anomalies: Investigate power supply
- No response for 2 checks: Mark as faulty, drain allocations
```

### Allocation Failure Handling

```
Scenario 1: Insufficient GPU Capacity
- Response: Return error to user
- Recommendation: Suggest alternative configurations (less GPUs, different type)
- Backoff: Requeue request after 5 minutes
- Alert: If consistent over 1 hour, notify operator

Scenario 2: GPU Faulty During Allocation
- Response: Skip faulty GPU, try next candidate
- If all exhausted: Return error
- Recovery: Mark GPU as faulty, trigger replacement workflow

Scenario 3: Allocation Timeout (> 10 seconds)
- Response: Return error
- Log: Full request state for debugging
- Action: Investigate scheduler performance
```

---

This strategy ensures:
1. **Predictable Performance**: GPU isolation and affinity configuration
2. **Efficient Utilization**: Bin-packing and policy-driven allocation
3. **High Availability**: Health monitoring and failover
4. **Operational Visibility**: Comprehensive metrics and alerting
