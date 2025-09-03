package ui

import (
	"faceit-cli/internal/config"
	"faceit-cli/internal/entity"
	"faceit-cli/internal/logger"
	"faceit-cli/internal/repository"

	"github.com/charmbracelet/lipgloss"
)

// AppState represents the current state of the application
type AppState int

const (
	StateSearch AppState = iota
	StateProfile
	StateMatches
	StateStats
	StateMatchDetail
	StatePlayerSwitch
	StateComparisonInput
	StateComparison
	StateLoading
	StateError
)

// PlayerStatsSummary represents aggregated statistics over recent matches
type PlayerStatsSummary struct {
	TotalMatches     int
	Wins             int
	Losses           int
	WinRate          float64
	TotalKills       int
	TotalDeaths      int
	TotalAssists     int
	AverageKDRatio   float64
	TotalKDA         float64 // Total K/D ratio (total kills / total deaths)
	AverageHS        float64
	BestKDRatio      float64
	WorstKDRatio     float64
	MostPlayedMap    string
	MapStats         map[string]int
	KDChartData      []float64 // K/D ratios for chart
	CurrentStreak    int    // positive for win streak, negative for loss streak
	StreakType       string // "win" or "loss"
	LongestWinStreak int
	LongestLossStreak int
}

// MatchDetail represents detailed statistics for a single match
type MatchDetail struct {
	MatchID             string
	Map                 string
	FinishedAt          int64
	Score               string
	Result              string
	PlayerStats         PlayerMatchStats
	TeamStats           TeamStats
	PerformanceMetrics  PerformanceMetrics
}

// PlayerMatchStats represents detailed player statistics for a match
type PlayerMatchStats struct {
	Kills               int
	Deaths              int
	Assists             int
	KDRatio             float64
	HeadshotsPercentage float64
	ADR                 float64
	HLTVRating          float64
	FirstKills          int
	FirstDeaths         int
	ClutchWins          int
	EntryFrags          int
	FlashAssists        int
	UtilityDamage       int
}

// TeamStats represents team-level statistics
type TeamStats struct {
	PlayerTeamScore  int
	EnemyTeamScore   int
	PlayerTeamRounds int
	EnemyTeamRounds  int
}

// PerformanceMetrics represents advanced performance metrics
type PerformanceMetrics struct {
	ConsistencyScore float64
	ImpactScore      float64
	ClutchScore      float64
	EntryScore       float64
	SupportScore     float64
}

// PlayerComparison represents comparison data between two players
type PlayerComparison struct {
	Player1Nickname string
	Player2Nickname string
	Player1Stats    PlayerStatsSummary
	Player2Stats    PlayerStatsSummary
	ComparisonData  ComparisonData
}

// ComparisonData represents detailed comparison metrics
type ComparisonData struct {
	KDRatioDiff      float64
	TotalKDADiff     float64 // Difference in total K/D ratio
	WinRateDiff      float64
	AverageHSDiff    float64
	TotalKillsDiff   int
	TotalDeathsDiff  int
	TotalAssistsDiff int
	BestKDDiff       float64
	WorstKDDiff      float64
	MostPlayedMap    string
	CommonMaps       []string
}

// AppModel represents the main application model
type AppModel struct {
	state              AppState
	repo               repository.FaceitRepository
	config             *config.Config
	logger             *logger.Logger
	searchInput        string
	player             *entity.PlayerProfile
	matches            []entity.PlayerMatchSummary
	stats              *PlayerStatsSummary
	lifetimeStats      *entity.PlayerStats
	matchDetail        *MatchDetail
	selectedMatchIndex int
	playerSwitchInput  string
	recentPlayers      []string
	comparison         *PlayerComparison
	comparisonInput    string
	error              string
	loading            bool
	width              int
	height             int
	// Pagination fields
	currentPage        int
	totalMatches       int
	matchesPerPage     int
	hasMoreMatches     bool
}

// Custom message types for async operations
type statsLoadedMsg struct {
	stats PlayerStatsSummary
}

type matchDetailLoadedMsg struct {
	matchDetail MatchDetail
}

type comparisonLoadedMsg struct {
	comparison PlayerComparison
}

type lifetimeStatsLoadedMsg struct {
	stats *entity.PlayerStats
}

type matchesPageLoadedMsg struct {
	matches      []entity.PlayerMatchSummary
	page         int
	hasMore      bool
	totalMatches int
}

type backgroundMatchesLoadedMsg struct {
	matches []entity.PlayerMatchSummary
}

// Styling constants
var (
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF5F87"))

	matchDetailStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#7D56F4")).
			Padding(1, 2).
			Margin(1, 0)

	comparisonStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#7D56F4")).
			Padding(1, 2).
			Margin(1, 0)

	player1Style = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#4ECDC4")).
			Bold(true)

	player2Style = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF6B6B")).
			Bold(true)

	betterStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#96CEB4")).
			Bold(true)

	worseStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF5F87")).
			Bold(true)
)
