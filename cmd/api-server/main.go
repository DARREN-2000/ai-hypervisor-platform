package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"net/http"
	"syscall"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/DARREN-2000/ai-hypervisor-platform/internal/api"
	"github.com/DARREN-2000/ai-hypervisor-platform/internal/app"
	"github.com/DARREN-2000/ai-hypervisor-platform/internal/config"
	"github.com/DARREN-2000/ai-hypervisor-platform/internal/logging"
	"github.com/DARREN-2000/ai-hypervisor-platform/internal/observability"
)

var rootCmd = &cobra.Command{
	Use:   "api-server",
	Short: "AI Hypervisor Platform API Server",
	Long:  "REST API server for VM and GPU orchestration",
	Run:   runAPIServer,
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		logrus.Fatal(err)
	}
}

func runAPIServer(cmd *cobra.Command, args []string) {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		logrus.WithError(err).Fatal("Failed to load configuration")
	}

	// Initialize logger
	log := logging.NewLogger("api-server", cfg.Logging)
	log.WithField("environment", cfg.Environment).Info("Starting API Server")

	ctx := context.Background()
	deps, err := app.Build(ctx, cfg, log)
	if err != nil {
		log.WithError(err).Fatal("Failed to initialize dependencies")
	}

	// Create server
	serverCfg := &api.Config{
		Address:        cfg.APIServer.Address,
		Port:           cfg.APIServer.Port,
		ReadTimeout:    cfg.APIServer.ReadTimeout,
		WriteTimeout:   cfg.APIServer.WriteTimeout,
		IdleTimeout:    cfg.APIServer.IdleTimeout,
		MaxHeaderBytes: cfg.APIServer.MaxHeaderBytes,
		TLSCert:        cfg.APIServer.TLSCert,
		TLSKey:         cfg.APIServer.TLSKey,
	}

	apiServer := api.NewAPIServer(serverCfg, log)

	obs, err := observability.NewManager(ctx, observability.Config{
		ServiceName:    "api-server",
		ServiceVersion: "1.0.0",
		Environment:    cfg.Environment,
		MetricsEnabled: cfg.Metrics.Enabled,
		MetricsAddr:    cfg.Metrics.PrometheusAddr,
		MetricsPort:    cfg.Metrics.PrometheusPort,
		TracingEnabled: true,
		OTLPEndpoint:   os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"),
		OTLPInsecure:   os.Getenv("OTEL_EXPORTER_OTLP_INSECURE") == "true",
	}, log)
	if err != nil {
		log.WithError(err).Fatal("Failed to initialize observability")
	}

	apiServer.SetMetrics(obs.Metrics())
	apiServer.SetHealthChecks(map[string]api.HealthCheck{
		"api": func(context.Context) error { return nil },
		"database": func(context.Context) error {
			if deps.DB == nil {
				return fmt.Errorf("database client not configured")
			}
			return deps.DB.PingContext(context.Background())
		},
		"redis": func(context.Context) error {
			if deps.Redis == nil {
				return fmt.Errorf("redis client not configured")
			}
			return deps.Redis.Ping(context.Background()).Err()
		},
		"nats": func(context.Context) error {
			if deps.NATS == nil {
				return fmt.Errorf("nats client not configured")
			}
			if deps.NATS.Status() != nats.CONNECTED {
				return fmt.Errorf("nats connection is %s", deps.NATS.Status())
			}
			return nil
		},
		"libvirt": func(context.Context) error {
			return nil
		},
	})

	// TODO: Initialize service implementations and inject them here.
	// apiServer.SetDependencies(vmMgr, scheduler, gpuOrch, taskExec, resMon, eventBus, auditLogger, stateStore)

	// Register routes
	apiServer.RegisterRoutes()

	if cfg.Metrics.Enabled {
		go func() {
			if err := obs.ServeMetrics(); err != nil && err != http.ErrServerClosed {
				log.WithError(err).Error("Metrics server error")
			}
		}()
	}

	// Start server in goroutine
	go func() {
		if err := apiServer.Start(); err != nil {
			log.WithError(err).Fatal("API server error")
		}
	}()

	log.Infof("API Server started on %s:%d", cfg.APIServer.Address, cfg.APIServer.Port)

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)

	<-sigChan
	log.Info("Shutdown signal received, gracefully stopping...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := apiServer.Stop(ctx); err != nil {
		log.WithError(err).Error("Error during shutdown")
	}
	if err := obs.Shutdown(ctx); err != nil {
		log.WithError(err).Error("Error stopping observability stack")
	}
	if err := deps.Close(ctx); err != nil {
		log.WithError(err).Error("Error closing dependencies")
	}

	log.Info("API Server stopped")
}
