# Realtime Leaderboard API

This project is a Go-based backend service for a realtime leaderboard system. It provides a simple RESTful API for user registration, authentication, score submission, and leaderboard retrieval.

This implementation is inspired by the [Realtime Leaderboard System project idea on roadmap.sh](https://roadmap.sh/projects/realtime-leaderboard-system).

## Features

*   **User Authentication**: Secure user registration and login using JWT (JSON Web Tokens).
*   **Score Submission**: Authenticated users can submit scores for various games.
*   **Leaderboard Retrieval**: Fetch game-specific leaderboards, sorted by score.
*   **Structured Logging**: Configurable JSON or text logging for better monitoring.
*   **Graceful Shutdown**: The server shuts down gracefully, ensuring all requests are completed.
*   **Configuration-driven**: All settings are managed via environment variables for easy deployment.

## Tech Stack

*   **Backend**: Go (Golang) with the standard `net/http` library.
*   **Database**: Redis for fast, in-memory data storage.
*   **Containerization**: Docker & Docker Compose.
*   **Configuration**: `godotenv` and `envconfig` for environment variable management.

## Getting Started

Follow these instructions to get the project up and running on your local machine.

### Prerequisites

*   Go (version 1.20+)
*   Docker and Docker Compose

### 1. Clone the Repository

```bash
git clone <your-repository-url>
cd leaderboard
```

### 2. Configure Environment Variables

The application requires a set of environment variables to run. A template is provided in `.env.example`.

Copy the example file to a new `.env` file:

```bash
cp .env.example .env
```

You can modify the values in `.env` as needed. The default values are suitable for local development.

### 3. Run the Application

You can run the application using Docker Compose (recommended) or directly with Go.

**Using Docker Compose:**

This is the easiest way to start the service, as it also manages the Redis container.

```bash
docker-compose up --build
```

The API will be available at `http://localhost:3001`.

## API Documentation

The API is fully documented using the OpenAPI 3.0 standard. The specification can be found in the `openapi.yaml` file.

You can use tools like the Swagger Editor or Postman to import the `openapi.yaml` file and interact with the API.