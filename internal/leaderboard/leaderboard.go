package leaderboard

import (
	"context"
	"fmt"
	"leaderboard/internal/models"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

// leaderboardKeyForGame generates a Redis key for a specific game's leaderboard.
func leaderboardKeyForGame(game string) string {
	return fmt.Sprintf("leaderboard:%s", game)
}

type LeaderboardService struct {
	redisClient *redis.Client
}

func NewLeaderboardService(redisClient *redis.Client) *LeaderboardService {
	return &LeaderboardService{redisClient: redisClient}
}

func (s *LeaderboardService) SubmitScore(userID string, game string, score float64) error {
	leaderboardKey := leaderboardKeyForGame(game)
	return s.redisClient.ZAdd(ctx, leaderboardKey, &redis.Z{
		Score:  score,
		Member: userID,
	}).Err()
}

func (s *LeaderboardService) GetLeaderboard(game string, limit int64) ([]models.LeaderboardEntry, error) {
	leaderboardKey := leaderboardKeyForGame(game)
	// ZRevRange returns members from highest to lowest score
	results, err := s.redisClient.ZRevRangeWithScores(ctx, leaderboardKey, 0, limit-1).Result()
	if err != nil {
		return nil, err
	}

	var leaderboard []models.LeaderboardEntry
	for _, r := range results {
		leaderboard = append(leaderboard, models.LeaderboardEntry{
			Username: r.Member.(string),
			Score:    r.Score,
		})
	}
	return leaderboard, nil
}

func (s *LeaderboardService) GetUserRank(userID string, game string) (int64, float64, error) {
	leaderboardKey := leaderboardKeyForGame(game)
	// ZRevRank returns the rank of a member (0-based) from highest to lowest score
	rank, err := s.redisClient.ZRevRank(ctx, leaderboardKey, userID).Result()
	if err == redis.Nil {
		// If user is not in the leaderboard, return rank -1 and score 0
		return -1, 0, nil
	}
	if err != nil {
		return 0, 0, err
	}
	score, err := s.redisClient.ZScore(ctx, leaderboardKey, userID).Result()
	return rank, score, err
}