package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/NostalgiaIm/Quokka-Lab/backend/go/internal/audio"
	"github.com/NostalgiaIm/Quokka-Lab/backend/go/internal/client/rails"
	"github.com/NostalgiaIm/Quokka-Lab/backend/go/internal/config"
	"github.com/NostalgiaIm/Quokka-Lab/backend/go/internal/handler"
	"github.com/NostalgiaIm/Quokka-Lab/backend/go/internal/service"
	"github.com/NostalgiaIm/Quokka-Lab/backend/go/internal/ws"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	cfg := config.Load()

	railsClient := rails.NewClient(cfg.RailsURL, logger)
	audioClient := audio.NewClient(cfg.AudioSocketPath, cfg.AudioBinaryPath, logger)
	quokkaService := service.NewQuokkaService(railsClient, audioClient, logger)
	hub := ws.NewHub(logger)

	go hub.Run()

	server := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           handler.NewRouter(handler.RouterConfig{Service: quokkaService, Hub: hub, Logger: logger}),
		ReadHeaderTimeout: 10 * time.Second,
	}

	go func() {
		logger.Info("go api gateway listening", "addr", cfg.HTTPAddr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("go api gateway failed", "error", err)
			os.Exit(1)
		}
	}()

	shutdownCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	<-shutdownCtx.Done()

	graceCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(graceCtx); err != nil {
		logger.Error("go api gateway shutdown failed", "error", err)
		os.Exit(1)
	}

	logger.Info("go api gateway stopped")
}
