package api

import (
	"encoding/json"
	"leaderboard/internal/auth"
	"leaderboard/internal/leaderboard"
	"leaderboard/internal/models"

	// "leaderboard/internal/models"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

type APIHandler struct {
	authService        *auth.AuthService
	leaderboardService *leaderboard.LeaderboardService
}

func NewAPIHandler(auth *auth.AuthService, leaderboard *leaderboard.LeaderboardService) *APIHandler {
	return &APIHandler{
		authService:        auth,
		leaderboardService: leaderboard,
	}
}

func (h *APIHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var creds models.Credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	user := &models.User{
		Username: creds.Username,
		Password: creds.Password,
	}

	if err := h.authService.RegisterUser(r.Context(), user); err != nil {
		if strings.Contains(err.Error(), "already exists") {
			respondWithError(w, http.StatusConflict, err.Error())
		} else {
			respondWithError(w, http.StatusInternalServerError, "Could not register user")
		}
		return
	}

	respondWithJSON(w, http.StatusCreated, map[string]string{"message": "User registered successfully"})
}

func (h *APIHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var creds models.Credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	token, err := h.authService.LoginUser(r.Context(), &creds)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"token": token})
}

func (h *APIHandler) SubmitScoreHandler(w http.ResponseWriter, r *http.Request) {
	var submission models.ScoreSubmission
	if err := json.NewDecoder(r.Body).Decode(&submission); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if submission.Game == "" {
		respondWithError(w, http.StatusBadRequest, "Game field is required")
		return
	}

	userID := r.Context().Value(auth.UserIDKey).(string)
	if err := h.leaderboardService.SubmitScore(userID, submission.Game, submission.Score); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not submit score")
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Score submitted successfully"})
}

func (h *APIHandler) GetLeaderboardHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	game, ok := vars["game"]
	if !ok {
		respondWithError(w, http.StatusBadRequest, "Game parameter is missing")
		return
	}

	limitStr := r.URL.Query().Get("limit")
	limit, err := strconv.ParseInt(limitStr, 10, 64)
	if err != nil || limit <= 0 {
		limit = 10 // Default limit
	}

	entries, err := h.leaderboardService.GetLeaderboard(game, limit)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not retrieve leaderboard")
		return
	}

	respondWithJSON(w, http.StatusOK, entries)
}

func (h *APIHandler) GetUserRankHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	game, ok := vars["game"]
	if !ok {
		respondWithError(w, http.StatusBadRequest, "Game parameter is missing")
		return
	}

	userID := r.Context().Value(auth.UserIDKey).(string)
	username := r.Context().Value(auth.UsernameKey).(string)

	rank, score, err := h.leaderboardService.GetUserRank(userID, game)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not retrieve user rank")
		return
	}

	var rankDisplay interface{}
	if rank == -1 {
		rankDisplay = "unranked"
	} else {
		rankDisplay = rank + 1 // Convert 0-based to 1-based for display
	}

	response := map[string]interface{}{
		"username": username,
		"rank":     rankDisplay,
		"score":    score,
	}

	respondWithJSON(w, http.StatusOK, response)
}

// respondWithError is a helper for sending error JSON responses.
func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

// respondWithJSON is a helper for sending JSON responses.
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}