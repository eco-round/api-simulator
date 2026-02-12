package aggregator

import (
	"time"

	"api-simulator/models"
)

const ConsensusMin = 2

// Aggregate collects results from all sources and determines consensus.
// 2-of-3 rule: if at least 2 sources agree on a winner, that's the final result.
func Aggregate(matchID string, teamA, teamB models.Team, sources []models.SourceResult) models.AggregatedResult {
	votes := make(map[string]int)
	for _, s := range sources {
		if s.Confident {
			votes[s.Winner]++
		}
	}

	result := models.AggregatedResult{
		MatchID:      matchID,
		TeamA:        teamA,
		TeamB:        teamB,
		Sources:      sources,
		ConsensusMin: ConsensusMin,
	}

	// Find winner with most votes
	bestWinner := ""
	bestCount := 0
	for winner, count := range votes {
		if count > bestCount {
			bestWinner = winner
			bestCount = count
		}
	}

	if bestCount >= ConsensusMin {
		result.Consensus = true
		result.AgreedSources = bestCount
		result.FinalWinner = bestWinner
		result.Status = "resolved"
		now := time.Now()
		result.ResolvedAt = &now
	} else {
		result.Consensus = false
		result.AgreedSources = bestCount
		result.FinalWinner = "disputed"
		result.Status = "disputed"
	}

	return result
}
