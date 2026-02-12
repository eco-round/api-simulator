package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"api-simulator/aggregator"
	"api-simulator/datasources"
	"api-simulator/models"
)

type MatchHandler struct {
	pandaScore *datasources.PandaScore
	vlr        *datasources.VLR
	liquipedia *datasources.Liquipedia
}

func NewMatchHandler() *MatchHandler {
	return &MatchHandler{
		pandaScore: datasources.NewPandaScore(),
		vlr:        datasources.NewVLR(),
		liquipedia: datasources.NewLiquipedia(),
	}
}

func (h *MatchHandler) RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api/v1")

	// PandaScore endpoints
	api.GET("/pandascore/matches", h.ListPandaScoreMatches)
	api.GET("/pandascore/matches/:id", h.GetPandaScoreMatch)

	// VLR endpoints
	api.GET("/vlr/matches/:id", h.GetVLRMatch)

	// Liquipedia endpoints
	api.GET("/liquipedia/matches/:id", h.GetLiquipediaMatch)

	// Aggregated oracle endpoint (what CRE workflow calls)
	api.GET("/aggregated/matches/:id", h.GetAggregatedResult)
}

// ── PandaScore ───────────────────────────────────────────────────────────

func (h *MatchHandler) ListPandaScoreMatches(c *gin.Context) {
	status := c.Query("status")
	matches := h.pandaScore.ListMatches(status)
	c.JSON(http.StatusOK, gin.H{"source": "pandascore", "matches": matches})
}

func (h *MatchHandler) GetPandaScoreMatch(c *gin.Context) {
	id := c.Param("id")
	result, err := h.pandaScore.GetResult(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

// ── VLR ──────────────────────────────────────────────────────────────────

func (h *MatchHandler) GetVLRMatch(c *gin.Context) {
	id := c.Param("id")
	result, err := h.vlr.GetResult(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

// ── Liquipedia ───────────────────────────────────────────────────────────

func (h *MatchHandler) GetLiquipediaMatch(c *gin.Context) {
	id := c.Param("id")
	result, err := h.liquipedia.GetResult(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

// ── Aggregated (CRE calls this) ──────────────────────────────────────────

func (h *MatchHandler) GetAggregatedResult(c *gin.Context) {
	id := c.Param("id")

	// Fetch from all 3 sources
	var sources []models.SourceResult

	if ps, err := h.pandaScore.GetResult(id); err == nil {
		sources = append(sources, *ps)
	}
	if vlr, err := h.vlr.GetResult(id); err == nil {
		sources = append(sources, *vlr)
	}
	if liq, err := h.liquipedia.GetResult(id); err == nil {
		sources = append(sources, *liq)
	}

	if len(sources) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "no sources returned data for this match"})
		return
	}

	// Get match info for team names
	match, _ := h.pandaScore.GetMatch(id)
	teamA := models.Team{Name: "Team A"}
	teamB := models.Team{Name: "Team B"}
	if match != nil {
		teamA = match.TeamA
		teamB = match.TeamB
	}

	result := aggregator.Aggregate(id, teamA, teamB, sources)
	c.JSON(http.StatusOK, result)
}
