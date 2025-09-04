package ui

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"faceit-cli/internal/entity"

	tea "github.com/charmbracelet/bubbletea"
)

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
	case "c":
		// Compare with friend
		m.state = StateComparisonInput
		m.comparisonInput = ""
		return m, nil
	case "p":
		// Switch player
		m.state = StatePlayerSwitch
		m.playerSwitchInput = ""
		return m, nil
	}
	return m, nil
}

// updatePlayerSwitch handles key events in the player switch state
func (m AppModel) updatePlayerSwitch(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit
	case "esc":
		m.state = StateProfile
		m.playerSwitchInput = ""
		return m, nil
	case "enter":
		if m.playerSwitchInput != "" {
			// Add current player to recent players if not already there
			if m.player != nil {
				m.addToRecentPlayers(m.player.Nickname)
			}
			// Switch to new player
			m.searchInput = m.playerSwitchInput
			m.state = StateLoading
			m.loading = true
			return m, m.loadPlayerProfile(m.playerSwitchInput)
		}
		return m, nil
	case "backspace":
		if len(m.playerSwitchInput) > 0 {
			m.playerSwitchInput = m.playerSwitchInput[:len(m.playerSwitchInput)-1]
		}
		return m, nil
	default:
		if len(msg.String()) == 1 {
			m.playerSwitchInput += msg.String()
		}
		return m, nil
	}
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
	case "left", "h":
		// Go to previous page
		if m.currentPage > 1 {
			m.currentPage--
			// Always reset cursor to first position on the new page
			m.selectedMatchIndex = (m.currentPage - 1) * m.matchesPerPage
		}
		return m, nil
	case "right", "l":
		// Go to next page
		totalPages := (len(m.matches) + m.matchesPerPage - 1) / m.matchesPerPage
		if m.currentPage < totalPages {
			m.currentPage++
			// Always reset cursor to first position on the new page
			m.selectedMatchIndex = (m.currentPage - 1) * m.matchesPerPage
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
		m.logger.Info("Loading player profile", map[string]interface{}{
			"nickname": nickname,
		})
		
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		profile, err := m.repo.GetPlayerByNickname(ctx, nickname)
		if err != nil {
			m.logger.Error("Failed to load player profile", map[string]interface{}{
				"nickname": nickname,
				"error":    err.Error(),
			})
			return errorMsg{err: err.Error()}
		}
		
		m.logger.Info("Player profile loaded successfully", map[string]interface{}{
			"nickname": nickname,
			"player_id": profile.ID,
		})
		return profileLoadedMsg{profile: *profile}
	}
}

// loadRecentMatches loads recent matches asynchronously
func (m AppModel) loadRecentMatches() tea.Cmd {
	return func() tea.Msg {
		// Load a smaller batch first for quick display
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// Load initial batch (first 20 matches for quick display)
		initialLimit := 20
		if initialLimit > m.config.MaxMatchesToLoad {
			initialLimit = m.config.MaxMatchesToLoad
		}
		
		matches, err := m.repo.GetPlayerRecentMatches(ctx, m.player.ID, "cs2", initialLimit)
		if err != nil {
			return errorMsg{err: err.Error()}
		}
		return matchesLoadedMsg{matches: matches}
	}
}


// loadBackgroundMatches loads matches in the background for better UX
func (m AppModel) loadBackgroundMatches() tea.Cmd {
	return func() tea.Msg {
		// Use a longer timeout for background loading
		ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
		defer cancel()

		// Load more matches in the background
		matches, err := m.repo.GetPlayerRecentMatches(ctx, m.player.ID, "cs2", m.config.MaxMatchesToLoad)
		if err != nil {
			// Don't return error for background loading, just return empty matches
			return backgroundMatchesLoadedMsg{matches: []entity.PlayerMatchSummary{}}
		}
		
		return backgroundMatchesLoadedMsg{matches: matches}
	}
}

// loadStatistics loads and calculates statistics from recent matches
func (m AppModel) loadStatistics() tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		matches, err := m.repo.GetPlayerRecentMatches(ctx, m.player.ID, "cs2", 20)
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

	// Create detailed match statistics
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
			ADR:                 m.calculateADR(baseMatch),
			HLTVRating:          m.calculateHLTVRating(baseMatch),
			FirstKills:          m.calculateFirstKills(baseMatch),
			FirstDeaths:         m.calculateFirstDeaths(baseMatch),
			ClutchWins:          m.calculateClutchWins(baseMatch),
			EntryFrags:          m.calculateEntryFrags(baseMatch),
			FlashAssists:        m.calculateFlashAssists(baseMatch),
			UtilityDamage:       m.calculateUtilityDamage(baseMatch),
		},
		TeamStats: TeamStats{
			PlayerTeamScore: m.extractPlayerTeamScore(baseMatch.Score),
			EnemyTeamScore:  m.extractEnemyTeamScore(baseMatch.Score),
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

// Helper functions for extracting team scores
func (m AppModel) extractPlayerTeamScore(score string) int {
	// Parse score like "16-14" and return the first number
	parts := strings.Split(score, "-")
	if len(parts) >= 2 {
		if playerScore, err := strconv.Atoi(strings.TrimSpace(parts[0])); err == nil {
			return playerScore
		}
	}
	return 0
}

func (m AppModel) extractEnemyTeamScore(score string) int {
	// Parse score like "16-14" and return the second number
	parts := strings.Split(score, "-")
	if len(parts) >= 2 {
		if enemyScore, err := strconv.Atoi(strings.TrimSpace(parts[1])); err == nil {
			return enemyScore
		}
	}
	return 0
}

// updateComparisonInput handles key events in the comparison input state
func (m AppModel) updateComparisonInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit
	case "esc":
		m.state = StateProfile
		m.comparisonInput = ""
		return m, nil
	case "enter":
		if strings.TrimSpace(m.comparisonInput) != "" {
			m.loading = true
			m.state = StateLoading
			return m, m.loadPlayerComparison(m.comparisonInput)
		}
	case "backspace":
		if len(m.comparisonInput) > 0 {
			m.comparisonInput = m.comparisonInput[:len(m.comparisonInput)-1]
		}
	default:
		if len(msg.String()) == 1 {
			m.comparisonInput += msg.String()
		}
	}
	return m, nil
}

// updateComparison handles key events in the comparison state
func (m AppModel) updateComparison(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit
	case "esc":
		m.state = StateProfile
		m.comparison = nil
		return m, nil
	}
	return m, nil
}

// loadPlayerComparison loads comparison data between current player and friend
func (m AppModel) loadPlayerComparison(friendNickname string) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		// Get friend's profile
		friendProfile, err := m.repo.GetPlayerByNickname(ctx, friendNickname)
		if err != nil {
			return errorMsg{err: fmt.Sprintf("Failed to load friend's profile: %v", err)}
		}

		// Get friend's recent matches
		friendMatches, err := m.repo.GetPlayerRecentMatches(ctx, friendProfile.ID, "cs2", m.config.ComparisonMatches)
		if err != nil {
			return errorMsg{err: fmt.Sprintf("Failed to load friend's matches: %v", err)}
		}

		// Get current player's recent matches for comparison (always load exactly the same number for fair comparison)
		currentMatches, err := m.repo.GetPlayerRecentMatches(ctx, m.player.ID, "cs2", m.config.ComparisonMatches)
		if err != nil {
			return errorMsg{err: fmt.Sprintf("Failed to load current player's matches: %v", err)}
		}

		// Calculate stats for both players
		currentStats := calculateStats(currentMatches)
		friendStats := calculateStats(friendMatches)

		// Create comparison data
		comparison := PlayerComparison{
			Player1Nickname: m.player.Nickname,
			Player2Nickname: friendProfile.Nickname,
			Player1Stats:    currentStats,
			Player2Stats:    friendStats,
			ComparisonData:  calculateComparisonData(currentStats, friendStats),
		}

		return comparisonLoadedMsg{comparison: comparison}
	}
}

// loadLifetimeStats loads lifetime statistics for the current player
func (m AppModel) loadLifetimeStats() tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Get lifetime stats from repository - try different game IDs
		var stats *entity.PlayerStats
		var err error
		
		// Try CS2 first
		stats, err = m.repo.GetPlayerStats(ctx, m.player.ID, "cs2")
		if err != nil || stats == nil || stats.Lifetime == nil || len(stats.Lifetime) == 0 {
			// Try CS:GO as fallback
			stats, err = m.repo.GetPlayerStats(ctx, m.player.ID, "csgo")
			if err != nil || stats == nil || stats.Lifetime == nil || len(stats.Lifetime) == 0 {
				// Try Counter-Strike 2
				stats, err = m.repo.GetPlayerStats(ctx, m.player.ID, "Counter-Strike 2")
			}
		}
		
		if err != nil {
			return errorMsg{err: fmt.Sprintf("Failed to load lifetime stats: %v", err)}
		}



		return lifetimeStatsLoadedMsg{stats: stats}
	}
}
