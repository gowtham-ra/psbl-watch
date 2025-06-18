FROM golang:1.22-bookworm AS builder
WORKDIR /src

# Cache Go modules first (leverages Docker layer caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build a statically-linked binary for linux/amd64
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build -ldflags="-s -w" -o /out/psbl-watch ./cmd/psbl-watch

# --- Runtime image ---
# Use a minimal distroless image that contains TLS certificates
FROM gcr.io/distroless/base-debian12

# Optional: set the timezone for cron scheduling
ENV TZ=America/Los_Angeles

# Copy the binary from the builder stage
COPY --from=builder /out/psbl-watch /psbl-watch

# Run as non-root user provided by distroless (UID 65532)
USER nonroot

ENTRYPOINT ["/psbl-watch"]