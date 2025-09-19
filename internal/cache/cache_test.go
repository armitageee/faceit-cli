package cache

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/armitageee/faceit-cli/internal/entity"
)

func TestCacheBasicOperations(t *testing.T) {
	cache := NewCache(1 * time.Minute)
	
	// Test Set and Get
	cache.Set("key1", "value1")
	
	value, found := cache.Get("key1")
	if !found {
		t.Error("Expected to find key1")
	}
	if value != "value1" {
		t.Errorf("Expected value1, got %v", value)
	}
	
	// Test non-existent key
	_, found = cache.Get("nonexistent")
	if found {
		t.Error("Expected not to find nonexistent key")
	}
}

func TestCacheExpiration(t *testing.T) {
	cache := NewCache(100 * time.Millisecond)
	
	// Set a value
	cache.Set("key1", "value1")
	
	// Should be found immediately
	_, found := cache.Get("key1")
	if !found {
		t.Error("Expected to find key1 immediately")
	}
	
	// Wait for expiration
	time.Sleep(150 * time.Millisecond)
	
	// Should not be found after expiration
	_, found = cache.Get("key1")
	if found {
		t.Error("Expected key1 to be expired")
	}
}

func TestCacheDelete(t *testing.T) {
	cache := NewCache(1 * time.Minute)
	
	// Set and verify
	cache.Set("key1", "value1")
	_, found := cache.Get("key1")
	if !found {
		t.Error("Expected to find key1")
	}
	
	// Delete and verify
	cache.Delete("key1")
	_, found = cache.Get("key1")
	if found {
		t.Error("Expected key1 to be deleted")
	}
}

func TestCacheClear(t *testing.T) {
	cache := NewCache(1 * time.Minute)
	
	// Set multiple values
	cache.Set("key1", "value1")
	cache.Set("key2", "value2")
	cache.Set("key3", "value3")
	
	// Verify all exist
	_, found := cache.Get("key1")
	if !found {
		t.Error("Expected to find key1")
	}
	
	// Clear and verify all are gone
	cache.Clear()
	_, found = cache.Get("key1")
	if found {
		t.Error("Expected key1 to be cleared")
	}
	_, found = cache.Get("key2")
	if found {
		t.Error("Expected key2 to be cleared")
	}
	_, found = cache.Get("key3")
	if found {
		t.Error("Expected key3 to be cleared")
	}
}

func TestGenerateKeys(t *testing.T) {
	// Test player matches key
	key := GeneratePlayerMatchesKey("player123", "cs2", 50)
	expected := "matches:player123:cs2:50"
	if key != expected {
		t.Errorf("Expected %s, got %s", expected, key)
	}
	
	// Test player profile key
	key = GeneratePlayerProfileKey("testplayer")
	expected = "profile:testplayer"
	if key != expected {
		t.Errorf("Expected %s, got %s", expected, key)
	}
	
	// Test player stats key
	key = GeneratePlayerStatsKey("player123", "cs2")
	expected = "stats:player123:cs2"
	if key != expected {
		t.Errorf("Expected %s, got %s", expected, key)
	}
	
	// Test match stats key
	key = GenerateMatchStatsKey("match123")
	expected = "match_stats:match123"
	if key != expected {
		t.Errorf("Expected %s, got %s", expected, key)
	}
}

// Mock repository for testing
type mockRepository struct {
	profiles   map[string]*entity.PlayerProfile
	stats      map[string]*entity.PlayerStats
	matches    map[string][]entity.PlayerMatchSummary
	matchStats map[string]*entity.MatchStats
}

func newMockRepository() *mockRepository {
	return &mockRepository{
		profiles:   make(map[string]*entity.PlayerProfile),
		stats:      make(map[string]*entity.PlayerStats),
		matches:    make(map[string][]entity.PlayerMatchSummary),
		matchStats: make(map[string]*entity.MatchStats),
	}
}

func (m *mockRepository) GetPlayerByNickname(ctx context.Context, nickname string) (*entity.PlayerProfile, error) {
	if profile, exists := m.profiles[nickname]; exists {
		return profile, nil
	}
	return nil, fmt.Errorf("player not found")
}

func (m *mockRepository) GetPlayerStats(ctx context.Context, playerID, gameID string) (*entity.PlayerStats, error) {
	key := playerID + ":" + gameID
	if stats, exists := m.stats[key]; exists {
		return stats, nil
	}
	return nil, fmt.Errorf("stats not found")
}

func (m *mockRepository) GetPlayerRecentMatches(ctx context.Context, playerID string, gameID string, limit int) ([]entity.PlayerMatchSummary, error) {
	key := playerID + ":" + gameID + ":" + string(rune(limit))
	if matches, exists := m.matches[key]; exists {
		return matches, nil
	}
	return nil, fmt.Errorf("matches not found")
}

func (m *mockRepository) GetMatchStats(ctx context.Context, matchID string) (*entity.MatchStats, error) {
	if stats, exists := m.matchStats[matchID]; exists {
		return stats, nil
	}
	return nil, fmt.Errorf("match stats not found")
}

func TestCachedRepository(t *testing.T) {
	mockRepo := newMockRepository()
	cachedRepo := NewCachedFaceitRepository(mockRepo, 1*time.Minute)
	
	// Test profile caching
	profile := &entity.PlayerProfile{
		ID:       "test123",
		Nickname: "testplayer",
		Country:  "US",
	}
	mockRepo.profiles["testplayer"] = profile
	
	ctx := context.Background()
	
	// First call should hit the repository
	result, err := cachedRepo.GetPlayerByNickname(ctx, "testplayer")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if result.Nickname != "testplayer" {
		t.Errorf("Expected testplayer, got %s", result.Nickname)
	}
	
	// Second call should hit the cache (we can't easily test this without modifying the mock)
	// But we can verify the cache has the item
	key := GeneratePlayerProfileKey("testplayer")
	_, found := cachedRepo.cache.Get(key)
	if !found {
		t.Error("Expected profile to be cached")
	}
}

func TestCachedMatchStats(t *testing.T) {
	mockRepo := newMockRepository()
	cachedRepo := NewCachedFaceitRepository(mockRepo, 1*time.Minute)
	
	// Test match stats caching
	matchStats := &entity.MatchStats{
		MatchID: "test-match-123",
		Map:     "de_dust2",
		Score:   "16-14",
		Result:  "finished",
		Team1: entity.TeamMatchStats{
			TeamID:   "team1",
			TeamName: "Team 1",
			Score:    16,
			Players:  []entity.PlayerMatchStats{},
		},
		Team2: entity.TeamMatchStats{
			TeamID:   "team2",
			TeamName: "Team 2",
			Score:    14,
			Players:  []entity.PlayerMatchStats{},
		},
		PlayerStats: []entity.PlayerMatchStats{},
	}
	mockRepo.matchStats["test-match-123"] = matchStats
	
	ctx := context.Background()
	
	// First call should hit the repository
	result, err := cachedRepo.GetMatchStats(ctx, "test-match-123")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if result.MatchID != "test-match-123" {
		t.Errorf("Expected test-match-123, got %s", result.MatchID)
	}
	
	// Verify the cache has the item
	key := GenerateMatchStatsKey("test-match-123")
	_, found := cachedRepo.cache.Get(key)
	if !found {
		t.Error("Expected match stats to be cached")
	}
}
