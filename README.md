# FACEIT CLI

A beautiful terminal user interface (TUI) for viewing FACEIT player profiles and match history, built with Go and Bubble Tea.

![Demo](assets/screen_recorder.gif)

## Features

- 🔍 Search for players by nickname
- 👤 View player profiles with CS2 stats (ELO, skill level, region)
- 🏆 Browse recent match history with detailed statistics and pagination
- 📊 View comprehensive statistics over last 20 matches
- 🔍 Detailed match analysis with advanced metrics
- ⚔️ Compare your stats with friends over last 20 matches
- 🔄 Switch between players without restarting
- 💾 Remember default player via environment variable
- 📝 Centralized logging with configurable levels
- 🚀 Kafka integration for log aggregation (optional)
- ⚡ Fast and responsive with API caching

## Quick Start

### Using Pre-built Binaries

Download the latest release from the [Releases page](https://github.com/armitageee/faceit-cli/releases) and extract the binary for your platform.

### Building from Source

```bash
# Clone the repository
git clone https://github.com/armitageee/faceit-cli.git
cd faceit-cli

# Build the application
make build

# Or build for all platforms
make build-all
```

### Development

```bash
# Install dependencies
make deps

# Run the application
make run

# Run with caching enabled
make run-cache

# Run with all optimizations
make run-optimized

# Run with Kafka logging
make run-kafka

# Run in production mode
make run-production
```

## Configuration

Create a `.env` file in the project root:

```bash
# Required
FACEIT_API_KEY=your_api_key_here

# Optional - Player Settings
FACEIT_DEFAULT_PLAYER=your_nickname_here
COMPARISON_MATCHES=20
MATCHES_PER_PAGE=10
MAX_MATCHES_TO_LOAD=100

# Optional - Logging
LOG_LEVEL=info
LOG_TO_STDOUT=true

# Optional - Caching
CACHE_ENABLED=true
CACHE_TTL=30

# Optional - Kafka Integration
KAFKA_ENABLED=false
KAFKA_BROKERS=localhost:9092
KAFKA_TOPIC=faceit-cli-logs

# Optional - Production Mode
PRODUCTION_MODE=false
```

### Environment Variables

**Required:**
- `FACEIT_API_KEY` (required): Your FACEIT API key

**Player Settings:**
- `FACEIT_DEFAULT_PLAYER` (optional): Default player nickname to load on startup
- `COMPARISON_MATCHES` (optional): Number of matches to use for player comparison (default: 20)
- `MATCHES_PER_PAGE` (optional): Matches per page (default: 10)
- `MAX_MATCHES_TO_LOAD` (optional): Maximum matches to load (default: 100)

**Logging:**
- `LOG_LEVEL` (optional): Log level - debug/info/warn/error (default: info)
- `LOG_TO_STDOUT` (optional): Log to stdout - true/false (default: true)

**Caching:**
- `CACHE_ENABLED` (optional): Enable API response caching - true/false (default: false)
- `CACHE_TTL` (optional): Cache TTL in minutes (default: 30)

**Kafka Integration:**
- `KAFKA_ENABLED` (optional): Enable Kafka logging - true/false (default: false)
- `KAFKA_BROKERS` (optional): Kafka brokers - comma-separated (default: localhost:9092)
- `KAFKA_TOPIC` (optional): Kafka topic for logs (default: faceit-cli-logs)

**Production Mode:**
- `PRODUCTION_MODE` (optional): Enable production mode - true/false (default: false)

## Usage

1. **Search for a player**: Enter a nickname and press Enter
2. **View profile**: See player stats, ELO, skill level, and lifetime statistics
3. **Browse matches**: Press `M` to view recent matches with pagination
4. **View statistics**: Press `S` to see comprehensive stats over last 20 matches
5. **Compare players**: Press `C` to compare with a friend
6. **Switch players**: Press `P` to switch to another player
7. **View match details**: Press `Enter` or `D` on any match for detailed analysis

## Controls

- `↑↓` or `KJ` - Navigate
- `←→` or `HL` - Change pages (in matches view)
- `Enter` or `D` - View details
- `Esc` - Go back
- `Ctrl+C` or `Q` - Quit

## Performance Optimizations

### 🚀 API Response Caching

Reduce API calls and improve response times with intelligent caching:

- **In-memory caching** with configurable TTL
- **Automatic expiration** of stale data
- **Background cleanup** of expired entries

### ⚡ Background Loading

Smart loading strategy for optimal user experience:

1. **Initial Load** - First 20 matches load quickly (30s timeout)
2. **Background Loading** - Remaining matches load in background (120s timeout)
3. **Seamless Updates** - UI updates automatically when more data arrives

## Kafka Integration

Optional centralized logging with Kafka:

```bash
# Start Kafka infrastructure
make kafka-up

# Run with Kafka logging
make run-kafka

# View Kafka UI
make kafka-ui

# Stop Kafka infrastructure
make kafka-down
```

## Development

### Available Make Commands

```bash
make help          # Show all available commands
make fmt           # Format code
make clean         # Clean build artifacts
make deps          # Install dependencies
make install-tools # Install development tools
make run           # Build and run the application
make run-cache     # Run with caching enabled
make run-optimized # Run with all optimizations
make test          # Run tests
make lint          # Run linter
make build         # Build binary
make build-all     # Build for all platforms
```

### Project Structure

```
├── internal/
│   ├── app/          # Application logic
│   ├── cache/        # API response caching
│   ├── config/       # Configuration management
│   ├── entity/       # Data models
│   ├── logger/       # Centralized logging
│   ├── repository/   # API client
│   └── ui/           # TUI components
├── assets/           # Demo GIF and other assets
├── .github/workflows/ # CI/CD pipelines
└── main.go           # Application entry point
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run tests: `make test`
5. Run linter: `make lint`
6. Submit a pull request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Built with [Bubble Tea](https://github.com/charmbracelet/bubbletea) TUI framework
- Uses [Lip Gloss](https://github.com/charmbracelet/lipgloss) for styling
- Powered by [FACEIT Data API](https://developers.faceit.com/)