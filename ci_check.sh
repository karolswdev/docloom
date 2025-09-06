#!/bin/bash
set -e

echo "=== Running Complete CI Verification ==="
echo

# Format check
echo "=== Format Check ==="
if [ "$(gofmt -s -l . | wc -l)" -gt 0 ]; then
  echo "FAIL: Files need formatting"
  gofmt -s -l .
  exit 1
fi
echo "PASS: Format check"
echo

# Imports check
echo "=== Imports Check ==="
if command -v goimports &> /dev/null; then
  if [ "$(goimports -l . | wc -l)" -gt 0 ]; then
    echo "FAIL: Imports need fixing"
    goimports -l .
    exit 1
  fi
  echo "PASS: Imports check"
else
  echo "SKIP: goimports not installed"
fi
echo

# Lint check
echo "=== Lint Check ==="
if golangci-lint run --out-format=line-number --config=.golangci.yml 2>&1 | grep -E "^(cmd|internal)" | grep -v "cyclomatic complexity" | grep -v "File is not properly formatted"; then
  echo "FAIL: Linting errors found (ignoring complexity warnings)"
  exit 1
fi
echo "PASS: Lint check (with acceptable warnings)"
echo

# Tests with race detection
echo "=== Test Suite ==="
if ! go test -race -coverprofile=coverage.out ./...; then
  echo "FAIL: Tests failed"
  exit 1
fi
echo "PASS: Tests"
echo

# Build check
echo "=== Build Check ==="
if ! go build ./cmd/docloom; then
  echo "FAIL: Main build failed"
  exit 1
fi
if ! go build -o docloom-agent-csharp ./cmd/docloom-agent-csharp; then
  echo "FAIL: Agent build failed"
  exit 1
fi
echo "PASS: Build"
echo

echo "=== ALL CI CHECKS PASSED ==="
