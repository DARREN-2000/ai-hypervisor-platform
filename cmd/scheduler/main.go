package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/DARREN-2000/ai-hypervisor-platform/internal/config"
	"github.com/DARREN-2000/ai-hypervisor-platform/internal/logging"
)

var configPath string

var rootCmd = &cobra.Command{
	Use:   "scheduler",
	Short: "AI Hypervisor Platform scheduler",
	Long:  "Scaffold entrypoint for the placement and scoring service.",
	RunE:  runScheduler,
}

func init() {
	rootCmd.Flags().StringVar(&configPath, "config", "", "optional config file path")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		logrus.Fatal(err)
	}
}

func runScheduler(cmd *cobra.Command, args []string) error {
	if configPath != "" {
		if err := os.Setenv("AIHYPERVISOR_CONFIG_PATH", configPath); err != nil {
			return err
		}
	}

	cfg, err := config.Load()
	if err != nil {
		return err
	}

	log := logging.NewLogger("scheduler", cfg.Logging)
	log.WithField("environment", cfg.Environment).Info("starting scheduler scaffold")

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.WithField("timeout", shutdownCtx.Err() == context.DeadlineExceeded).Info("stopping scheduler scaffold")
	return nil
}
