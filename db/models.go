package db

import "time"

// Match represents a Valorant match in the database.
type Match struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	TeamAName string    `gorm:"not null" json:"team_a_name"`
	TeamATag  string    `json:"team_a_tag"`
	TeamBName string    `gorm:"not null" json:"team_b_name"`
	TeamBTag  string    `json:"team_b_tag"`
	Status    string    `gorm:"default:upcoming" json:"status"` // upcoming, live, finished, cancelled
	BestOf    int       `gorm:"default:3" json:"best_of"`
	Event     string    `json:"event"`
	StartTime time.Time `json:"start_time"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Results []MatchResult `gorm:"foreignKey:MatchID" json:"results,omitempty"`
}

// MatchResult stores a result reported by a specific source.
// Each match can have up to 3 results (one per source).
type MatchResult struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	MatchID   uint      `gorm:"not null;index" json:"match_id"`
	Source    string    `gorm:"not null" json:"source"` // pandascore, vlr, liquipedia
	Winner   string    `gorm:"not null" json:"winner"` // TeamA or TeamB
	ScoreA   int       `json:"score_a"`
	ScoreB   int       `json:"score_b"`
	MapCount int       `json:"map_count"`
	ReportedAt time.Time `gorm:"autoCreateTime" json:"reported_at"`
}
