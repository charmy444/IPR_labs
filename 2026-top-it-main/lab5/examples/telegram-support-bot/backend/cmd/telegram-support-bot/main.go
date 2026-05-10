package main

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"telegram-support-bot/internal/bot"
	"telegram-support-bot/internal/config"
	"telegram-support-bot/internal/repository"
	"telegram-support-bot/internal/server"
	"telegram-support-bot/internal/telemetry"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	_ "github.com/lib/pq"
)

func main() {
	cfg := config.Load()

	// Initialize logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	tracerShutdown, err := telemetry.InitTracer(context.Background())
	if err != nil {
		slog.Error("OpenTelemetry init failed", "error", err)
		os.Exit(1)
	}
	defer func() {
		_ = tracerShutdown(context.Background())
	}()

	// Connect to database
	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	// Test database connection
	if err := db.Ping(); err != nil {
		slog.Error("Failed to ping database", "error", err)
		os.Exit(1)
	}
	slog.Info("Connected to database")

	// Run migrations
	if err := runMigrations(db); err != nil {
		slog.Error("Failed to run migrations", "error", err)
		os.Exit(1)
	}
	slog.Info("Migrations completed successfully")

	// Initialize repository
	repo := repository.New(db)

	// Initialize bot
	botAPI, err := tgbotapi.NewBotAPI(cfg.BotToken)
	if err != nil {
		slog.Error("Failed to create bot", "error", err)
		slog.Error("Please check your BOT_TOKEN in .env file. Get a valid token from @BotFather in Telegram.")
		os.Exit(1)
	}

	botAPI.Debug = false

	slog.Info("Bot authorized", "username", botAPI.Self.UserName)

	// Create bot instance
	botInstance := bot.New(botAPI, repo)

	// Create server instance
	srv := server.NewServer(cfg, repo, botInstance)

	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start bot in goroutine
	go func() {
		slog.Info("Starting bot")
		if err := botInstance.Start(); err != nil {
			slog.Error("Bot stopped with error", "error", err)
			cancel()
		}
	}()

	// Start HTTP server in goroutine
	go func() {
		slog.Info("Starting HTTP server", "port", cfg.ServerPort)
		if err := srv.Start(); err != nil {
			slog.Error("HTTP server stopped with error", "error", err)
			cancel()
		}
	}()

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	slog.Info("Shutting down...")

	// Graceful shutdown
	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("Server shutdown error", "error", err)
	}

	slog.Info("Application stopped")
}

func runMigrations(db *sql.DB) error {
	migrationsDir := "./migrations"

	entries, err := os.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		filename := entry.Name()
		if filepath.Ext(filename) != ".sql" {
			continue
		}

		content, err := os.ReadFile(filepath.Join(migrationsDir, filename))
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", filename, err)
		}

		if _, err := db.Exec(string(content)); err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", filename, err)
		}

		slog.Info("Migration applied", "file", filename)
	}

	return nil
}
