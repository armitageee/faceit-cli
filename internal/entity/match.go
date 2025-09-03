package entity

// PlayerMatchSummary captures a player's performance in a single
// match.  It is used by the delivery layer to present a concise
// overview of recent matches to end users.  Fields are optional
// depending on the availability of data from the FACEIT API.
type PlayerMatchSummary struct {
	// MatchID uniquely identifies the match on FACEIT.
	MatchID string
	// Map identifies the game mode or map name.  The FACEIT API uses
	// this field to describe the map played (for example
	// "de_inferno").
	Map string
	// FinishedAt holds the UNIX timestamp at which the match
	// concluded.  A value of zero means the timestamp was
	// unavailable.
	FinishedAt int64
	// Score is a string representation of the final score of the
	// match, for example "16-14".  When unavailable this field
	// will be empty.
	Score string
	// Kills is the number of kills the player achieved.
	Kills int
	// Deaths is the number of deaths the player suffered.
	Deaths int
	// Assists is the number of assists the player recorded.
	Assists int
	// KDRatio is the kill/death ratio for the match.  A zero value
	// indicates it could not be parsed from the API response.
	KDRatio float64
	// HeadshotsPercentage represents the percentage of kills that were
	// headshots.  A value of zero means the data was unavailable.
	HeadshotsPercentage float64
	// Result is "Win" when the player's team won the match and
	// "Loss" otherwise.
	Result string
}
