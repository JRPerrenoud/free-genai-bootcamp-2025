package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"lang_portal_go/internal/models"
)

// GroupHandler handles group-related routes
type GroupHandler struct {
	db *models.DB
}

// NewGroupHandler creates a new group handler
func NewGroupHandler(db *models.DB) *GroupHandler {
	return &GroupHandler{db: db}
}

// ListGroups returns a paginated list of groups
func (h *GroupHandler) ListGroups(c *gin.Context) {
	page, pageSize := getPaginationParams(c, 20)

	groups, total, err := h.db.ListGroups(page, pageSize)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Error fetching groups")
		return
	}

	response := wrapWithPagination(groups, page, total, pageSize)
	respondWithSuccess(c, response)
}

// GetGroup returns details for a specific group
func (h *GroupHandler) GetGroup(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid group ID")
		return
	}

	group, err := h.db.GetGroup(id)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Error fetching group")
		return
	}

	if group == nil {
		respondWithError(c, http.StatusNotFound, "Group not found")
		return
	}

	// Get words in this group
	words, total, err := h.db.GetGroupWords(id, 1, 1000)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Error fetching group words")
		return
	}

	response := gin.H{
		"group": group,
		"words": words,
		"total_words": total,
	}

	respondWithSuccess(c, response)
}

// CreateGroup creates a new group
func (h *GroupHandler) CreateGroup(c *gin.Context) {
	var group models.Group
	if err := c.ShouldBindJSON(&group); err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid group data")
		return
	}

	if err := h.db.CreateGroup(&group); err != nil {
		respondWithError(c, http.StatusInternalServerError, "Error creating group")
		return
	}

	respondWithSuccess(c, group)
}

// UpdateGroup updates an existing group
func (h *GroupHandler) UpdateGroup(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid group ID")
		return
	}

	var group models.Group
	if err := c.ShouldBindJSON(&group); err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid group data")
		return
	}

	group.ID = id
	if err := h.db.UpdateGroup(&group); err != nil {
		respondWithError(c, http.StatusInternalServerError, "Error updating group")
		return
	}

	respondWithSuccess(c, group)
}

// DeleteGroup deletes a group
func (h *GroupHandler) DeleteGroup(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid group ID")
		return
	}

	if err := h.db.DeleteGroup(id); err != nil {
		respondWithError(c, http.StatusInternalServerError, "Error deleting group")
		return
	}

	respondWithSuccess(c, gin.H{"message": "Group deleted successfully"})
}

// AddWordToGroup adds a word to a group
func (h *GroupHandler) AddWordToGroup(c *gin.Context) {
	groupIDStr := c.Param("id")
	groupID, err := strconv.Atoi(groupIDStr)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid group ID")
		return
	}

	var req struct {
		WordID int `json:"word_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid request data")
		return
	}

	if err := h.db.AddWordToGroup(req.WordID, groupID); err != nil {
		respondWithError(c, http.StatusInternalServerError, "Error adding word to group")
		return
	}

	respondWithSuccess(c, gin.H{"message": "Word added to group successfully"})
}

// RemoveWordFromGroup removes a word from a group
func (h *GroupHandler) RemoveWordFromGroup(c *gin.Context) {
	groupIDStr := c.Param("id")
	groupID, err := strconv.Atoi(groupIDStr)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid group ID")
		return
	}

	wordIDStr := c.Param("word_id")
	wordID, err := strconv.Atoi(wordIDStr)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid word ID")
		return
	}

	if err := h.db.RemoveWordFromGroup(wordID, groupID); err != nil {
		respondWithError(c, http.StatusInternalServerError, "Error removing word from group")
		return
	}

	respondWithSuccess(c, gin.H{"message": "Word removed from group successfully"})
}
