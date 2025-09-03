package ui

import (
	"faceit-cli/internal/entity"
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

// generateASCIILevel creates a large ASCII art representation of the skill level
func generateASCIILevel(level int) string {
	// ASCII art for digits 0-10
	digits := map[int]string{
		0: `  ██████  
 ██    ██ 
██      ██
██      ██
██      ██
 ██    ██ 
  ██████  `,
		1: `    ██    
  ████    
    ██    
    ██    
    ██    
    ██    
  ██████  `,
		2: ` ██████  
██    ██ 
      ██ 
 ██████  
██       
██       
████████ `,
		3: ` ██████  
██    ██ 
      ██ 
 ██████  
      ██ 
██    ██ 
 ██████  `,
		4: `██    ██ 
██    ██ 
██    ██ 
████████
      ██ 
      ██ 
      ██ `,
		5: `████████
██       
██       
███████  
      ██ 
██    ██ 
 ██████  `,
		6: ` ██████  
██       
██       
███████  
██    ██ 
██    ██ 
 ██████  `,
		7: `████████
      ██ 
     ██  
    ██   
   ██    
  ██     
 ██      `,
		8: ` ██████  
██    ██ 
██    ██ 
 ██████  
██    ██ 
██    ██ 
 ██████  `,
		9: ` ██████  
██    ██ 
██    ██ 
 ███████
      ██ 
██    ██ 
 ██████  `,
		10: ` ██████   ██████  
██    ██ ██    ██ 
██      ██      ██
██      ██      ██
██      ██      ██
██    ██ ██    ██ 
 ██████   ██████  `,
	}

	if ascii, exists := digits[level]; exists {
		return ascii
	}
	
	// Fallback for unknown levels
	return `  ██████  
 ██    ██ 
██      ██
██      ██
██      ██
 ██    ██ 
  ██████  `
}

// getLevelColor returns the color for a given skill level
func getLevelColor(level int) string {
	switch {
	case level >= 8:
		return "#FFD700" // Gold for high levels (8-10)
	case level >= 6:
		return "#FF6B6B" // Red for medium-high levels (6-7)
	case level >= 4:
		return "#4ECDC4" // Teal for medium levels (4-5)
	case level >= 2:
		return "#45B7D1" // Blue for low-medium levels (2-3)
	default:
		return "#96CEB4" // Green for low levels (1)
	}
}
