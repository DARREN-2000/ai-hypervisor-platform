package main

import (
    "context"
    "os"
    "os/signal"
    "syscall"

    "github.com/DARREN-2000/ai-hypervisor-platform/internal/collectors"
    "github.com/DARREN-2000/ai-hypervisor-platform/internal/config"
    "github.com/DARREN-2000/ai-hypervisor-platform/internal/logging"
    "github.com/DARREN-2000/ai-hypervisor-platform/internal/observability"
)

func main() {
    ctx := context.Background()

    cfg := config.DefaultConfig()

    log := logging.NewLogger("resource-monitor", cfg.Logging)
    log.WithField("environment", cfg.Environment).Info("starting resource-monitor (demo collectors)")

    obs, err := observability.NewManager(ctx, observability.Config{
        ServiceName:    "resource-monitor",
        ServiceVersion: "0.1.0",
        Environment:    cfg.Environment,
        MetricsEnabled: cfg.Metrics.Enabled,
        MetricsAddr:    cfg.Metrics.PrometheusAddr,
        MetricsPort:    cfg.Metrics.PrometheusPort,
        TracingEnabled: true,
        OTLPEndpoint:   os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"),
        OTLPInsecure:   os.Getenv("OTEL_EXPORTER_OTLP_INSECURE") == "true",
    }, log)
    if err != nil {
        log.WithError(err).Fatal("failed to init observability")
    }

    if cfg.Metrics.Enabled {
        go func() {
            if err := obs.ServeMetrics(); err != nil {
                log.WithError(err).Error("metrics server failed")
            }
        }()
    }

    // Start synthetic collectors for demo/testing. Replace with real fetchers.
    cancelVM := collectors.StartVMCollector(ctx, cfg.ResourceMonitor.MetricsInterval, obs.Metrics(), collectors.SyntheticVMFetcher(5))
    cancelGPU := collectors.StartGPUCollector(ctx, cfg.ResourceMonitor.MetricsInterval, obs.Metrics(), collectors.SyntheticGPUFetcher(2))

    // Wait for termination
    sig := make(chan os.Signal, 1)
    signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
    <-sig

    cancelVM()
    cancelGPU()

    _ = obs.Shutdown(context.Background())
    log.Info("resource-monitor stopped")
}
