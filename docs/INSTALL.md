# Installation Guide

## Quick Install (Recommended)

### 1. Install via Go

```bash
# Install the latest version
go install github.com/armitageee/faceit-cli@latest

# Initialize configuration
faceit-cli init

# Edit the config file and add your API key
# ~/.config/faceit-cli/config.yml

# Run the application
faceit-cli
```

### 2. Configure

After running `faceit-cli init`, edit the configuration file:

```bash
# Edit the config file
nano ~/.config/faceit-cli/config.yml
# or
code ~/.config/faceit-cli/config.yml
```

Set your FACEIT API key:

```yaml
api_key: "your_actual_faceit_api_key_here"
```

## Alternative Installation Methods

### Pre-built Binaries

1. Download the latest release from [GitHub Releases](https://github.com/armitageee/faceit-cli/releases)
2. Extract the binary for your platform
3. Move to a directory in your PATH (e.g., `/usr/local/bin` on macOS/Linux)
4. Run `faceit-cli init` to create configuration

### Docker

```bash
# Clone the repository
git clone https://github.com/armitageee/faceit-cli.git
cd faceit-cli

# Set up environment
cp .env.example .env
# Edit .env and add your FACEIT_API_KEY

# Build and run with Docker
make docker-build
make docker-run
```

### Building from Source

```bash
# Clone the repository
git clone https://github.com/armitageee/faceit-cli.git
cd faceit-cli

# Build the application
make build

# Initialize configuration
./faceit-cli init

# Run the application
./faceit-cli
```

## Configuration

The application uses a flexible configuration system with priority order:

1. **Environment Variables** (highest priority)
2. **YAML Configuration** (`~/.config/faceit-cli/config.yml`)
3. **Default Values** (lowest priority)

### 1. YAML Configuration

Location: `~/.config/faceit-cli/config.yml`

```yaml
# Required
api_key: "your_faceit_api_key_here"

# Optional settings
default_player: "your_nickname"
log_level: "info"
matches_per_page: 10
max_matches_to_load: 100
cache_enabled: false
telemetry_enabled: false
# ... and more
```

### 2. Environment Variables Override

Environment variables always override YAML configuration. This is useful for:
- Production deployments
- CI/CD pipelines
- Temporary overrides
- Security-sensitive values

If no YAML config is found, the application falls back to environment variables:

```bash
export FACEIT_API_KEY="your_api_key_here"
export FACEIT_DEFAULT_PLAYER="your_nickname"
export LOG_LEVEL="info"
# ... and more
```

## Platform-specific Notes

### macOS

- Config location: `~/.config/faceit-cli/config.yml`
- Install via Homebrew: `go install github.com/armitageee/faceit-cli@latest`

### Linux

- Config location: `~/.config/faceit-cli/config.yml`
- Make sure `~/.local/bin` is in your PATH for `go install`

### Windows

- Config location: `%USERPROFILE%\.config\faceit-cli\config.yml`
- Use PowerShell or Command Prompt for installation

## Troubleshooting

### "command not found: faceit-cli"

Make sure the Go binary directory is in your PATH:

```bash
# Add to ~/.bashrc or ~/.zshrc
export PATH=$PATH:$(go env GOPATH)/bin
```

### "API key not configured"

1. Run `faceit-cli init` to create the config file
2. Edit `~/.config/faceit-cli/config.yml`
3. Set your `api_key` value

### "config file not found"

The application will create the config directory automatically when you run `faceit-cli init`.

## Getting Your FACEIT API Key

1. Go to [FACEIT Developer Portal](https://developers.faceit.com/)
2. Sign in with your FACEIT account
3. Create a new application
4. Copy the API key
5. Add it to your configuration file

## Updating

To update to the latest version:

```bash
go install github.com/armitageee/faceit-cli@latest
```

Your configuration will be preserved.
