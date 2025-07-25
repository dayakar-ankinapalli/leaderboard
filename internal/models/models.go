package models

// User represents a user in the system.
// The password hash is stored, not the plaintext password.
type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"-"` // The password hash should not be exposed
}

// Credentials is used for login and registration requests.
type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// ScoreSubmission represents a score submitted by a user.
type ScoreSubmission struct {
	Game  string  `json:"game"`
	Score float64 `json:"score"`
}

type LeaderboardEntry struct {
	Username string  `json:"username"`
	Score    float64 `json:"score"`
}