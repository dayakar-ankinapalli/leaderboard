# Stage 1: Build the Go application
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum files to download dependencies
COPY go.mod ./
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o leaderboard-app ./cmd/api

# Stage 2: Create a minimal final image
FROM alpine:latest

WORKDIR /root/
COPY --from=builder /app/leaderboard-app .

CMD ["./leaderboard-app"]