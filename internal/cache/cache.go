package cache

import (
	"context"
	"fmt"
	"sync"
	"time"

	"faceit-cli/internal/entity"
)

// CacheEntry represents a cached item with expiration
type CacheEntry struct {
	Data      interface{}
	ExpiresAt time.Time
}

// IsExpired checks if the cache entry has expired
func (e *CacheEntry) IsExpired() bool {
	return time.Now().After(e.ExpiresAt)
}

// Cache provides in-memory caching with TTL support
type Cache struct {
	mu    sync.RWMutex
	items map[string]*CacheEntry
	ttl   time.Duration
}

// NewCache creates a new cache instance with the specified TTL
func NewCache(ttl time.Duration) *Cache {
	c := &Cache{
		items: make(map[string]*CacheEntry),
		ttl:   ttl,
	}
	
	// Start cleanup goroutine
	go c.cleanup()
	
	return c
}

// Set stores a value in the cache with the default TTL
func (c *Cache) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	c.items[key] = &CacheEntry{
		Data:      value,
		ExpiresAt: time.Now().Add(c.ttl),
	}
}

// Get retrieves a value from the cache
func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	entry, exists := c.items[key]
	if !exists || entry.IsExpired() {
		return nil, false
	}
	
	return entry.Data, true
}

// Delete removes a value from the cache
func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	delete(c.items, key)
}

// Clear removes all items from the cache
func (c *Cache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	c.items = make(map[string]*CacheEntry)
}

// cleanup removes expired entries periodically
func (c *Cache) cleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	
	for range ticker.C {
		c.mu.Lock()
		for key, entry := range c.items {
			if entry.IsExpired() {
				delete(c.items, key)
			}
		}
		c.mu.Unlock()
	}
}

// GenerateKey creates a cache key for player matches
func GeneratePlayerMatchesKey(playerID, gameID string, limit int) string {
	return fmt.Sprintf("matches:%s:%s:%d", playerID, gameID, limit)
}

// GenerateKey creates a cache key for player profile
func GeneratePlayerProfileKey(nickname string) string {
	return fmt.Sprintf("profile:%s", nickname)
}

// GenerateKey creates a cache key for player stats
func GeneratePlayerStatsKey(playerID, gameID string) string {
	return fmt.Sprintf("stats:%s:%s", playerID, gameID)
}

// GenerateKey creates a cache key for match stats
func GenerateMatchStatsKey(matchID string) string {
	return fmt.Sprintf("match_stats:%s", matchID)
}

// CachedFaceitRepository wraps a FaceitRepository with caching
type CachedFaceitRepository struct {
	repo  FaceitRepository
	cache *Cache
}

// FaceitRepository interface for dependency injection
type FaceitRepository interface {
	GetPlayerByNickname(ctx context.Context, nickname string) (*entity.PlayerProfile, error)
	GetPlayerStats(ctx context.Context, playerID, gameID string) (*entity.PlayerStats, error)
	GetPlayerRecentMatches(ctx context.Context, playerID string, gameID string, limit int) ([]entity.PlayerMatchSummary, error)
	GetMatchStats(ctx context.Context, matchID string) (*entity.MatchStats, error)
}

// NewCachedFaceitRepository creates a new cached repository
func NewCachedFaceitRepository(repo FaceitRepository, cacheTTL time.Duration) *CachedFaceitRepository {
	return &CachedFaceitRepository{
		repo:  repo,
		cache: NewCache(cacheTTL),
	}
}

// GetPlayerByNickname implements FaceitRepository interface with caching
func (c *CachedFaceitRepository) GetPlayerByNickname(ctx context.Context, nickname string) (*entity.PlayerProfile, error) {
	key := GeneratePlayerProfileKey(nickname)
	
	// Try to get from cache
	if cached, found := c.cache.Get(key); found {
		if profile, ok := cached.(*entity.PlayerProfile); ok {
			return profile, nil
		}
	}
	
	// Get from repository
	profile, err := c.repo.GetPlayerByNickname(ctx, nickname)
	if err != nil {
		return nil, err
	}
	
	// Cache the result
	c.cache.Set(key, profile)
	
	return profile, nil
}

// GetPlayerStats implements FaceitRepository interface with caching
func (c *CachedFaceitRepository) GetPlayerStats(ctx context.Context, playerID, gameID string) (*entity.PlayerStats, error) {
	key := GeneratePlayerStatsKey(playerID, gameID)
	
	// Try to get from cache
	if cached, found := c.cache.Get(key); found {
		if stats, ok := cached.(*entity.PlayerStats); ok {
			return stats, nil
		}
	}
	
	// Get from repository
	stats, err := c.repo.GetPlayerStats(ctx, playerID, gameID)
	if err != nil {
		return nil, err
	}
	
	// Cache the result
	c.cache.Set(key, stats)
	
	return stats, nil
}

// GetPlayerRecentMatches implements FaceitRepository interface with caching
func (c *CachedFaceitRepository) GetPlayerRecentMatches(ctx context.Context, playerID string, gameID string, limit int) ([]entity.PlayerMatchSummary, error) {
	key := GeneratePlayerMatchesKey(playerID, gameID, limit)
	
	// Try to get from cache
	if cached, found := c.cache.Get(key); found {
		if matches, ok := cached.([]entity.PlayerMatchSummary); ok {
			return matches, nil
		}
	}
	
	// Get from repository
	matches, err := c.repo.GetPlayerRecentMatches(ctx, playerID, gameID, limit)
	if err != nil {
		return nil, err
	}
	
	// Cache the result
	c.cache.Set(key, matches)
	
	return matches, nil
}

// ClearCache clears all cached data
func (c *CachedFaceitRepository) ClearCache() {
	c.cache.Clear()
}

// GetCacheStats returns cache statistics
func (c *CachedFaceitRepository) GetCacheStats() map[string]interface{} {
	c.cache.mu.RLock()
	defer c.cache.mu.RUnlock()
	
	return map[string]interface{}{
		"total_items": len(c.cache.items),
		"ttl":         c.cache.ttl.String(),
	}
}

// GetMatchStats implements FaceitRepository interface with caching
func (c *CachedFaceitRepository) GetMatchStats(ctx context.Context, matchID string) (*entity.MatchStats, error) {
	key := GenerateMatchStatsKey(matchID)
	
	// Try to get from cache
	if cached, found := c.cache.Get(key); found {
		if stats, ok := cached.(*entity.MatchStats); ok {
			return stats, nil
		}
	}
	
	// Get from repository
	stats, err := c.repo.GetMatchStats(ctx, matchID)
	if err != nil {
		return nil, err
	}
	
	// Cache the result
	c.cache.Set(key, stats)
	
	return stats, nil
}
