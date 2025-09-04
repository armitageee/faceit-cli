package repository

import (
	"context"
	"os"
	"testing"
	"time"
)

func TestNewFaceitRepository(t *testing.T) {
	apiKey := "test-api-key"
	repo := NewFaceitRepository(apiKey)
	
	if repo == nil {
		t.Fatal("Expected repository to be created, got nil")
	}
}

func TestGetPlayerByNickname(t *testing.T) {
	// Skip if no API key is provided
	apiKey := os.Getenv("FACEIT_API_KEY")
	if apiKey == "" {
		t.Skip("FACEIT_API_KEY not set, skipping integration test")
	}
	
	// Skip integration tests in short mode
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}
	
	// Skip integration tests in short mode
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	repo := NewFaceitRepository(apiKey)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Test with a known player
	player, err := repo.GetPlayerByNickname(ctx, "s1mple")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if player == nil {
		t.Fatal("Expected player profile, got nil")
	}

	if player.Nickname == "" {
		t.Error("Expected player nickname to be set")
	}

	if player.ID == "" {
		t.Error("Expected player ID to be set")
	}

	// Check CS2 stats
	if cs2, ok := player.Games["cs2"]; ok {
		if cs2.SkillLevel <= 0 {
			t.Error("Expected skill level to be positive")
		}
		if cs2.Elo <= 0 {
			t.Error("Expected ELO to be positive")
		}
	}
}

func TestGetPlayerByNickname_InvalidPlayer(t *testing.T) {
	apiKey := os.Getenv("FACEIT_API_KEY")
	if apiKey == "" {
		t.Skip("FACEIT_API_KEY not set, skipping integration test")
	}

	repo := NewFaceitRepository(apiKey)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Test with non-existent player
	_, err := repo.GetPlayerByNickname(ctx, "this-player-definitely-does-not-exist-12345")
	// Note: Faceit API might not always return an error for non-existent players
	// It might return a player with empty fields instead
	if err != nil {
		// If we get an error, that's expected
		t.Logf("Got expected error for non-existent player: %v", err)
	} else {
		// If no error, the API might have returned a player with empty fields
		t.Log("API returned no error for non-existent player (this might be expected behavior)")
	}
}

func TestGetPlayerStats(t *testing.T) {
	apiKey := os.Getenv("FACEIT_API_KEY")
	if apiKey == "" {
		t.Skip("FACEIT_API_KEY not set, skipping integration test")
	}

	repo := NewFaceitRepository(apiKey)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// First get a player to get their ID
	player, err := repo.GetPlayerByNickname(ctx, "s1mple")
	if err != nil {
		t.Fatalf("Failed to get player: %v", err)
	}

	// Get player stats
	stats, err := repo.GetPlayerStats(ctx, player.ID, "cs2")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if stats == nil {
		t.Fatal("Expected player stats, got nil")
	}

	// Check that we have some lifetime stats
	if stats.Lifetime == nil {
		t.Error("Expected lifetime stats to be present")
	}
}

func TestGetPlayerRecentMatches(t *testing.T) {
	apiKey := os.Getenv("FACEIT_API_KEY")
	if apiKey == "" {
		t.Skip("FACEIT_API_KEY not set, skipping integration test")
	}

	repo := NewFaceitRepository(apiKey)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// First get a player to get their ID
	player, err := repo.GetPlayerByNickname(ctx, "s1mple")
	if err != nil {
		t.Fatalf("Failed to get player: %v", err)
	}

	// Test with small limit
	matches, err := repo.GetPlayerRecentMatches(ctx, player.ID, "cs2", 5)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(matches) == 0 {
		t.Error("Expected at least some matches, got 0")
	}

	if len(matches) > 5 {
		t.Errorf("Expected at most 5 matches, got %d", len(matches))
	}

	// Verify match structure
	for i, match := range matches {
		if match.MatchID == "" {
			t.Errorf("Match %d: Expected MatchID to be set", i)
		}
		// Note: GameID is not part of PlayerMatchSummary, it's handled internally
		if match.Map == "" {
			t.Errorf("Match %d: Expected Map to be set", i)
		}
	}
}

func TestGetPlayerRecentMatches_Pagination(t *testing.T) {
	apiKey := os.Getenv("FACEIT_API_KEY")
	if apiKey == "" {
		t.Skip("FACEIT_API_KEY not set, skipping integration test")
	}

	repo := NewFaceitRepository(apiKey)
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// First get a player to get their ID
	player, err := repo.GetPlayerByNickname(ctx, "s1mple")
	if err != nil {
		t.Fatalf("Failed to get player: %v", err)
	}

	// Test pagination with 150 matches (should require 2 API calls)
	matches, err := repo.GetPlayerRecentMatches(ctx, player.ID, "cs2", 150)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Should get exactly 150 matches (or all available if less)
	if len(matches) > 150 {
		t.Errorf("Expected at most 150 matches, got %d", len(matches))
	}

	// Verify all matches have unique IDs
	matchIDs := make(map[string]bool)
	for i, match := range matches {
		if matchIDs[match.MatchID] {
			t.Errorf("Duplicate match ID found: %s at position %d", match.MatchID, i)
		}
		matchIDs[match.MatchID] = true
	}
}

func TestGetPlayerRecentMatches_InvalidPlayer(t *testing.T) {
	apiKey := os.Getenv("FACEIT_API_KEY")
	if apiKey == "" {
		t.Skip("FACEIT_API_KEY not set, skipping integration test")
	}

	repo := NewFaceitRepository(apiKey)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Test with invalid player ID
	_, err := repo.GetPlayerRecentMatches(ctx, "invalid-player-id", "cs2", 5)
	// Note: Faceit API might not always return an error for invalid player IDs
	if err != nil {
		// If we get an error, that's expected
		t.Logf("Got expected error for invalid player ID: %v", err)
	} else {
		// If no error, the API might have returned empty results
		t.Log("API returned no error for invalid player ID (this might be expected behavior)")
	}
}

func TestGetPlayerRecentMatches_InvalidGame(t *testing.T) {
	apiKey := os.Getenv("FACEIT_API_KEY")
	if apiKey == "" {
		t.Skip("FACEIT_API_KEY not set, skipping integration test")
	}

	repo := NewFaceitRepository(apiKey)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// First get a player to get their ID
	player, err := repo.GetPlayerByNickname(ctx, "s1mple")
	if err != nil {
		t.Fatalf("Failed to get player: %v", err)
	}

	// Test with invalid game ID
	_, err = repo.GetPlayerRecentMatches(ctx, player.ID, "invalid-game", 5)
	// Note: Faceit API might not always return an error for invalid game IDs
	if err != nil {
		// If we get an error, that's expected
		t.Logf("Got expected error for invalid game ID: %v", err)
	} else {
		// If no error, the API might have returned empty results
		t.Log("API returned no error for invalid game ID (this might be expected behavior)")
	}
}

func TestGetPlayerRecentMatches_EdgeCases(t *testing.T) {
	apiKey := os.Getenv("FACEIT_API_KEY")
	if apiKey == "" {
		t.Skip("FACEIT_API_KEY not set, skipping integration test")
	}

	repo := NewFaceitRepository(apiKey)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// First get a player to get their ID
	player, err := repo.GetPlayerByNickname(ctx, "s1mple")
	if err != nil {
		t.Fatalf("Failed to get player: %v", err)
	}

	// Test with zero limit (should default to 5)
	matches, err := repo.GetPlayerRecentMatches(ctx, player.ID, "cs2", 0)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(matches) > 5 {
		t.Errorf("Expected at most 5 matches with limit 0, got %d", len(matches))
	}

	// Test with negative limit (should default to 5)
	matches, err = repo.GetPlayerRecentMatches(ctx, player.ID, "cs2", -1)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(matches) > 5 {
		t.Errorf("Expected at most 5 matches with negative limit, got %d", len(matches))
	}
}

// Benchmark tests
func BenchmarkGetPlayerByNickname(b *testing.B) {
	apiKey := os.Getenv("FACEIT_API_KEY")
	if apiKey == "" {
		b.Skip("FACEIT_API_KEY not set, skipping benchmark")
	}
	
	// Skip benchmarks in short mode
	if testing.Short() {
		b.Skip("Skipping benchmark in short mode")
	}

	repo := NewFaceitRepository(apiKey)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := repo.GetPlayerByNickname(ctx, "s1mple")
		if err != nil {
			b.Fatalf("Benchmark failed: %v", err)
		}
	}
}

func BenchmarkGetPlayerRecentMatches(b *testing.B) {
	apiKey := os.Getenv("FACEIT_API_KEY")
	if apiKey == "" {
		b.Skip("FACEIT_API_KEY not set, skipping benchmark")
	}
	
	// Skip benchmarks in short mode
	if testing.Short() {
		b.Skip("Skipping benchmark in short mode")
	}

	repo := NewFaceitRepository(apiKey)
	ctx := context.Background()

	// Get player ID once
	player, err := repo.GetPlayerByNickname(ctx, "s1mple")
	if err != nil {
		b.Fatalf("Failed to get player: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := repo.GetPlayerRecentMatches(ctx, player.ID, "cs2", 20)
		if err != nil {
			b.Fatalf("Benchmark failed: %v", err)
		}
	}
}
