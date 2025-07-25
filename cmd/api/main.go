package main

import (
	"context"
	"fmt"
	"leaderboard/internal/api"
	"leaderboard/internal/auth"
	"leaderboard/internal/leaderboard"
	"leaderboard/internal/config"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-redis/redis/v8"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}

	// Set up structured logger
	logger := setupLogger(cfg.LogLevel, cfg.LogFormat)
	slog.SetDefault(logger)

	slog.Info("Starting Leaderboard Service...")

	// Set up Redis connection
	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
	})

	if _, err := redisClient.Ping(context.Background()).Result(); err != nil {
		slog.Error("could not connect to Redis", "error", err)
		os.Exit(1)
	}
	slog.Info("Successfully connected to Redis")

	// Initialize services
	authService := auth.NewAuthService(redisClient, cfg.JWTSecret)
	leaderboardService := leaderboard.NewLeaderboardService(redisClient)

	// Set up router and server
	router := api.NewRouter(authService, leaderboardService)
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Port),
		Handler: router,
	}

	// Graceful shutdown
	go func() {
		slog.Info("Server is running", "port", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server failed to start", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("server shutdown failed", "error", err)
		os.Exit(1)
	}

	slog.Info("Server gracefully stopped")
}

func setupLogger(level, format string) *slog.Logger {
	var logLevel slog.Level
	switch level {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{Level: logLevel}
	var handler slog.Handler = slog.NewTextHandler(os.Stdout, opts)
	if format == "json" {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	}

	return slog.New(handler)
}