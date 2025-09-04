# Testing Guide

This document explains how to run tests for the FACEIT CLI application.

## Test Types

### Unit Tests
Fast tests that don't require external dependencies:
```bash
go test -v -short ./...
```

### Integration Tests
Tests that require FACEIT API access:
```bash
FACEIT_API_KEY=your_api_key go test -v ./internal/repository/
```

## Setting up GitHub Actions

To enable integration tests in GitHub Actions, you need to add a repository secret:

### 1. Get your FACEIT API Key
1. Go to [FACEIT Developers](https://developers.faceit.com/)
2. Create an account or log in
3. Create a new application
4. Copy your API key

### 2. Add Secret to GitHub Repository
1. Go to your GitHub repository
2. Click **Settings** → **Secrets and variables** → **Actions**
3. Click **New repository secret**
4. Name: `FACEIT_API_KEY`
5. Value: Your FACEIT API key
6. Click **Add secret**

### 3. Verify CI is Working
After adding the secret, push to `main` or `develop` branch to trigger:
- **Unit Tests**: Run on every push/PR
- **Integration Tests**: Run only on pushes to main/develop branches

## Local Testing

### Running All Tests
```bash
# Unit tests only (fast)
make test

# Integration tests (requires API key)
FACEIT_API_KEY=your_key make test-integration
```

### Running Specific Tests
```bash
# Test specific package
go test -v ./internal/repository/

# Test specific function
go test -v ./internal/repository/ -run TestGetPlayerByNickname

# Run benchmarks
go test -v -bench=. ./internal/repository/
```

## Test Coverage

View test coverage:
```bash
go test -v -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
open coverage.html
```

## Troubleshooting

### Tests Skipped
If you see "FACEIT_API_KEY not set, skipping integration test":
- Make sure you've set the environment variable
- Check that your API key is valid
- Verify you have internet connection

### API Rate Limits
- Integration tests make real API calls
- FACEIT has rate limits, so tests may fail if run too frequently
- Consider running integration tests less frequently in CI

### Timeouts
- Integration tests have timeouts to prevent hanging
- If tests timeout, check your internet connection
- Consider increasing timeouts for slow connections
