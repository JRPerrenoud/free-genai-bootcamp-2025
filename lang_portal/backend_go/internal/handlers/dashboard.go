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
	correct, wrong, err := h.db.GetStudySessionStats(session.ID)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Error fetching session statistics")
		return
	}

	response := gin.H{
		"id":                session.ID,
		"group_id":          session.GroupID,
		"created_at":        session.CreatedAt,
		"study_activity_id": session.StudyActivityID,
		"group_name":        session.GroupName,
		"correct_count":     correct,
		"wrong_count":       wrong,
	}

	respondWithSuccess(c, response)
}

// StudyProgress returns the overall study progress
func (h *DashboardHandler) StudyProgress(c *gin.Context) {
	studied, total, err := h.db.GetStudyProgress()
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Error fetching study progress")
		return
	}

	response := gin.H{
		"total_words_studied":    studied,
		"total_available_words": total,
	}

	respondWithSuccess(c, response)
}

// QuickStats returns quick overview statistics
func (h *DashboardHandler) QuickStats(c *gin.Context) {
	// Get total words
	_, totalWords, err := h.db.GetStudyProgress()
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Error fetching word count")
		return
	}

	// Get total groups
	groups, _, err := h.db.ListGroups(1, 1)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Error fetching group count")
		return
	}

	// Get total study sessions
	sessions, _, err := h.db.ListStudySessions(1, 1)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Error fetching session count")
		return
	}

	// Get overall accuracy
	accuracy, err := h.db.GetOverallAccuracy()
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Error fetching accuracy")
		return
	}

	// Get study streak
	streak, err := h.db.GetStudyStreak()
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Error fetching study streak")
		return
	}

	response := gin.H{
		"total_words":         totalWords,
		"total_groups":        len(groups),
		"total_study_sessions": len(sessions),
		"overall_accuracy":    accuracy,
		"study_streak":        streak,
	}

	respondWithSuccess(c, response)
}
