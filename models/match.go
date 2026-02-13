package models

import "time"

type Team struct {
	Name    string `json:"name"`
	Tag     string `json:"tag"`
	LogoURL string `json:"logo_url,omitempty"`
}

type MatchResult struct {
	Winner   string `json:"winner"`
	ScoreA   int    `json:"score_a"`
	ScoreB   int    `json:"score_b"`
	MapCount int    `json:"map_count"`
}

type Match struct {
	ID         string       `json:"id"`
	ExternalID string       `json:"external_id,omitempty"`
	TeamA      Team         `json:"team_a"`
	TeamB      Team         `json:"team_b"`
	Status     string       `json:"status"` // upcoming, live, finished, cancelled
	BestOf     int          `json:"best_of"`
	Result     *MatchResult `json:"result,omitempty"`
	StartTime  time.Time    `json:"start_time"`
	Event      string       `json:"event"`
}

type SourceResult struct {
	Source    string      `json:"source"`
	Winner   string      `json:"winner"`
	Result   MatchResult `json:"result"`
	Timestamp time.Time  `json:"timestamp"`
	Confident bool       `json:"confident"`
}
