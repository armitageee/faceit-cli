package entity

import (
	"testing"
)

func TestPlayerMatchSummary(t *testing.T) {
	match := PlayerMatchSummary{
		MatchID:             "test-match-id",
		Map:                 "de_dust2",
		FinishedAt:          1640995200, // 2022-01-01 00:00:00 UTC
		Score:               "16-14",
		Kills:               20,
		Deaths:              15,
		Assists:             5,
		KDRatio:             1.33,
		HeadshotsPercentage: 60.0,
		ADR:                 75.5,
		Result:              "Win",
	}

	// Test basic fields
	if match.MatchID != "test-match-id" {
		t.Errorf("MatchID = %v, want %v", match.MatchID, "test-match-id")
	}
	if match.Map != "de_dust2" {
		t.Errorf("Map = %v, want %v", match.Map, "de_dust2")
	}
	if match.Score != "16-14" {
		t.Errorf("Score = %v, want %v", match.Score, "16-14")
	}
	if match.Kills != 20 {
		t.Errorf("Kills = %v, want %v", match.Kills, 20)
	}
	if match.Deaths != 15 {
		t.Errorf("Deaths = %v, want %v", match.Deaths, 15)
	}
	if match.Assists != 5 {
		t.Errorf("Assists = %v, want %v", match.Assists, 5)
	}
	if match.KDRatio != 1.33 {
		t.Errorf("KDRatio = %v, want %v", match.KDRatio, 1.33)
	}
	if match.HeadshotsPercentage != 60.0 {
		t.Errorf("HeadshotsPercentage = %v, want %v", match.HeadshotsPercentage, 60.0)
	}
	if match.ADR != 75.5 {
		t.Errorf("ADR = %v, want %v", match.ADR, 75.5)
	}
	if match.Result != "Win" {
		t.Errorf("Result = %v, want %v", match.Result, "Win")
	}
}

func TestPlayerProfile(t *testing.T) {
	profile := PlayerProfile{
		ID:        "test-player-id",
		Nickname:  "testplayer",
		Country:   "US",
		Avatar:    "https://example.com/avatar.jpg",
		FaceitURL: "https://www.faceit.com/en/players/testplayer",
		Games: map[string]GameDetail{
			"cs2": {
				Elo:        2500,
				SkillLevel: 8,
				Region:     "NA",
			},
		},
	}

	// Test basic fields
	if profile.ID != "test-player-id" {
		t.Errorf("ID = %v, want %v", profile.ID, "test-player-id")
	}
	if profile.Nickname != "testplayer" {
		t.Errorf("Nickname = %v, want %v", profile.Nickname, "testplayer")
	}
	if profile.Country != "US" {
		t.Errorf("Country = %v, want %v", profile.Country, "US")
	}

	// Test games
	cs2Game, exists := profile.Games["cs2"]
	if !exists {
		t.Error("CS2 game should exist in Games map")
	}
	if cs2Game.Elo != 2500 {
		t.Errorf("CS2 Elo = %v, want %v", cs2Game.Elo, 2500)
	}
	if cs2Game.SkillLevel != 8 {
		t.Errorf("CS2 SkillLevel = %v, want %v", cs2Game.SkillLevel, 8)
	}
	if cs2Game.Region != "NA" {
		t.Errorf("CS2 Region = %v, want %v", cs2Game.Region, "NA")
	}
}

func TestGameDetail(t *testing.T) {
	game := GameDetail{
		Elo:        2000,
		SkillLevel: 6,
		Region:     "EU",
	}

	if game.Elo != 2000 {
		t.Errorf("Elo = %v, want %v", game.Elo, 2000)
	}
	if game.SkillLevel != 6 {
		t.Errorf("SkillLevel = %v, want %v", game.SkillLevel, 6)
	}
	if game.Region != "EU" {
		t.Errorf("Region = %v, want %v", game.Region, "EU")
	}
}
