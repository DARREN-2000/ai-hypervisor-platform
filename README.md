# AI Hypervisor Platform

[![Build Status](https://github.com/ai-hypervisor/platform/workflows/Go/badge.svg)](https://github.com/ai-hypervisor/platform/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/ai-hypervisor/platform)](https://goreportcard.com/report/github.com/ai-hypervisor/platform)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

**An opinionated, production-grade virtualization control plane focused on GPU-accelerated AI workloads.**

AI Hypervisor Platform provides service operators with a unified control plane to provision, schedule, and observe virtual machines tailored for GPU workloads. It combines secure VM lifecycle management, intelligent GPU orchestration, and a robust observability stack to run inference services at scale.

<p align="center">
	<img src="docs/site/animations/ai-hypervisor-hero.svg" alt="AI Hypervisor Platform animated control plane" width="100%" />
</p>

## Product Vision

To become the standard infrastructure layer for heterogeneous GPU clusters, enabling zero-trust, multi-tenant AI inference and training environments with native observability and deterministic performance guarantees.

## Key Features

- **Intelligent GPU Orchestration**: Bin-packing, spread, and NUMA-aware scheduling for optimal GPU utilization.
- **Hardware Virtualization**: Secure VM lifecycle management using KVM/QEMU via Libvirt.
- **Unified Observability**: Prometheus metrics, OpenTelemetry tracing, and structured logging built-in.
- **Kubernetes-Native Operations**: Deployable via Helm, exposing standard Kubernetes APIs and resources.
- **Multi-Tenant Isolation**: Strict resource boundaries, RBAC, and namespace isolation.

## Architecture Overview

The system is composed of modular microservices communicating asynchronously over NATS, with PostgreSQL as the authoritative store:

*   **API Server**: External REST/WebSocket API and UI ingress.
*   **VM Manager**: VM lifecycle orchestration.
*   **GPU Orchestrator**: Allocation engine for GPU resources.
*   **Scheduler**: Pluggable policy engine for node selection.
*   **Task Executor**: Reliable asynchronous job execution.
*   **Resource Monitor**: Telemetry aggregation.
*   **Host Agent**: Node-level daemon interfacing with libvirt and NVML.

*(See [Architecture Guide](docs/architecture/system.md) for detailed diagrams.)*

## Why This Exists

Modern AI workloads require direct hardware access (GPUs, NVLink) while maintaining strict multi-tenant isolation. Traditional virtualization layers add unacceptable overhead, while bare-metal lacks the necessary flexibility. The AI Hypervisor Platform bridges this gap by providing bare-metal performance with cloud-native orchestration primitives.

## Comparison With Alternatives

| Feature | AI Hypervisor | Kubernetes (Native) | Proxmox VE | OpenStack |
| :--- | :--- | :--- | :--- | :--- |
| **Primary Workload** | GPU VMs | Containers | VMs/LXC | VMs |
| **Complexity** | Low | Medium | Low | High |
| **GPU Scheduling** | Advanced (NUMA, MIG) | Basic (Device Plugins) | Manual | Basic |
| **Observability** | Built-in (OTEL) | Add-ons required | Basic | Add-ons required |

## Technology Stack

- Language: Go 1.21
- Messaging: NATS
- Datastore: PostgreSQL (primary), Redis (caching/coordination)
- Container orchestration: Kubernetes (manifests & Helm charts included)
- Virtualization: KVM/QEMU via Libvirt
- GPU telemetry: NVML (NVIDIA), vendor-specific tools for AMD/Intel
- Observability: Prometheus, Grafana, OpenTelemetry (OTLP)

## Project Structure

- `cmd/` – Service entry points (API, VM manager, GPU orchestrator, resource monitor, host agent).
- `internal/` – Core implementation packages and integrations.
- `pkg/` – Reusable public packages (telemetry, errors, common utilities).
- `deploy/` – Kubernetes manifests, Helm charts, and provisioning scripts.
- `docs/` – Architecture and operations documentation.
- `docs/site/` – Modern static frontend published to GitHub Pages.
- `docs/site/animations/` – Animated SVG assets used by the README and Pages demo.

## Quick Start

### Installation

```bash
# Clone the repository
git clone https://github.com/ai-hypervisor/platform.git
cd platform

# Set up development environment
make setup-dev

# Build services
make build
```

### Configuration

Edit `config/sample-config.yaml` to match your environment (networking, metrics ports, GPU policies).

### Running Locally

```bash
# Start infrastructure dependencies (requires Docker Compose or Kubernetes)
# See docs/getting-started/installation.md for full instructions

# Run the API server
./bin/api-server --config config/sample-config.yaml
```

## Performance & Benchmarks

The AI Hypervisor Platform is designed for minimal overhead. Bare-metal performance is retained with zero CPU throttling during VM provisioning and sub-millisecond API response times. Detailed benchmarking data will be available in future releases.

## Limitations

- Currently limited to KVM/QEMU virtualization via libvirt.
- Does not yet support live migration.
- Tested and verified primarily on NVIDIA GPUs.

## Roadmap

- ✅ Core VM orchestration and scheduler
- ✅ GPU allocation primitives and monitoring
- ✅ Observability: Prometheus + OTLP traces
- [ ] Live VM migration
- [ ] Multi-cluster federation and cross-region scheduling
- [ ] ML-driven scheduler recommendations
- [ ] Advanced GPU virtualization features (fine-grained MIG profiles)

## Documentation

Full documentation is available in the `docs/` directory:

- [Getting Started](docs/getting-started/README.md)
- [Architecture Guide](docs/architecture/system.md)
- [API Reference](docs/api/endpoints.md)
- [Developer Experience](docs/developer/onboarding.md)
- [Security](docs/security/threat-model.md)

## Contributing

We welcome contributions! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for details.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
