# Multi-stage Dockerfile for docloom
# Stage 1: Build the Go binary
FROM golang:1.22-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git make

# Set working directory
WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build arguments for version information
ARG VERSION=dev
ARG GIT_COMMIT=unknown
ARG BUILD_DATE=unknown

# Build the binary with version information
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo \
    -ldflags "-X github.com/karolswdev/docloom/internal/version.Version=${VERSION} \
    -X github.com/karolswdev/docloom/internal/version.GitCommit=${GIT_COMMIT} \
    -X github.com/karolswdev/docloom/internal/version.BuildDate=${BUILD_DATE}" \
    -o docloom ./cmd/docloom

# Run tests (with CI flag to skip tests that need bash)
ENV CI=true
RUN go test ./...

# Stage 2: Create minimal final image
FROM alpine:3.19

# Install runtime dependencies
RUN apk add --no-cache ca-certificates

# Create non-root user
RUN adduser -D -u 1000 docloom

# Copy binary from builder
COPY --from=builder /build/docloom /usr/local/bin/docloom

# Copy default templates (when we have them)
# COPY --from=builder /build/templates /usr/local/share/docloom/templates

# Switch to non-root user
USER docloom

# Set working directory
WORKDIR /workspace

# Expose any ports if needed (currently none)
# EXPOSE 8080

# Default command
ENTRYPOINT ["docloom"]
CMD ["--help"]