package entity

import (
	"testing"
)

func TestPlayerMatchSummaryStruct(t *testing.T) {
	tests := []struct {
		name     string
		match    PlayerMatchSummary
		expected PlayerMatchSummary
	}{
		{
			name: "complete match data",
			match: PlayerMatchSummary{
				MatchID:             "match-123",
				Map:                 "de_dust2",
				FinishedAt:          1640995200,
				Score:               "16-14",
				Kills:               25,
				Deaths:              20,
				Assists:             5,
				KDRatio:             1.25,
				HeadshotsPercentage: 45.5,
				ADR:                 85.2,
				Result:              "Win",
			},
			expected: PlayerMatchSummary{
				MatchID:             "match-123",
				Map:                 "de_dust2",
				FinishedAt:          1640995200,
				Score:               "16-14",
				Kills:               25,
				Deaths:              20,
				Assists:             5,
				KDRatio:             1.25,
				HeadshotsPercentage: 45.5,
				ADR:                 85.2,
				Result:              "Win",
			},
		},
		{
			name: "minimal match data",
			match: PlayerMatchSummary{
				MatchID: "match-456",
				Map:     "de_inferno",
				Result:  "Loss",
			},
			expected: PlayerMatchSummary{
				MatchID: "match-456",
				Map:     "de_inferno",
				Result:  "Loss",
			},
		},
		{
			name: "zero values",
			match: PlayerMatchSummary{
				MatchID:             "match-789",
				Map:                 "de_mirage",
				FinishedAt:          0,
				Score:               "",
				Kills:               0,
				Deaths:              0,
				Assists:             0,
				KDRatio:             0.0,
				HeadshotsPercentage: 0.0,
				ADR:                 0.0,
				Result:              "",
			},
			expected: PlayerMatchSummary{
				MatchID:             "match-789",
				Map:                 "de_mirage",
				FinishedAt:          0,
				Score:               "",
				Kills:               0,
				Deaths:              0,
				Assists:             0,
				KDRatio:             0.0,
				HeadshotsPercentage: 0.0,
				ADR:                 0.0,
				Result:              "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.match != tt.expected {
				t.Errorf("PlayerMatchSummary = %+v, want %+v", tt.match, tt.expected)
			}
		})
	}
}

func TestMatchStats(t *testing.T) {
	matchStats := MatchStats{
		MatchID:    "match-123",
		Map:        "de_dust2",
		FinishedAt: 1640995200,
		Score:      "16-14",
		Result:     "finished",
		Team1: TeamMatchStats{
			TeamID:   "team1",
			TeamName: "Team Alpha",
			Score:    16,
			Players: []PlayerMatchStats{
				{
					PlayerID:            "player1",
					Nickname:            "Player1",
					Team:                "team1",
					Kills:               25,
					Deaths:              20,
					Assists:             5,
					KDRatio:             1.25,
					HeadshotsPercentage: 45.5,
					ADR:                 85.2,
				},
			},
		},
		Team2: TeamMatchStats{
			TeamID:   "team2",
			TeamName: "Team Beta",
			Score:    14,
			Players: []PlayerMatchStats{
				{
					PlayerID:            "player2",
					Nickname:            "Player2",
					Team:                "team2",
					Kills:               20,
					Deaths:              25,
					Assists:             3,
					KDRatio:             0.8,
					HeadshotsPercentage: 35.0,
					ADR:                 70.5,
				},
			},
		},
		PlayerStats: []PlayerMatchStats{
			{
				PlayerID:            "player1",
				Nickname:            "Player1",
				Team:                "team1",
				Kills:               25,
				Deaths:              20,
				Assists:             5,
				KDRatio:             1.25,
				HeadshotsPercentage: 45.5,
				ADR:                 85.2,
			},
			{
				PlayerID:            "player2",
				Nickname:            "Player2",
				Team:                "team2",
				Kills:               20,
				Deaths:              25,
				Assists:             3,
				KDRatio:             0.8,
				HeadshotsPercentage: 35.0,
				ADR:                 70.5,
			},
		},
	}

	// Test basic fields
	if matchStats.MatchID != "match-123" {
		t.Errorf("MatchStats.MatchID = %s, want match-123", matchStats.MatchID)
	}

	if matchStats.Map != "de_dust2" {
		t.Errorf("MatchStats.Map = %s, want de_dust2", matchStats.Map)
	}

	if matchStats.Score != "16-14" {
		t.Errorf("MatchStats.Score = %s, want 16-14", matchStats.Score)
	}

	// Test team data
	if matchStats.Team1.TeamName != "Team Alpha" {
		t.Errorf("MatchStats.Team1.TeamName = %s, want Team Alpha", matchStats.Team1.TeamName)
	}

	if matchStats.Team2.TeamName != "Team Beta" {
		t.Errorf("MatchStats.Team2.TeamName = %s, want Team Beta", matchStats.Team2.TeamName)
	}

	// Test player data
	if len(matchStats.PlayerStats) != 2 {
		t.Errorf("MatchStats.PlayerStats length = %d, want 2", len(matchStats.PlayerStats))
	}

	if matchStats.PlayerStats[0].Nickname != "Player1" {
		t.Errorf("MatchStats.PlayerStats[0].Nickname = %s, want Player1", matchStats.PlayerStats[0].Nickname)
	}
}

func TestTeamMatchStats(t *testing.T) {
	teamStats := TeamMatchStats{
		TeamID:   "team1",
		TeamName: "Team Alpha",
		Score:    16,
		Players: []PlayerMatchStats{
			{
				PlayerID: "player1",
				Nickname: "Player1",
				Team:     "team1",
				Kills:    25,
				Deaths:   20,
				Assists:  5,
				KDRatio:  1.25,
			},
			{
				PlayerID: "player2",
				Nickname: "Player2",
				Team:     "team1",
				Kills:    20,
				Deaths:   18,
				Assists:  3,
				KDRatio:  1.11,
			},
		},
	}

	// Test basic fields
	if teamStats.TeamID != "team1" {
		t.Errorf("TeamMatchStats.TeamID = %s, want team1", teamStats.TeamID)
	}

	if teamStats.TeamName != "Team Alpha" {
		t.Errorf("TeamMatchStats.TeamName = %s, want Team Alpha", teamStats.TeamName)
	}

	if teamStats.Score != 16 {
		t.Errorf("TeamMatchStats.Score = %d, want 16", teamStats.Score)
	}

	// Test players
	if len(teamStats.Players) != 2 {
		t.Errorf("TeamMatchStats.Players length = %d, want 2", len(teamStats.Players))
	}

	if teamStats.Players[0].Nickname != "Player1" {
		t.Errorf("TeamMatchStats.Players[0].Nickname = %s, want Player1", teamStats.Players[0].Nickname)
	}
}

func TestPlayerMatchStats(t *testing.T) {
	playerStats := PlayerMatchStats{
		PlayerID:            "player1",
		Nickname:            "Player1",
		Team:                "team1",
		Kills:               25,
		Deaths:              20,
		Assists:             5,
		KDRatio:             1.25,
		HeadshotsPercentage: 45.5,
		ADR:                 85.2,
		HLTVRating:          1.15,
		FirstKills:          3,
		FirstDeaths:         2,
		ClutchWins:          1,
		EntryFrags:          4,
		FlashAssists:        2,
		UtilityDamage:       150,
	}

	// Test basic fields
	if playerStats.PlayerID != "player1" {
		t.Errorf("PlayerMatchStats.PlayerID = %s, want player1", playerStats.PlayerID)
	}

	if playerStats.Nickname != "Player1" {
		t.Errorf("PlayerMatchStats.Nickname = %s, want Player1", playerStats.Nickname)
	}

	if playerStats.Team != "team1" {
		t.Errorf("PlayerMatchStats.Team = %s, want team1", playerStats.Team)
	}

	// Test statistics
	if playerStats.Kills != 25 {
		t.Errorf("PlayerMatchStats.Kills = %d, want 25", playerStats.Kills)
	}

	if playerStats.Deaths != 20 {
		t.Errorf("PlayerMatchStats.Deaths = %d, want 20", playerStats.Deaths)
	}

	if playerStats.Assists != 5 {
		t.Errorf("PlayerMatchStats.Assists = %d, want 5", playerStats.Assists)
	}

	if playerStats.KDRatio != 1.25 {
		t.Errorf("PlayerMatchStats.KDRatio = %f, want 1.25", playerStats.KDRatio)
	}

	// Test advanced statistics
	if playerStats.HLTVRating != 1.15 {
		t.Errorf("PlayerMatchStats.HLTVRating = %f, want 1.15", playerStats.HLTVRating)
	}

	if playerStats.FirstKills != 3 {
		t.Errorf("PlayerMatchStats.FirstKills = %d, want 3", playerStats.FirstKills)
	}

	if playerStats.ClutchWins != 1 {
		t.Errorf("PlayerMatchStats.ClutchWins = %d, want 1", playerStats.ClutchWins)
	}
}

// Benchmark tests
func BenchmarkPlayerMatchSummary_Creation(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = PlayerMatchSummary{
			MatchID:             "match-123",
			Map:                 "de_dust2",
			FinishedAt:          1640995200,
			Score:               "16-14",
			Kills:               25,
			Deaths:              20,
			Assists:             5,
			KDRatio:             1.25,
			HeadshotsPercentage: 45.5,
			ADR:                 85.2,
			Result:              "Win",
		}
	}
}

func BenchmarkMatchStats_Creation(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = MatchStats{
			MatchID:    "match-123",
			Map:        "de_dust2",
			FinishedAt: 1640995200,
			Score:      "16-14",
			Result:     "finished",
			Team1: TeamMatchStats{
				TeamID:   "team1",
				TeamName: "Team Alpha",
				Score:    16,
			},
			Team2: TeamMatchStats{
				TeamID:   "team2",
				TeamName: "Team Beta",
				Score:    14,
			},
		}
	}
}
