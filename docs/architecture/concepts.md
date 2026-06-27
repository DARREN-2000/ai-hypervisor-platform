# Core Concepts

Understanding the core concepts of the AI Hypervisor Platform is crucial for effective operation and development.

## Virtual Machine (VM)

A Virtual Machine is the fundamental unit of compute. It is an isolated environment running its own operating system, provisioned with specific CPU, memory, and GPU resources.

## Host Node

A Host Node is a physical, bare-metal server in your cluster. It runs the `host-agent` and provides the underlying hardware resources (CPU, RAM, GPUs) for your VMs.

## GPU Orchestrator

The GPU Orchestrator is the component responsible for managing the lifecycle and allocation of GPU resources across the cluster. It ensures that VMs receive the GPUs they requested and manages advanced features like Multi-Instance GPU (MIG).

## Scheduler

The Scheduler decides *where* to place a VM. It evaluates all available Host Nodes and selects the best fit based on configured policies (e.g., bin-packing or spreading).

## Task

A Task represents an asynchronous background operation, such as provisioning a VM or creating a volume. Tasks are tracked and managed by the `task-executor`.
