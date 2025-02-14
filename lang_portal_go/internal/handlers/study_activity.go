package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"lang_portal_go/internal/models"
)

// StudyActivityHandler handles study activity-related routes
type StudyActivityHandler struct {
	db *models.DB
}

// NewStudyActivityHandler creates a new study activity handler
func NewStudyActivityHandler(db *models.DB) *StudyActivityHandler {
	return &StudyActivityHandler{db: db}
}

// ListStudyActivities returns a paginated list of study activities
func (h *StudyActivityHandler) ListStudyActivities(c *gin.Context) {
	page, pageSize := getPaginationParams(c, 20)

	activities, total, err := h.db.ListStudyActivities(page, pageSize)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Error fetching study activities")
		return
	}

	response := wrapWithPagination(activities, page, total, pageSize)
	respondWithSuccess(c, response)
}

// GetStudyActivity returns details for a specific study activity
func (h *StudyActivityHandler) GetStudyActivity(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid study activity ID")
		return
	}

	activity, err := h.db.GetStudyActivity(id)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Error fetching study activity")
		return
	}

	if activity == nil {
		respondWithError(c, http.StatusNotFound, "Study activity not found")
		return
	}

	respondWithSuccess(c, activity)
}

// CreateStudyActivity creates a new study activity
func (h *StudyActivityHandler) CreateStudyActivity(c *gin.Context) {
	var activity models.StudyActivity
	if err := c.ShouldBindJSON(&activity); err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid study activity data")
		return
	}

	if err := h.db.CreateStudyActivity(&activity); err != nil {
		respondWithError(c, http.StatusInternalServerError, "Error creating study activity")
		return
	}

	respondWithSuccess(c, activity)
}

// UpdateStudyActivity updates an existing study activity
func (h *StudyActivityHandler) UpdateStudyActivity(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid study activity ID")
		return
	}

	var activity models.StudyActivity
	if err := c.ShouldBindJSON(&activity); err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid study activity data")
		return
	}

	activity.ID = id
	if err := h.db.UpdateStudyActivity(&activity); err != nil {
		respondWithError(c, http.StatusInternalServerError, "Error updating study activity")
		return
	}

	respondWithSuccess(c, activity)
}

// DeleteStudyActivity deletes a study activity
func (h *StudyActivityHandler) DeleteStudyActivity(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid study activity ID")
		return
	}

	if err := h.db.DeleteStudyActivity(id); err != nil {
		respondWithError(c, http.StatusInternalServerError, "Error deleting study activity")
		return
	}

	respondWithSuccess(c, gin.H{"message": "Study activity deleted successfully"})
}

// GetStudyActivitySessions returns all study sessions for a specific activity
func (h *StudyActivityHandler) GetStudyActivitySessions(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid study activity ID")
		return
	}

	page, pageSize := getPaginationParams(c, 20)

	sessions, total, err := h.db.GetStudyActivitySessions(id, page, pageSize)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Error fetching study sessions")
		return
	}

	response := wrapWithPagination(sessions, page, total, pageSize)
	respondWithSuccess(c, response)
}

// StartStudySession starts a new study session for an activity
func (h *StudyActivityHandler) StartStudySession(c *gin.Context) {
	idStr := c.Param("id")
	activityID, err := strconv.Atoi(idStr)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid study activity ID")
		return
	}

	var req struct {
		GroupID int `json:"group_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid request data")
		return
	}

	session, err := h.db.CreateStudySession(req.GroupID, activityID)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Error starting study session")
		return
	}

	respondWithSuccess(c, session)
}

// RecordStudyResult records a study result for a session
func (h *StudyActivityHandler) RecordStudyResult(c *gin.Context) {
	sessionIDStr := c.Param("session_id")
	sessionID, err := strconv.Atoi(sessionIDStr)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid session ID")
		return
	}

	var req struct {
		WordID  int  `json:"word_id" binding:"required"`
		Correct bool `json:"correct" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid request data")
		return
	}

	reviewItem, err := h.db.CreateWordReviewItem(req.WordID, sessionID, req.Correct)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Error recording study result")
		return
	}

	respondWithSuccess(c, reviewItem)
}
