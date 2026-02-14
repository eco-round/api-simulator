package datasources

import (
	"fmt"
	"time"

	"api-simulator/db"
	"api-simulator/models"
)

// PandaScore simulates the PandaScore REST API for Valorant matches.
type PandaScore struct{}

func NewPandaScore() *PandaScore { return &PandaScore{} }

func (p *PandaScore) Name() string { return "pandascore" }

func (p *PandaScore) ListMatches(status string) []models.Match {
	var dbMatches []db.Match
	query := db.DB.Order("id DESC")
	if status != "" {
		query = query.Where("status = ?", status)
	}
	query.Find(&dbMatches)

	var results []models.Match
	for _, m := range dbMatches {
		results = append(results, toAPIMatch(m, "pandascore"))
	}
	return results
}

func (p *PandaScore) GetMatch(id string) (*models.Match, error) {
	var m db.Match
	if err := db.DB.First(&m, id).Error; err != nil {
		return nil, fmt.Errorf("match %s not found", id)
	}
	match := toAPIMatch(m, "pandascore")
	return &match, nil
}

func (p *PandaScore) GetResult(id string) (*models.SourceResult, error) {
	var m db.Match
	if err := db.DB.First(&m, id).Error; err != nil {
		return nil, fmt.Errorf("match %s not found", id)
	}

	var result db.MatchResult
	if err := db.DB.Where("match_id = ? AND source = ?", m.ID, "pandascore").First(&result).Error; err != nil {
		return nil, fmt.Errorf("no pandascore result for match %s", id)
	}

	return &models.SourceResult{
		Source:      "pandascore",
		MatchStatus: result.MatchStatus,
		Winner:      result.Winner,
		Result:      models.MatchResult{MatchStatus: result.MatchStatus, Winner: result.Winner, ScoreA: result.ScoreA, ScoreB: result.ScoreB, MapCount: result.MapCount},
		Timestamp:   time.Now(),
		Confident:   true,
	}, nil
}
