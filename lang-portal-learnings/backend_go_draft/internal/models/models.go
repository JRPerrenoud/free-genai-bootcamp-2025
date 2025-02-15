package models

import (
	"database/sql"
	"time"
	_ "github.com/mattn/go-sqlite3"
)

// Word represents a vocabulary word with its translations
type Word struct {
	ID      int    `json:"id"`
	Spanish string `json:"spanish"`
	English string `json:"english"`
}

// Group represents a thematic group of words
type Group struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// LaunchButton represents a launch button for a study activity
type LaunchButton struct {
	Text string `json:"text"`
	URL  string `json:"url"`
}

// StudyActivity represents a study activity
type StudyActivity struct {
	ID             int            `json:"id"`
	Name           string         `json:"name"`
	Description    string         `json:"description"`
	ThumbnailURL   string         `json:"thumbnail_url"`
	LaunchURL      string         `json:"-"`
	LaunchButton   *LaunchButton  `json:"launch_button"`
	StudySessions  []StudySession `json:"study_sessions,omitempty"`
	CreatedAt    time.Time `json:"-"`
}

// StudySession represents a learning session
type StudySession struct {
	ID               int       `json:"id"`
	GroupID         int       `json:"-"`
	StudyActivityID int       `json:"-"`
	ActivityName    string    `json:"activity_name"`
	GroupName       string    `json:"group_name"`
	StartTime       time.Time `json:"start_time"`
	EndTime         time.Time `json:"end_time"`
	ReviewItemsCount int      `json:"review_items_count"`
}

// WordReviewItem represents a practice record for a word
type WordReviewItem struct {
	ID             int       `json:"id"`
	WordID         int       `json:"word_id"`
	StudySessionID int       `json:"study_session_id"`
	Correct        bool      `json:"correct"`
	CreatedAt      time.Time `json:"created_at"`
}

// DB represents our database connection
type DB struct {
	*sql.DB
}

// NewDB creates a new database connection
func NewDB(dataSourceName string) (*DB, error) {
	db, err := sql.Open("sqlite3", dataSourceName)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return &DB{db}, nil
}
