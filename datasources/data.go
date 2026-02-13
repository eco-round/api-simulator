package datasources

import (
	"fmt"

	"api-simulator/db"
	"api-simulator/models"
)

// toAPIMatch converts a DB match to an API response match.
func toAPIMatch(m db.Match, source string) models.Match {
	apiMatch := models.Match{
		ID:         fmt.Sprintf("%d", m.ID),
		ExternalID: fmt.Sprintf("%s-%d", source, m.ID),
		TeamA:      models.Team{Name: m.TeamAName, Tag: m.TeamATag},
		TeamB:      models.Team{Name: m.TeamBName, Tag: m.TeamBTag},
		Status:     m.Status,
		BestOf:     m.BestOf,
		Event:      m.Event,
		StartTime:  m.StartTime,
	}

	// Attach result from this source if available
	var result db.MatchResult
	if err := db.DB.Where("match_id = ? AND source = ?", m.ID, source).First(&result).Error; err == nil {
		apiMatch.Result = &models.MatchResult{
			Winner:   result.Winner,
			ScoreA:   result.ScoreA,
			ScoreB:   result.ScoreB,
			MapCount: result.MapCount,
		}
	}

	return apiMatch
}
