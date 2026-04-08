# Build stage
FROM golang:1.23-alpine AS builder

# Set working directory
WORKDIR /build

# Install git and ca-certificates for private repos and HTTPS
RUN apk add --no-cache git ca-certificates tzdata

# Copy dependency files first for better layer caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download && go mod verify

# Copy source code
COPY . .

# Build the application with optimizations
# -ldflags "-w -s" removes debug info and symbol table for smaller binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s -extldflags '-static'" \
    -a -installsuffix cgo \
    -o server \
    ./cmd/server/main.go

# Final stage - using distroless for minimal attack surface and size
FROM gcr.io/distroless/static:nonroot

# Copy timezone data
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Copy the binary from builder stage
COPY --from=builder /build/server /server

# Use nonroot user (provided by distroless)
USER nonroot:nonroot

# Expose port
EXPOSE 8080

# Health check endpoint is available at /health
# Run the server
ENTRYPOINT ["/server"] 