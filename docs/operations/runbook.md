# Operations Guide

This guide covers the day-to-day operations, monitoring, and scaling of the AI Hypervisor Platform.

## Observability

The platform exports rich telemetry using open standards.

### Monitoring (Prometheus)

All services expose a `/metrics` endpoint. The `resource-monitor` service aggregates cluster-wide metrics (VM counts, host capacity, GPU utilization).

**Key Metrics:**
- `aihypervisor_vms_total`: Total number of VMs by status.
- `aihypervisor_gpu_utilization`: Percentage utilization of GPUs.
- `aihypervisor_host_memory_bytes`: Available vs used memory per node.

### Tracing (OpenTelemetry)

Distributed tracing is implemented via OTLP. Traces track the lifecycle of a request from the API Server, through NATS, to the Host Agent.

- Set `OtelEndpoint` in the configuration to forward traces to Jaeger or an OpenTelemetry Collector.

### Logging

Services use structured JSON logging via Logrus.
- Log levels can be configured per service.
- Correlation IDs (`trace_id`, `request_id`) are injected into log entries automatically.

## Scaling

### Control Plane Scaling

- The `api-server` and `scheduler` are stateless and can be scaled horizontally behind a load balancer.
- The `task-executor` uses NATS JetStream consumer groups. Scaling it increases the concurrency of background operations.
- `gpu-orchestrator` should typically be run as a singleton or with leader election (handled via Redis) to prevent race conditions during allocation.

### Infrastructure Scaling

- **PostgreSQL**: Must be scaled using standard HA techniques (e.g., streaming replication, Patroni).
- **Redis**: Can be clustered for high availability.
- **NATS**: Deploy as a highly available cluster.

## Maintenance and Backups

- Regularly backup the PostgreSQL database as it contains the authoritative state of all VMs and cluster topologies.
- Redis and NATS states are ephemeral and can be lost without data corruption, though active tasks might need to be retried.