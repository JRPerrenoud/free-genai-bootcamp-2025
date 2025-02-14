package models

import (
	"database/sql"
	"fmt"
	"time"
)

// CreateStudySession creates a new study session
func (db *DB) CreateStudySession(groupID, studyActivityID int) (*StudySession, error) {
	query := `
		INSERT INTO study_sessions (group_id, study_activity_id, created_at) 
		VALUES (?, ?, ?)
	`
	result, err := db.Exec(query, groupID, studyActivityID, time.Now())
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
		SELECT s.id, s.group_id, s.created_at, s.study_activity_id, g.name as group_name
		FROM study_sessions s
		LEFT JOIN groups g ON s.group_id = g.id
		WHERE s.id = ?
	`
	err := db.QueryRow(query, id).Scan(
		&session.ID,
		&session.GroupID,
		&session.CreatedAt,
		&session.StudyActivityID,
		&session.GroupName,
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
		SELECT s.id, s.group_id, s.created_at, s.study_activity_id, g.name as group_name
		FROM study_sessions s
		LEFT JOIN groups g ON s.group_id = g.id
		ORDER BY s.created_at DESC
		LIMIT 1
	`
	err := db.QueryRow(query).Scan(
		&session.ID,
		&session.GroupID,
		&session.CreatedAt,
		&session.StudyActivityID,
		&session.GroupName,
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
	var sessions []StudySession
	var total int

	// Get total count
	err := db.QueryRow("SELECT COUNT(*) FROM study_sessions").Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("error getting total count: %v", err)
	}

	// Get sessions
	query := `
		SELECT s.id, s.group_id, s.created_at, s.study_activity_id, g.name as group_name
		FROM study_sessions s
		LEFT JOIN groups g ON s.group_id = g.id
		ORDER BY s.created_at DESC
		LIMIT ? OFFSET ?
	`
	offset := (page - 1) * pageSize
	rows, err := db.Query(query, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("error querying study sessions: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var session StudySession
		err := rows.Scan(
			&session.ID,
			&session.GroupID,
			&session.CreatedAt,
			&session.StudyActivityID,
			&session.GroupName,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("error scanning study session: %v", err)
		}
		sessions = append(sessions, session)
	}

	return sessions, total, nil
}

// GetStudyProgress retrieves overall study progress
func (db *DB) GetStudyProgress() (int, int, error) {
	var totalWords, studiedWords int

	err := db.QueryRow(`
		SELECT COUNT(DISTINCT id) 
		FROM words
	`).Scan(&totalWords)
	if err != nil {
		return 0, 0, fmt.Errorf("error getting total words: %v", err)
	}

	err = db.QueryRow(`
		SELECT COUNT(DISTINCT word_id) 
		FROM word_review_items
	`).Scan(&studiedWords)
	if err != nil {
		return 0, 0, fmt.Errorf("error getting studied words: %v", err)
	}

	return totalWords, studiedWords, nil
}

// GetStudySessionStats retrieves statistics for a study session
func (db *DB) GetStudySessionStats(sessionID int) (int, int, error) {
	var correct, wrong int

	query := `
		SELECT 
			SUM(CASE WHEN correct = 1 THEN 1 ELSE 0 END) as correct_count,
			SUM(CASE WHEN correct = 0 THEN 1 ELSE 0 END) as wrong_count
		FROM word_review_items
		WHERE study_session_id = ?
	`
	err := db.QueryRow(query, sessionID).Scan(&correct, &wrong)
	if err != nil {
		return 0, 0, fmt.Errorf("error getting study session stats: %v", err)
	}

	return correct, wrong, nil
}
