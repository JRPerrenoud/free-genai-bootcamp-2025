package handlers

import (
	"net/http"

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

// ListGroups returns a list of all part of speech groups
func (h *GroupHandler) ListGroups(c *gin.Context) {
	groups, err := h.db.ListGroups()
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Error fetching groups")
		return
	}

	respondWithSuccess(c, gin.H{
		"items": groups,
	})
}

// GetGroup returns details for a specific part of speech group
func (h *GroupHandler) GetGroup(c *gin.Context) {
	partOfSpeech := c.Param("part_of_speech")

	group, err := h.db.GetGroupStats(partOfSpeech)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Error fetching group")
		return
	}

	if group == nil {
		respondWithError(c, http.StatusNotFound, "Group not found")
		return
	}

	respondWithSuccess(c, group)
}

// GetGroupWords returns all words for a specific part of speech
func (h *GroupHandler) GetGroupWords(c *gin.Context) {
	partOfSpeech := c.Param("part_of_speech")
	page, pageSize := getPaginationParams(c, 20)

	words, total, err := h.db.GetGroupWords(partOfSpeech, page, pageSize)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Error fetching group words")
		return
	}

	response := wrapWithPagination(words, page, total, pageSize)
	respondWithSuccess(c, response)
}
