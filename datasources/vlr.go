package datasources

import (
	"fmt"
	"time"

	"api-simulator/db"
	"api-simulator/models"
)

// VLR simulates a VLR.gg scraper API for Valorant match results.
type VLR struct{}

func NewVLR() *VLR { return &VLR{} }

func (v *VLR) Name() string { return "vlr" }

func (v *VLR) GetMatch(id string) (*models.Match, error) {
	var m db.Match
	if err := db.DB.First(&m, id).Error; err != nil {
		return nil, fmt.Errorf("match %s not found", id)
	}
	match := toAPIMatch(m, "vlr")
	return &match, nil
}

func (v *VLR) GetResult(id string) (*models.SourceResult, error) {
	var m db.Match
	if err := db.DB.First(&m, id).Error; err != nil {
		return nil, fmt.Errorf("match %s not found", id)
	}

	var result db.MatchResult
	if err := db.DB.Where("match_id = ? AND source = ?", m.ID, "vlr").First(&result).Error; err != nil {
		return nil, fmt.Errorf("no vlr result for match %s", id)
	}

	return &models.SourceResult{
		Source:    "vlr",
		Winner:   result.Winner,
		Result:   models.MatchResult{Winner: result.Winner, ScoreA: result.ScoreA, ScoreB: result.ScoreB, MapCount: result.MapCount},
		Timestamp: time.Now(),
		Confident: true,
	}, nil
}
