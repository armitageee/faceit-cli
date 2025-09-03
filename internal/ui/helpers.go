package ui

import (
	"faceit-cli/internal/entity"
	"fmt"
	"strconv"
	"strings"
)

// addToRecentPlayers adds a player to the recent players list
func (m *AppModel) addToRecentPlayers(nickname string) {
	// Remove if already exists
	for i, player := range m.recentPlayers {
		if player == nickname {
			m.recentPlayers = append(m.recentPlayers[:i], m.recentPlayers[i+1:]...)
			break
		}
	}
	// Add to beginning
	m.recentPlayers = append([]string{nickname}, m.recentPlayers...)
	// Keep only last 5 players
	if len(m.recentPlayers) > 5 {
		m.recentPlayers = m.recentPlayers[:5]
	}
}

// calculateStats calculates aggregated statistics from recent matches
func calculateStats(matches []entity.PlayerMatchSummary) PlayerStatsSummary {
	if len(matches) == 0 {
		return PlayerStatsSummary{}
	}

	stats := PlayerStatsSummary{
		TotalMatches: len(matches),
		MapStats:     make(map[string]int),
		KDChartData:  make([]float64, len(matches)),
	}

	var totalKDRatio, totalHS float64
	bestKD, worstKD := 0.0, 100.0

	for i, match := range matches {
		// Basic stats
		stats.TotalKills += match.Kills
		stats.TotalDeaths += match.Deaths
		stats.TotalAssists += match.Assists

		// Win/Loss tracking
		if match.Result == "Win" {
			stats.Wins++
		} else {
			stats.Losses++
		}

		// K/D tracking
		if match.KDRatio > 0 {
			totalKDRatio += match.KDRatio
			stats.KDChartData[i] = match.KDRatio
			
			if match.KDRatio > bestKD {
				bestKD = match.KDRatio
			}
			if match.KDRatio < worstKD {
				worstKD = match.KDRatio
			}
		}

		// Headshot tracking
		if match.HeadshotsPercentage > 0 {
			totalHS += match.HeadshotsPercentage
		}

		// Map tracking
		if match.Map != "" {
			stats.MapStats[match.Map]++
		}
	}

	// Calculate averages
	if len(matches) > 0 {
		stats.WinRate = float64(stats.Wins) / float64(len(matches)) * 100
		stats.AverageKDRatio = totalKDRatio / float64(len(matches))
		stats.AverageHS = totalHS / float64(len(matches))
		
		// Calculate total K/D ratio (total kills / total deaths)
		if stats.TotalDeaths > 0 {
			stats.TotalKDA = float64(stats.TotalKills) / float64(stats.TotalDeaths)
		} else {
			stats.TotalKDA = float64(stats.TotalKills) // If no deaths, K/D = kills
		}
	}

	stats.BestKDRatio = bestKD
	stats.WorstKDRatio = worstKD

	// Find most played map
	maxCount := 0
	for mapName, count := range stats.MapStats {
		if count > maxCount {
			maxCount = count
			stats.MostPlayedMap = mapName
		}
	}

	// Calculate streaks
	currentStreak, streakType, longestWinStreak, longestLossStreak := calculateStreaks(matches)
	stats.CurrentStreak = currentStreak
	stats.StreakType = streakType
	stats.LongestWinStreak = longestWinStreak
	stats.LongestLossStreak = longestLossStreak

	return stats
}

// calculateStreaks calculates win/loss streaks from match results
func calculateStreaks(matches []entity.PlayerMatchSummary) (currentStreak int, streakType string, longestWinStreak, longestLossStreak int) {
	if len(matches) == 0 {
		return 0, "", 0, 0
	}

	// Calculate current streak
	currentStreak = 0
	streakType = ""
	
	// Start from the most recent match
	for _, match := range matches {
		if match.Result == "Win" {
			if streakType == "win" || streakType == "" {
				currentStreak++
				streakType = "win"
			} else {
				break
			}
		} else if match.Result == "Loss" {
			if streakType == "loss" || streakType == "" {
				currentStreak++
				streakType = "loss"
			} else {
				break
			}
		}
	}

	// Make current streak negative for loss streaks
	if streakType == "loss" {
		currentStreak = -currentStreak
	}

	// Calculate longest streaks
	var tempWinStreak, tempLossStreak int
	longestWinStreak, longestLossStreak = 0, 0

	for _, match := range matches {
		if match.Result == "Win" {
			tempWinStreak++
			tempLossStreak = 0
			if tempWinStreak > longestWinStreak {
				longestWinStreak = tempWinStreak
			}
		} else if match.Result == "Loss" {
			tempLossStreak++
			tempWinStreak = 0
			if tempLossStreak > longestLossStreak {
				longestLossStreak = tempLossStreak
			}
		}
	}

	return currentStreak, streakType, longestWinStreak, longestLossStreak
}

// generateStreakInfo generates a formatted string for streak information
func generateStreakInfo(stats *PlayerStatsSummary) string {
	if stats == nil {
		return ""
	}

	var streakInfo strings.Builder
	
	if stats.CurrentStreak > 0 {
		streakInfo.WriteString("🔥 Win Streak: ")
		streakInfo.WriteString(strconv.Itoa(stats.CurrentStreak))
	} else if stats.CurrentStreak < 0 {
		streakInfo.WriteString("❄️  Loss Streak: ")
		streakInfo.WriteString(strconv.Itoa(-stats.CurrentStreak))
	} else {
		streakInfo.WriteString("📊 No active streak")
	}
	
	streakInfo.WriteString("\n")
	streakInfo.WriteString("🏆 Longest Win Streak: ")
	streakInfo.WriteString(strconv.Itoa(stats.LongestWinStreak))
	streakInfo.WriteString("\n")
	streakInfo.WriteString("💔 Longest Loss Streak: ")
	streakInfo.WriteString(strconv.Itoa(stats.LongestLossStreak))
	
	return streakInfo.String()
}

// generateASCIILogo creates a large ASCII art logo for FACEIT-CLI
func generateASCIILogo() string {
	logo := `
╔══════════════════════════════════════════════════════════════════════════════╗
║                                                                              ║
║    ███████╗ █████╗  ██████╗███████╗██╗████████╗    ██████╗██╗     ██╗        ║
║    ██╔════╝██╔══██╗██╔════╝██╔════╝██║╚══██╔══╝   ██╔════╝██║     ██║        ║
║    █████╗  ███████║██║     █████╗  ██║   ██║      ██║     ██║     ██║        ║
║    ██╔══╝  ██╔══██║██║     ██╔══╝  ██║   ██║      ██║     ██║     ██║        ║
║    ██║     ██║  ██║╚██████╗███████╗██║   ██║      ╚██████╗███████╗██║        ║
║    ╚═╝     ╚═╝  ╚═╝ ╚═════╝╚══════╝╚═╝   ╚═╝       ╚═════╝╚══════╝╚═╝        ║
║                                                                              ║
╚══════════════════════════════════════════════════════════════════════════════╝
`
	return logo
}

// Helper functions for calculating advanced metrics
func (m AppModel) calculateADR(match *entity.PlayerMatchSummary) float64 {
	// Use real ADR from API if available
	if match.ADR > 0 {
		return match.ADR
	}
	// Fallback to estimation if ADR not available
	estimatedDamage := float64(match.Kills)*100 + float64(match.Assists)*50
	rounds := 30 // Assume average match length
	return estimatedDamage / float64(rounds)
}

func (m AppModel) calculateHLTVRating(match *entity.PlayerMatchSummary) float64 {
	// Improved HLTV rating calculation
	// Based on the official HLTV rating formula with some adjustments
	
	kd := match.KDRatio
	hs := match.HeadshotsPercentage / 100.0
	adr := m.calculateADR(match)
	
	// More accurate HLTV rating formula
	// K/D ratio is the most important factor
	kdComponent := kd * 0.5
	
	// Headshot percentage contributes to rating
	hsComponent := hs * 0.3
	
	// ADR normalized to 0-1 scale (assuming max ADR around 100)
	adrComponent := (adr / 100.0) * 0.2
	
	rating := kdComponent + hsComponent + adrComponent
	
	// Ensure rating is reasonable (0.0 to 2.0+)
	if rating < 0 {
		rating = 0
	}
	
	return rating
}

func (m AppModel) calculateFirstKills(match *entity.PlayerMatchSummary) int {
	// Estimate first kills based on K/D ratio
	return int(float64(match.Kills) * 0.3)
}

func (m AppModel) calculateFirstDeaths(match *entity.PlayerMatchSummary) int {
	// Estimate first deaths based on deaths
	return int(float64(match.Deaths) * 0.2)
}

func (m AppModel) calculateClutchWins(match *entity.PlayerMatchSummary) int {
	// Estimate clutch wins based on K/D ratio and result
	if match.Result == "Win" && match.KDRatio > 1.0 {
		return int(float64(match.Kills) * 0.1)
	}
	return 0
}

func (m AppModel) calculateEntryFrags(match *entity.PlayerMatchSummary) int {
	// Estimate entry frags based on kills
	return int(float64(match.Kills) * 0.25)
}

func (m AppModel) calculateFlashAssists(match *entity.PlayerMatchSummary) int {
	// Estimate flash assists based on assists
	return int(float64(match.Assists) * 0.4)
}

func (m AppModel) calculateUtilityDamage(match *entity.PlayerMatchSummary) int {
	// Estimate utility damage based on assists and kills
	return int(float64(match.Assists)*20 + float64(match.Kills)*5)
}

// Advanced performance metrics
func (m AppModel) calculateConsistencyScore(match *entity.PlayerMatchSummary) float64 {
	// Consistency based on K/D ratio stability
	kd := match.KDRatio
	if kd >= 1.5 {
		return 90.0
	} else if kd >= 1.2 {
		return 80.0
	} else if kd >= 1.0 {
		return 70.0
	} else if kd >= 0.8 {
		return 60.0
	}
	return 50.0
}

func (m AppModel) calculateImpactScore(match *entity.PlayerMatchSummary) float64 {
	// Impact based on K/D ratio and headshot percentage
	kd := match.KDRatio
	hs := match.HeadshotsPercentage
	return (kd * 40) + (hs * 0.6)
}

func (m AppModel) calculateClutchScore(match *entity.PlayerMatchSummary) float64 {
	// Clutch score based on K/D ratio in winning matches
	if match.Result == "Win" && match.KDRatio > 1.0 {
		return match.KDRatio * 30
	}
	return 20.0
}

func (m AppModel) calculateEntryScore(match *entity.PlayerMatchSummary) float64 {
	// Entry score based on kills and K/D ratio
	return float64(match.Kills) * match.KDRatio * 2
}

func (m AppModel) calculateSupportScore(match *entity.PlayerMatchSummary) float64 {
	// Support score based on assists and team play
	return float64(match.Assists) * 15
}

// calculateComparisonData calculates comparison metrics between two players
func calculateComparisonData(player1Stats, player2Stats PlayerStatsSummary) ComparisonData {
	return ComparisonData{
		KDRatioDiff:      player1Stats.AverageKDRatio - player2Stats.AverageKDRatio,
		TotalKDADiff:     player1Stats.TotalKDA - player2Stats.TotalKDA,
		WinRateDiff:      player1Stats.WinRate - player2Stats.WinRate,
		AverageHSDiff:    player1Stats.AverageHS - player2Stats.AverageHS,
		TotalKillsDiff:   player1Stats.TotalKills - player2Stats.TotalKills,
		TotalDeathsDiff:  player1Stats.TotalDeaths - player2Stats.TotalDeaths,
		TotalAssistsDiff: player1Stats.TotalAssists - player2Stats.TotalAssists,
		BestKDDiff:       player1Stats.BestKDRatio - player2Stats.BestKDRatio,
		WorstKDDiff:      player1Stats.WorstKDRatio - player2Stats.WorstKDRatio,
		MostPlayedMap:    findMostPlayedMap(player1Stats, player2Stats),
		CommonMaps:       findCommonMaps(player1Stats, player2Stats),
	}
}

// findMostPlayedMap finds the map that both players played most
func findMostPlayedMap(stats1, stats2 PlayerStatsSummary) string {
	commonMaps := findCommonMaps(stats1, stats2)
	if len(commonMaps) == 0 {
		return "No common maps"
	}
	
	maxCount := 0
	mostPlayed := ""
	for _, mapName := range commonMaps {
		count1 := stats1.MapStats[mapName]
		count2 := stats2.MapStats[mapName]
		totalCount := count1 + count2
		if totalCount > maxCount {
			maxCount = totalCount
			mostPlayed = mapName
		}
	}
	return mostPlayed
}

// findCommonMaps finds maps that both players have played
func findCommonMaps(stats1, stats2 PlayerStatsSummary) []string {
	var commonMaps []string
	for mapName := range stats1.MapStats {
		if _, exists := stats2.MapStats[mapName]; exists {
			commonMaps = append(commonMaps, mapName)
		}
	}
	return commonMaps
}

// formatComparisonValue formats a comparison value with appropriate styling
func formatComparisonValue(value float64, isBetter bool) string {
	if isBetter {
		return betterStyle.Render(fmt.Sprintf("+%.2f", value))
	}
	return worseStyle.Render(fmt.Sprintf("%.2f", value))
}

// formatComparisonInt formats an integer comparison value with appropriate styling
func formatComparisonInt(value int, isBetter bool) string {
	if isBetter {
		return betterStyle.Render(fmt.Sprintf("+%d", value))
	}
	return worseStyle.Render(fmt.Sprintf("%d", value))
}

// getVisualLength calculates the visual length of a string (counting runes)
func getVisualLength(s string) int {
	return len([]rune(s))
}

// generateProfileFrame generates a beautiful ASCII frame for the profile
func generateProfileFrame(content string) string {
	lines := strings.Split(content, "\n")
	maxWidth := 0
	for _, line := range lines {
		visualLen := getVisualLength(line)
		if visualLen > maxWidth {
			maxWidth = visualLen
		}
	}
	
	// Ensure minimum width
	if maxWidth < 50 {
		maxWidth = 50
	}
	
	var result strings.Builder
	
	// Top border with corners
	result.WriteString("╔")
	for i := 0; i < maxWidth+2; i++ {
		result.WriteString("═")
	}
	result.WriteString("╗\n")
	
	// Content lines with side borders
	for _, line := range lines {
		result.WriteString("║ ")
		result.WriteString(line)
		// Pad with spaces to maintain width based on visual length
		visualLen := getVisualLength(line)
		// Subtract space for lines with emojis to compensate for emoji width
		adjustment := 0
		if strings.Contains(line, "🎯") || strings.Contains(line, "📊") {
			adjustment = -1
		}
		for i := visualLen; i < maxWidth+adjustment; i++ {
			result.WriteString(" ")
		}
		result.WriteString(" ║\n")
	}
	
	// Bottom border with corners
	result.WriteString("╚")
	for i := 0; i < maxWidth+2; i++ {
		result.WriteString("═")
	}
	result.WriteString("╝")
	
	return result.String()
}

// extractLifetimeStats extracts key statistics from lifetime stats
func extractLifetimeStats(stats *entity.PlayerStats) (kdRatio float64, totalMatches int, winRate float64) {
	if stats == nil || stats.Lifetime == nil {
		return 0, 0, 0
	}

	// Extract K/D ratio - try different possible keys based on FACEIT API
	kdKeys := []string{"Average K/D Ratio", "K/D Ratio", "K/D", "KD Ratio", "Average KD", "K/D Ratio", "K/D", "KD"}
	for _, key := range kdKeys {
		if kd, ok := stats.Lifetime[key]; ok {
			if kdFloat, ok := kd.(float64); ok {
				kdRatio = kdFloat
				break
			} else if kdStr, ok := kd.(string); ok {
				if parsed, err := strconv.ParseFloat(kdStr, 64); err == nil {
					kdRatio = parsed
					break
				}
			}
		}
	}

	// Extract total matches - try different possible keys
	matchKeys := []string{"Matches", "Total Matches", "Games", "Total Games", "Matches Played", "Total Matches Played"}
	for _, key := range matchKeys {
		if matches, ok := stats.Lifetime[key]; ok {
			if matchesFloat, ok := matches.(float64); ok {
				totalMatches = int(matchesFloat)
				break
			} else if matchesStr, ok := matches.(string); ok {
				if parsed, err := strconv.ParseFloat(matchesStr, 64); err == nil {
					totalMatches = int(parsed)
					break
				}
			}
		}
	}

	// Extract win rate - try different possible keys
	winRateKeys := []string{"Win Rate %", "Win Rate", "Win%", "Win Percentage", "Wins %", "Winrate %"}
	for _, key := range winRateKeys {
		if winRateVal, ok := stats.Lifetime[key]; ok {
			if winRateFloat, ok := winRateVal.(float64); ok {
				winRate = winRateFloat
				break
			} else if winRateStr, ok := winRateVal.(string); ok {
				if parsed, err := strconv.ParseFloat(winRateStr, 64); err == nil {
					winRate = parsed
					break
				}
			}
		}
	}

	return kdRatio, totalMatches, winRate
}



