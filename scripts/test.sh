#!/bin/bash

# Hospital EMR System - Test Script
# Runs comprehensive test suite

set -e

echo "=========================================="
echo "Hospital EMR System - Test Suite"
echo "=========================================="
echo ""

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# Test results
TESTS_PASSED=0
TESTS_FAILED=0

# Function to run tests
run_test() {
    local test_name=$1
    local test_command=$2
    
    echo -e "${BLUE}Running: ${test_name}${NC}"
    if eval $test_command; then
        echo -e "${GREEN}✓ ${test_name} passed${NC}"
        ((TESTS_PASSED++))
    else
        echo -e "${RED}✗ ${test_name} failed${NC}"
        ((TESTS_FAILED++))
    fi
    echo ""
}

# Unit Tests
echo "1. Unit Tests"
echo "----------------------------------------"
run_test "Unit Tests" "go test -v -short -race ./..."

# Integration Tests
echo "2. Integration Tests"
echo "----------------------------------------"
run_test "Integration Tests" "go test -v -run Integration ./..."

# Code Coverage
echo "3. Code Coverage"
echo "----------------------------------------"
run_test "Coverage Analysis" "go test -coverprofile=coverage.out ./... && go tool cover -func=coverage.out"

# Linting
echo "4. Code Linting"
echo "----------------------------------------"
if command -v golangci-lint &> /dev/null; then
    run_test "Linting" "golangci-lint run ./..."
else
    echo -e "${YELLOW}⚠ golangci-lint not installed, skipping...${NC}"
fi

# Security Scan
echo "5. Security Scan"
echo "----------------------------------------"
if command -v gosec &> /dev/null; then
    run_test "Security Scan" "gosec -quiet ./..."
else
    echo -e "${YELLOW}⚠ gosec not installed, skipping...${NC}"
fi

# Vet
echo "6. Go Vet"
echo "----------------------------------------"
run_test "Go Vet" "go vet ./..."

# Format Check
echo "7. Format Check"
echo "----------------------------------------"
run_test "Format Check" "test -z \$(gofmt -l .)"

# Dependency Check
echo "8. Dependency Check"
echo "----------------------------------------"
run_test "Module Verification" "go mod verify"

# Summary
echo "=========================================="
echo "Test Summary"
echo "=========================================="
echo -e "Tests Passed: ${GREEN}${TESTS_PASSED}${NC}"
echo -e "Tests Failed: ${RED}${TESTS_FAILED}${NC}"
echo ""

if [ $TESTS_FAILED -eq 0 ]; then
    echo -e "${GREEN}All tests passed! ✓${NC}"
    exit 0
else
    echo -e "${RED}Some tests failed! ✗${NC}"
    exit 1
fi
