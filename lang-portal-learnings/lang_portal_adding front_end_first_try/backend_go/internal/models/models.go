package models

import (
	"database/sql"
	"time"
	_ "github.com/mattn/go-sqlite3"
)

// Word represents a vocabulary word with its translations and parts of speech
type Word struct {
	ID           int    `json:"id"`
	Spanish      string `json:"spanish"`
	English      string `json:"english"`
	PartOfSpeech string `json:"part_of_speech"`
}

// StudyActivity represents a specific study activity
type StudyActivity struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

// StudySession represents a learning session
type StudySession struct {
	ID              int       `json:"id"`
	PartOfSpeech    string    `json:"part_of_speech"`
	StudyActivityID int       `json:"study_activity_id"`
	CreatedAt       time.Time `json:"created_at"`
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
