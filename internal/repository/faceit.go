package repository

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"sort"

	"faceit-cli/internal/entity"

	"github.com/antihax/optional"
	faceit "github.com/mconnat/go-faceit"
)

// FaceitRepository defines the interface for FACEIT API operations
type FaceitRepository interface {
	GetPlayerByNickname(ctx context.Context, nickname string) (*entity.PlayerProfile, error)
	GetPlayerStats(ctx context.Context, playerID, gameID string) (*entity.PlayerStats, error)
	GetPlayerRecentMatches(ctx context.Context, playerID string, gameID string, limit int) ([]entity.PlayerMatchSummary, error)
}

// faceitRepository is a concrete implementation of FaceitRepository that
// delegates calls to the generated FACEIT API client. It stores the API
// key and applies it to each request via the context.
type faceitRepository struct {
	client *faceit.APIClient
	apiKey string
}

// NewFaceitRepository constructs a repository backed by the FACEIT API.
// It takes an API key which will be sent with each request. The
// underlying API client is created with default configuration –
// including the base URL "https://open.faceit.com/data/v4".
func NewFaceitRepository(apiKey string) FaceitRepository {
	cfg := faceit.NewConfiguration()
	client := faceit.NewAPIClient(cfg)
	return &faceitRepository{
		client: client,
		apiKey: apiKey,
	}
}

// contextWithAPIKey injects the API key into the provided context. The
// FACEIT client looks up this value under the ContextAPIKey key to set
// the Authorization header on outgoing requests.
func (r *faceitRepository) contextWithAPIKey(ctx context.Context) context.Context {
	return context.WithValue(ctx, faceit.ContextAccessToken, r.apiKey)
}

// GetPlayerByNickname resolves a player's FACEIT nickname to their unique
// player ID and returns detailed profile information. It performs two
// API calls: first a search to find the player ID and then a fetch of
// the player's profile.
func (r *faceitRepository) GetPlayerByNickname(ctx context.Context, nickname string) (*entity.PlayerProfile, error) {
	if nickname == "" {
		return nil, fmt.Errorf("nickname must not be empty")
	}

	ctx = r.contextWithAPIKey(ctx)
	// Search players by nickname. We don't filter by game or country
	// because a nickname should uniquely identify a user. The API returns
	// a list – we take the first item. Optionally, callers could handle
	// ambiguities here by matching the exact nickname or by returning
	// multiple results.
	opts := &faceit.SearchApiSearchPlayersOpts{}
	list, _, err := r.client.SearchApi.SearchPlayers(ctx, nickname, opts)
	if err != nil {
		return nil, fmt.Errorf("search players: %w", err)
	}
	if len(list.Items) == 0 {
		return nil, fmt.Errorf("player not found: %s", nickname)
	}
	playerID := list.Items[0].PlayerId
	// Retrieve full player details using the resolved ID. A separate
	// endpoint exists to fetch details directly by nickname, but the
	// search call ensures we have a valid ID before proceeding.
	player, _, err := r.client.PlayersApi.GetPlayer(ctx, playerID)
	if err != nil {
		return nil, fmt.Errorf("get player: %w", err)
	}
	profile := &entity.PlayerProfile{
		ID:        player.PlayerId,
		Nickname:  player.Nickname,
		Country:   player.Country,
		Avatar:    player.Avatar,
		FaceitURL: player.FaceitUrl,
		Games:     make(map[string]entity.GameDetail),
	}
	// Populate per‑game details. Each entry in the map corresponds to
	// a registered game (e.g. "cs2", "dota2").
	if player.Games != nil {
		for k, g := range player.Games {
			profile.Games[k] = entity.GameDetail{
				Elo:        int(g.FaceitElo),
				SkillLevel: int(g.SkillLevel),
				Region:     g.Region,
			}
		}
	}

	return profile, nil
}

// GetPlayerStats fetches lifetime statistics for a given player and game.
// The Faceit API exposes two methods: one that returns statistics for a
// given number of recent matches and another that returns aggregated
// lifetime statistics. This implementation uses the latter via
// GetPlayerStats_1.
func (r *faceitRepository) GetPlayerStats(ctx context.Context, playerID, gameID string) (*entity.PlayerStats, error) {
	if playerID == "" {
		return nil, fmt.Errorf("playerID must not be empty")
	}
	if gameID == "" {
		return nil, fmt.Errorf("gameID must not be empty")
	}

	ctx = r.contextWithAPIKey(ctx)
	stats, _, err := r.client.PlayersApi.GetPlayerStats_1(ctx, playerID, gameID)
	if err != nil {
		return nil, fmt.Errorf("get player stats: %w", err)
	}

	result := &entity.PlayerStats{
		GameID:   stats.GameId,
		PlayerID: stats.PlayerId,
		Lifetime: stats.Lifetime,
		Segments: stats.Segments,
	}

	return result, nil
}

// GetPlayerRecentMatches implements the FaceitRepository interface.  It
// fetches the most recent matches for the given player and game and
// enriches the data with individual player statistics and match
// outcomes.  When the limit is zero or negative a default of five
// matches is used.  Any errors encountered during history or stats
// retrieval are returned immediately.
func (r *faceitRepository) GetPlayerRecentMatches(ctx context.Context, playerID string, gameID string, limit int) ([]entity.PlayerMatchSummary, error) {
	if playerID == "" {
		return nil, fmt.Errorf("playerID must not be empty")
	}
	if gameID == "" {
		return nil, fmt.Errorf("gameID must not be empty")
	}
	if limit <= 0 {
		limit = 5
	}

	ctx = r.contextWithAPIKey(ctx)

	// Prepare optional parameters for the history call.  The
	// PlayersApiGetPlayerHistoryOpts type supports specifying a
	// maximum number of records to return via the Limit field.
	opts := &faceit.PlayersApiGetPlayerHistoryOpts{}
	opts.Limit = optional.NewInt32(int32(limit))

	history, _, err := r.client.PlayersApi.GetPlayerHistory(ctx, playerID, gameID, opts)
	if err != nil {
		return nil, fmt.Errorf("get player history: %w", err)
	}
	// MatchHistoryList is a struct, not a pointer, so it cannot be compared to nil.
	// When the history has no items, simply return an empty slice.
	if len(history.Items) == 0 {
		return nil, nil
	}
	
	// Debug logging (can be removed in production)
	// fmt.Printf("DEBUG: Requested limit: %d, Got matches: %d\n", limit, len(history.Items))
	results := make([]entity.PlayerMatchSummary, 0, len(history.Items))

	// Iterate through the returned matches.  The API lists matches
	// from newest to oldest so we preserve the order provided.
	for _, item := range history.Items {

		// Identify the team the player belonged to by scanning the
		// Teams map for the player's ID.
		var playerTeamID string
		if item.Teams != nil {
			for teamID, faction := range item.Teams {
				if faction.Players != nil {
					for _, p := range faction.Players {
						// MatchHistoryPlayer.PlayerId is a string
						if p.PlayerId == playerID {
							playerTeamID = teamID
							break
						}
					}
				}
				if playerTeamID != "" {
					break
				}
			}
		}
		// Determine the match result (win/loss) by comparing the
		// player's team with the winner reported in the results.
		result := "Loss"
		scoreStr := ""
		if item.Results != nil {
			if item.Results.Winner == playerTeamID {
				result = "Win"
			}
			// Format the score as a hyphen separated string (e.g. "16-14").
			if item.Results.Score != nil {
				// To maintain a deterministic order we iterate over
				// faction IDs, sort them alphabetically and then join
				// their scores.  There are usually two teams.
				var teamIDs []string
				for tID := range item.Results.Score {
					teamIDs = append(teamIDs, tID)
				}
				if len(teamIDs) > 0 {
					// Simple lexicographic sort so that order is stable
					// across calls.
					sort.Strings(teamIDs)
					var parts []string
					for _, tID := range teamIDs {
						v := item.Results.Score[tID]
						parts = append(parts, strconv.FormatInt(v, 10))
					}
					scoreStr = strings.Join(parts, "-")
				}
			}
		}

		// Extract per‑player statistics by requesting the match stats
		// endpoint.  If this call fails we log the error and skip
		// statistics for this match rather than aborting the whole
		// request.  This ensures that at least minimal information is
		// returned to the user.
		var kills, deaths, assists int
		var kdRatio, hsPerc, adr float64
		// Create a separate context with longer timeout for match stats
		statsCtx, statsCancel := context.WithTimeout(context.Background(), 30*time.Second)
		statsCtx = r.contextWithAPIKey(statsCtx)
		
		stats, _, err := r.client.MatchesApi.GetMatchStats(statsCtx, item.MatchId)
		statsCancel()
		
		if err != nil {
			// Log the error and continue without stats.
			// fmt.Printf("DEBUG: Failed to get stats for match %s: %v\n", item.MatchId, err)
		} else if len(stats.Rounds) > 0 {
			// Aggregate statistics across all rounds for the player. The API may
			// provide per‑round data; summing kills, deaths and assists over
			// rounds yields the total for the match. Headshot percentage is
			// computed based on the proportion of headshot kills over total
			// kills.
			var hsKills float64
			for _, round := range stats.Rounds {
				if round.Teams == nil {
					continue
				}
				found := false
				for _, team := range round.Teams {
					if team.Players == nil {
						continue
					}
					for _, ps := range team.Players {
						var pid string
						switch v := ps.PlayerId.(type) {
						case string:
							pid = v
						}
						if pid != playerID {
							continue
						}
						if ps.PlayerStats != nil {
							// Helpers to convert arbitrary values. Defined inside
							// loop to capture strconv and strings from outer scope.
							toInt := func(x interface{}) int {
								switch v := x.(type) {
								case float64:
									return int(v)
								case float32:
									return int(v)
								case int64:
									return int(v)
								case int32:
									return int(v)
								case int:
									return v
								case string:
									i, _ := strconv.Atoi(v)
									return i
								default:
									return 0
								}
							}
							toFloat := func(x interface{}) float64 {
								switch v := x.(type) {
								case float64:
									return v
								case float32:
									return float64(v)
								case int64:
									return float64(v)
								case int32:
									return float64(v)
								case int:
									return float64(v)
								case string:
									f, _ := strconv.ParseFloat(strings.TrimSuffix(v, "%"), 64)
									return f
								default:
									return 0
								}
							}
							// Sum kills, deaths and assists across rounds.
							if v, ok := ps.PlayerStats["Kills"]; ok {
								kills += toInt(v)
								// Track headshot kills proportionally when available.
								var hsRound float64
								if vhs, ok2 := ps.PlayerStats["Headshots %"]; ok2 {
									hsRound = toFloat(vhs)
								} else if vhs, ok2 := ps.PlayerStats["HS %"]; ok2 {
									hsRound = toFloat(vhs)
								}
								// Add headshot kills based on percentage.
								hsKills += float64(toInt(v)) * hsRound / 100.0
							}
							if v, ok := ps.PlayerStats["Deaths"]; ok {
								deaths += toInt(v)
							}
							if v, ok := ps.PlayerStats["Assists"]; ok {
								assists += toInt(v)
							}
							// Extract ADR (Average Damage per Round)
							if v, ok := ps.PlayerStats["Average Damage per Round"]; ok {
								adr += toFloat(v)
							} else if v, ok := ps.PlayerStats["ADR"]; ok {
								adr += toFloat(v)
							} else if v, ok := ps.PlayerStats["Avg Damage"]; ok {
								adr += toFloat(v)
							}
						}
						// Break out after finding the player in this round.
						found = true
						break
					}
					if found {
						break
					}
				}
			}
			// Compute K/D ratio and headshot percentage. Fall back to 0 when
			// values are not meaningful.
			if deaths > 0 {
				kdRatio = float64(kills) / float64(deaths)
			} else if kills > 0 {
				kdRatio = float64(kills)
			}
			if kills > 0 {
				hsPerc = (hsKills / float64(kills)) * 100.0
			}
		}

		// Get the map name - we already have match stats, so use them directly
		mapName := item.GameMode

		// If GameMode is generic, try to get map name from match stats (which we already have)
		if mapName == "5v5" || mapName == "" {
			if len(stats.Rounds) > 0 {
				firstRound := stats.Rounds[0]
				if firstRound.RoundStats != nil {
					if mapVal, exists := firstRound.RoundStats["Map"]; exists {
						if mapStr, ok := mapVal.(string); ok {
							mapName = mapStr
						}
					}
				}
			}
		}

		// Compose the summary.  FinishedAt is provided as int64.  Map
		// corresponds to the actual map name from the match.  Not all
		// fields are guaranteed to be populated so defaults apply.
		summary := entity.PlayerMatchSummary{
			MatchID:             item.MatchId,
			Map:                 mapName,
			FinishedAt:          item.FinishedAt,
			Score:               scoreStr,
			Kills:               kills,
			Deaths:              deaths,
			Assists:             assists,
			KDRatio:             kdRatio,
			HeadshotsPercentage: hsPerc,
			ADR:                 adr,
			Result:              result,
		}
		results = append(results, summary)
	}

	return results, nil
}
