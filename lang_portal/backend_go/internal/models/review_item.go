package models

import (
	"database/sql"
	"fmt"
	"time"
)

// CreateWordReviewItem creates a new word review item
func (db *DB) CreateWordReviewItem(wordID, studySessionID int, correct bool) (*WordReviewItem, error) {
	query := `
		INSERT INTO word_review_items (word_id, study_session_id, correct, created_at)
		VALUES (?, ?, ?, ?)
	`
	result, err := db.Exec(query, wordID, studySessionID, correct, time.Now())
	if err != nil {
		return nil, fmt.Errorf("error creating word review item: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("error getting last insert id: %v", err)
	}

	return db.GetWordReviewItem(int(id))
}

// GetWordReviewItem retrieves a word review item by ID
func (db *DB) GetWordReviewItem(id int) (*WordReviewItem, error) {
	item := &WordReviewItem{}
	query := `
		SELECT id, word_id, study_session_id, correct, created_at
		FROM word_review_items
		WHERE id = ?
	`
	err := db.QueryRow(query, id).Scan(
		&item.ID,
		&item.WordID,
		&item.StudySessionID,
		&item.Correct,
		&item.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error getting word review item: %v", err)
	}

	return item, nil
}

// ListSessionReviewItems retrieves all review items for a study session
func (db *DB) ListSessionReviewItems(sessionID int) ([]WordReviewItem, error) {
	var items []WordReviewItem

	query := `
		SELECT id, word_id, study_session_id, correct, created_at
		FROM word_review_items
		WHERE study_session_id = ?
		ORDER BY created_at DESC
	`
	rows, err := db.Query(query, sessionID)
	if err != nil {
		return nil, fmt.Errorf("error querying session review items: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var item WordReviewItem
		err := rows.Scan(
			&item.ID,
			&item.WordID,
			&item.StudySessionID,
			&item.Correct,
			&item.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning word review item: %v", err)
		}
		items = append(items, item)
	}

	return items, nil
}

// GetOverallAccuracy retrieves the overall accuracy across all review items
func (db *DB) GetOverallAccuracy() (float64, error) {
	var accuracy float64
	err := db.QueryRow(
		`SELECT 
			CAST(COUNT(CASE WHEN correct = 1 THEN 1 END) AS FLOAT) / 
			CAST(COUNT(*) AS FLOAT) as accuracy
		FROM word_review_items`,
	).Scan(&accuracy)

	if err != nil {
		return 0, err
	}

	return accuracy, nil
}

// GetStudyStreak returns the number of consecutive days with study sessions
func (db *DB) GetStudyStreak() (int, error) {
	rows, err := db.Query(
		`WITH RECURSIVE dates AS (
			SELECT date(created_at) as study_date
			FROM study_sessions
			GROUP BY date(created_at)
		),
		streak_calc AS (
			SELECT study_date, 
				   row_number() OVER (ORDER BY study_date DESC) as row_num
			FROM dates
		)
		SELECT COUNT(*) as streak
		FROM streak_calc
		WHERE study_date = date('now', '-' || (row_num - 1) || ' days')`,
	)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	var streak int
	if rows.Next() {
		if err := rows.Scan(&streak); err != nil {
			return 0, err
		}
	}

	return streak, nil
}
