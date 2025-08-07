# file: Dockerfile
# version: 1.0.0
# guid: a1b2c3d4-e5f6-7890-abcd-ef1234567890

# Build stage
FROM golang:1.21-alpine AS builder

# Install git (needed for go modules)
RUN apk add --no-cache git

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o copilot-agent-util ./cmd/copilot-agent-util

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS calls
RUN apk --no-cache add ca-certificates git

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/copilot-agent-util .

# Create logs directory
RUN mkdir -p logs

# Set the entrypoint
ENTRYPOINT ["./copilot-agent-util"]
