#!/bin/bash

# Test Coverage Script for Asset Pulse API
# This script runs all tests and generates coverage reports

set -e

echo "🧪 Running Asset Pulse API Tests with Coverage..."

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if we're in the right directory
if [ ! -f "go.mod" ]; then
    print_error "Please run this script from the asset-pulse-api directory"
    exit 1
fi

# Clean previous coverage files
print_status "Cleaning previous coverage files..."
rm -f coverage.out coverage.html

# Run tests with coverage
print_status "Running tests with coverage..."
go test -v -coverprofile=coverage.out -covermode=atomic ./...

# Check if coverage file was created
if [ ! -f "coverage.out" ]; then
    print_error "Coverage file not created. Tests may have failed."
    exit 1
fi

# Generate HTML coverage report
print_status "Generating HTML coverage report..."
go tool cover -html=coverage.out -o coverage.html

# Calculate coverage percentage
COVERAGE=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//')
print_status "Total coverage: ${COVERAGE}%"

# Check if coverage meets minimum requirement (70%)
if (( $(echo "$COVERAGE >= 70" | bc -l) )); then
    print_status "✅ Coverage meets minimum requirement (70%)"
else
    print_warning "⚠️  Coverage is below minimum requirement (70%)"
fi

# Show detailed coverage by package
print_status "Coverage by package:"
go tool cover -func=coverage.out | grep -v "total:"

# Open coverage report if on macOS
if [[ "$OSTYPE" == "darwin"* ]]; then
    print_status "Opening coverage report in browser..."
    open coverage.html
fi

print_status "Test coverage analysis complete!"
print_status "Coverage report saved as: coverage.html"
print_status "Raw coverage data saved as: coverage.out"
