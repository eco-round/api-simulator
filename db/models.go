package db

import "time"

// Match represents a Valorant match in the database.
type Match struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	OnChainMatchID uint      `json:"on_chain_match_id"`
	VaultAddress   string    `json:"vault_address"`
	TeamAName      string    `gorm:"not null" json:"team_a_name"`
	TeamATag       string    `json:"team_a_tag"`
	TeamBName      string    `gorm:"not null" json:"team_b_name"`
	TeamBTag       string    `json:"team_b_tag"`
	Status         string    `gorm:"default:open" json:"status"` // open, locked, finished, cancelled
	BestOf         int       `gorm:"default:3" json:"best_of"`
	Event          string    `json:"event"`
	StartTime      time.Time `json:"start_time"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`

	Results []MatchResult `gorm:"foreignKey:MatchID" json:"results,omitempty"`
}

// MatchResult stores a result reported by a specific source.
// Each match can have up to 3 results (one per source).
type MatchResult struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	MatchID     uint      `gorm:"not null;index" json:"match_id"`
	Source      string    `gorm:"not null" json:"source"`       // pandascore, vlr, liquipedia
	MatchStatus string    `gorm:"default:upcoming" json:"match_status"` // upcoming, started, ended
	Winner      string    `json:"winner"`                       // TeamA or TeamB (only set when ended)
	ScoreA      int       `json:"score_a"`
	ScoreB      int       `json:"score_b"`
	MapCount    int       `json:"map_count"`
	ReportedAt  time.Time `gorm:"autoCreateTime" json:"reported_at"`
}
