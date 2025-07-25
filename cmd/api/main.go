package main

import (
	"context"
	"leaderboard/internal/api"
	"leaderboard/internal/auth"
	"leaderboard/internal/leaderboard"
	"log"
	"net/http"
	"os"

	"github.com/go-redis/redis/v8"
)

func main() {
	log.Println("Starting Leaderboard Service...")

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "default_secret_key"
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	if _, err := redisClient.Ping(context.Background()).Result(); err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}

	authService := auth.NewAuthService(redisClient, jwtSecret)
	leaderboardService := leaderboard.NewLeaderboardService(redisClient)

	router := api.NewRouter(authService, leaderboardService)

	log.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}