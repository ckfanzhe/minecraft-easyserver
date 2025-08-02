# Multi-stage build Dockerfile for Minecraft Server Manager

# Build stage
FROM golang:1.23.11 AS builder

# # Install build dependencies
# RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /app

# Configure Go proxy for China network environment
ENV GOPROXY=https://goproxy.cn,https://goproxy.io,direct
ENV GOSUMDB=sum.golang.google.cn
ENV GO111MODULE=on

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application with embedded web assets
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o minecraft-easyserver .

# Runtime stage
FROM ubuntu:22.04

# Install runtime dependencies
RUN apt-get update && apt-get install -y \
    ca-certificates \
    wget \
    unzip \
    && rm -rf /var/lib/apt/lists/*

# Create non-root user
RUN useradd -m -u 1000 minecraft

# Create directories
RUN mkdir -p /data/bedrock-server && \
    chown -R minecraft:minecraft /data

# Copy the binary from builder stage
COPY --from=builder /app/minecraft-easyserver /data/minecraft-easyserver

# Set permissions
RUN chmod +x /data/minecraft-easyserver && \
    chown minecraft:minecraft /data/minecraft-easyserver

# Switch to non-root user
USER minecraft

# Set working directory
WORKDIR /data

# Expose ports
EXPOSE 8080 19132/udp 19133/udp

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/ || exit 1

# Run the application
CMD ["./minecraft-easyserver"]