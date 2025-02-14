package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"lang_portal_go/internal/config"
	"lang_portal_go/internal/models"
)

func main() {
	// Initialize configuration
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("Failed to initialize config: %v", err)
	}
	defer cfg.Close()

	// Create DB instance
	db := models.NewDB(cfg.DB)

	// Initialize Gin router
	r := gin.Default()

	// Enable CORS for all origins
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

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// API routes
	api := r.Group("/api")
	{
		// Words endpoints
		words := api.Group("/words")
		{
			words.GET("", func(c *gin.Context) {
				page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
				perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "10"))
				
				words, err := db.GetWords(page, perPage)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				c.JSON(http.StatusOK, words)
			})

			words.POST("", func(c *gin.Context) {
				var word models.Word
				if err := c.ShouldBindJSON(&word); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}

				if err := db.CreateWord(&word); err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				c.JSON(http.StatusCreated, word)
			})

			words.PUT("/:id", func(c *gin.Context) {
				id, err := strconv.ParseInt(c.Param("id"), 10, 64)
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": "invalid word ID"})
					return
				}

				var word models.Word
				if err := c.ShouldBindJSON(&word); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}
				word.ID = id

				if err := db.UpdateWord(&word); err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				c.JSON(http.StatusOK, word)
			})

			words.DELETE("/:id", func(c *gin.Context) {
				id, err := strconv.ParseInt(c.Param("id"), 10, 64)
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": "invalid word ID"})
					return
				}

				if err := db.DeleteWord(id); err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				c.Status(http.StatusNoContent)
			})
		}

		// Groups endpoints
		groups := api.Group("/groups")
		{
			groups.GET("", func(c *gin.Context) {
				page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
				perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "10"))
				
				groups, err := db.GetGroups(page, perPage)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				c.JSON(http.StatusOK, groups)
			})

			groups.POST("", func(c *gin.Context) {
				var group models.Group
				if err := c.ShouldBindJSON(&group); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}

				if err := db.CreateGroup(&group); err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				c.JSON(http.StatusCreated, group)
			})

			groups.PUT("/:id", func(c *gin.Context) {
				id, err := strconv.ParseInt(c.Param("id"), 10, 64)
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": "invalid group ID"})
					return
				}

				var group models.Group
				if err := c.ShouldBindJSON(&group); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}
				group.ID = id

				if err := db.UpdateGroup(&group); err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				c.JSON(http.StatusOK, group)
			})

			groups.DELETE("/:id", func(c *gin.Context) {
				id, err := strconv.ParseInt(c.Param("id"), 10, 64)
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": "invalid group ID"})
					return
				}

				if err := db.DeleteGroup(id); err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				c.Status(http.StatusNoContent)
			})

			groups.GET("/:id/words", func(c *gin.Context) {
				groupID, err := strconv.ParseInt(c.Param("id"), 10, 64)
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": "invalid group ID"})
					return
				}

				page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
				perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "10"))
				
				words, err := db.GetGroupWords(groupID, page, perPage)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				c.JSON(http.StatusOK, words)
			})

			groups.POST("/:id/words", func(c *gin.Context) {
				groupID, err := strconv.ParseInt(c.Param("id"), 10, 64)
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": "invalid group ID"})
					return
				}

				var wordIDs []int64
				if err := c.ShouldBindJSON(&wordIDs); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}

				if err := db.AddWordsToGroup(groupID, wordIDs); err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				c.Status(http.StatusNoContent)
			})

			groups.DELETE("/:id/words", func(c *gin.Context) {
				groupID, err := strconv.ParseInt(c.Param("id"), 10, 64)
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": "invalid group ID"})
					return
				}

				var wordIDs []int64
				if err := c.ShouldBindJSON(&wordIDs); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}

				if err := db.RemoveWordsFromGroup(groupID, wordIDs); err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				c.Status(http.StatusNoContent)
			})
		}

		// Study activities endpoints
		activities := api.Group("/study-activities")
		{
			activities.GET("", func(c *gin.Context) {
				page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
				perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "10"))
				
				activities, err := db.GetStudyActivities(page, perPage)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				c.JSON(http.StatusOK, activities)
			})
		}

		// Study sessions endpoints
		sessions := api.Group("/study-sessions")
		{
			sessions.GET("", func(c *gin.Context) {
				page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
				perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "10"))
				
				sessions, err := db.GetStudySessions(page, perPage)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				c.JSON(http.StatusOK, sessions)
			})

			sessions.POST("", func(c *gin.Context) {
				var session models.StudySession
				if err := c.ShouldBindJSON(&session); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}

				if err := db.CreateStudySession(&session); err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				c.JSON(http.StatusCreated, session)
			})

			sessions.POST("/:id/reviews", func(c *gin.Context) {
				sessionID, err := strconv.ParseInt(c.Param("id"), 10, 64)
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": "invalid session ID"})
					return
				}

				var reviews []models.WordReview
				if err := c.ShouldBindJSON(&reviews); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}

				for i := range reviews {
					reviews[i].StudySessionID = sessionID
				}

				if err := db.CreateWordReviews(reviews); err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				c.JSON(http.StatusCreated, reviews)
			})
		}
	}

	// Start the server
	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("Server starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
