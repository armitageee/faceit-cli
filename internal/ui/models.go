package ui

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"faceit-cli/internal/entity"
	"faceit-cli/internal/repository"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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

// PlayerMatchStats represents player's detailed stats for a match
type PlayerMatchStats struct {
	Kills               int
	Deaths              int
	Assists             int
	KDRatio             float64
	HeadshotsPercentage float64
	ADR                 float64 // Average Damage per Round
	Rating              float64 // HLTV Rating
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
	PlayerTeamID     string
	EnemyTeamID      string
	PlayerTeamPlayers []string
	EnemyTeamPlayers []string
}



// PerformanceMetrics represents advanced performance metrics
type PerformanceMetrics struct {
	ConsistencyScore float64 // How consistent the performance was
	ImpactScore      float64 // How much impact the player had
	ClutchScore      float64 // Performance in clutch situations
	EntryScore       float64 // Performance in entry situations
	SupportScore     float64 // Support/utility usage score
}

// AppState represents the current state of the application
type AppState int

const (
	StateSearch AppState = iota
	StateProfile
	StateMatches
	StateStats
	StateMatchDetail
	StateLoading
	StateError
)

// AppModel is the main application model
type AppModel struct {
	state        AppState
	searchInput  string
	profile      *entity.PlayerProfile
	matches      []entity.PlayerMatchSummary
	stats        *PlayerStatsSummary
	matchDetail  *MatchDetail
	selectedMatchIndex int
	error        string
	loading      bool
	repo         repository.FaceitRepository
	width        int
	height       int
}

// InitialModel returns the initial state of the application
func InitialModel(repo repository.FaceitRepository) AppModel {
	return AppModel{
		state:       StateSearch,
		searchInput: "",
		repo:        repo,
		width:       80,
		height:      24,
		selectedMatchIndex: 0,
	}
}

// Init implements the tea.Model interface
func (m AppModel) Init() tea.Cmd {
	return nil
}

// Update implements the tea.Model interface
func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch m.state {
		case StateSearch:
			return m.updateSearch(msg)
		case StateProfile:
			return m.updateProfile(msg)
		case StateMatches:
			return m.updateMatches(msg)
		case StateStats:
			return m.updateStats(msg)
		case StateMatchDetail:
			return m.updateMatchDetail(msg)
		case StateError:
			return m.updateError(msg)
		}

	case profileLoadedMsg:
		m.profile = &msg.profile
		m.state = StateProfile
		m.loading = false
		return m, nil

	case matchesLoadedMsg:
		m.matches = msg.matches
		m.selectedMatchIndex = 0 // Reset to first match
		m.state = StateMatches
		m.loading = false
		return m, nil

	case statsLoadedMsg:
		m.stats = &msg.stats
		m.state = StateStats
		m.loading = false
		return m, nil

	case matchDetailLoadedMsg:
		m.matchDetail = &msg.matchDetail
		m.state = StateMatchDetail
		m.loading = false
		return m, nil

	case errorMsg:
		m.error = msg.err
		m.state = StateError
		m.loading = false
		return m, nil
	}

	return m, nil
}

// View implements the tea.Model interface
func (m AppModel) View() string {
	switch m.state {
	case StateSearch:
		return m.viewSearch()
	case StateProfile:
		return m.viewProfile()
	case StateMatches:
		return m.viewMatches()
	case StateStats:
		return m.viewStats()
	case StateMatchDetail:
		return m.viewMatchDetail()
	case StateLoading:
		return m.viewLoading()
	case StateError:
		return m.viewError()
	default:
		return "Unknown state"
	}
}

// updateSearch handles key events in the search state
func (m AppModel) updateSearch(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit
	case "enter":
		if strings.TrimSpace(m.searchInput) != "" {
			m.loading = true
			m.state = StateLoading
			return m, m.loadPlayerProfile(m.searchInput)
		}
	case "backspace":
		if len(m.searchInput) > 0 {
			m.searchInput = m.searchInput[:len(m.searchInput)-1]
		}
	default:
		if len(msg.String()) == 1 {
			m.searchInput += msg.String()
		}
	}
	return m, nil
}

// updateProfile handles key events in the profile state
func (m AppModel) updateProfile(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit
	case "esc":
		m.state = StateSearch
		m.searchInput = ""
		return m, nil
	case "m":
		// Load recent matches
		m.loading = true
		m.state = StateLoading
		return m, m.loadRecentMatches()
	case "s":
		// Load statistics
		m.loading = true
		m.state = StateLoading
		return m, m.loadStatistics()
	}
	return m, nil
}

// updateMatches handles key events in the matches state
func (m AppModel) updateMatches(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit
	case "esc":
		m.state = StateProfile
		return m, nil
	case "up", "k":
		// Navigate to previous match
		if m.selectedMatchIndex > 0 {
			m.selectedMatchIndex--
		}
		return m, nil
	case "down", "j":
		// Navigate to next match
		if m.selectedMatchIndex < len(m.matches)-1 {
			m.selectedMatchIndex++
		}
		return m, nil
	case "enter", "d":
		// Load detailed view of the selected match
		if len(m.matches) > 0 && m.selectedMatchIndex < len(m.matches) {
			m.loading = true
			m.state = StateLoading
			return m, m.loadMatchDetail(m.matches[m.selectedMatchIndex].MatchID)
		}
	}
	return m, nil
}

// updateStats handles key events in the stats state
func (m AppModel) updateStats(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit
	case "esc":
		m.state = StateProfile
		return m, nil
	}
	return m, nil
}

// updateMatchDetail handles key events in the match detail state
func (m AppModel) updateMatchDetail(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit
	case "esc":
		m.state = StateMatches
		return m, nil
	}
	return m, nil
}

// updateError handles key events in the error state
func (m AppModel) updateError(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit
	case "esc", "enter":
		m.state = StateSearch
		m.error = ""
		return m, nil
	}
	return m, nil
}

// loadPlayerProfile loads a player profile asynchronously
func (m AppModel) loadPlayerProfile(nickname string) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		profile, err := m.repo.GetPlayerByNickname(ctx, nickname)
		if err != nil {
			return errorMsg{err: err.Error()}
		}
		return profileLoadedMsg{profile: *profile}
	}
}

// loadRecentMatches loads recent matches asynchronously
func (m AppModel) loadRecentMatches() tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		matches, err := m.repo.GetPlayerRecentMatches(ctx, m.profile.ID, "cs2", 10)
		if err != nil {
			return errorMsg{err: err.Error()}
		}
		return matchesLoadedMsg{matches: matches}
	}
}

// loadStatistics loads and calculates statistics from recent matches
func (m AppModel) loadStatistics() tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		matches, err := m.repo.GetPlayerRecentMatches(ctx, m.profile.ID, "cs2", 20)
		if err != nil {
			return errorMsg{err: err.Error()}
		}
		
		stats := calculateStats(matches)
		return statsLoadedMsg{stats: stats}
	}
}

// loadMatchDetail loads detailed statistics for a specific match
func (m AppModel) loadMatchDetail(matchID string) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Get detailed match stats from repository
		matchDetail, err := m.getDetailedMatchStats(ctx, matchID)
		if err != nil {
			return errorMsg{err: err.Error()}
		}
		return matchDetailLoadedMsg{matchDetail: matchDetail}
	}
}

// calculateStats calculates aggregated statistics from matches
func calculateStats(matches []entity.PlayerMatchSummary) PlayerStatsSummary {
	if len(matches) == 0 {
		return PlayerStatsSummary{}
	}

	stats := PlayerStatsSummary{
		TotalMatches: len(matches),
		MapStats:     make(map[string]int),
		BestKDRatio:  -1,
		WorstKDRatio: 999,
	}

	var totalKDRatio, totalHS float64
	mapCount := make(map[string]int)

	for _, match := range matches {
		// Win/Loss counting
		if match.Result == "Win" {
			stats.Wins++
		} else {
			stats.Losses++
		}

		// K/D/A totals
		stats.TotalKills += match.Kills
		stats.TotalDeaths += match.Deaths
		stats.TotalAssists += match.Assists

		// K/D ratio tracking
		if match.KDRatio > 0 {
			totalKDRatio += match.KDRatio
			if match.KDRatio > stats.BestKDRatio {
				stats.BestKDRatio = match.KDRatio
			}
			if match.KDRatio < stats.WorstKDRatio {
				stats.WorstKDRatio = match.KDRatio
			}
		}

		// Headshot percentage
		if match.HeadshotsPercentage > 0 {
			totalHS += match.HeadshotsPercentage
		}

		// Map statistics
		if match.Map != "" {
			mapCount[match.Map]++
		}
	}

	// Calculate averages
	if stats.TotalMatches > 0 {
		stats.WinRate = float64(stats.Wins) / float64(stats.TotalMatches) * 100
		stats.AverageKDRatio = totalKDRatio / float64(stats.TotalMatches)
		stats.AverageHS = totalHS / float64(stats.TotalMatches)
	}

	// Find most played map
	maxCount := 0
	for mapName, count := range mapCount {
		stats.MapStats[mapName] = count
		if count > maxCount {
			maxCount = count
			stats.MostPlayedMap = mapName
		}
	}

	// Handle edge cases
	if stats.BestKDRatio == -1 {
		stats.BestKDRatio = 0
	}
	if stats.WorstKDRatio == 999 {
		stats.WorstKDRatio = 0
	}

	return stats
}

// getDetailedMatchStats retrieves and processes detailed match statistics
func (m AppModel) getDetailedMatchStats(ctx context.Context, matchID string) (MatchDetail, error) {
	// For now, we'll create a mock detailed match with enhanced statistics
	// In a real implementation, this would fetch detailed match data from the API
	
	// Find the match in our recent matches
	var baseMatch *entity.PlayerMatchSummary
	for _, match := range m.matches {
		if match.MatchID == matchID {
			baseMatch = &match
			break
		}
	}
	
	if baseMatch == nil {
		return MatchDetail{}, fmt.Errorf("match not found")
	}
	
	// Create enhanced match detail with calculated metrics
	matchDetail := MatchDetail{
		MatchID:    baseMatch.MatchID,
		Map:        baseMatch.Map,
		FinishedAt: baseMatch.FinishedAt,
		Score:      baseMatch.Score,
		Result:     baseMatch.Result,
		PlayerStats: PlayerMatchStats{
			Kills:               baseMatch.Kills,
			Deaths:              baseMatch.Deaths,
			Assists:             baseMatch.Assists,
			KDRatio:             baseMatch.KDRatio,
			HeadshotsPercentage: baseMatch.HeadshotsPercentage,
			// Calculate additional metrics
			ADR:                 m.calculateADR(baseMatch),
			Rating:              m.calculateHLTVRating(baseMatch),
			FirstKills:          m.calculateFirstKills(baseMatch),
			FirstDeaths:         m.calculateFirstDeaths(baseMatch),
			ClutchWins:          m.calculateClutchWins(baseMatch),
			EntryFrags:          m.calculateEntryFrags(baseMatch),
			FlashAssists:        m.calculateFlashAssists(baseMatch),
			UtilityDamage:       m.calculateUtilityDamage(baseMatch),
		},
		TeamStats: TeamStats{
			PlayerTeamScore:  m.extractPlayerTeamScore(baseMatch.Score),
			EnemyTeamScore:   m.extractEnemyTeamScore(baseMatch.Score),
			PlayerTeamID:     "team1", // Mock data
			EnemyTeamID:      "team2", // Mock data
		},
		PerformanceMetrics: PerformanceMetrics{
			ConsistencyScore: m.calculateConsistencyScore(baseMatch),
			ImpactScore:      m.calculateImpactScore(baseMatch),
			ClutchScore:      m.calculateClutchScore(baseMatch),
			EntryScore:       m.calculateEntryScore(baseMatch),
			SupportScore:     m.calculateSupportScore(baseMatch),
		},
	}
	
	return matchDetail, nil
}

// Helper functions for calculating advanced metrics
func (m AppModel) calculateADR(match *entity.PlayerMatchSummary) float64 {
	// Estimate ADR based on kills and assists (rough calculation)
	estimatedDamage := float64(match.Kills)*100 + float64(match.Assists)*50
	rounds := 30 // Assume average match length
	return estimatedDamage / float64(rounds)
}

func (m AppModel) calculateHLTVRating(match *entity.PlayerMatchSummary) float64 {
	// Simplified HLTV rating calculation
	kd := match.KDRatio
	hs := match.HeadshotsPercentage / 100.0
	adr := m.calculateADR(match)
	
	// Basic rating formula (simplified)
	rating := (kd * 0.4) + (hs * 0.2) + (adr/100 * 0.4)
	if rating > 2.0 {
		rating = 2.0
	}
	return rating
}

func (m AppModel) calculateFirstKills(match *entity.PlayerMatchSummary) int {
	// Estimate first kills as 20% of total kills
	return int(float64(match.Kills) * 0.2)
}

func (m AppModel) calculateFirstDeaths(match *entity.PlayerMatchSummary) int {
	// Estimate first deaths as 15% of total deaths
	return int(float64(match.Deaths) * 0.15)
}

func (m AppModel) calculateClutchWins(match *entity.PlayerMatchSummary) int {
	// Estimate clutch wins based on performance
	if match.KDRatio > 1.5 {
		return 2
	} else if match.KDRatio > 1.0 {
		return 1
	}
	return 0
}

func (m AppModel) calculateEntryFrags(match *entity.PlayerMatchSummary) int {
	// Estimate entry frags as 30% of total kills
	return int(float64(match.Kills) * 0.3)
}

func (m AppModel) calculateFlashAssists(match *entity.PlayerMatchSummary) int {
	// Estimate flash assists as 40% of total assists
	return int(float64(match.Assists) * 0.4)
}

func (m AppModel) calculateUtilityDamage(match *entity.PlayerMatchSummary) int {
	// Estimate utility damage based on assists
	return match.Assists * 25
}

func (m AppModel) extractPlayerTeamScore(score string) int {
	// Parse score string like "16-14"
	parts := strings.Split(score, "-")
	if len(parts) >= 1 {
		if val, err := strconv.Atoi(parts[0]); err == nil {
			return val
		}
	}
	return 0
}

func (m AppModel) extractEnemyTeamScore(score string) int {
	// Parse score string like "16-14"
	parts := strings.Split(score, "-")
	if len(parts) >= 2 {
		if val, err := strconv.Atoi(parts[1]); err == nil {
			return val
		}
	}
	return 0
}



func (m AppModel) calculateConsistencyScore(match *entity.PlayerMatchSummary) float64 {
	// Calculate consistency based on K/D ratio stability
	if match.KDRatio > 1.5 {
		return 0.9
	} else if match.KDRatio > 1.0 {
		return 0.7
	} else if match.KDRatio > 0.8 {
		return 0.5
	}
	return 0.3
}

func (m AppModel) calculateImpactScore(match *entity.PlayerMatchSummary) float64 {
	// Calculate impact based on kills and assists
	impact := float64(match.Kills)*0.6 + float64(match.Assists)*0.4
	return impact / 20.0 // Normalize to 0-1 scale
}

func (m AppModel) calculateClutchScore(match *entity.PlayerMatchSummary) float64 {
	// Calculate clutch performance
	if match.Result == "Win" && match.KDRatio > 1.2 {
		return 0.8
	} else if match.Result == "Win" {
		return 0.6
	}
	return 0.4
}

func (m AppModel) calculateEntryScore(match *entity.PlayerMatchSummary) float64 {
	// Calculate entry performance
	entryKills := m.calculateEntryFrags(match)
	return float64(entryKills) / 10.0 // Normalize
}

func (m AppModel) calculateSupportScore(match *entity.PlayerMatchSummary) float64 {
	// Calculate support performance
	support := float64(match.Assists) + float64(m.calculateFlashAssists(match))
	return support / 15.0 // Normalize
}

// Message types for async operations
type profileLoadedMsg struct {
	profile entity.PlayerProfile
}

type matchesLoadedMsg struct {
	matches []entity.PlayerMatchSummary
}

type statsLoadedMsg struct {
	stats PlayerStatsSummary
}

type matchDetailLoadedMsg struct {
	matchDetail MatchDetail
}

type errorMsg struct {
	err string
}

// Styles
var (
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1)

	searchStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#874BFD")).
			Padding(1, 2)

	profileStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#04B575")).
			Padding(1, 2)

	matchesStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#F25D94")).
			Padding(1, 2)

	statsStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#FFA500")).
			Padding(1, 2)

	matchDetailStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#9D4EDD")).
			Padding(1, 2)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF5F87")).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#FF5F87")).
			Padding(1, 2)

	loadingStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7D56F4")).
			Padding(1, 2)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262")).
			Italic(true)

	winStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#04B575")).
			Bold(true)

	lossStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF5F87")).
			Bold(true)
)

// viewSearch renders the search screen
func (m AppModel) viewSearch() string {
	title := titleStyle.Render("🎮 FACEIT CLI")
	search := searchStyle.Render(fmt.Sprintf("Enter player nickname:\n\n%s", m.searchInput))
	help := helpStyle.Render("Press Enter to search • Ctrl+C or Q to quit")

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center,
		lipgloss.JoinVertical(lipgloss.Center, title, search, help))
}

// viewProfile renders the profile screen
func (m AppModel) viewProfile() string {
	if m.profile == nil {
		return "No profile data"
	}

	title := titleStyle.Render("👤 " + m.profile.Nickname)
	
	var content strings.Builder
	content.WriteString(fmt.Sprintf("Country: %s\n", m.profile.Country))
	content.WriteString(fmt.Sprintf("ID: %s\n", m.profile.ID))
	
	if cs2, ok := m.profile.Games["cs2"]; ok {
		content.WriteString("\n🎯 CS2 Stats:\n")
		content.WriteString(fmt.Sprintf("  ELO: %d\n", cs2.Elo))
		content.WriteString(fmt.Sprintf("  Skill Level: %d\n", cs2.SkillLevel))
		content.WriteString(fmt.Sprintf("  Region: %s\n", cs2.Region))
	}

	profile := profileStyle.Render(content.String())
	help := helpStyle.Render("M - Recent matches • S - Statistics (20 matches) • Esc - Back to search • Ctrl+C or Q to quit")

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center,
		lipgloss.JoinVertical(lipgloss.Center, title, profile, help))
}

// viewMatches renders the matches screen
func (m AppModel) viewMatches() string {
	if len(m.matches) == 0 {
		return "No matches found"
	}

	title := titleStyle.Render("🏆 Recent Matches - " + m.profile.Nickname)
	
	var content strings.Builder
	for i, match := range m.matches {
		// Highlight selected match
		prefix := "  "
		if i == m.selectedMatchIndex {
			prefix = "▶ "
		}
		
		resultStyle := lossStyle
		if match.Result == "Win" {
			resultStyle = winStyle
		}
		
		finishedAt := time.Unix(match.FinishedAt, 0).Format("2006-01-02 15:04")
		
		content.WriteString(fmt.Sprintf("%s%s %s | %s | %s\n", 
			prefix,
			resultStyle.Render(match.Result),
			match.Map,
			match.Score,
			finishedAt))
		content.WriteString(fmt.Sprintf("    K/D/A: %d/%d/%d (%.2f) | HS: %.1f%%\n\n",
			match.Kills, match.Deaths, match.Assists, match.KDRatio, match.HeadshotsPercentage))
	}

	matches := matchesStyle.Render(content.String())
	help := helpStyle.Render("↑↓/KJ - Navigate • Enter/D - View details • Esc - Back to profile • Ctrl+C or Q to quit")

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center,
		lipgloss.JoinVertical(lipgloss.Center, title, matches, help))
}

// viewStats renders the statistics screen
func (m AppModel) viewStats() string {
	if m.stats == nil {
		return "No statistics data"
	}

	title := titleStyle.Render("📊 Statistics (20 matches) - " + m.profile.Nickname)
	
	var content strings.Builder
	content.WriteString("📈 Overall Performance:\n")
	content.WriteString(fmt.Sprintf("  Matches: %d | Wins: %d | Losses: %d\n", 
		m.stats.TotalMatches, m.stats.Wins, m.stats.Losses))
	content.WriteString(fmt.Sprintf("  Win Rate: %.1f%%\n\n", m.stats.WinRate))
	
	content.WriteString("🎯 Combat Statistics:\n")
	content.WriteString(fmt.Sprintf("  Total K/D/A: %d/%d/%d\n", 
		m.stats.TotalKills, m.stats.TotalDeaths, m.stats.TotalAssists))
	content.WriteString(fmt.Sprintf("  Average K/D: %.2f\n", m.stats.AverageKDRatio))
	content.WriteString(fmt.Sprintf("  Best K/D: %.2f | Worst K/D: %.2f\n", 
		m.stats.BestKDRatio, m.stats.WorstKDRatio))
	content.WriteString(fmt.Sprintf("  Average HS%%: %.1f%%\n\n", m.stats.AverageHS))
	
	content.WriteString("🗺️  Map Statistics:\n")
	content.WriteString(fmt.Sprintf("  Most Played: %s\n", m.stats.MostPlayedMap))
	content.WriteString("  Map Breakdown:\n")
	
	// Sort maps by count (most played first)
	type mapStat struct {
		name  string
		count int
	}
	var sortedMaps []mapStat
	for mapName, count := range m.stats.MapStats {
		sortedMaps = append(sortedMaps, mapStat{mapName, count})
	}
	
	// Simple bubble sort for small number of maps
	for i := 0; i < len(sortedMaps)-1; i++ {
		for j := 0; j < len(sortedMaps)-i-1; j++ {
			if sortedMaps[j].count < sortedMaps[j+1].count {
				sortedMaps[j], sortedMaps[j+1] = sortedMaps[j+1], sortedMaps[j]
			}
		}
	}
	
	for _, mapStat := range sortedMaps {
		content.WriteString(fmt.Sprintf("    %s: %d matches\n", mapStat.name, mapStat.count))
	}

	stats := statsStyle.Render(content.String())
	help := helpStyle.Render("Esc - Back to profile • Ctrl+C or Q to quit")

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center,
		lipgloss.JoinVertical(lipgloss.Center, title, stats, help))
}

// viewMatchDetail renders the detailed match statistics screen
func (m AppModel) viewMatchDetail() string {
	if m.matchDetail == nil {
		return "No match detail data"
	}

	title := titleStyle.Render("🔍 Match Details - " + m.matchDetail.Map)
	
	var content strings.Builder
	
	// Match overview
	finishedAt := time.Unix(m.matchDetail.FinishedAt, 0).Format("2006-01-02 15:04")
	resultStyle := lossStyle
	if m.matchDetail.Result == "Win" {
		resultStyle = winStyle
	}
	
	content.WriteString("📊 Match Overview:\n")
	content.WriteString(fmt.Sprintf("  %s | %s | %s\n", 
		resultStyle.Render(m.matchDetail.Result),
		m.matchDetail.Score,
		finishedAt))
	content.WriteString(fmt.Sprintf("  Map: %s | Match ID: %s\n\n", 
		m.matchDetail.Map, m.matchDetail.MatchID))
	
	// Player statistics
	content.WriteString("🎯 Player Performance:\n")
	content.WriteString(fmt.Sprintf("  K/D/A: %d/%d/%d (%.2f)\n", 
		m.matchDetail.PlayerStats.Kills,
		m.matchDetail.PlayerStats.Deaths,
		m.matchDetail.PlayerStats.Assists,
		m.matchDetail.PlayerStats.KDRatio))
	content.WriteString(fmt.Sprintf("  HS%%: %.1f%% | ADR: %.1f | Rating: %.2f\n\n", 
		m.matchDetail.PlayerStats.HeadshotsPercentage,
		m.matchDetail.PlayerStats.ADR,
		m.matchDetail.PlayerStats.Rating))
	
	// Advanced metrics
	content.WriteString("⚡ Advanced Metrics:\n")
	content.WriteString(fmt.Sprintf("  First Kills: %d | First Deaths: %d\n", 
		m.matchDetail.PlayerStats.FirstKills,
		m.matchDetail.PlayerStats.FirstDeaths))
	content.WriteString(fmt.Sprintf("  Clutch Wins: %d | Entry Frags: %d\n", 
		m.matchDetail.PlayerStats.ClutchWins,
		m.matchDetail.PlayerStats.EntryFrags))
	content.WriteString(fmt.Sprintf("  Flash Assists: %d | Utility Damage: %d\n\n", 
		m.matchDetail.PlayerStats.FlashAssists,
		m.matchDetail.PlayerStats.UtilityDamage))
	
	// Performance scores
	content.WriteString("📈 Performance Scores:\n")
	content.WriteString(fmt.Sprintf("  Consistency: %.1f%% | Impact: %.1f%%\n", 
		m.matchDetail.PerformanceMetrics.ConsistencyScore*100,
		m.matchDetail.PerformanceMetrics.ImpactScore*100))
	content.WriteString(fmt.Sprintf("  Clutch: %.1f%% | Entry: %.1f%% | Support: %.1f%%\n\n", 
		m.matchDetail.PerformanceMetrics.ClutchScore*100,
		m.matchDetail.PerformanceMetrics.EntryScore*100,
		m.matchDetail.PerformanceMetrics.SupportScore*100))
	
	// Team statistics
	content.WriteString("👥 Team Statistics:\n")
	content.WriteString(fmt.Sprintf("  Your Team: %d | Enemy Team: %d\n", 
		m.matchDetail.TeamStats.PlayerTeamScore,
		m.matchDetail.TeamStats.EnemyTeamScore))
	


	matchDetail := matchDetailStyle.Render(content.String())
	help := helpStyle.Render("Esc - Back to matches • Ctrl+C or Q to quit")

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center,
		lipgloss.JoinVertical(lipgloss.Center, title, matchDetail, help))
}

// viewLoading renders the loading screen
func (m AppModel) viewLoading() string {
	loading := loadingStyle.Render("⏳ Loading...")
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, loading)
}

// viewError renders the error screen
func (m AppModel) viewError() string {
	error := errorStyle.Render(fmt.Sprintf("❌ Error: %s", m.error))
	help := helpStyle.Render("Esc or Enter - Back to search • Ctrl+C or Q to quit")
	
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center,
		lipgloss.JoinVertical(lipgloss.Center, error, help))
}
