package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

// viewSearch renders the search screen
func (m AppModel) viewSearch() string {
	asciiTitle := generateASCIILogo()
	title := titleStyle.Render("🎮 FACEIT CLI")
	search := searchStyle.Render(fmt.Sprintf("Enter player nickname:\n\n%s", m.searchInput))
	help := helpStyle.Render("Press Enter to search • Ctrl+C or Q to quit")

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center,
		lipgloss.JoinVertical(lipgloss.Center, asciiTitle, title, search, help))
}

// viewProfile renders the profile screen
func (m AppModel) viewProfile() string {
	if m.player == nil {
		return "No profile data"
	}

	asciiTitle := generateASCIILogo()
	title := titleStyle.Render("👤 " + m.player.Nickname)
	
	var content strings.Builder
	content.WriteString(fmt.Sprintf("Country: %s\n", m.player.Country))
	content.WriteString(fmt.Sprintf("ID: %s\n", m.player.ID))
	
	if cs2, ok := m.player.Games["cs2"]; ok {
		content.WriteString("\n🎯 CS2 Stats:\n")
		content.WriteString(fmt.Sprintf("  ELO: %d\n", cs2.Elo))
		content.WriteString(fmt.Sprintf("  Skill Level: %d\n", cs2.SkillLevel))
		content.WriteString(fmt.Sprintf("  Region: %s\n", cs2.Region))
	}

	profile := profileStyle.Render(content.String())
	help := helpStyle.Render("M - Recent matches • S - Statistics (20 matches) • C - Compare with friend • P - Switch player • Esc - Back to search • Ctrl+C or Q to quit")

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center,
		lipgloss.JoinVertical(lipgloss.Center, asciiTitle, title, profile, help))
}

// viewMatches renders the matches screen
func (m AppModel) viewMatches() string {
	if len(m.matches) == 0 {
		return "No matches found"
	}

	asciiTitle := generateASCIILogo()
	title := titleStyle.Render("🏆 Recent Matches - " + m.player.Nickname)
	
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
		lipgloss.JoinVertical(lipgloss.Center, asciiTitle, title, matches, help))
}

// viewStats renders the statistics screen
func (m AppModel) viewStats() string {
	if m.stats == nil {
		return "No statistics data"
	}

	asciiTitle := generateASCIILogo()
	title := titleStyle.Render("📊 Statistics (20 matches) - " + m.player.Nickname)
	
	// Left side - Statistics
	var statsContent strings.Builder
	statsContent.WriteString("📈 Overall Performance:\n")
	statsContent.WriteString(fmt.Sprintf("  Matches: %d | Wins: %d | Losses: %d\n", 
		m.stats.TotalMatches, m.stats.Wins, m.stats.Losses))
	statsContent.WriteString(fmt.Sprintf("  Win Rate: %.1f%%\n\n", m.stats.WinRate))
	
	statsContent.WriteString("🎯 Combat Statistics:\n")
	statsContent.WriteString(fmt.Sprintf("  Total K/D/A: %d/%d/%d\n", 
		m.stats.TotalKills, m.stats.TotalDeaths, m.stats.TotalAssists))
	statsContent.WriteString(fmt.Sprintf("  Average K/D: %.2f\n", m.stats.AverageKDRatio))
	statsContent.WriteString(fmt.Sprintf("  Best K/D: %.2f | Worst K/D: %.2f\n", 
		m.stats.BestKDRatio, m.stats.WorstKDRatio))
	statsContent.WriteString(fmt.Sprintf("  Average HS%%: %.1f%%\n\n", m.stats.AverageHS))
	
	statsContent.WriteString("🗺️  Map Statistics:\n")
	statsContent.WriteString(fmt.Sprintf("  Most Played: %s\n", m.stats.MostPlayedMap))
	statsContent.WriteString("  Map Breakdown:\n")
	
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
		statsContent.WriteString(fmt.Sprintf("    %s: %d matches\n", mapStat.name, mapStat.count))
	}

	// Right side - Streak information
	streakInfo := generateStreakInfo(m.stats)
	
	// Create styled boxes
	statsBox := statsStyle.Render(statsContent.String())
	streakBox := statsStyle.Render(streakInfo)
	
	// Combine stats and streak info side by side
	combinedContent := lipgloss.JoinHorizontal(lipgloss.Top, statsBox, "  ", streakBox)
	
	help := helpStyle.Render("Esc - Back to profile • Ctrl+C or Q to quit")

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center,
		lipgloss.JoinVertical(lipgloss.Center, asciiTitle, title, combinedContent, help))
}

// viewMatchDetail renders the detailed match statistics screen
func (m AppModel) viewMatchDetail() string {
	if m.matchDetail == nil {
		return "No match detail data"
	}

	asciiTitle := generateASCIILogo()
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
		m.matchDetail.PlayerStats.HLTVRating))
	
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
		m.matchDetail.PerformanceMetrics.ConsistencyScore,
		m.matchDetail.PerformanceMetrics.ImpactScore))
	content.WriteString(fmt.Sprintf("  Clutch: %.1f%% | Entry: %.1f%% | Support: %.1f%%\n\n", 
		m.matchDetail.PerformanceMetrics.ClutchScore,
		m.matchDetail.PerformanceMetrics.EntryScore,
		m.matchDetail.PerformanceMetrics.SupportScore))
	
	// Team statistics
	content.WriteString("👥 Team Statistics:\n")
	content.WriteString(fmt.Sprintf("  Your Team: %d | Enemy Team: %d\n", 
		m.matchDetail.TeamStats.PlayerTeamScore,
		m.matchDetail.TeamStats.EnemyTeamScore))

	matchDetail := matchDetailStyle.Render(content.String())
	help := helpStyle.Render("Esc - Back to matches • Ctrl+C or Q to quit")

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center,
		lipgloss.JoinVertical(lipgloss.Center, asciiTitle, title, matchDetail, help))
}

// viewLoading renders the loading screen
func (m AppModel) viewLoading() string {
	loading := loadingStyle.Render("⏳ Loading...")
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, loading)
}

// viewPlayerSwitch renders the player switch screen
func (m AppModel) viewPlayerSwitch() string {
	asciiTitle := generateASCIILogo()
	title := titleStyle.Render("🔄 Switch Player")
	
	var content strings.Builder
	content.WriteString("Enter player nickname:\n\n")
	content.WriteString(fmt.Sprintf("> %s", m.playerSwitchInput))
	
	// Show recent players if any
	if len(m.recentPlayers) > 0 {
		content.WriteString("\n\nRecent players:\n")
		for i, player := range m.recentPlayers {
			content.WriteString(fmt.Sprintf("  %d. %s\n", i+1, player))
		}
	}
	
	playerSwitch := profileStyle.Render(content.String())
	help := helpStyle.Render("Enter - Switch player • Esc - Back to profile • Ctrl+C or Q to quit")
	
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center,
		lipgloss.JoinVertical(lipgloss.Center, asciiTitle, title, playerSwitch, help))
}

// viewComparison renders the player comparison screen
func (m AppModel) viewComparison() string {
	if m.comparison == nil {
		return "No comparison data"
	}

	asciiTitle := generateASCIILogo()
	title := titleStyle.Render("⚔️ Player Comparison")
	
	var content strings.Builder
	
	// Player names
	content.WriteString(fmt.Sprintf("%s vs %s\n\n", 
		player1Style.Render(m.comparison.Player1Nickname),
		player2Style.Render(m.comparison.Player2Nickname)))
	
	// Basic stats comparison
	content.WriteString("📊 Basic Statistics:\n")
	content.WriteString(fmt.Sprintf("  K/D Ratio: %.2f vs %.2f (%s)\n", 
		m.comparison.Player1Stats.AverageKDRatio,
		m.comparison.Player2Stats.AverageKDRatio,
		formatComparisonValue(m.comparison.ComparisonData.KDRatioDiff, m.comparison.ComparisonData.KDRatioDiff > 0)))
	
	content.WriteString(fmt.Sprintf("  Win Rate: %.1f%% vs %.1f%% (%s%%)\n", 
		m.comparison.Player1Stats.WinRate,
		m.comparison.Player2Stats.WinRate,
		formatComparisonValue(m.comparison.ComparisonData.WinRateDiff, m.comparison.ComparisonData.WinRateDiff > 0)))
	
	content.WriteString(fmt.Sprintf("  Headshots: %.1f%% vs %.1f%% (%s%%)\n", 
		m.comparison.Player1Stats.AverageHS,
		m.comparison.Player2Stats.AverageHS,
		formatComparisonValue(m.comparison.ComparisonData.AverageHSDiff, m.comparison.ComparisonData.AverageHSDiff > 0)))
	
	content.WriteString("\n🎯 Kills & Deaths:\n")
	content.WriteString(fmt.Sprintf("  Total Kills: %d vs %d (%s)\n", 
		m.comparison.Player1Stats.TotalKills,
		m.comparison.Player2Stats.TotalKills,
		formatComparisonInt(m.comparison.ComparisonData.TotalKillsDiff, m.comparison.ComparisonData.TotalKillsDiff > 0)))
	
	content.WriteString(fmt.Sprintf("  Total Deaths: %d vs %d (%s)\n", 
		m.comparison.Player1Stats.TotalDeaths,
		m.comparison.Player2Stats.TotalDeaths,
		formatComparisonInt(m.comparison.ComparisonData.TotalDeathsDiff, m.comparison.ComparisonData.TotalDeathsDiff < 0)))
	
	content.WriteString(fmt.Sprintf("  Total Assists: %d vs %d (%s)\n", 
		m.comparison.Player1Stats.TotalAssists,
		m.comparison.Player2Stats.TotalAssists,
		formatComparisonInt(m.comparison.ComparisonData.TotalAssistsDiff, m.comparison.ComparisonData.TotalAssistsDiff > 0)))
	
	content.WriteString("\n🏆 Performance:\n")
	content.WriteString(fmt.Sprintf("  Best K/D: %.2f vs %.2f (%s)\n", 
		m.comparison.Player1Stats.BestKDRatio,
		m.comparison.Player2Stats.BestKDRatio,
		formatComparisonValue(m.comparison.ComparisonData.BestKDDiff, m.comparison.ComparisonData.BestKDDiff > 0)))
	
	content.WriteString(fmt.Sprintf("  Worst K/D: %.2f vs %.2f (%s)\n", 
		m.comparison.Player1Stats.WorstKDRatio,
		m.comparison.Player2Stats.WorstKDRatio,
		formatComparisonValue(m.comparison.ComparisonData.WorstKDDiff, m.comparison.ComparisonData.WorstKDDiff > 0)))
	
	content.WriteString("\n🗺️ Maps:\n")
	content.WriteString(fmt.Sprintf("  Most Played Together: %s\n", m.comparison.ComparisonData.MostPlayedMap))
	content.WriteString(fmt.Sprintf("  Common Maps: %d\n", len(m.comparison.ComparisonData.CommonMaps)))
	
	comparison := comparisonStyle.Render(content.String())
	help := helpStyle.Render("Esc - Back to profile • Ctrl+C or Q to quit")
	
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center,
		lipgloss.JoinVertical(lipgloss.Center, asciiTitle, title, comparison, help))
}

// viewComparisonInput renders the comparison input screen
func (m AppModel) viewComparisonInput() string {
	asciiTitle := generateASCIILogo()
	title := titleStyle.Render("⚔️ Compare with Friend")
	search := searchStyle.Render(fmt.Sprintf("Enter friend's nickname:\n\n%s", m.comparisonInput))
	help := helpStyle.Render("Press Enter to compare • Esc - Back to profile • Ctrl+C or Q to quit")

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center,
		lipgloss.JoinVertical(lipgloss.Center, asciiTitle, title, search, help))
}

// viewError renders the error screen
func (m AppModel) viewError() string {
	error := errorStyle.Render(fmt.Sprintf("❌ Error: %s", m.error))
	help := helpStyle.Render("Esc or Enter - Back to search • Ctrl+C or Q to quit")
	
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center,
		lipgloss.JoinVertical(lipgloss.Center, error, help))
}
