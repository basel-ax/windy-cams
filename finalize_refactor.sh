#!/bin/bash

# This script helps finalize the refactoring from "Platforms" to "Webcams".

set -e

# 1. File System Cleanup
echo "Verifying file system cleanup..."
if [ -d "internal/platform" ]; then
    echo "Error: The directory 'internal/platform' still exists. Please remove it."
    exit 1
fi

if [ -f "internal/domain/platform.go" ]; then
    echo "Error: The file 'internal/domain/platform.go' still exists. Please remove it."
    exit 1
fi
echo "File system checks passed."
echo ""

# 2. Dependency Management
echo "Tidying go modules..."
go mod tidy
echo "go mod tidy finished."
echo ""

# 3. Code Formatting
echo "Formatting Go code..."
go fmt ./...
echo "go fmt finished."
echo ""

# 4. Verification
echo "Running tests..."
go test ./...
echo "Tests passed."
echo ""

echo "Build and run verification:"
echo "You can now run the application to perform a final check:"
echo "go run cmd/app/main.go"
echo ""

echo "Refactoring finalization is complete."
