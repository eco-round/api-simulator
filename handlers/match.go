package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"api-simulator/datasources"
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
