package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"lang_portal_go/internal/models"
)

// DashboardHandler handles dashboard-related routes
type DashboardHandler struct {
	db *models.DB
}

// NewDashboardHandler creates a new dashboard handler
func NewDashboardHandler(db *models.DB) *DashboardHandler {
	return &DashboardHandler{db: db}
}

// LastStudySession returns information about the most recent study session
func (h *DashboardHandler) LastStudySession(c *gin.Context) {
	session, err := h.db.GetLastStudySession()
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Error fetching last study session")
		return
	}

	if session == nil {
		respondWithSuccess(c, nil)
		return
	}

	// Get statistics for this session
	totalWords, correctWords, err := h.db.GetStudySessionStats(session.ID)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Error fetching session statistics")
		return
	}

	response := gin.H{
		"id":                session.ID,
		"part_of_speech":    session.PartOfSpeech,
		"created_at":        session.CreatedAt,
		"study_activity_id": session.StudyActivityID,
		"total_words":       totalWords,
		"correct_words":     correctWords,
	}

	respondWithSuccess(c, response)
}

// StudyProgress returns the overall study progress
func (h *DashboardHandler) StudyProgress(c *gin.Context) {
	totalWords, reviewedWords, err := h.db.GetStudyProgress()
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Error fetching study progress")
		return
	}

	response := gin.H{
		"total_words":     totalWords,
		"reviewed_words":  reviewedWords,
		"completion_rate": float64(reviewedWords) / float64(totalWords),
	}

	respondWithSuccess(c, response)
}

// QuickStats returns quick overview statistics
func (h *DashboardHandler) QuickStats(c *gin.Context) {
	// Get total words count
	totalWords, reviewedWords, err := h.db.GetStudyProgress()
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Error fetching word counts")
		return
	}

	// Get groups count and stats
	groups, err := h.db.ListGroups()
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Error fetching groups")
		return
	}

	// Get last study session
	lastSession, err := h.db.GetLastStudySession()
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Error fetching last session")
		return
	}

	var lastSessionStats gin.H
	if lastSession != nil {
		totalWords, correctWords, err := h.db.GetStudySessionStats(lastSession.ID)
		if err != nil {
			respondWithError(c, http.StatusInternalServerError, "Error fetching session stats")
			return
		}

		lastSessionStats = gin.H{
			"id":                lastSession.ID,
			"part_of_speech":    lastSession.PartOfSpeech,
			"created_at":        lastSession.CreatedAt,
			"study_activity_id": lastSession.StudyActivityID,
			"total_words":       totalWords,
			"correct_words":     correctWords,
		}
	}

	response := gin.H{
		"total_words":      totalWords,
		"reviewed_words":   reviewedWords,
		"completion_rate":  float64(reviewedWords) / float64(totalWords),
		"groups":          groups,
		"last_session":    lastSessionStats,
	}

	respondWithSuccess(c, response)
}
