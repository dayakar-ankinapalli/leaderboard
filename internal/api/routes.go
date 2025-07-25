package api

import (
	"leaderboard/internal/auth"
	"leaderboard/internal/leaderboard"
	"net/http"

	"github.com/gorilla/mux"
)

func NewRouter(authService *auth.AuthService, leaderboardService *leaderboard.LeaderboardService) *mux.Router {
	router := mux.NewRouter()
	apiHandler := NewAPIHandler(authService, leaderboardService)

	apiRouter := router.PathPrefix("/api").Subrouter()

	// Auth routes
	apiRouter.HandleFunc("/register", apiHandler.RegisterHandler).Methods(http.MethodPost)
	apiRouter.HandleFunc("/login", apiHandler.LoginHandler).Methods(http.MethodPost)

	// Leaderboard routes (protected)
	protected := apiRouter.PathPrefix("").Subrouter()
	protected.Use(authService.AuthMiddleware)
	protected.HandleFunc("/scores", apiHandler.SubmitScoreHandler).Methods(http.MethodPost)
	protected.HandleFunc("/leaderboard/{game}", apiHandler.GetLeaderboardHandler).Methods(http.MethodGet)
	protected.HandleFunc("/rank/{game}", apiHandler.GetUserRankHandler).Methods(http.MethodGet)
	protected.HandleFunc("/report/top-players/{game}", apiHandler.GetLeaderboardHandler).Methods(http.MethodGet) // Re-using handler for simplicity

	return router
}