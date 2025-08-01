#!/bin/bash

# Minecraft Easy Server Build Script
# Used to package single-file releases with embedded frontend files

set -e

echo "🚀 Starting Minecraft Easy Server build..."

# Clean old build files
echo "🧹 Cleaning old build files..."
rm -f minecraft-server-manager-*

# Build current platform version
echo "📦 Building current platform version..."
go build -ldflags "-s -w" -o minecraft-server-manager-embedded

# Build Windows version
echo "🪟 Building Windows version..."
GOOS=windows GOARCH=amd64 go build -ldflags "-s -w" -o minecraft-server-manager-windows.exe

# Build Linux version
echo "🐧 Building Linux version..."
GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o minecraft-server-manager-linux

# Build macOS version
echo "🍎 Building macOS version..."
GOOS=darwin GOARCH=amd64 go build -ldflags "-s -w" -o minecraft-server-manager-macos

# Display build results
echo "✅ Build completed! Generated files:"
ls -lh minecraft-server-manager-*

echo ""
echo "📋 Usage:"
echo "  - minecraft-server-manager-embedded: Current platform version"
echo "  - minecraft-server-manager-windows.exe: Windows version"
echo "  - minecraft-server-manager-linux: Linux version"
echo "  - minecraft-server-manager-macos: macOS version"
echo ""
echo "🎉 All versions include complete frontend files and can run independently!"