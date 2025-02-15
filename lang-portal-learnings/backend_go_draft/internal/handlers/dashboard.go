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

	response := gin.H{
		"id":                session.ID,
		"group_id":          session.GroupID,
		"study_activity_id": session.StudyActivityID,
		"activity_name":     session.ActivityName,
		"group_name":        session.GroupName,
		"start_time":        session.StartTime,
		"end_time":          session.EndTime,
		"review_items_count": session.ReviewItemsCount,
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
		"total_words_studied":   studied,
		"total_available_words": total,
		"progress_percentage":   float64(0),
	}

	if total > 0 {
		response["progress_percentage"] = float64(studied) / float64(total) * 100
	}

	respondWithSuccess(c, response)
}

// QuickStats returns quick overview statistics
func (h *DashboardHandler) QuickStats(c *gin.Context) {
	// Get total words
	studied, totalWords, err := h.db.GetStudyProgress()
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Error fetching word count")
		return
	}

	// Get total groups
	_, totalGroups, err := h.db.ListGroups(1, 1)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Error fetching group count")
		return
	}

	// Get total study sessions
	_, totalSessions, err := h.db.ListStudySessions(1, 1)
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
		respondWithError(c, http.StatusInternalServerError, "Error fetching streak")
		return
	}

	response := gin.H{
		"total_words_studied":   studied,
		"total_words":          totalWords,
		"total_groups":         totalGroups,
		"total_study_sessions": totalSessions,
		"overall_accuracy":     accuracy,
		"study_streak":         streak,
	}

	respondWithSuccess(c, response)
}
