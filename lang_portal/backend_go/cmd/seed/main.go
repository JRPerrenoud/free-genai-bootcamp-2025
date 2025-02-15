package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/JPerreno/free-genai-bootcamp-2025/lang_portal/backend_go/internal/models"
)

type Word struct {
	Spanish string `json:"spanish"`
	English string `json:"english"`
	Type    string `json:"type"`
}

type Group struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type GroupWords struct {
	GroupName string   `json:"group_name"`
	Words     []string `json:"words"`
}

type StudyActivity struct {
	Name         string `json:"name"`
	ThumbnailURL string `json:"thumbnail_url"`
	LaunchURL    string `json:"launch_url"`
}

type ReviewItem struct {
	Word    string `json:"word"`
	Correct bool   `json:"correct"`
}

type StudySession struct {
	GroupName    string       `json:"group_name"`
	ActivityName string       `json:"activity_name"`
	CreatedAt    string       `json:"created_at"`
	ReviewItems  []ReviewItem `json:"review_items"`
}

type SeedData struct {
	Words           []Word           `json:"words"`
	Groups          []Group          `json:"groups"`
	GroupWords      []GroupWords     `json:"group_words"`
	StudyActivities []StudyActivity  `json:"study_activities"`
	StudySessions   []StudySession   `json:"study_sessions"`
}

func main() {
	// Get the absolute path to the database file
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "db/lang_portal.db"
	}
	absDbPath, err := filepath.Abs(dbPath)
	if err != nil {
		log.Fatalf("Failed to get absolute path to database: %v", err)
	}

	// Initialize database
	db, err := models.NewDB(absDbPath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Read seed data from JSON file
	seedPath := filepath.Join("db", "seeds", "initial_seed_data.json")
	seedFile, err := os.ReadFile(seedPath)
	if err != nil {
		log.Fatalf("Failed to read seed file: %v", err)
	}

	var seedData SeedData
	if err := json.Unmarshal(seedFile, &seedData); err != nil {
		log.Fatalf("Failed to parse seed data: %v", err)
	}

	// Begin transaction
	tx, err := db.Begin()
	if err != nil {
		log.Fatalf("Failed to begin transaction: %v", err)
	}

	// Insert words
	wordStmt, err := tx.Prepare(`
		INSERT INTO words (spanish, english, type)
		VALUES (?, ?, ?)
		RETURNING id
	`)
	if err != nil {
		tx.Rollback()
		log.Fatalf("Failed to prepare word statement: %v", err)
	}
	defer wordStmt.Close()

	wordIDs := make(map[string]int)
	for _, word := range seedData.Words {
		var id int
		err = wordStmt.QueryRow(word.Spanish, word.English, word.Type).Scan(&id)
		if err != nil {
			tx.Rollback()
			log.Fatalf("Failed to insert word: %v", err)
		}
		wordIDs[word.Spanish] = id
	}

	// Insert groups
	groupStmt, err := tx.Prepare(`
		INSERT INTO groups (name, description)
		VALUES (?, ?)
		RETURNING id
	`)
	if err != nil {
		tx.Rollback()
		log.Fatalf("Failed to prepare group statement: %v", err)
	}
	defer groupStmt.Close()

	groupIDs := make(map[string]int)
	for _, group := range seedData.Groups {
		var id int
		err = groupStmt.QueryRow(group.Name, group.Description).Scan(&id)
		if err != nil {
			tx.Rollback()
			log.Fatalf("Failed to insert group: %v", err)
		}
		groupIDs[group.Name] = id
	}

	// Insert words_groups
	wordsGroupsStmt, err := tx.Prepare(`
		INSERT INTO words_groups (word_id, group_id)
		VALUES (?, ?)
	`)
	if err != nil {
		tx.Rollback()
		log.Fatalf("Failed to prepare words_groups statement: %v", err)
	}
	defer wordsGroupsStmt.Close()

	for _, groupWords := range seedData.GroupWords {
		groupID := groupIDs[groupWords.GroupName]
		for _, word := range groupWords.Words {
			wordID := wordIDs[word]
			_, err = wordsGroupsStmt.Exec(wordID, groupID)
			if err != nil {
				tx.Rollback()
				log.Fatalf("Failed to insert words_groups: %v", err)
			}
		}
	}

	// Insert study activities
	activityStmt, err := tx.Prepare(`
		INSERT INTO study_activities (name, thumbnail_url, launch_url)
		VALUES (?, ?, ?)
		RETURNING id
	`)
	if err != nil {
		tx.Rollback()
		log.Fatalf("Failed to prepare activity statement: %v", err)
	}
	defer activityStmt.Close()

	activityIDs := make(map[string]int)
	for _, activity := range seedData.StudyActivities {
		var id int
		err = activityStmt.QueryRow(activity.Name, activity.ThumbnailURL, activity.LaunchURL).Scan(&id)
		if err != nil {
			tx.Rollback()
			log.Fatalf("Failed to insert activity: %v", err)
		}
		activityIDs[activity.Name] = id
	}

	// Insert study sessions and review items
	sessionStmt, err := tx.Prepare(`
		INSERT INTO study_sessions (group_id, study_activity_id, created_at)
		VALUES (?, ?, ?)
		RETURNING id
	`)
	if err != nil {
		tx.Rollback()
		log.Fatalf("Failed to prepare session statement: %v", err)
	}
	defer sessionStmt.Close()

	reviewStmt, err := tx.Prepare(`
		INSERT INTO word_review_items (word_id, study_session_id, correct, created_at)
		VALUES (?, ?, ?, ?)
	`)
	if err != nil {
		tx.Rollback()
		log.Fatalf("Failed to prepare review statement: %v", err)
	}
	defer reviewStmt.Close()

	for _, session := range seedData.StudySessions {
		groupID := groupIDs[session.GroupName]
		activityID := activityIDs[session.ActivityName]
		createdAt, err := time.Parse(time.RFC3339, session.CreatedAt)
		if err != nil {
			tx.Rollback()
			log.Fatalf("Failed to parse created_at: %v", err)
		}

		var sessionID int
		err = sessionStmt.QueryRow(groupID, activityID, createdAt).Scan(&sessionID)
		if err != nil {
			tx.Rollback()
			log.Fatalf("Failed to insert session: %v", err)
		}

		for _, review := range session.ReviewItems {
			wordID := wordIDs[review.Word]
			_, err = reviewStmt.Exec(wordID, sessionID, review.Correct, createdAt)
			if err != nil {
				tx.Rollback()
				log.Fatalf("Failed to insert review item: %v", err)
			}
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		log.Fatalf("Failed to commit transaction: %v", err)
	}

	fmt.Println("Database seeded successfully!")
}
