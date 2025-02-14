package models

import (
	"database/sql"
	"fmt"
	"time"
)

// CreateStudySession creates a new study session
func (db *DB) CreateStudySession(partOfSpeech string, studyActivityID int) (*StudySession, error) {
	query := `
		INSERT INTO study_sessions (part_of_speech, study_activity_id, created_at) 
		VALUES (?, ?, ?)
	`
	result, err := db.Exec(query, partOfSpeech, studyActivityID, time.Now())
	if err != nil {
		return nil, fmt.Errorf("error creating study session: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("error getting last insert id: %v", err)
	}

	return db.GetStudySession(int(id))
}

// GetStudySession retrieves a study session by its ID
func (db *DB) GetStudySession(id int) (*StudySession, error) {
	session := &StudySession{}
	query := `
		SELECT id, part_of_speech, created_at, study_activity_id
		FROM study_sessions
		WHERE id = ?
	`
	err := db.QueryRow(query, id).Scan(
		&session.ID,
		&session.PartOfSpeech,
		&session.CreatedAt,
		&session.StudyActivityID,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error getting study session: %v", err)
	}

	return session, nil
}

// GetLastStudySession retrieves the most recent study session
func (db *DB) GetLastStudySession() (*StudySession, error) {
	session := &StudySession{}
	query := `
		SELECT id, part_of_speech, created_at, study_activity_id
		FROM study_sessions
		ORDER BY created_at DESC
		LIMIT 1
	`
	err := db.QueryRow(query).Scan(
		&session.ID,
		&session.PartOfSpeech,
		&session.CreatedAt,
		&session.StudyActivityID,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error getting last study session: %v", err)
	}

	return session, nil
}

// ListStudySessions retrieves a paginated list of study sessions
func (db *DB) ListStudySessions(page, pageSize int) ([]StudySession, int, error) {
	var total int
	err := db.QueryRow(`SELECT COUNT(*) FROM study_sessions`).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("error getting total count: %v", err)
	}

	offset := (page - 1) * pageSize
	query := `
		SELECT id, part_of_speech, created_at, study_activity_id
		FROM study_sessions
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`
	rows, err := db.Query(query, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("error querying study sessions: %v", err)
	}
	defer rows.Close()

	var sessions []StudySession
	for rows.Next() {
		var session StudySession
		err := rows.Scan(
			&session.ID,
			&session.PartOfSpeech,
			&session.CreatedAt,
			&session.StudyActivityID,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("error scanning study session: %v", err)
		}
		sessions = append(sessions, session)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating study sessions: %v", err)
	}

	return sessions, total, nil
}

// GetStudyProgress retrieves overall study progress
func (db *DB) GetStudyProgress() (int, int, error) {
	var totalWords, reviewedWords int

	err := db.QueryRow(`SELECT COUNT(*) FROM words`).Scan(&totalWords)
	if err != nil {
		return 0, 0, fmt.Errorf("error getting total words: %v", err)
	}

	err = db.QueryRow(`
		SELECT COUNT(DISTINCT word_id) 
		FROM word_review_items
	`).Scan(&reviewedWords)
	if err != nil {
		return 0, 0, fmt.Errorf("error getting reviewed words: %v", err)
	}

	return totalWords, reviewedWords, nil
}

// GetStudySessionStats retrieves statistics for a study session
func (db *DB) GetStudySessionStats(sessionID int) (int, int, error) {
	var totalWords, correctWords int

	err := db.QueryRow(`
		SELECT COUNT(*), SUM(CASE WHEN correct THEN 1 ELSE 0 END)
		FROM word_review_items
		WHERE study_session_id = ?
	`, sessionID).Scan(&totalWords, &correctWords)
	if err != nil {
		return 0, 0, fmt.Errorf("error getting session stats: %v", err)
	}

	return totalWords, correctWords, nil
}
