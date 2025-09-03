package entity

// GameDetail represents a subset of the information contained in the
// Faceit API's GameDetail model. It focuses on the attributes that are
// useful for presenting a player's standing within a specific game. The
// number of games tracked on FACEIT grows over time, so the map of games
// held by PlayerProfile may contain keys such as "cs2", "dota2" and
// others. Each key points to a GameDetail describing the player's ELO,
// skill level and region for that game.
type GameDetail struct {
	// Elo is the player's rating for the game on FACEIT. Higher values
	// represent better relative skill. When the Faceit API does not return
	// an elo value, this field is zero.
	Elo int
	// SkillLevel is a one‑to‑ten scale used by FACEIT to group players of
	// similar ability. A value of zero means the level was missing from
	// the response.
	SkillLevel int
	// Region is the geographic region associated with the player's
	// preferred servers for this game. Regions are typically codes like
	// "EU", "NA" or "SA".
	Region string
}

// PlayerProfile aggregates profile data returned from FACEIT for a user.
// It contains high‑level identity fields alongside per‑game details.
//
// The Games map uses the FACEIT game identifier (for example "cs2") as
// the key. Each value holds the player's rating and skill information for
// that game. If the player has not registered a particular game on
// FACEIT, the corresponding key will be absent from the map.
type PlayerProfile struct {
	ID        string
	Nickname  string
	Country   string
	Avatar    string
	FaceitURL string
	Games     map[string]GameDetail
}

// PlayerStats wraps the statistics returned from the Faceit API for a
// particular player and game. The API exposes lifetime statistics as a
// dynamic map whose keys and values depend on the game in question.
// Segments provide more granular statistics (for example per‑map stats)
// but are left unprocessed here to give consumers flexibility.
type PlayerStats struct {
	GameID   string
	PlayerID string
	Lifetime map[string]interface{}
	Segments []map[string]interface{}
}
