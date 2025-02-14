package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"lang_portal_go/internal/handlers"
	"lang_portal_go/internal/models"
)

func setupRouter(db *models.DB) *gin.Engine {
	router := gin.Default()

	// CORS middleware
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		
		c.Next()
	})

	// Create handlers
	dashboardHandler := handlers.NewDashboardHandler(db)
	wordHandler := handlers.NewWordHandler(db)
	groupHandler := handlers.NewGroupHandler(db)
	studyActivityHandler := handlers.NewStudyActivityHandler(db)

	// API routes will be grouped under /api
	api := router.Group("/api")
	{
		// Dashboard routes
		api.GET("/dashboard/last_study_session", dashboardHandler.LastStudySession)
		api.GET("/dashboard/study_progress", dashboardHandler.StudyProgress)
		api.GET("/dashboard/quick-stats", dashboardHandler.QuickStats)

		// Word routes
		api.GET("/words", wordHandler.ListWords)
		api.GET("/words/:id", wordHandler.GetWord)
		api.POST("/words", wordHandler.CreateWord)
		api.PUT("/words/:id", wordHandler.UpdateWord)
		api.DELETE("/words/:id", wordHandler.DeleteWord)

		// Group routes
		api.GET("/groups", groupHandler.ListGroups)
		api.GET("/groups/:id", groupHandler.GetGroup)
		api.POST("/groups", groupHandler.CreateGroup)
		api.PUT("/groups/:id", groupHandler.UpdateGroup)
		api.DELETE("/groups/:id", groupHandler.DeleteGroup)
		api.POST("/groups/:id/words", groupHandler.AddWordToGroup)
		api.DELETE("/groups/:id/words/:word_id", groupHandler.RemoveWordFromGroup)

		// Study activity routes
		api.GET("/study_activities", studyActivityHandler.ListStudyActivities)
		api.GET("/study_activities/:id", studyActivityHandler.GetStudyActivity)
		api.POST("/study_activities", studyActivityHandler.CreateStudyActivity)
		api.PUT("/study_activities/:id", studyActivityHandler.UpdateStudyActivity)
		api.DELETE("/study_activities/:id", studyActivityHandler.DeleteStudyActivity)
		api.GET("/study_activities/:id/sessions", studyActivityHandler.GetStudyActivitySessions)
		api.POST("/study_activities/:id/start", studyActivityHandler.StartStudySession)
		api.POST("/study_sessions/:session_id/results", studyActivityHandler.RecordStudyResult)
	}

	return router
}

func main() {
	// Initialize database connection
	db, err := models.NewDB("words.db")
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}
	
	router := setupRouter(db)
	
	// Start the server on port 8080
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}
