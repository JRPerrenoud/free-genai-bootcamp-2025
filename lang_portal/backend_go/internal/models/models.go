package models

import (
	"time"
)

type Word struct {
	ID      int    `json:"id"`
	Spanish string `json:"spanish"`
	English string `json:"english"`
	Type    string `json:"type"` // verb, noun, adjective, etc.
}

type Group struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	WordCount   int    `json:"word_count,omitempty"`
}

type StudyActivity struct {
	ID           int            `json:"id"`
	Name         string         `json:"name"`
	ThumbnailURL *string        `json:"thumbnail_url,omitempty"`
	LaunchURL    string         `json:"launch_url"`
	Description  *string        `json:"description,omitempty"`
	CreatedAt    time.Time      `json:"created_at"`
	Sessions     []StudySession `json:"sessions,omitempty"`
}

type StudySession struct {
	ID               int              `json:"id"`
	GroupID          int              `json:"group_id"`
	StudyActivityID  int              `json:"study_activity_id"`
	CreatedAt        time.Time        `json:"created_at"`
	Group            *Group           `json:"group,omitempty"`
	Activity         *StudyActivity   `json:"activity,omitempty"`
	ReviewItems      []WordReviewItem `json:"review_items,omitempty"`
	ReviewItemsCount int              `json:"review_items_count"`
}

type WordReviewItem struct {
	ID             int       `json:"id"`
	WordID         int       `json:"word_id"`
	StudySessionID int       `json:"study_session_id"`
	Correct        bool      `json:"correct"`
	CreatedAt      time.Time `json:"created_at"`
	Word           *Word     `json:"word,omitempty"`
}

type PaginatedResponse struct {
	TotalItems   int         `json:"total_items"`
	CurrentPage  int         `json:"current_page"`
	TotalPages   int         `json:"total_pages"`
	ItemsPerPage int         `json:"items_per_page"`
	Items        interface{} `json:"items"`
}

type DashboardStats struct {
	TotalWords     int `json:"total_words"`
	TotalGroups    int `json:"total_groups"`
	TotalSessions  int `json:"total_sessions"`
	TotalReviews   int `json:"total_reviews"`
	CorrectReviews int `json:"correct_reviews"`
	WrongReviews   int `json:"wrong_reviews"`
}

type StudyProgress struct {
	CorrectCount int `json:"correct_count"`
	WrongCount   int `json:"wrong_count"`
}

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}
