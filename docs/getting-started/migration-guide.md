# Migration Guide

This guide will help you migrate your workloads from existing platforms to the AI Hypervisor Platform.

## Migrating from Kubernetes Device Plugins

If you are currently running AI workloads in Kubernetes using native GPU device plugins, moving to the AI Hypervisor Platform provides stronger isolation and better control over GPU features like MIG.

### Key Differences

1. **Isolation Level:** You will be transitioning from container-level isolation to full virtual machine isolation.
2. **Resource Requests:** Instead of requesting `nvidia.com/gpu: 1` in a Pod specification, you will define `gpu_count: 1` in your VM request payload.

### Migration Steps

1. Export your existing workload configurations from Kubernetes.
2. Create base VM images that mirror the environment of your containers (e.g., Ubuntu + CUDA 12).
3. Update your deployment pipelines to use the AI Hypervisor Platform API or CLI instead of `kubectl`.

## Migrating from Proxmox VE

If you are using Proxmox VE, you are already familiar with VM-based isolation, but you will gain advanced, automated GPU scheduling capabilities.

### Key Differences

1. **Automation:** Proxmox often requires manual configuration for GPU passthrough. The AI Hypervisor Platform automates this via the `gpu-orchestrator`.
2. **API-First:** Our platform is designed around a modern RESTful API, making it easier to integrate with CI/CD systems.

### Migration Steps

1. Review your Proxmox VM configurations.
2. Ensure your VM images are compatible with the AI Hypervisor Platform (KVM/QEMU).
3. Use our API to script the provisioning of your VMs on the new platform.
