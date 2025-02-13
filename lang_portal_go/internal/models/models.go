package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type PaginatedResponse[T any] struct {
	Items []T   `json:"items"`
	Total int   `json:"total"`
	Page  int   `json:"page"`
}

type WordParts []string

func (p *WordParts) MarshalJSON() ([]byte, error) {
	if p == nil {
		return []byte("[]"), nil
	}
	return json.Marshal([]string(*p))
}

func (p *WordParts) UnmarshalJSON(data []byte) error {
	var parts []string
	if err := json.Unmarshal(data, &parts); err != nil {
		return err
	}
	*p = WordParts(parts)
	return nil
}

func (p *WordParts) Scan(value interface{}) error {
	if value == nil {
		*p = WordParts{}
		return nil
	}
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("expected string, got %T", value)
	}
	*p = WordParts(strings.Split(str, ","))
	return nil
}

func (p WordParts) Value() (driver.Value, error) {
	if len(p) == 0 {
		return "", nil
	}
	return strings.Join([]string(p), ","), nil
}

type Word struct {
	ID           int64     `json:"id"`
	Spanish      string    `json:"spanish"`
	English      string    `json:"english"`
	Parts        WordParts `json:"parts"`
	CorrectCount int       `json:"correct_count"`
	WrongCount   int       `json:"wrong_count"`
}

type Group struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	WordCount int    `json:"word_count"`
}

type StudyActivity struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type StudySession struct {
	ID           int64     `json:"id"`
	ActivityID   int64     `json:"activity_id"`
	StartTime    time.Time `json:"start_time"`
	EndTime      time.Time `json:"end_time"`
	ReviewCount  int       `json:"review_count"`
	CorrectCount int       `json:"correct_count"`
}

type WordReview struct {
	ID             int64 `json:"id"`
	StudySessionID int64 `json:"study_session_id"`
	WordID         int64 `json:"word_id"`
	IsCorrect      bool  `json:"is_correct"`
}
