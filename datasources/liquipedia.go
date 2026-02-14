package datasources

import (
	"fmt"
	"time"

	"api-simulator/db"
	"api-simulator/models"
)

// Liquipedia simulates a Liquipedia API for Valorant match results.
type Liquipedia struct{}

func NewLiquipedia() *Liquipedia { return &Liquipedia{} }

func (l *Liquipedia) Name() string { return "liquipedia" }

func (l *Liquipedia) GetMatch(id string) (*models.Match, error) {
	var m db.Match
	if err := db.DB.First(&m, id).Error; err != nil {
		return nil, fmt.Errorf("match %s not found", id)
	}
	match := toAPIMatch(m, "liquipedia")
	return &match, nil
}

func (l *Liquipedia) GetResult(id string) (*models.SourceResult, error) {
	var m db.Match
	if err := db.DB.First(&m, id).Error; err != nil {
		return nil, fmt.Errorf("match %s not found", id)
	}

	var result db.MatchResult
	if err := db.DB.Where("match_id = ? AND source = ?", m.ID, "liquipedia").First(&result).Error; err != nil {
		return nil, fmt.Errorf("no liquipedia result for match %s", id)
	}

	return &models.SourceResult{
		Source:      "liquipedia",
		MatchStatus: result.MatchStatus,
		Winner:      result.Winner,
		Result:      models.MatchResult{MatchStatus: result.MatchStatus, Winner: result.Winner, ScoreA: result.ScoreA, ScoreB: result.ScoreB, MapCount: result.MapCount},
		Timestamp:   time.Now(),
		Confident:   true,
	}, nil
}
