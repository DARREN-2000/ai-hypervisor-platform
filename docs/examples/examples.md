# Examples

## Beginner: Provisioning a Single GPU VM

This example shows how to provision a basic VM for inference.

```json
{
  "name": "basic-inference",
  "cpu": 8,
  "memory_mb": 32768,
  "gpu_count": 1,
  "image": "ubuntu-22.04-cuda12"
}
```

## Advanced: NUMA-Aware Multi-GPU Node

For training jobs requiring low-latency memory access, request strict NUMA affinity.

```json
{
  "name": "training-node",
  "cpu": 32,
  "memory_mb": 131072,
  "gpu_count": 4,
  "image": "ubuntu-22.04-cuda12",
  "metadata": {
    "numa_policy": "strict"
  }
}
```