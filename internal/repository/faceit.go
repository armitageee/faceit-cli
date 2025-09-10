package repository

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"sort"

	"faceit-cli/internal/entity"
	"faceit-cli/internal/logger"

	"github.com/antihax/optional"
	faceit "github.com/mconnat/go-faceit"
)

// FaceitRepository defines the interface for FACEIT API operations
type FaceitRepository interface {
	GetPlayerByNickname(ctx context.Context, nickname string) (*entity.PlayerProfile, error)
	GetPlayerStats(ctx context.Context, playerID, gameID string) (*entity.PlayerStats, error)
	GetPlayerRecentMatches(ctx context.Context, playerID string, gameID string, limit int) ([]entity.PlayerMatchSummary, error)
	GetMatchStats(ctx context.Context, matchID string) (*entity.MatchStats, error)
}

// faceitRepository is a concrete implementation of FaceitRepository that
// delegates calls to the generated FACEIT API client. It stores the API
// key and applies it to each request via the context.
type faceitRepository struct {
	client *faceit.APIClient
	apiKey string
	logger *logger.Logger
}

// NewFaceitRepository constructs a repository backed by the FACEIT API.
// It takes an API key which will be sent with each request. The
// underlying API client is created with default configuration –
// including the base URL "https://open.faceit.com/data/v4".
func NewFaceitRepository(apiKey string) FaceitRepository {
	cfg := faceit.NewConfiguration()
	client := faceit.NewAPIClient(cfg)
	
	// Create logger with default config
	loggerConfig := logger.Config{
		Level:          logger.LogLevelInfo,
		KafkaEnabled:   false,
		ServiceName:    "faceit-repository",
		ProductionMode: false,
		LogToStdout:    true,
	}
	appLogger, _ := logger.New(loggerConfig)
	
	return &faceitRepository{
		client: client,
		apiKey: apiKey,
		logger: appLogger,
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
	r.logger.Debug("Starting GetPlayerByNickname", map[string]interface{}{
		"nickname": nickname,
	})

	if nickname == "" {
		r.logger.Debug("Empty nickname provided", nil)
		return nil, fmt.Errorf("nickname must not be empty")
	}

	ctx = r.contextWithAPIKey(ctx)
	// Search players by nickname. We don't filter by game or country
	// because a nickname should uniquely identify a user. The API returns
	// a list – we take the first item. Optionally, callers could handle
	// ambiguities here by matching the exact nickname or by returning
	// multiple results.
	opts := &faceit.SearchApiSearchPlayersOpts{}
	
	r.logger.Debug("Searching players by nickname", map[string]interface{}{
		"nickname": nickname,
		"options":  opts,
	})
	
	list, _, err := r.client.SearchApi.SearchPlayers(ctx, nickname, opts)
	if err != nil {
		r.logger.Error("Failed to search players", map[string]interface{}{
			"nickname": nickname,
			"error":    err.Error(),
		})
		return nil, fmt.Errorf("search players: %w", err)
	}
	
	r.logger.Debug("Search players response", map[string]interface{}{
		"nickname":     nickname,
		"results_count": len(list.Items),
	})
	
	if len(list.Items) == 0 {
		r.logger.Debug("No players found", map[string]interface{}{
			"nickname": nickname,
		})
		return nil, fmt.Errorf("player not found: %s", nickname)
	}
	playerID := list.Items[0].PlayerId
	
	r.logger.Debug("Found player ID, fetching details", map[string]interface{}{
		"nickname": nickname,
		"player_id": playerID,
	})
	
	// Retrieve full player details using the resolved ID. A separate
	// endpoint exists to fetch details directly by nickname, but the
	// search call ensures we have a valid ID before proceeding.
	player, _, err := r.client.PlayersApi.GetPlayer(ctx, playerID)
	if err != nil {
		r.logger.Error("Failed to get player details", map[string]interface{}{
			"nickname":  nickname,
			"player_id": playerID,
			"error":     err.Error(),
		})
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

	var allMatches []entity.PlayerMatchSummary
	offset := 0
	maxPerRequest := 100 // Faceit API maximum per request

	for len(allMatches) < limit {
		// Calculate how many matches to request in this batch
		remaining := limit - len(allMatches)
		batchSize := maxPerRequest
		if remaining < maxPerRequest {
			batchSize = remaining
		}

		// Prepare optional parameters for the history call
		opts := &faceit.PlayersApiGetPlayerHistoryOpts{}
		opts.Limit = optional.NewInt32(int32(batchSize))
		opts.Offset = optional.NewInt32(int32(offset))

		history, _, err := r.client.PlayersApi.GetPlayerHistory(ctx, playerID, gameID, opts)
		if err != nil {
			return nil, fmt.Errorf("get player history: %w", err)
		}

		// If no more matches, break
		if len(history.Items) == 0 {
			break
		}

		// If we got fewer matches than requested, we've reached the end
		if len(history.Items) < batchSize {
			// Process this final batch and break
			allMatches = append(allMatches, r.processMatches(history.Items, playerID)...)
			break
		}

		// Process this batch
		allMatches = append(allMatches, r.processMatches(history.Items, playerID)...)

		// Move to next batch
		offset += len(history.Items)
	}

	// Debug logging (can be removed in production)
	// fmt.Printf("DEBUG: Requested limit: %d, Got matches: %d\n", limit, len(allMatches))
	return allMatches, nil
}

// processMatches processes a batch of matches and returns PlayerMatchSummary slice
func (r *faceitRepository) processMatches(items []faceit.MatchHistory, playerID string) []entity.PlayerMatchSummary {
	results := make([]entity.PlayerMatchSummary, 0, len(items))
	// Iterate through the returned matches.  The API lists matches
	// from newest to oldest so we preserve the order provided.
	for _, item := range items {

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

	return results
}

// GetMatchStats retrieves detailed match statistics by match ID
func (r *faceitRepository) GetMatchStats(ctx context.Context, matchID string) (*entity.MatchStats, error) {
	r.logger.Debug("Starting GetMatchStats", map[string]interface{}{
		"match_id": matchID,
	})

	if matchID == "" {
		r.logger.Debug("Empty matchID provided", nil)
		return nil, fmt.Errorf("matchID must not be empty")
	}

	ctx = r.contextWithAPIKey(ctx)
	
	r.logger.Debug("Fetching match details", map[string]interface{}{
		"match_id": matchID,
	})
	
	// Try to get match details first
	match, _, err := r.client.MatchesApi.GetMatch(ctx, matchID)
	if err != nil {
		r.logger.Error("Failed to get match details", map[string]interface{}{
			"match_id": matchID,
			"error":    err.Error(),
		})
		// If match not found, return a helpful error
		return nil, fmt.Errorf("match not found: %s. Please check the Match ID and try again", matchID)
	}

	r.logger.Debug("Match details retrieved", map[string]interface{}{
		"match_id":    matchID,
		"status":      match.Status,
		"finished_at": match.FinishedAt,
	})

	r.logger.Debug("Fetching match statistics", map[string]interface{}{
		"match_id": matchID,
	})

	// Get match statistics
	stats, _, err := r.client.MatchesApi.GetMatchStats(ctx, matchID)
	if err != nil {
		r.logger.Error("Failed to get match statistics", map[string]interface{}{
			"match_id": matchID,
			"error":    err.Error(),
		})
		// If stats not available, return basic match info
		return &entity.MatchStats{
			MatchID:    matchID,
			Map:        "Unknown",
			FinishedAt: match.FinishedAt,
			Score:      "N/A",
			Result:     match.Status,
			Team1: entity.TeamMatchStats{
				TeamID:   "team1",
				TeamName: "Team 1",
				Score:    0,
				Players:  []entity.PlayerMatchStats{},
			},
			Team2: entity.TeamMatchStats{
				TeamID:   "team2", 
				TeamName: "Team 2",
				Score:    0,
				Players:  []entity.PlayerMatchStats{},
			},
			PlayerStats: []entity.PlayerMatchStats{},
		}, nil
	}

	// Extract map name from RoundStats
	mapName := "Unknown"
	scoreStr := "0-0"
	
	if len(stats.Rounds) > 0 {
		roundStats := stats.Rounds[0].RoundStats
		if roundStats != nil {
			if mapVal, ok := roundStats["Map"]; ok {
				if mapStr, ok := mapVal.(string); ok {
					mapName = mapStr
				}
			}
			if scoreVal, ok := roundStats["Score"]; ok {
				if scoreStrVal, ok := scoreVal.(string); ok {
					// Convert "6 / 13" format to "6-13"
					scoreStr = strings.ReplaceAll(scoreStrVal, " / ", "-")
				}
			}
		}
	}

	// Initialize match stats with basic info
	matchStats := &entity.MatchStats{
		MatchID:    matchID,
		Map:        mapName,
		FinishedAt: match.FinishedAt,
		Score:      scoreStr,
		Result:     match.Status,
		Team1: entity.TeamMatchStats{
			TeamID:   "team1",
			TeamName: "Team 1",
			Score:    0,
			Players:  []entity.PlayerMatchStats{},
		},
		Team2: entity.TeamMatchStats{
			TeamID:   "team2",
			TeamName: "Team 2", 
			Score:    0,
			Players:  []entity.PlayerMatchStats{},
		},
		PlayerStats: []entity.PlayerMatchStats{},
	}

	// Process teams and players with real data
	if len(stats.Rounds) > 0 && len(stats.Rounds[0].Teams) >= 2 {
		round := stats.Rounds[0]
		
		for i, team := range round.Teams {
			if i >= 2 {
				break // Only process first 2 teams
			}
			
			teamID := fmt.Sprintf("team%d", i+1)
			teamName := fmt.Sprintf("Team %d", i+1)
			teamScore := 0
			
			// Try to get team name and score from TeamStats
			if team.TeamStats != nil {
				if nameVal, ok := team.TeamStats["Team"]; ok {
					if nameStr, ok := nameVal.(string); ok {
						teamName = nameStr
					}
				}
				if scoreVal, ok := team.TeamStats["Final Score"]; ok {
					switch v := scoreVal.(type) {
					case float64:
						teamScore = int(v)
					case int:
						teamScore = v
					case string:
						if parsed, err := strconv.Atoi(v); err == nil {
							teamScore = parsed
						}
					}
				}
			}
			
			teamStats := entity.TeamMatchStats{
				TeamID:   teamID,
				TeamName: teamName,
				Score:    teamScore,
				Players:  []entity.PlayerMatchStats{},
			}

			// Process players in team
			for _, player := range team.Players {
				playerID := "unknown"
				nickname := "Unknown Player"
				
				// Extract player info
				if pid, ok := player.PlayerId.(string); ok {
					playerID = pid
				}
				if nick, ok := player.Nickname.(string); ok {
					nickname = nick
				}

				// Extract stats from PlayerStats
				kills := 0
				deaths := 0
				assists := 0
				hsPerc := 0.0
				adr := 0.0
				
				if player.PlayerStats != nil {
					// Extract stats with proper type handling
					if k, ok := player.PlayerStats["Kills"]; ok {
						switch v := k.(type) {
						case float64:
							kills = int(v)
						case int:
							kills = v
						case string:
							if parsed, err := strconv.Atoi(v); err == nil {
								kills = parsed
							}
						}
					}
					if d, ok := player.PlayerStats["Deaths"]; ok {
						switch v := d.(type) {
						case float64:
							deaths = int(v)
						case int:
							deaths = v
						case string:
							if parsed, err := strconv.Atoi(v); err == nil {
								deaths = parsed
							}
						}
					}
					if a, ok := player.PlayerStats["Assists"]; ok {
						switch v := a.(type) {
						case float64:
							assists = int(v)
						case int:
							assists = v
						case string:
							if parsed, err := strconv.Atoi(v); err == nil {
								assists = parsed
							}
						}
					}
					if h, ok := player.PlayerStats["Headshots %"]; ok {
						switch v := h.(type) {
						case float64:
							hsPerc = v
						case int:
							hsPerc = float64(v)
						case string:
							if parsed, err := strconv.ParseFloat(v, 64); err == nil {
								hsPerc = parsed
							}
						}
					}
					if r, ok := player.PlayerStats["ADR"]; ok {
						switch v := r.(type) {
						case float64:
							adr = v
						case int:
							adr = float64(v)
						case string:
							if parsed, err := strconv.ParseFloat(v, 64); err == nil {
								adr = parsed
							}
						}
					}
				}

				// Calculate K/D ratio
				kdRatio := 0.0
				if deaths > 0 {
					kdRatio = float64(kills) / float64(deaths)
				} else if kills > 0 {
					kdRatio = float64(kills)
				}

				playerStats := entity.PlayerMatchStats{
					PlayerID:            playerID,
					Nickname:            nickname,
					Team:                teamName,
					Kills:               kills,
					Deaths:              deaths,
					Assists:             assists,
					KDRatio:             kdRatio,
					HeadshotsPercentage: hsPerc,
					ADR:                 adr,
					HLTVRating:          0.0, // Not available in basic stats
					FirstKills:          0,   // Not available in basic stats
					FirstDeaths:         0,   // Not available in basic stats
					ClutchWins:          0,   // Not available in basic stats
					EntryFrags:          0,   // Not available in basic stats
					FlashAssists:        0,   // Not available in basic stats
					UtilityDamage:       0,   // Not available in basic stats
				}

				teamStats.Players = append(teamStats.Players, playerStats)
				matchStats.PlayerStats = append(matchStats.PlayerStats, playerStats)
			}

			if i == 0 {
				matchStats.Team1 = teamStats
			} else {
				matchStats.Team2 = teamStats
			}
		}
	}

	return matchStats, nil
}

