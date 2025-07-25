package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"leaderboard/internal/models"
	"net/http"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

// A private key type is used to prevent collisions with context keys from other packages.
type contextKey string

// Context keys for user information.
const (
	UserIDKey   = contextKey("user_id")
	UsernameKey = contextKey("username")
)

type AuthService struct {
	redisClient *redis.Client
	jwtSecret   []byte
}

func NewAuthService(redisClient *redis.Client, jwtSecret string) *AuthService {
	return &AuthService{
		redisClient: redisClient,
		jwtSecret:   []byte(jwtSecret),
	}
}

func (s *AuthService) RegisterUser(ctx context.Context, user *models.User) error {
	// Check if user already exists
	exists, err := s.redisClient.HExists(ctx, "users", user.Username).Result()
	if err != nil {
		return fmt.Errorf("failed to check user existence: %w", err)
	}
	if exists {
		return fmt.Errorf("user '%s' already exists", user.Username)
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	user.ID = user.Username // Using username as ID for simplicity

	userJSON, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("failed to marshal user data: %w", err)
	}

	// Store user data in a Redis hash
	return s.redisClient.HSet(ctx, "users", user.ID, userJSON).Err()
}

func (s *AuthService) LoginUser(ctx context.Context, creds *models.Credentials) (string, error) {
	val, err := s.redisClient.HGet(ctx, "users", creds.Username).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("user not found")
	} else if err != nil {
		return "", err
	}

	var user models.User
	if err := json.Unmarshal([]byte(val), &user); err != nil {
		return "", fmt.Errorf("failed to unmarshal user data: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password)); err != nil {
		return "", fmt.Errorf("invalid credentials")
	}

	// Create JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	})

	tokenString, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *AuthService) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return s.jwtSecret, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			// Add user info to the request context
			ctx := context.WithValue(r.Context(), UserIDKey, claims["user_id"])
			ctx = context.WithValue(ctx, UsernameKey, claims["username"])
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
		}
	})
}