package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/alarm-agent/internal/adapters/http"
	"github.com/alarm-agent/internal/adapters/llm"
	"github.com/alarm-agent/internal/adapters/repo"
	"github.com/alarm-agent/internal/adapters/whatsapp"
	"github.com/alarm-agent/internal/config"
	"github.com/alarm-agent/internal/domain"
	"github.com/alarm-agent/internal/infra"
	"github.com/alarm-agent/internal/usecase"
	"github.com/alarm-agent/internal/workers"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	logger, err := infra.NewLogger(cfg.App.Environment)
	if err != nil {
		return fmt.Errorf("failed to create logger: %w", err)
	}
	defer logger.Sync()

	logger.Info("Starting Alarm Agent",
		zap.String("environment", cfg.App.Environment),
		zap.String("port", cfg.App.Port),
	)

	repos, err := repo.NewPostgresRepositories(cfg.Database.DSN)
	if err != nil {
		return fmt.Errorf("failed to create repositories: %w", err)
	}
	defer repos.Close()

	if err := initializeWhitelist(repos, cfg, logger); err != nil {
		logger.Warn("Failed to initialize whitelist", zap.Error(err))
	}

	// LLM client is now created per-request from database configuration

	whatsappSender := whatsapp.NewInfobipClient(
		cfg.Infobip.BaseURL,
		cfg.Infobip.APIKey,
		cfg.Infobip.WhatsAppSender,
	)

	webhookVerifier := whatsapp.NewInfobipWebhookVerifier(cfg.Infobip.WebhookSecret)
	timeProvider := infra.NewRealTimeProvider()

	eventUseCase := usecase.NewEventUseCase(repos)
	messageUseCase := usecase.NewMessageUseCase(
		repos,
		whatsappSender,
		eventUseCase,
		cfg.App.DefaultTimezone,
		cfg,
	)

	reminderWorker := workers.NewReminderWorker(
		repos,
		whatsappSender,
		timeProvider,
		logger,
		cfg.Worker.ReminderTickInterval,
	)

	server := http.NewServer(
		cfg,
		repos,
		messageUseCase,
		eventUseCase,
		webhookVerifier,
		logger,
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := server.Start(); err != nil {
			logger.Error("HTTP server error", zap.Error(err))
			cancel()
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := reminderWorker.Start(ctx); err != nil && err != context.Canceled {
			logger.Error("Reminder worker error", zap.Error(err))
			cancel()
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-quit:
		logger.Info("Received shutdown signal", zap.String("signal", sig.String()))
	case <-ctx.Done():
		logger.Info("Context cancelled, shutting down")
	}

	logger.Info("Shutting down gracefully...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	reminderWorker.Stop()

	if err := server.Stop(shutdownCtx); err != nil {
		logger.Error("Error shutting down HTTP server", zap.Error(err))
	}

	cancel()
	wg.Wait()

	logger.Info("Shutdown complete")
	return nil
}

// getLLMAPIKey function removed - API keys are now stored in database

func initializeWhitelist(repos *repo.PostgresRepositories, cfg *config.Config, logger *zap.Logger) error {
	ctx := context.Background()

	for _, number := range cfg.Security.WhitelistNumbers {
		exists, err := repos.Whitelist().IsWhitelisted(ctx, number)
		if err != nil {
			return err
		}

		if !exists {
			whitelist := &domain.WhitelistNumber{
				Number: number,
				Note:   nil, // Could add a note like "Added from config"
			}
			if err := repos.Whitelist().Add(ctx, whitelist); err != nil {
				logger.Error("Failed to add number to whitelist", 
					zap.String("number", number), 
					zap.Error(err))
			} else {
				logger.Info("Added number to whitelist", zap.String("number", number))
			}
		}
	}

	return nil
}