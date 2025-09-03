# FACEIT CLI

A beautiful terminal user interface (TUI) for viewing FACEIT player profiles and match history, built with Go and Bubble Tea.

## Features

- 🔍 Search for players by nickname
- 👤 View player profiles with CS2 stats (ELO, skill level, region)
- 🏆 Browse recent match history with detailed statistics
- 📊 View comprehensive statistics over last 20 matches with streak information
- 🔍 Detailed match analysis with advanced metrics and performance scores
- 🎯 Navigate through all matches and view detailed stats for any match
- 🔄 Switch between players without restarting the application
- 💾 Remember default player via environment variable
- 🎨 Beautiful ASCII logo and terminal interface with colors and styling
- ⚡ Fast and responsive

## Prerequisites

- Go 1.23 or later
- FACEIT API key (get one from [FACEIT Developers](https://developers.faceit.com/))

## Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd faceit-cli
```

2. Install dependencies:
```bash
go mod tidy
```

3. Set up your API key:
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

Example `.env` file:
```bash
FACEIT_API_KEY=your_api_key_here
FACEIT_DEFAULT_PLAYER=your_nickname_here
```

## Usage

Run the application:
```bash
go run main.go
```

### Controls

- **Search Screen**: Type a player nickname and press Enter
- **Profile Screen**: 
  - `M` - View recent matches (10 matches)
  - `S` - View statistics (20 matches)
  - `P` - Switch to another player
  - `Esc` - Back to search
  - `Ctrl+C` or `Q` - Quit
- **Matches Screen**:
  - `↑`/`↓` or `K`/`J` - Navigate through matches
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
- **Error Screen**:
  - `Esc` or `Enter` - Back to search
  - `Ctrl+C` or `Q` - Quit

## Building

Build the binary:
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

This CLI uses the official FACEIT Data API v4. It fetches:
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



## Match Navigation

The matches screen now supports full navigation through all recent matches:

- **Visual Selection**: The currently selected match is highlighted with a `▶` arrow
- **Keyboard Navigation**: Use arrow keys (`↑`/`↓`) or Vim-style keys (`K`/`J`) to navigate
- **Detailed View**: Press `Enter` or `D` to view detailed statistics for the selected match
- **All Matches**: Navigate through all 20 recent matches, not just the first 10

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test thoroughly
5. Submit a pull request

## License

This project is licensed under the MIT License.
