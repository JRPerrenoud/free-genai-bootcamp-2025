package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/JPerreno/free-genai-bootcamp-2025/lang_portal/backend_go/internal/handlers"
	"github.com/JPerreno/free-genai-bootcamp-2025/lang_portal/backend_go/internal/models"
	"github.com/gin-gonic/gin"
)

func main() {
	// Get database path from environment variable or use default
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "db/lang_portal.db"
	}

	// Get absolute path
	absDbPath, err := filepath.Abs(dbPath)
	if err != nil {
		log.Fatalf("Failed to get absolute path to database: %v", err)
	}

	// Initialize database
	db, err := models.NewDB(absDbPath)
	if err != nil {
		log.Fatal(err)
	}

	// Create handler
	h := handlers.NewHandler(db)

	// Initialize router
	r := gin.Default()

	// Add CORS middleware
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		
		c.Next()
	})

	// API routes
	api := r.Group("/api")
	{
		// Words endpoints
		api.GET("/words", h.GetWords)
		api.GET("/words/:id", h.GetWord)
		api.POST("/words", h.CreateWord)
		api.PUT("/words/:id", h.UpdateWord)
		api.DELETE("/words/:id", h.DeleteWord)

		// Groups endpoints
		api.GET("/groups", h.GetGroups)
		api.GET("/groups/:id", h.GetGroup)
		api.POST("/groups", h.CreateGroup)
		api.PUT("/groups/:id", h.UpdateGroup)
		api.DELETE("/groups/:id", h.DeleteGroup)
		api.POST("/groups/:id/words", h.AddWordsToGroup)
		api.DELETE("/groups/:id/words", h.RemoveWordsFromGroup)

		// Study activities endpoints
		api.GET("/study_activities", h.GetStudyActivities)
		api.GET("/study_activities/:id", h.GetStudyActivityByID)
		api.GET("/study_activities/:id/sessions", h.GetStudyActivitySessions)
		api.POST("/study_activities", h.CreateStudyActivity)

		// Study sessions endpoints
		api.POST("/study_sessions", h.CreateStudySession)
		api.GET("/study_sessions/:id", h.GetStudySession)
		api.POST("/study_sessions/:id/review_items", h.AddReviewItems)

		// Dashboard endpoints
		api.GET("/dashboard/last_study_session", h.GetLastStudySession)
		api.GET("/dashboard/study_progress", h.GetStudyProgress)
		api.GET("/dashboard/quick_stats", h.GetQuickStats)
	}

	// Start server
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
