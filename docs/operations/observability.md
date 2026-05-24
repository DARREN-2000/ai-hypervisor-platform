# Observability Architecture

## Overview

The AI Hypervisor Platform exposes observability through a shared pipeline:

1. API and control-plane handlers emit request, resource, and health events.
2. A Prometheus-backed recorder stores metrics in a dedicated registry.
3. The API server exposes `/metrics` on the metrics port for Prometheus scraping.
4. OpenTelemetry spans are created per request and exported through OTLP when configured.
5. Structured JSON logs include request IDs and trace IDs for correlation.
6. Health endpoints return liveness and readiness state for Kubernetes and operators.

## Metrics Flow

```mermaid
flowchart LR
    A[API / Control Plane Requests] --> B[Observability Middleware]
    B --> C[Prometheus Registry]
    B --> D[OpenTelemetry Tracer]
    C --> E[/metrics]
    D --> F[OTLP Collector]
    E --> G[Prometheus]
    G --> H[Grafana Dashboards]
    F --> I[Distributed Tracing Backend]
```

## Exported Metrics

The platform records the following metric families:

- HTTP request duration and request totals by method, route, and status class.
- VM resource metrics for CPU, memory, disk I/O, and network I/O.
- GPU usage metrics for utilization, memory usage, temperature, and power.
- Health status indicators for core components.

## Endpoint Map

- `GET /metrics` on the metrics port: Prometheus scrape endpoint.
- `GET /health`: aggregated health status.
- `GET /ready`: readiness check for dependencies.
- `GET /live`: liveness check for the API process.
- `GET /api/v1/metrics`: JSON cluster metrics for API clients.

## Tracing

Tracing uses OpenTelemetry spans around each HTTP request. If `OTEL_EXPORTER_OTLP_ENDPOINT` is set, spans are exported over OTLP. The middleware also forwards or generates `X-Request-ID` so logs and traces can be correlated.

## Logging

Logs are emitted in structured JSON through Logrus. Each request log includes:

- `request_id`
- `trace_id`
- `method`
- `path`
- `status`
- `duration_ms`
- `remote_addr`

## Grafana Dashboard Structure

The dashboard layout is intentionally split into three views:

### Cluster Overview

- API request rate and latency
- Readiness and health state
- Cluster-level CPU and memory utilization
- VM count and GPU count

### GPU Utilization

- GPU utilization by device
- GPU memory used/free
- GPU temperature and power draw
- Faulty or degraded GPU counts

### VM Lifecycle

- VM state transitions over time
- Provisioning and scheduling activity
- VM resource consumption
- Failure and retry rates

## Kubernetes Notes

- The API server service already exposes the metrics port on `8081`.
- Prometheus can scrape the pod using the existing annotations in `deploy/kubernetes/manifests.yaml`.
- The readiness probe should target `/ready` and the liveness probe should target `/live` or `/health` depending on deployment preference.

## Environment Variables

- `OTEL_EXPORTER_OTLP_ENDPOINT`: OTLP collector endpoint, for example `otel-collector.monitoring.svc.cluster.local:4317`.
- `OTEL_EXPORTER_OTLP_INSECURE`: set to `true` for local insecure collector connections.
