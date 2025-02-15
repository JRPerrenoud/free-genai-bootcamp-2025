package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Response represents a standard API response
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string     `json:"error,omitempty"`
}

// PaginationResponse wraps paginated data
type PaginationResponse struct {
	Items         interface{} `json:"items"`
	CurrentPage   int         `json:"current_page"`
	TotalPages    int         `json:"total_pages"`
	TotalItems    int         `json:"total_items"`
	ItemsPerPage  int         `json:"items_per_page"`
}

// respondWithError sends an error response
func respondWithError(c *gin.Context, code int, message string) {
	c.JSON(code, Response{
		Success: false,
		Error:   message,
	})
}

// respondWithSuccess sends a success response
func respondWithSuccess(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    data,
	})
}

// getPaginationParams extracts pagination parameters from the request
func getPaginationParams(c *gin.Context, defaultPageSize int) (page, pageSize int) {
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", strconv.Itoa(defaultPageSize))

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err = strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		pageSize = defaultPageSize
	}

	return page, pageSize
}

// wrapWithPagination wraps data with pagination metadata
func wrapWithPagination(items interface{}, currentPage, totalItems, pageSize int) PaginationResponse {
	totalPages := (totalItems + pageSize - 1) / pageSize
	if totalPages < 1 {
		totalPages = 1
	}

	return PaginationResponse{
		Items:         items,
		CurrentPage:   currentPage,
		TotalPages:    totalPages,
		TotalItems:    totalItems,
		ItemsPerPage:  pageSize,
	}
}
