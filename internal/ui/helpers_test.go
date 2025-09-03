package ui

import (
	"faceit-cli/internal/entity"
	"testing"
)

func TestCalculateStats(t *testing.T) {
	tests := []struct {
		name     string
		matches  []entity.PlayerMatchSummary
		expected PlayerStatsSummary
	}{
		{
			name:     "empty matches",
			matches:  []entity.PlayerMatchSummary{},
			expected: PlayerStatsSummary{},
		},
		{
			name: "single match",
			matches: []entity.PlayerMatchSummary{
				{
					Kills:               20,
					Deaths:              15,
					Assists:             5,
					KDRatio:             1.33,
					HeadshotsPercentage: 60.0,
					Result:              "Win",
					Map:                 "de_dust2",
				},
			},
			expected: PlayerStatsSummary{
				TotalMatches:   1,
				Wins:           1,
				Losses:         0,
				WinRate:        100.0,
				TotalKills:     20,
				TotalDeaths:    15,
				TotalAssists:   5,
				AverageKDRatio: 1.33,
				AverageHS:      60.0,
				BestKDRatio:    1.33,
				WorstKDRatio:   1.33,
				MostPlayedMap:  "de_dust2",
				MapStats:       map[string]int{"de_dust2": 1},
				CurrentStreak:  1,
				StreakType:     "win",
				LongestWinStreak: 1,
				LongestLossStreak: 0,
			},
		},
		{
			name: "multiple matches with streaks",
			matches: []entity.PlayerMatchSummary{
				{Result: "Win", KDRatio: 1.5, HeadshotsPercentage: 70.0, Map: "de_dust2"},
				{Result: "Win", KDRatio: 1.2, HeadshotsPercentage: 60.0, Map: "de_dust2"},
				{Result: "Loss", KDRatio: 0.8, HeadshotsPercentage: 50.0, Map: "de_inferno"},
				{Result: "Loss", KDRatio: 0.9, HeadshotsPercentage: 55.0, Map: "de_inferno"},
				{Result: "Win", KDRatio: 1.1, HeadshotsPercentage: 65.0, Map: "de_mirage"},
			},
			expected: PlayerStatsSummary{
				TotalMatches:   5,
				Wins:           3,
				Losses:         2,
				WinRate:        60.0,
				AverageKDRatio: 1.1,
				AverageHS:      60.0,
				BestKDRatio:    1.5,
				WorstKDRatio:   0.8,
				MostPlayedMap:  "de_dust2",
				MapStats:       map[string]int{"de_dust2": 2, "de_inferno": 2, "de_mirage": 1},
				CurrentStreak:  2, // The last two matches are Win, Win (from newest to oldest)
				StreakType:     "win",
				LongestWinStreak: 2,
				LongestLossStreak: 2,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calculateStats(tt.matches)
			
			if result.TotalMatches != tt.expected.TotalMatches {
				t.Errorf("TotalMatches = %v, want %v", result.TotalMatches, tt.expected.TotalMatches)
			}
			if result.Wins != tt.expected.Wins {
				t.Errorf("Wins = %v, want %v", result.Wins, tt.expected.Wins)
			}
			if result.Losses != tt.expected.Losses {
				t.Errorf("Losses = %v, want %v", result.Losses, tt.expected.Losses)
			}
			if result.WinRate != tt.expected.WinRate {
				t.Errorf("WinRate = %v, want %v", result.WinRate, tt.expected.WinRate)
			}
			if result.BestKDRatio != tt.expected.BestKDRatio {
				t.Errorf("BestKDRatio = %v, want %v", result.BestKDRatio, tt.expected.BestKDRatio)
			}
			if result.WorstKDRatio != tt.expected.WorstKDRatio {
				t.Errorf("WorstKDRatio = %v, want %v", result.WorstKDRatio, tt.expected.WorstKDRatio)
			}
			if result.MostPlayedMap != tt.expected.MostPlayedMap {
				t.Errorf("MostPlayedMap = %v, want %v", result.MostPlayedMap, tt.expected.MostPlayedMap)
			}
			if result.CurrentStreak != tt.expected.CurrentStreak {
				t.Errorf("CurrentStreak = %v, want %v", result.CurrentStreak, tt.expected.CurrentStreak)
			}
			if result.StreakType != tt.expected.StreakType {
				t.Errorf("StreakType = %v, want %v", result.StreakType, tt.expected.StreakType)
			}
		})
	}
}

func TestCalculateStreaks(t *testing.T) {
	tests := []struct {
		name                string
		matches             []entity.PlayerMatchSummary
		expectedCurrent     int
		expectedType        string
		expectedLongestWin  int
		expectedLongestLoss int
	}{
		{
			name:                "empty matches",
			matches:             []entity.PlayerMatchSummary{},
			expectedCurrent:     0,
			expectedType:        "",
			expectedLongestWin:  0,
			expectedLongestLoss: 0,
		},
		{
			name: "win streak",
			matches: []entity.PlayerMatchSummary{
				{Result: "Win"},
				{Result: "Win"},
				{Result: "Win"},
			},
			expectedCurrent:     3,
			expectedType:        "win",
			expectedLongestWin:  3,
			expectedLongestLoss: 0,
		},
		{
			name: "loss streak",
			matches: []entity.PlayerMatchSummary{
				{Result: "Loss"},
				{Result: "Loss"},
			},
			expectedCurrent:     -2,
			expectedType:        "loss",
			expectedLongestWin:  0,
			expectedLongestLoss: 2,
		},
		{
			name: "mixed results",
			matches: []entity.PlayerMatchSummary{
				{Result: "Win"},
				{Result: "Win"},
				{Result: "Loss"},
				{Result: "Loss"},
				{Result: "Loss"},
				{Result: "Win"},
			},
			expectedCurrent:     2, // The first two matches are Win, Win
			expectedType:        "win",
			expectedLongestWin:  2,
			expectedLongestLoss: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			current, streakType, longestWin, longestLoss := calculateStreaks(tt.matches)
			
			if current != tt.expectedCurrent {
				t.Errorf("Current streak = %v, want %v", current, tt.expectedCurrent)
			}
			if streakType != tt.expectedType {
				t.Errorf("Streak type = %v, want %v", streakType, tt.expectedType)
			}
			if longestWin != tt.expectedLongestWin {
				t.Errorf("Longest win streak = %v, want %v", longestWin, tt.expectedLongestWin)
			}
			if longestLoss != tt.expectedLongestLoss {
				t.Errorf("Longest loss streak = %v, want %v", longestLoss, tt.expectedLongestLoss)
			}
		})
	}
}

func TestGenerateStreakInfo(t *testing.T) {
	tests := []struct {
		name     string
		stats    *PlayerStatsSummary
		expected string
	}{
		{
			name:     "nil stats",
			stats:    nil,
			expected: "",
		},
		{
			name: "win streak",
			stats: &PlayerStatsSummary{
				CurrentStreak:     3,
				StreakType:        "win",
				LongestWinStreak:  5,
				LongestLossStreak: 2,
			},
			expected: "🔥 Win Streak: 3\n🏆 Longest Win Streak: 5\n💔 Longest Loss Streak: 2",
		},
		{
			name: "loss streak",
			stats: &PlayerStatsSummary{
				CurrentStreak:     -2,
				StreakType:        "loss",
				LongestWinStreak:  3,
				LongestLossStreak: 4,
			},
			expected: "❄️  Loss Streak: 2\n🏆 Longest Win Streak: 3\n💔 Longest Loss Streak: 4",
		},
		{
			name: "no active streak",
			stats: &PlayerStatsSummary{
				CurrentStreak:     0,
				StreakType:        "",
				LongestWinStreak:  2,
				LongestLossStreak: 1,
			},
			expected: "📊 No active streak\n🏆 Longest Win Streak: 2\n💔 Longest Loss Streak: 1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := generateStreakInfo(tt.stats)
			if result != tt.expected {
				t.Errorf("generateStreakInfo() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGenerateASCIILogo(t *testing.T) {
	logo := generateASCIILogo()
	
	// Check that logo is not empty
	if logo == "" {
		t.Error("generateASCIILogo() returned empty string")
	}
	
	// Check that logo contains expected elements (Unicode box drawing characters)
	if !contains(logo, "╔") {
		t.Error("Logo should contain box drawing characters")
	}
	
	if !contains(logo, "║") {
		t.Error("Logo should contain box drawing characters")
	}
	
	if !contains(logo, "╚") {
		t.Error("Logo should contain box drawing characters")
	}
	
	// Check that logo has reasonable length (not too short, not too long)
	if len(logo) < 100 {
		t.Error("Logo seems too short")
	}
	if len(logo) > 2000 {
		t.Error("Logo seems too long")
	}
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 || 
		(len(s) > len(substr) && (s[:len(substr)] == substr || 
		s[len(s)-len(substr):] == substr || 
		contains(s[1:], substr))))
}


