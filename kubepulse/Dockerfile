# Build stage
FROM golang:1.25-bookworm AS builder

WORKDIR /build

# Cache dependencies
# COPY go.mod go.sum ./
COPY . .
RUN go mod download

# Copy source

RUN go build -o /build/kubepulse ./cmd

# Final stage
FROM debian:bookworm-slim

# Install CA certs
RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates && rm -rf /var/lib/apt/lists/*

# Copy binary
COPY --from=builder /build/kubepulse /app

EXPOSE 8081

ENTRYPOINT [ "/app"]
