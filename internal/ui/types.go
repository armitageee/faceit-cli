package ui

import (
	"faceit-cli/internal/config"
	"faceit-cli/internal/entity"
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

// AppModel represents the main application model
type AppModel struct {
	state              AppState
	repo               repository.FaceitRepository
	config             *config.Config
	searchInput        string
	player             *entity.PlayerProfile
	matches            []entity.PlayerMatchSummary
	stats              *PlayerStatsSummary
	matchDetail        *MatchDetail
	selectedMatchIndex int
	playerSwitchInput  string
	recentPlayers      []string
	error              string
	loading            bool
	width              int
	height             int
}

// Custom message types for async operations
type statsLoadedMsg struct {
	stats PlayerStatsSummary
}

type matchDetailLoadedMsg struct {
	matchDetail MatchDetail
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
)
