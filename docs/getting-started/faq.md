# Frequently Asked Questions

## General

**What is the AI Hypervisor Platform?**
It is an open-source, production-grade virtualization control plane designed to run GPU-accelerated workloads (like AI inference and training) with strict isolation and high performance.

**Why not just use Kubernetes device plugins?**
Kubernetes container-level isolation is often insufficient for secure multi-tenancy. AI Hypervisor provides VM-level isolation via KVM/QEMU while retaining a cloud-native API and orchestration model.

## Architecture

**Which hypervisor does this use?**
The platform manages VMs using KVM/QEMU, controlled via Libvirt by the `host-agent`.

**How does GPU scheduling work?**
The `gpu-orchestrator` handles advanced placement decisions, supporting bin-packing, spreading, and NUMA-aware affinity to ensure optimal latency between CPU and GPU memory.

## Operations

**Is it safe to restart the API Server?**
Yes. The API Server is stateless. Ongoing VM provisioning tasks are handled asynchronously by the `task-executor`.

**How do I view metrics?**
All services expose a `/metrics` endpoint for Prometheus. You can use the provided Grafana dashboards in `deploy/grafana/dashboards` to visualize cluster health and GPU utilization.