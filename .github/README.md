# GitHub Actions Workflows

This repository includes automated CI/CD workflows for testing, building, and releasing the faceit-cli application.

## Workflows

### 1. CI/CD Pipeline (`.github/workflows/ci.yml`)

**Triggers:**
- Push to `main` or `develop` branches
- Pull requests to `main` branch
- Release events

**Jobs:**

#### Test Job
- Runs on Ubuntu latest
- Sets up Go 1.24
- Runs tests with race detection
- Generates coverage reports
- Uploads coverage to Codecov

#### Lint Job
- Runs on Ubuntu latest
- Uses golangci-lint for code quality checks
- Ensures code follows Go best practices

#### Build Job
- Runs after test and lint jobs pass
- Builds for multiple platforms:
  - Linux (amd64, arm64)
  - Windows (amd64)
  - macOS (amd64)
- Uploads build artifacts

### 2. Release Workflow (`.github/workflows/release.yml`)

**Triggers:**
- Push of version tags (e.g., `v1.0.0`)

**Features:**
- Automatically creates GitHub releases
- Builds binaries for all supported platforms
- Generates checksums for verification
- Includes release notes

## Usage

### Running Tests Locally

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run specific test
go test -v ./internal/ui
```

### Building Locally

```bash
# Build for current platform
make build

# Build for all platforms
make build-all

# Build with specific version
VERSION=v1.0.0 make build
```

### Creating a Release

1. Create and push a version tag:
   ```bash
   git tag v1.0.0
   git push origin v1.0.0
   ```

2. GitHub Actions will automatically:
   - Run tests
   - Build binaries for all platforms
   - Create a GitHub release
   - Upload all artifacts

### Development Tools

```bash
# Install development tools
make install-tools

# Format code
make fmt

# Run linter
make lint

# Clean build artifacts
make clean
```

## Environment Variables

The workflows use the following environment variables:

- `GO_VERSION`: Go version to use (default: 1.24)
- `GITHUB_TOKEN`: Automatically provided by GitHub Actions

## Artifacts

### Build Artifacts
- `faceit-cli-linux-amd64`
- `faceit-cli-linux-arm64`
- `faceit-cli-windows-amd64.exe`
- `faceit-cli-darwin-amd64`
- `faceit-cli-darwin-arm64`

### Release Assets
- All platform binaries
- `checksums.txt` with SHA256 checksums
- Source code archive

## Coverage

Code coverage is automatically generated and uploaded to Codecov. You can view coverage reports at:
- Codecov dashboard
- Local HTML report: `coverage.html`

## Troubleshooting

### Common Issues

1. **Tests failing**: Check that all dependencies are properly imported
2. **Build failing**: Ensure Go version compatibility
3. **Release failing**: Verify tag format (must start with 'v')

### Local Development

```bash
# Install dependencies
make deps

# Run the application
make run

# Check version
./faceit-cli --version
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run tests: `make test`
5. Run linter: `make lint`
6. Submit a pull request

The CI/CD pipeline will automatically test your changes before they can be merged.
