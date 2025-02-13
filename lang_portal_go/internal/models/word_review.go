package models

import "time"

// WordReviewItem represents a record of word practice
type WordReviewItem struct {
	ID             int64     `json:"id"`
	WordID         int64     `json:"word_id"`
	StudySessionID int64     `json:"study_session_id"`
	Correct        bool      `json:"correct"`
	CreatedAt      time.Time `json:"created_at"`
}

// CreateWordReviewItem creates a new word review item
func (db *DB) CreateWordReviewItem(wordID, studySessionID int64, correct bool) (*WordReviewItem, error) {
	result, err := db.Exec(`
		INSERT INTO word_review_items (word_id, study_session_id, correct, created_at)
		VALUES (?, ?, ?, CURRENT_TIMESTAMP)`,
		wordID, studySessionID, correct)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	// Retrieve the created review item
	var review WordReviewItem
	err = db.QueryRow(`
		SELECT id, word_id, study_session_id, correct, created_at
		FROM word_review_items
		WHERE id = ?`, id).Scan(
		&review.ID, &review.WordID, &review.StudySessionID,
		&review.Correct, &review.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &review, nil
}

// GetStudyProgress retrieves overall study progress
func (db *DB) GetStudyProgress() (map[string]int, error) {
	query := `
		WITH WordStats AS (
			SELECT 
				COUNT(DISTINCT word_id) as words_studied
			FROM word_review_items
		),
		TotalWords AS (
			SELECT COUNT(*) as total_words
			FROM words
		)
		SELECT 
			COALESCE(ws.words_studied, 0) as total_words_studied,
			tw.total_words as total_available_words
		FROM WordStats ws
		CROSS JOIN TotalWords tw`

	var progress struct {
		TotalWordsStudied    int
		TotalAvailableWords int
	}

	err := db.QueryRow(query).Scan(&progress.TotalWordsStudied, &progress.TotalAvailableWords)
	if err != nil {
		return nil, err
	}

	return map[string]int{
		"total_words_studied":    progress.TotalWordsStudied,
		"total_available_words": progress.TotalAvailableWords,
	}, nil
}

// GetQuickStats retrieves quick overview statistics
func (db *DB) GetQuickStats() (map[string]interface{}, error) {
	query := `
		WITH Stats AS (
			SELECT 
				(SELECT COUNT(*) FROM words) as total_words,
				(SELECT COUNT(*) FROM groups) as total_groups,
				(SELECT COUNT(*) FROM study_sessions) as total_study_sessions,
				(SELECT 
					CAST(SUM(CASE WHEN correct THEN 1 ELSE 0 END) AS FLOAT) / 
					CAST(COUNT(*) AS FLOAT)
				FROM word_review_items) as overall_accuracy
		)
		SELECT 
			total_words,
			total_groups,
			total_study_sessions,
			COALESCE(overall_accuracy, 0) as overall_accuracy
		FROM Stats`

	var stats struct {
		TotalWords         int
		TotalGroups        int
		TotalStudySessions int
		OverallAccuracy    float64
	}

	err := db.QueryRow(query).Scan(
		&stats.TotalWords,
		&stats.TotalGroups,
		&stats.TotalStudySessions,
		&stats.OverallAccuracy,
	)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"total_words":          stats.TotalWords,
		"total_groups":         stats.TotalGroups,
		"total_study_sessions": stats.TotalStudySessions,
		"overall_accuracy":     stats.OverallAccuracy,
	}, nil
}
