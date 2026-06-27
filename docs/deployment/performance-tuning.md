# Performance Tuning

The AI Hypervisor Platform is built for high-performance workloads, but you can optimize it further based on your specific requirements.

## GPU Scheduling

### NUMA Awareness

For latency-sensitive workloads like training large language models, enabling NUMA-aware scheduling can significantly reduce memory access times. Ensure that your GPU allocation policies in the scheduler are configured to prioritize locality.

## Database Optimization

### PostgreSQL Connections

If your cluster scales to hundreds of hosts or thousands of VMs, you may need to increase the `max_connections` setting in PostgreSQL to handle the load from the `vm-manager` and `gpu-orchestrator`.

## Monitoring and Observability

### Prometheus Scrape Intervals

If you are generating too much telemetry data, consider increasing the Prometheus scrape interval from the default 15 seconds to 30 or 60 seconds.

## Memory Management

### Disabling Swap

To ensure deterministic performance, it is highly recommended to disable swap on your host nodes. This prevents the OS from paging out VM memory, which can lead to unpredictable latency spikes.
