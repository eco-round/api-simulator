package handlers

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"api-simulator/datasources"
)

type MatchHandler struct {
	pandaScore    *datasources.PandaScore
	vlr           *datasources.VLR
	liquipedia    *datasources.Liquipedia
	pandaKey      string
	vlrKey        string
	liquipediaKey string
}

func NewMatchHandler() *MatchHandler {
	return &MatchHandler{
		pandaScore: datasources.NewPandaScore(),
		vlr:        datasources.NewVLR(),
		liquipedia: datasources.NewLiquipedia(),
		// API keys loaded from env — must match the keys stored in CRE DON Vault
		pandaKey:      os.Getenv("PANDASCORE_API_KEY"),
		vlrKey:        os.Getenv("VLR_API_KEY"),
		liquipediaKey: os.Getenv("LIQUIPEDIA_API_KEY"),
	}
}

func (h *MatchHandler) RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api/v1")

	// PandaScore endpoints — protected by API key (simulates real PandaScore auth)
	ps := api.Group("/pandascore", h.requireKey("PANDASCORE", h.pandaKey))
	ps.GET("/matches", h.ListPandaScoreMatches)
	ps.GET("/matches/:id", h.GetPandaScoreMatch)

	// VLR endpoints — protected by API key
	vlr := api.Group("/vlr", h.requireKey("VLR", h.vlrKey))
	vlr.GET("/matches/:id", h.GetVLRMatch)

	// Liquipedia endpoints — protected by API key
	lq := api.Group("/liquipedia", h.requireKey("LIQUIPEDIA", h.liquipediaKey))
	lq.GET("/matches/:id", h.GetLiquipediaMatch)
}

// requireKey returns a Gin middleware that validates the X-Api-Key header.
// If no key is configured (empty env var), validation is skipped — dev mode.
// This simulates real API authentication that the CRE oracle bypasses using
// Chainlink Confidential HTTP with secrets injected inside the enclave.
func (h *MatchHandler) requireKey(source, expectedKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if expectedKey == "" {
			// No key configured — open access (dev/local mode without auth)
			c.Next()
			return
		}

		provided := c.GetHeader("X-Api-Key")
		if provided == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":  "missing X-Api-Key header",
				"source": source,
			})
			c.Abort()
			return
		}

		if provided != expectedKey {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":  "invalid API key",
				"source": source,
			})
			c.Abort()
			return
		}

		c.Next()
	}
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
