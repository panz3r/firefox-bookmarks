#!/bin/bash

# Build script for Firefox Bookmarks Converter (Go version)
# Creates binaries for multiple platforms

set -e

echo "Building Firefox Bookmarks Converter for multiple platforms..."

# Create builds directory
mkdir -p builds

# Version info
VERSION=${1:-"1.0.0"}
echo "Building version: $VERSION"

# Build for different platforms
echo "Building for Windows (amd64)..."
GOOS=windows GOARCH=amd64 go build -ldflags "-s -w" -o builds/firefox-bookmarks_windows_amd64.exe

echo "Building for Windows (arm64)..."
GOOS=windows GOARCH=arm64 go build -ldflags "-s -w" -o builds/firefox-bookmarks_windows_arm64.exe

echo "Building for macOS (Intel)..."
GOOS=darwin GOARCH=amd64 go build -ldflags "-s -w" -o builds/firefox-bookmarks_macos_intel

echo "Building for macOS (Apple Silicon)..."
GOOS=darwin GOARCH=arm64 go build -ldflags "-s -w" -o builds/firefox-bookmarks_macos_arm64

echo "Building for Linux (amd64)..."
GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o builds/firefox-bookmarks_linux_amd64

echo "Building for Linux (arm64)..."
GOOS=linux GOARCH=arm64 go build -ldflags "-s -w" -o builds/firefox-bookmarks_linux_arm64

echo "Building for current platform..."
go build -ldflags "-s -w" -o firefox-bookmarks

echo ""
echo "Build complete! Binaries created in builds/ directory:"
ls -la builds/

echo ""
echo "Current platform binary: firefox-bookmarks"
ls -la firefox-bookmarks

echo ""
echo "To test the current platform binary:"
echo "  ./firefox-bookmarks example/test_bookmarks.json"
