# Dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install git (required by go get), and tzdata (for timezone support)
RUN apk add --no-cache git tzdata

# Copy go mod and sum first, then download deps
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire source
COPY . .

# Build the binary
RUN go build -o timezone-utils ./cmd/server

# Final stage: minimal image
FROM alpine:latest

# Copy tzdata for timezone conversion
RUN apk add --no-cache tzdata ca-certificates

WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/timezone-utils .

# Copy .env if needed (only for dev, optional)
COPY .env .

# Expose port (default from .env or fallback)
EXPOSE 8080

CMD ["./timezone-api"]
