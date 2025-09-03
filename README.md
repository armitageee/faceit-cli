# FACEIT CLI

A beautiful terminal user interface (TUI) for viewing FACEIT player profiles and match history, built with Go and Bubble Tea.

## Features

- 🔍 Search for players by nickname
- 👤 View player profiles with CS2 stats (ELO, skill level, region)
- 🏆 Browse recent match history with detailed statistics and pagination
- 📊 View comprehensive statistics over last 20 matches with streak information
- 🔍 Detailed match analysis with advanced metrics and performance scores
- 🎯 Navigate through all matches and view detailed stats for any match
- ⚔️ Compare your stats with friends over last 20 matches
- 🔄 Switch between players without restarting the application
- 💾 Remember default player via environment variable
- 🎨 Beautiful ASCII logo and terminal interface with colors and styling
- 📝 Centralized logging with configurable levels (debug, info, warn, error)
- 🚀 Kafka integration for log aggregation (optional)
- ⚡ Fast and responsive

## Prerequisites

- Go 1.22 or later
- FACEIT API key (get one from [FACEIT Developers](https://developers.faceit.com/))

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

# Run tests
make test

# Run with coverage
make test-coverage

# Format code
make fmt

# Run linter
make lint

# Install development tools
make install-tools
```

## Installation

1. **Clone the repository:**
```bash
git clone https://github.com/armitageee/faceit-cli.git
cd faceit-cli
```

2. **Install dependencies:**
```bash
go mod tidy
```

3. **Set up your API key:**
```bash
cp .env.example .env
# Edit .env and add your FACEIT_API_KEY
```

Or set the environment variable directly:
```bash
export FACEIT_API_KEY=your_api_key_here
```

### Environment Variables

- `FACEIT_API_KEY` (required): Your FACEIT API key
- `FACEIT_DEFAULT_PLAYER` (optional): Default player nickname to load automatically on startup
- `LOG_LEVEL` (optional): Log level - debug, info, warn, error (default: info)
- `KAFKA_ENABLED` (optional): Enable Kafka logging - true/false (default: false)
- `KAFKA_BROKERS` (optional): Comma-separated Kafka brokers (default: localhost:9092)
- `KAFKA_TOPIC` (optional): Kafka topic for logs (default: faceit-cli-logs)
- `PRODUCTION_MODE` (optional): Enable production mode - true/false (default: false)
- `LOG_TO_STDOUT` (optional): Log to stdout - true/false (default: true)
- `MATCHES_PER_PAGE` (optional): Matches per page (default: 10)
- `MAX_MATCHES_TO_LOAD` (optional): Maximum matches to load (default: 100)
- `CACHE_ENABLED` (optional): Enable API response caching - true/false (default: false)
- `CACHE_TTL` (optional): Cache TTL in minutes (default: 30)
- `COMPARISON_MATCHES` (optional): Number of matches to use for player comparison (default: 20)

**Example `.env` file:**
```bash
FACEIT_API_KEY=your_api_key_here
FACEIT_DEFAULT_PLAYER=your_nickname_here
LOG_LEVEL=debug
KAFKA_ENABLED=false
PRODUCTION_MODE=false
LOG_TO_STDOUT=true
MATCHES_PER_PAGE=10
MAX_MATCHES_TO_LOAD=100
CACHE_ENABLED=false
CACHE_TTL=30
```

## Usage

**Run the application:**
```bash
go run main.go
```

### Controls

- **Search Screen**: Type a player nickname and press Enter
- **Profile Screen**: 
  - `M` - View recent matches (10 matches)
  - `S` - View statistics (20 matches)
  - `C` - Compare with friend
  - `P` - Switch to another player
  - `Esc` - Back to search
  - `Ctrl+C` or `Q` - Quit
- **Matches Screen**:
  - `↑`/`↓` or `K`/`J` - Navigate through matches
  - `←`/`→` or `H`/`L` - Navigate between pages
  - `Enter` or `D` - View detailed analysis of selected match
  - `Esc` - Back to profile
  - `Ctrl+C` or `Q` - Quit
- **Match Detail Screen**:
  - `Esc` - Back to matches
  - `Ctrl+C` or `Q` - Quit
- **Statistics Screen**:
  - `Esc` - Back to profile
  - `Ctrl+C` or `Q` - Quit
- **Player Switch Screen**:
  - Type player nickname and press `Enter` to switch
  - `Esc` - Back to profile
  - `Ctrl+C` or `Q` - Quit
- **Comparison Input Screen**:
  - Type friend's nickname and press `Enter` to compare
  - `Esc` - Back to profile
  - `Ctrl+C` or `Q` - Quit
- **Comparison Screen**:
  - `Esc` - Back to profile
  - `Ctrl+C` or `Q` - Quit
- **Error Screen**:
  - `Esc` or `Enter` - Back to search
  - `Ctrl+C` or `Q` - Quit

## Building

**Build the binary:**
```bash
go build -o faceit-cli main.go
```

## Project Structure

```
faceit-cli/
├── main.go                 # Application entry point
├── internal/
│   ├── app/               # Application logic
│   ├── config/            # Configuration management
│   ├── entity/            # Data models
│   ├── repository/        # FACEIT API client
│   └── ui/                # TUI models and views
├── go.mod                 # Go module file
├── .env.example          # Environment variables template
└── README.md             # This file
```

## API Integration

This CLI uses the official **FACEIT Data API v4**. It fetches:

- Player profiles and basic information
- CS2 game statistics (ELO, skill level, region)
- Recent match history with detailed performance metrics

## Advanced Match Analysis

The detailed match view includes:

### 📊 Basic Statistics
- K/D/A ratio and headshot percentage
- ADR (Average Damage per Round)
- **HLTV Rating**: Standard CS:GO/CS2 performance metric
  - 2.0+ = Incredible performance
  - 1.5-2.0 = Excellent performance
  - 1.2-1.5 = Good performance
  - 1.0-1.2 = Average performance
  - 0.8-1.0 = Below average
  - < 0.8 = Poor performance

### ⚡ Advanced Metrics
- **First Kills/Deaths**: Opening duels won/lost
- **Clutch Wins**: 1vX situations won
- **Entry Frags**: Opening kills for team
- **Flash Assists**: Flashbang assists
- **Utility Damage**: Damage from grenades/utility

### 📈 Performance Scores
- **Consistency Score**: Performance stability
- **Impact Score**: Overall match impact
- **Clutch Score**: Performance in clutch situations
- **Entry Score**: Entry fragging effectiveness
- **Support Score**: Support/utility usage

### 🔥 Streak Information
- **Current Streak**: Shows current win or loss streak
- **Longest Streaks**: Displays record win and loss streaks
- **Recent Performance**: Average K/D over last 5 matches
- **Performance Tracking**: Monitor improvement or decline trends

## Player Comparison

Compare your performance with friends over the last 20 matches (configurable via `COMPARISON_MATCHES` environment variable)! This feature provides detailed head-to-head statistics to see who's performing better. Both players are compared using exactly the same number of recent matches for fair comparison.

### ⚔️ How to Use

1. **Open your profile** (search for your nickname)
2. **Press `C`** to start comparison
3. **Enter friend's nickname** and press Enter
4. **View detailed comparison** with color-coded results

### 📊 Comparison Metrics

The comparison shows side-by-side statistics with **color-coded differences**:

#### 🎯 Basic Statistics
- **K/D Ratio**: Average kills per death
- **Win Rate**: Percentage of matches won
- **Headshots**: Average headshot percentage

#### 🔫 Kills & Deaths
- **Total Kills**: Total kills across 20 matches
- **Total Deaths**: Total deaths across 20 matches  
- **Total Assists**: Total assists across 20 matches

#### 🏆 Performance
- **Best K/D**: Highest K/D ratio achieved
- **Worst K/D**: Lowest K/D ratio recorded

#### 🗺️ Maps
- **Most Played Together**: Map both players played most
- **Common Maps**: Number of maps both players have played

### 🎨 Color Coding

- **🟢 Green (+X.XX)**: You're performing better
- **🔴 Red (-X.XX)**: Friend is performing better
- **🔵 Teal**: Your statistics
- **🔴 Red**: Friend's statistics

### 💡 Example Comparison

```
wUwunchik vs s1mple

📊 Basic Statistics:
  K/D Ratio: 1.25 vs 1.45 (-0.20)
  Win Rate: 65.0% vs 70.0% (-5.0%)
  Headshots: 55.2% vs 58.1% (-2.9%)

🎯 Kills & Deaths:
  Total Kills: 245 vs 289 (-44)
  Total Deaths: 196 vs 199 (+3)
  Total Assists: 89 vs 95 (-6)

🏆 Performance:
  Best K/D: 2.1 vs 2.8 (-0.7)
  Worst K/D: 0.6 vs 0.8 (+0.2)

🗺️ Maps:
  Most Played Together: de_dust2
  Common Maps: 8
```

This shows that while s1mple has better overall stats, you have fewer deaths and a better worst performance!



## Match Navigation

The matches screen now supports full navigation through all recent matches:

- **Visual Selection**: The currently selected match is highlighted with a `▶` arrow
- **Keyboard Navigation**: Use arrow keys (`↑`/`↓`) or Vim-style keys (`K`/`J`) to navigate
- **Detailed View**: Press `Enter` or `D` to view detailed statistics for the selected match
- **All Matches**: Navigate through all 20 recent matches, not just the first 10

## Match History Pagination

Browse through your complete match history with intuitive pagination controls! No more limitations on viewing your recent matches.

### 🏆 How Pagination Works

- **Page Navigation**: Use `←`/`→` or `H`/`L` keys to move between pages
- **Configurable Page Size**: Set `MATCHES_PER_PAGE` environment variable (default: 10)
- **Smart Loading**: Matches are loaded on-demand as you navigate
- **Visual Indicators**: See current page, match range, and availability of more pages

### 📊 Pagination Features

#### **Visual Feedback**
- **Page Information**: "Page 2 | Matches 11-20"
- **Navigation Hints**: "More available (→)" and "Previous (←)"
- **Match Range**: Shows which matches you're currently viewing

#### **Keyboard Controls**
- `←`/`H` - Go to previous page
- `→`/`L` - Go to next page (if available)
- `↑`/`↓`/`K`/`J` - Navigate within current page
- `Enter`/`D` - View detailed match analysis

#### **Configuration**
```bash
# Set matches per page (default: 10)
MATCHES_PER_PAGE=15

# Set maximum matches to load (default: 100)
MAX_MATCHES_TO_LOAD=200
```

### 🎯 Benefits

- **Complete History**: Access all your recent matches, not just the first 10
- **Performance**: Load matches efficiently without overwhelming the interface
- **Flexibility**: Configure page size based on your preference
- **Intuitive**: Familiar arrow key navigation for page switching

## Performance Optimizations

The application includes several performance optimizations to provide a smooth user experience:

### 🚀 API Response Caching

Reduce API calls and improve response times with intelligent caching:

#### **Features**
- **In-memory caching** with configurable TTL
- **Automatic expiration** of stale data
- **Background cleanup** of expired entries
- **Cache statistics** for monitoring

#### **Configuration**
```bash
# Enable caching (default: false)
CACHE_ENABLED=true

# Cache TTL in minutes (default: 30)
CACHE_TTL=30
```

#### **Benefits**
- **Faster navigation** - cached data loads instantly
- **Reduced API calls** - less load on Faceit servers
- **Better reliability** - works even with slow API responses
- **Cost efficiency** - fewer API requests

### ⚡ Background Loading

Smart loading strategy for optimal user experience:

#### **How it Works**
1. **Initial Load** - First 20 matches load quickly (30s timeout)
2. **Background Loading** - Remaining matches load in background (120s timeout)
3. **Seamless Updates** - UI updates automatically when more data arrives
4. **No Blocking** - User can navigate while data loads

#### **Benefits**
- **Instant Display** - First matches appear quickly
- **Progressive Enhancement** - More data appears as it loads
- **Non-blocking** - User can interact while loading
- **Graceful Degradation** - Works even if background loading fails

### 🎯 Cursor Management

Improved navigation experience:

- **Smart Positioning** - Cursor always starts at top of new page
- **Consistent Behavior** - Predictable navigation across pages
- **Visual Feedback** - Clear indication of current position

## Logging and Monitoring

The application features a comprehensive logging system with the following capabilities:

### Log Levels
- **debug**: Detailed information for debugging
- **info**: General information about application flow  
- **warn**: Warning messages for potential issues
- **error**: Error messages for failures

### Kafka Integration
Optional Kafka integration for log aggregation and monitoring:

```bash
# Start Kafka infrastructure
make kafka-up

# Run with Kafka logging enabled
make run-kafka

# View Kafka UI
make kafka-ui  # Opens http://localhost:8080
```

### Log Examples
```bash
# Development mode
LOG_LEVEL=debug go run main.go

# Run with Kafka integration
KAFKA_ENABLED=true LOG_LEVEL=info go run main.go

# Production mode (no stdout logs)
make run-production

# Production mode with Kafka
make run-prod-kafka
```

For detailed logging configuration and Kafka setup, see [LOGGING.md](LOGGING.md).

## CI/CD

This project uses **GitHub Actions** for automated testing, building, and releasing:

- **Automated Testing**: Runs tests on every push and pull request
- **Multi-platform Builds**: Builds binaries for Linux, Windows, and macOS
- **Automatic Releases**: Creates releases when version tags are pushed
- **Code Quality**: Runs linters and coverage reports

See [.github/README.md](.github/README.md) for detailed information about the CI/CD pipeline.

## Contributing

1. **Fork the repository**
2. **Create a feature branch**
3. **Make your changes**
4. **Run tests:** `make test`
5. **Run linter:** `make lint`
6. **Submit a pull request**

The CI/CD pipeline will automatically test your changes before they can be merged.

## License

This project is licensed under the **MIT License**.
