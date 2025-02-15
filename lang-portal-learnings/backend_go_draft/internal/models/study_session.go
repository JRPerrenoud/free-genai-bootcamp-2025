package models

import (
	"database/sql"
	"fmt"
	"time"
)

// CreateStudySession creates a new study session
func (db *DB) CreateStudySession(groupID, studyActivityID int) (*StudySession, error) {
	// First, get the activity name and group name
	var activityName, groupName string
	err := db.QueryRow(`SELECT name FROM study_activities WHERE id = ?`, studyActivityID).Scan(&activityName)
	if err != nil {
		return nil, fmt.Errorf("error getting activity name: %v", err)
	}

	err = db.QueryRow(`SELECT name FROM groups WHERE id = ?`, groupID).Scan(&groupName)
	if err != nil {
		return nil, fmt.Errorf("error getting group name: %v", err)
	}

	startTime := time.Now()
	session := &StudySession{
		GroupID:         groupID,
		StudyActivityID: studyActivityID,
		ActivityName:    activityName,
		GroupName:       groupName,
		StartTime:       startTime,
	}

	query := `
		INSERT INTO study_sessions (group_id, study_activity_id, start_time)
		VALUES (?, ?, ?)
		RETURNING id
	`
	err = db.QueryRow(query, session.GroupID, session.StudyActivityID, session.StartTime).Scan(&session.ID)
	if err != nil {
		return nil, fmt.Errorf("error creating study session: %v", err)
	}

	return session, nil
}

// GetStudySession retrieves a study session by its ID
func (db *DB) GetStudySession(id int) (*StudySession, error) {
	session := &StudySession{}
	query := `
		SELECT s.id, s.group_id, s.study_activity_id, s.start_time, s.end_time,
			   g.name as group_name, a.name as activity_name,
			   COUNT(r.id) as review_items_count
		FROM study_sessions s
		JOIN groups g ON s.group_id = g.id
		JOIN study_activities a ON s.study_activity_id = a.id
		LEFT JOIN word_review_items r ON s.id = r.study_session_id
		WHERE s.id = ?
		GROUP BY s.id
	`
	err := db.QueryRow(query, id).Scan(
		&session.ID,
		&session.GroupID,
		&session.StudyActivityID,
		&session.StartTime,
		&session.EndTime,
		&session.GroupName,
		&session.ActivityName,
		&session.ReviewItemsCount,
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
		SELECT s.id, s.group_id, s.study_activity_id, s.start_time, s.end_time,
			   g.name as group_name, a.name as activity_name,
			   COUNT(r.id) as review_items_count
		FROM study_sessions s
		JOIN groups g ON s.group_id = g.id
		JOIN study_activities a ON s.study_activity_id = a.id
		LEFT JOIN word_review_items r ON s.id = r.study_session_id
		GROUP BY s.id
		ORDER BY s.start_time DESC
		LIMIT 1
	`
	err := db.QueryRow(query).Scan(
		&session.ID,
		&session.GroupID,
		&session.StudyActivityID,
		&session.StartTime,
		&session.EndTime,
		&session.GroupName,
		&session.ActivityName,
		&session.ReviewItemsCount,
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
	offset := (page - 1) * pageSize

	// Get total count
	var total int
	err := db.QueryRow(`SELECT COUNT(*) FROM study_sessions`).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("error getting total count: %v", err)
	}

	// Get paginated sessions
	query := `
		SELECT s.id, s.group_id, s.study_activity_id, s.start_time, s.end_time,
			   g.name as group_name, a.name as activity_name,
			   COUNT(r.id) as review_items_count
		FROM study_sessions s
		JOIN groups g ON s.group_id = g.id
		JOIN study_activities a ON s.study_activity_id = a.id
		LEFT JOIN word_review_items r ON s.id = r.study_session_id
		GROUP BY s.id
		ORDER BY s.start_time DESC
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
			&session.GroupID,
			&session.StudyActivityID,
			&session.StartTime,
			&session.EndTime,
			&session.GroupName,
			&session.ActivityName,
			&session.ReviewItemsCount,
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

	// Get total words available
	err := db.QueryRow(`
		SELECT COUNT(*) 
		FROM words
	`).Scan(&totalWords)
	if err != nil {
		return 0, 0, fmt.Errorf("error getting total words: %v", err)
	}
	fmt.Printf("DEBUG: Total words from database: %d\n", totalWords)

	// Get total words studied
	err = db.QueryRow(`
		SELECT COUNT(DISTINCT word_id) 
		FROM word_review_items 
		WHERE word_id IN (SELECT id FROM words)
	`).Scan(&studiedWords)
	if err != nil {
		return 0, 0, fmt.Errorf("error getting studied words: %v", err)
	}
	fmt.Printf("DEBUG: Studied words from database: %d\n", studiedWords)

	fmt.Printf("DEBUG: Returning totalWords=%d, studiedWords=%d\n", totalWords, studiedWords)
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
