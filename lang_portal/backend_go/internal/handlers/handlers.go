package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/JPerreno/free-genai-bootcamp-2025/lang_portal/backend_go/internal/models"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	db *models.DB
}

func NewHandler(db *models.DB) *Handler {
	return &Handler{db: db}
}

// GetLastStudySession handles GET /api/dashboard/last_study_session
func (h *Handler) GetLastStudySession(c *gin.Context) {
	// Get groupID from query parameter, default to 0 if not provided
	groupID := 0
	groupIDStr := c.Query("group_id")
	if groupIDStr != "" {
		var err error
		groupID, err = strconv.Atoi(groupIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, models.Response{
				Success: false,
				Error:   "invalid group_id",
			})
			return
		}
	}

	session, err := h.db.GetLastStudySession(groupID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	if session == nil {
		c.JSON(http.StatusNotFound, models.Response{
			Success: false,
			Error:   "no study sessions found",
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Data:    session,
	})
}

// GetStudyProgress handles GET /api/dashboard/study_progress
func (h *Handler) GetStudyProgress(c *gin.Context) {
	// Get groupID from query parameter, default to 0 if not provided
	groupID := 0
	groupIDStr := c.Query("group_id")
	if groupIDStr != "" {
		var err error
		groupID, err = strconv.Atoi(groupIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, models.Response{
				Success: false,
				Error:   "invalid group_id",
			})
			return
		}
	}

	progress, err := h.db.GetStudyProgress(groupID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Data:    progress,
	})
}

// GetQuickStats handles GET /api/dashboard/quick_stats
func (h *Handler) GetQuickStats(c *gin.Context) {
	// Get groupID from query parameter, default to 0 if not provided
	groupID := 0
	groupIDStr := c.Query("group_id")
	if groupIDStr != "" {
		var err error
		groupID, err = strconv.Atoi(groupIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, models.Response{
				Success: false,
				Error:   "invalid group_id",
			})
			return
		}
	}

	stats, err := h.db.GetDashboardStats(groupID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Error:   "Failed to get dashboard stats",
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Data:    stats,
	})
}

// GetStudyActivities handles GET /api/study_activities
func (h *Handler) GetStudyActivities(c *gin.Context) {
	// Get groupID from query parameter, default to 0 if not provided
	groupID := 0
	groupIDStr := c.Query("group_id")
	if groupIDStr != "" {
		var err error
		groupID, err = strconv.Atoi(groupIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, models.Response{
				Success: false,
				Error:   "invalid group_id",
			})
			return
		}
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}

	response, err := h.db.GetStudyActivities(groupID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Error:   "Failed to get study activities",
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Data:    response,
	})
}

// GetWords handles GET /api/words
func (h *Handler) GetWords(c *gin.Context) {
	// Get groupID from query parameter, default to 0 if not provided
	groupID := 0
	groupIDStr := c.Query("group_id")
	if groupIDStr != "" {
		var err error
		groupID, err = strconv.Atoi(groupIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, models.Response{
				Success: false,
				Error:   "invalid group_id",
			})
			return
		}
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "100"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 100
	}

	response, err := h.db.GetWords(groupID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Error:   "Failed to get words",
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Data:    response,
	})
}

// GetWord handles GET /api/words/:id
func (h *Handler) GetWord(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Error:   "Invalid word ID",
		})
		return
	}

	word, err := h.db.GetWordByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Error:   "Failed to get word",
		})
		return
	}

	if word == nil {
		c.JSON(http.StatusNotFound, models.Response{
			Success: false,
			Error:   "Word not found",
		})
		return
	}

	// Get groupID from query parameter, default to 0 if not provided
	groupID := 0
	groupIDStr := c.Query("group_id")
	if groupIDStr != "" {
		var err error
		groupID, err = strconv.Atoi(groupIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, models.Response{
				Success: false,
				Error:   "invalid group_id",
			})
			return
		}
	}

	stats, err := h.db.GetWordStats(id, groupID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Error:   "Failed to get word stats",
		})
		return
	}

	response := map[string]interface{}{
		"word":  word,
		"stats": stats,
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Data:    response,
	})
}

// CreateWord handles POST /api/words
func (h *Handler) CreateWord(c *gin.Context) {
	var word models.Word
	if err := c.ShouldBindJSON(&word); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Error:   "Invalid request body: " + err.Error(),
		})
		return
	}

	if err := h.db.CreateWord(&word); err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Error:   "Failed to create word: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, models.Response{
		Success: true,
		Data:    word,
	})
}

// UpdateWord handles PUT /api/words/:id
func (h *Handler) UpdateWord(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Error:   "Invalid word ID",
		})
		return
	}

	var word models.Word
	if err := c.ShouldBindJSON(&word); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Error:   "Invalid request body: " + err.Error(),
		})
		return
	}
	word.ID = id

	if err := h.db.UpdateWord(&word); err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Error:   "Failed to update word: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Data:    word,
	})
}

// DeleteWord handles DELETE /api/words/:id
func (h *Handler) DeleteWord(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Error:   "Invalid word ID",
		})
		return
	}

	if err := h.db.DeleteWord(id); err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Error:   "Failed to delete word: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
	})
}

// GetGroups handles GET /api/groups
func (h *Handler) GetGroups(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}

	response, err := h.db.GetGroups(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Error:   "Failed to get groups",
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Data:    response,
	})
}

// GetGroup handles GET /api/groups/:id
func (h *Handler) GetGroup(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Error:   "Invalid group ID",
		})
		return
	}

	group, err := h.db.GetGroupByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Error:   "Failed to get group",
		})
		return
	}

	if group == nil {
		c.JSON(http.StatusNotFound, models.Response{
			Success: false,
			Error:   "Group not found",
		})
		return
	}

	words, err := h.db.GetGroupWords(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Error:   "Failed to get group words",
		})
		return
	}

	response := map[string]interface{}{
		"group": group,
		"words": words,
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Data:    response,
	})
}

// CreateGroup handles POST /api/groups
func (h *Handler) CreateGroup(c *gin.Context) {
	var group models.Group
	if err := c.ShouldBindJSON(&group); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Error:   "Invalid request body: " + err.Error(),
		})
		return
	}

	if err := h.db.CreateGroup(&group); err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Error:   "Failed to create group: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, models.Response{
		Success: true,
		Data:    group,
	})
}

// UpdateGroup handles PUT /api/groups/:id
func (h *Handler) UpdateGroup(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Error:   "Invalid group ID",
		})
		return
	}

	var group models.Group
	if err := c.ShouldBindJSON(&group); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Error:   "Invalid request body: " + err.Error(),
		})
		return
	}
	group.ID = id

	if err := h.db.UpdateGroup(&group); err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Error:   "Failed to update group: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Data:    group,
	})
}

// DeleteGroup handles DELETE /api/groups/:id
func (h *Handler) DeleteGroup(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Error:   "Invalid group ID",
		})
		return
	}

	if err := h.db.DeleteGroup(id); err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Error:   "Failed to delete group: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
	})
}

// AddWordsToGroup handles POST /api/groups/:id/words
func (h *Handler) AddWordsToGroup(c *gin.Context) {
	groupID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Error:   "Invalid group ID",
		})
		return
	}

	var req struct {
		WordIDs []int `json:"word_ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Error:   "Invalid request body: " + err.Error(),
		})
		return
	}

	if err := h.db.AddWordsToGroup(groupID, req.WordIDs); err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Error:   "Failed to add words to group: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
	})
}

// RemoveWordsFromGroup handles DELETE /api/groups/:id/words
func (h *Handler) RemoveWordsFromGroup(c *gin.Context) {
	groupID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Error:   "Invalid group ID",
		})
		return
	}

	var req struct {
		WordIDs []int `json:"word_ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Error:   "Invalid request body: " + err.Error(),
		})
		return
	}

	if err := h.db.RemoveWordsFromGroup(groupID, req.WordIDs); err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Error:   "Failed to remove words from group: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
	})
}

// CreateStudySession handles POST /api/study_sessions
func (h *Handler) CreateStudySession(c *gin.Context) {
	var session models.StudySession
	if err := c.ShouldBindJSON(&session); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	err := h.db.CreateStudySession(&session)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, models.Response{
		Success: true,
		Data:    session,
	})
}

// GetStudySession handles GET /api/study_sessions/:id
func (h *Handler) GetStudySession(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Error:   "invalid session ID",
		})
		return
	}

	session, err := h.db.GetStudySession(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}
	if session == nil {
		c.JSON(http.StatusNotFound, models.Response{
			Success: false,
			Error:   "session not found",
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Data:    session,
	})
}

// AddReviewItems handles POST /api/study_sessions/:id/review_items
func (h *Handler) AddReviewItems(c *gin.Context) {
	sessionID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Error:   "invalid session ID",
		})
		return
	}

	var req struct {
		WordIDs []int `json:"word_ids" binding:"required"`
		Correct bool  `json:"correct"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	err = h.db.AddReviewItems(sessionID, req.WordIDs, req.Correct)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
	})
}

// GetStudyActivityByID handles GET /api/study_activities/:id
func (h *Handler) GetStudyActivityByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Error:   "invalid activity id",
		})
		return
	}

	activity, err := h.db.GetStudyActivityByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Error:   fmt.Sprintf("error getting study activity: %v", err),
		})
		return
	}

	if activity == nil {
		c.JSON(http.StatusNotFound, models.Response{
			Success: false,
			Error:   "study activity not found",
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Data:    activity,
	})
}

// GetStudyActivitySessions handles GET /api/study_activities/:id/sessions
func (h *Handler) GetStudyActivitySessions(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Error:   "invalid activity id",
		})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}

	response, err := h.db.GetStudyActivitySessions(id, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Error:   fmt.Sprintf("error getting study sessions: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Data:    response,
	})
}

// CreateStudyActivity handles POST /api/study_activities
func (h *Handler) CreateStudyActivity(c *gin.Context) {
	var activity models.StudyActivity
	if err := c.ShouldBindJSON(&activity); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Error:   "invalid request body",
		})
		return
	}

	if activity.Name == "" {
		c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Error:   "activity name is required",
		})
		return
	}

	if activity.LaunchURL == "" {
		c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Error:   "launch URL is required",
		})
		return
	}

	err := h.db.CreateStudyActivity(&activity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Error:   fmt.Sprintf("error creating study activity: %v", err),
		})
		return
	}

	c.JSON(http.StatusCreated, models.Response{
		Success: true,
		Data:    activity,
	})
}
