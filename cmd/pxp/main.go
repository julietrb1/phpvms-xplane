package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/julietrb1/phpvms-xplane/internal/api"
	"github.com/julietrb1/phpvms-xplane/internal/config"
	"github.com/julietrb1/phpvms-xplane/internal/logging"
	"github.com/julietrb1/phpvms-xplane/internal/service"
	"github.com/julietrb1/phpvms-xplane/internal/tui"
	"github.com/julietrb1/phpvms-xplane/internal/udp"
)

func main() {
	var configFile string
	var enableTUI bool
	flag.StringVar(&configFile, "config", "", "Path to config file")
	flag.BoolVar(&enableTUI, "tui", true, "Enable Terminal User Interface")
	flag.Parse()

	cfg := config.DefaultConfig()
	if err := cfg.LoadFromDotEnv(configFile); err != nil {
		fmt.Fprintf(os.Stderr, "Error loading configuration from .env file: %v\n", err)
		os.Exit(1)
	}

	if err := cfg.LoadPreferences(""); err != nil {
		fmt.Fprintf(os.Stderr, "Error loading user preferences: %v\n", err)
	}

	cfg.TUIEnabled = enableTUI

	if err := cfg.Validate(); err != nil {
		fmt.Fprintf(os.Stderr, "Invalid configuration: %v\n", err)
		os.Exit(1)
	}

	logger := logging.SetupLogger(cfg.LogLevel)
	logger.Info("Starting phpVMS ACARS client",
		"udp_bind", fmt.Sprintf("%s:%d", cfg.UDPBindHost, cfg.UDPBindPort),
	)

	apiClient := api.NewClient(cfg.PhpVMSBaseURL, cfg.PhpVMSAPIKey, logger)
	flightService := service.NewFlightService(apiClient, logger)
	udpListener, err := udp.NewListener(cfg.UDPBindHost, cfg.UDPBindPort, flightService, logger)
	if err != nil {
		logger.Error("Failed to create UDP listener", "error", err)
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigCh
		logger.Info("Received signal, shutting down", "signal", sig)
		cancel()
	}()

	if cfg.TUIEnabled {
		logger.Info("Starting Terminal User Interface")
		go func() {
			if err := tui.Run(ctx, cancel, udpListener.GetMetrics(), flightService, cfg, logger); err != nil {
				logger.Error("TUI error", "error", err)
				os.Exit(1)
			}
		}()
	}

	logger.Info("Starting UDP listener", "addr", fmt.Sprintf("%s:%d", cfg.UDPBindHost, cfg.UDPBindPort))
	if err := udpListener.Start(ctx); err != nil && err != context.Canceled {
		logger.Error("UDP listener error", "error", err)
		os.Exit(1)
	}

	logger.Info("Shutdown complete")
}
