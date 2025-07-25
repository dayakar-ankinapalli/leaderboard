# Real-Time Leaderboard Service

This project is a backend system for a real-time leaderboard service built with Go and Redis. It allows users to register, log in, submit scores for games, and view a global leaderboard in real-time.

## Features

-   User Registration and Login (JWT-based authentication)
-   Score Submission for different games
-   Real-time Global Leaderboard
-   Get User's Rank
-   Top N Players Report
-   Dockerized setup for easy deployment

## Tech Stack

-   **Backend:** Go (Golang)
-   **Database/Cache:** Redis (for leaderboards and user storage)
-   **Router:** `gorilla/mux`
-   **Containerization:** Docker & Docker Compose

## Prerequisites

-   Docker
-   Docker Compose

## Getting Started

1.  **Clone the repository:**

    ```sh
    git clone <repository-url>
    cd leaderboard
    ```

2.  **Run the application:**

    Use Docker Compose to build and run the Go application and the Redis container.

    ```sh
    docker-compose up --build
    ```

    The API will be available at `http://localhost:8080`.

## API Endpoints

All endpoints are prefixed with `/api`.

### Authentication

#### `POST /api/register`

Register a new user.

**Request Body:**

```json
{
    "username": "newuser",
    "password": "password123"
}
```

**Response (Success 201):**

```json
{
    "message": "User registered successfully"
}
```

#### `POST /api/login`

Log in an existing user and get a JWT token.

**Request Body:**

```json
{
    "username": "newuser",
    "password": "password123"
}
```

**Response (Success 200):**

```json
{
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### Leaderboard

*Note: The following endpoints require authentication. Include the JWT token in the `Authorization` header as a Bearer token.*

`Authorization: Bearer <your-jwt-token>`

#### `POST /api/scores`

Submit a score for the logged-in user. The system uses the user's ID from the JWT token.

**Request Body:**

```json
{
    "game": "space-invaders",
    "score": 150
}
```

**Response (Success 200):**

```json
{
    "message": "Score submitted successfully"
}
```

#### `GET /api/leaderboard`

Get the global leaderboard. By default, it returns the top 10 players.

**Query Parameters:**

-   `limit` (optional): Number of top players to return. Default is 10.

**Example:** `GET /api/leaderboard?limit=5`

**Response (Success 200):**

```json
[
    {
        "username": "user1",
        "score": 500
    },
    {
        "username": "user2",
        "score": 450
    }
]
```

#### `GET /api/rank`

Get the rank of the currently logged-in user.

**Response (Success 200):**

```json
{
    "username": "newuser",
    "rank": 5,
    "score": 150
}
```

#### `GET /api/report/top-players`

This is an example of a report endpoint. It behaves similarly to the main leaderboard endpoint but could be extended for different time periods or games.

**Query Parameters:**

-   `limit` (optional): Number of top players to return. Default is 10.

**Response (Success 200):**

```json
[
    {
        "username": "user1",
        "score": 500
    },
    {
        "username": "user2",
        "score": 450
    }
]
```# leaderboard
