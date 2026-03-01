package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"api-simulator/db"
)

type AdminHandler struct{}

func NewAdminHandler() *AdminHandler {
	return &AdminHandler{}
}

func (h *AdminHandler) RegisterRoutes(r *gin.Engine) {
	admin := r.Group("/api/v1/admin")

	admin.POST("/matches", h.CreateMatch)
	admin.GET("/matches", h.ListMatches)
	admin.GET("/matches/:id", h.GetMatch)
	admin.PATCH("/matches/:id", h.UpdateMatchStatus)
	admin.PATCH("/matches/:id/vault", h.UpdateMatchVault)
	admin.POST("/matches/:id/result", h.SetResult)
}

// ── Request Bodies ──────────────────────────────────────────────────────

type CreateMatchRequest struct {
	TeamAName      string    `json:"team_a_name" binding:"required"`
	TeamATag       string    `json:"team_a_tag"`
	TeamBName      string    `json:"team_b_name" binding:"required"`
	TeamBTag       string    `json:"team_b_tag"`
	BestOf         int       `json:"best_of"`
	Event          string    `json:"event"`
	StartTime      time.Time `json:"start_time"`
	OnChainMatchID uint      `json:"on_chain_match_id"`
	VaultAddress   string    `json:"vault_address"`
}

type UpdateStatusRequest struct {
	Status string `json:"status" binding:"required"` // upcoming, live, finished, cancelled
}

type UpdateVaultRequest struct {
	OnChainMatchID uint   `json:"on_chain_match_id" binding:"required"`
	VaultAddress   string `json:"vault_address" binding:"required"`
}

type SetResultRequest struct {
	Source      string `json:"source" binding:"required"`       // pandascore, vlr, liquipedia
	MatchStatus string `json:"match_status" binding:"required"` // upcoming, started, ended
	Winner      string `json:"winner"`                          // TeamA, TeamB (required when ended)
	ScoreA      int    `json:"score_a"`
	ScoreB      int    `json:"score_b"`
	MapCount    int    `json:"map_count"`
}

// ── Handlers ────────────────────────────────────────────────────────────

func (h *AdminHandler) CreateMatch(c *gin.Context) {
	var req CreateMatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	bestOf := req.BestOf
	if bestOf == 0 {
		bestOf = 3
	}

	match := db.Match{
		TeamAName:      req.TeamAName,
		TeamATag:       req.TeamATag,
		TeamBName:      req.TeamBName,
		TeamBTag:       req.TeamBTag,
		BestOf:         bestOf,
		Event:          req.Event,
		Status:         "open",
		StartTime:      req.StartTime,
		OnChainMatchID: req.OnChainMatchID,
		VaultAddress:   req.VaultAddress,
	}

	if err := db.DB.Create(&match).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, match)
}

func (h *AdminHandler) ListMatches(c *gin.Context) {
	var matches []db.Match
	query := db.DB.Preload("Results").Order("id DESC")

	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Find(&matches).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"matches": matches})
}

func (h *AdminHandler) GetMatch(c *gin.Context) {
	id := c.Param("id")
	var match db.Match
	if err := db.DB.Preload("Results").First(&match, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "match not found"})
		return
	}
	c.JSON(http.StatusOK, match)
}

func (h *AdminHandler) UpdateMatchStatus(c *gin.Context) {
	id := c.Param("id")
	var req UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	validStatuses := map[string]bool{
		"open": true, "locked": true, "finished": true, "cancelled": true,
	}
	if !validStatuses[req.Status] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid status, must be: open, locked, finished, cancelled"})
		return
	}

	var match db.Match
	if err := db.DB.First(&match, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "match not found"})
		return
	}

	match.Status = req.Status
	db.DB.Save(&match)

	c.JSON(http.StatusOK, match)
}

func (h *AdminHandler) SetResult(c *gin.Context) {
	id := c.Param("id")
	var req SetResultRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate source
	// Validate match_status
	validStatuses := map[string]bool{
		"upcoming": true, "started": true, "ended": true,
	}
	if !validStatuses[req.MatchStatus] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid match_status, must be: upcoming, started, ended"})
		return
	}

	// Validate winner only if ended
	if req.MatchStatus == "ended" {
		if req.Winner != "TeamA" && req.Winner != "TeamB" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "winner must be TeamA or TeamB when match_status is ended"})
			return
		}
	}

	// Check match exists
	var match db.Match
	if err := db.DB.First(&match, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "match not found"})
		return
	}

	// Upsert: delete existing result for this source, then create new
	db.DB.Where("match_id = ? AND source = ?", match.ID, req.Source).Delete(&db.MatchResult{})

	result := db.MatchResult{
		MatchID:     match.ID,
		Source:      req.Source,
		MatchStatus: req.MatchStatus,
		Winner:      req.Winner,
		ScoreA:      req.ScoreA,
		ScoreB:      req.ScoreB,
		MapCount:    req.MapCount,
	}

	if err := db.DB.Create(&result).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, result)
}

func (h *AdminHandler) UpdateMatchVault(c *gin.Context) {
	id := c.Param("id")
	var req UpdateVaultRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var match db.Match
	if err := db.DB.First(&match, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "match not found"})
		return
	}

	match.OnChainMatchID = req.OnChainMatchID
	match.VaultAddress = req.VaultAddress
	db.DB.Save(&match)

	c.JSON(http.StatusOK, match)
}
