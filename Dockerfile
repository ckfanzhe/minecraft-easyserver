# Multi-stage build Dockerfile for Minecraft Server Manager

# Build stage
FROM node:22.18-alpine AS frontend-builder

# Set working directory for frontend
WORKDIR /app/minecraft-easyserver-web

# Copy frontend package files
COPY minecraft-easyserver-web/package*.json ./

# Install frontend dependencies
RUN npm ci --only=production

# Copy frontend source code
COPY minecraft-easyserver-web/ ./

# Build frontend
RUN npm run build

# Go build stage
FROM golang:1.23.11 AS go-builder


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

# Copy frontend build artifacts
COPY --from=frontend-builder /app/minecraft-easyserver-web/dist ./web/dist

# Build the application with embedded web assets
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o minecraft-easyserver .

# Runtime stage
FROM ubuntu:22.04

# Install runtime dependencies
RUN apt-get update && apt-get install -y \
    ca-certificates libcurl4-openssl-dev \
    wget \
    unzip \
    && rm -rf /var/lib/apt/lists/*

# Create directories
RUN mkdir -p /data/bedrock-server

# Copy the binary from go-builder stage
COPY --from=go-builder /app/minecraft-easyserver /data/minecraft-easyserver

# Set permissions
RUN chmod +x /data/minecraft-easyserver

# Set working directory
WORKDIR /data

# Expose ports
EXPOSE 8080 19132/udp 19133/udp

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/ || exit 1

# Run the application
CMD ["./minecraft-easyserver"]