package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"lang_portal_go/internal/models"
)

// WordHandler handles word-related routes
type WordHandler struct {
	db *models.DB
}

// NewWordHandler creates a new word handler
func NewWordHandler(db *models.DB) *WordHandler {
	return &WordHandler{db: db}
}

// ListWords returns a paginated list of words
func (h *WordHandler) ListWords(c *gin.Context) {
	page, pageSize := getPaginationParams(c, 20)
	
	// Get optional group filter
	groupIDStr := c.Query("group_id")
	var groupID *int
	if groupIDStr != "" {
		if id, err := strconv.Atoi(groupIDStr); err == nil {
			groupID = &id
		}
	}

	words, total, err := h.db.ListWords(page, pageSize, groupID)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Error fetching words")
		return
	}

	response := wrapWithPagination(words, page, total, pageSize)
	respondWithSuccess(c, response)
}

// GetWord returns details for a specific word
func (h *WordHandler) GetWord(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid word ID")
		return
	}

	word, err := h.db.GetWord(id)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Error fetching word")
		return
	}

	if word == nil {
		respondWithError(c, http.StatusNotFound, "Word not found")
		return
	}

	// Get study statistics for this word
	totalReviews, correctReviews, err := h.db.GetWordStats(id)
	if err != nil {
		if err.Error() == "word not found" {
			respondWithError(c, http.StatusNotFound, "Word not found")
			return
		}
		respondWithError(c, http.StatusInternalServerError, "Error fetching word statistics")
		return
	}

	response := gin.H{
		"word": word,
		"study_stats": gin.H{
			"total_reviews":   totalReviews,
			"correct_reviews": correctReviews,
		},
	}

	respondWithSuccess(c, response)
}

// CreateWord creates a new word
func (h *WordHandler) CreateWord(c *gin.Context) {
	var word models.Word
	if err := c.ShouldBindJSON(&word); err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid word data")
		return
	}

	if err := h.db.CreateWord(&word); err != nil {
		respondWithError(c, http.StatusInternalServerError, "Error creating word")
		return
	}

	respondWithSuccess(c, word)
}

// UpdateWord updates an existing word
func (h *WordHandler) UpdateWord(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid word ID")
		return
	}

	var word models.Word
	if err := c.ShouldBindJSON(&word); err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid word data")
		return
	}

	word.ID = id
	if err := h.db.UpdateWord(&word); err != nil {
		respondWithError(c, http.StatusInternalServerError, "Error updating word")
		return
	}

	respondWithSuccess(c, word)
}

// DeleteWord deletes a word
func (h *WordHandler) DeleteWord(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid word ID")
		return
	}

	if err := h.db.DeleteWord(id); err != nil {
		respondWithError(c, http.StatusInternalServerError, "Error deleting word")
		return
	}

	respondWithSuccess(c, gin.H{"message": "Word deleted successfully"})
}
