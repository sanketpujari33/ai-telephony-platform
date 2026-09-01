FROM golang:alpine AS builder

# Set the working directory
WORKDIR /app

# Install necessary build tools
RUN apk add --no-cache git gcc musl-dev

# Copy go.mod and go.sum first to cache dependencies
COPY go.mod go.sum ./
RUN sed -i 's/^go 1.*/go 1.23/' go.mod
RUN go mod download

# Copy the rest of the application code
COPY . .

# Pass the service name as a build argument
ARG SERVICE_NAME

# Build the Go binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/service ./cmd/${SERVICE_NAME}

# Use a minimal alpine image for the final stage
FROM alpine:latest

# Install necessary certificates
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/bin/service /app/service

# Copy the API documentation
COPY --from=builder /app/docs /app/docs

# Add any default configuration directories if needed
RUN mkdir -p /etc/ai-telephony-platform

# Expose common ports (specifics handled by Kubernetes manifests)
EXPOSE 8080 9090-9095 5060-5063

# Command to run the executable
ENTRYPOINT ["/app/service"]
