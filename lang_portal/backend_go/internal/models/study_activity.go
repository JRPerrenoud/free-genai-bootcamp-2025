package models

import (
	"database/sql"
	"fmt"
	"time"
)

// CreateStudyActivity creates a new study activity
func (db *DB) CreateStudyActivity(activity *StudyActivity) error {
	query := `
		INSERT INTO study_activities (name, description, created_at)
		VALUES (?, ?, ?)
	`

	activity.CreatedAt = time.Now()
	result, err := db.Exec(query, activity.Name, activity.Description, activity.CreatedAt)
	if err != nil {
		return fmt.Errorf("error creating study activity: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("error getting last insert id: %v", err)
	}

	activity.ID = int(id)
	return nil
}

// GetStudyActivity retrieves a study activity by ID
func (db *DB) GetStudyActivity(id int) (*StudyActivity, error) {
	activity := &StudyActivity{}
	query := `SELECT id, name, description, created_at FROM study_activities WHERE id = ?`
	
	err := db.QueryRow(query, id).Scan(&activity.ID, &activity.Name, &activity.Description, &activity.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error getting study activity: %v", err)
	}

	return activity, nil
}

// UpdateStudyActivity updates an existing study activity
func (db *DB) UpdateStudyActivity(activity *StudyActivity) error {
	query := `
		UPDATE study_activities
		SET name = ?, description = ?
		WHERE id = ?
	`

	result, err := db.Exec(query, activity.Name, activity.Description, activity.ID)
	if err != nil {
		return fmt.Errorf("error updating study activity: %v", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %v", err)
	}
	if rows == 0 {
		return fmt.Errorf("study activity not found")
	}

	return nil
}

// DeleteStudyActivity deletes a study activity by ID
func (db *DB) DeleteStudyActivity(id int) error {
	query := `DELETE FROM study_activities WHERE id = ?`
	
	result, err := db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error deleting study activity: %v", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %v", err)
	}
	if rows == 0 {
		return fmt.Errorf("study activity not found")
	}

	return nil
}

// ListStudyActivities returns a paginated list of study activities
func (db *DB) ListStudyActivities(page, pageSize int) ([]StudyActivity, int, error) {
	var activities []StudyActivity
	var total int

	// Get total count
	err := db.QueryRow(`SELECT COUNT(*) FROM study_activities`).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("error getting total count: %v", err)
	}

	// Get paginated activities
	query := `
		SELECT id, name, description, created_at
		FROM study_activities
		LIMIT ? OFFSET ?
	`

	offset := (page - 1) * pageSize
	rows, err := db.Query(query, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("error querying study activities: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var activity StudyActivity
		err := rows.Scan(&activity.ID, &activity.Name, &activity.Description, &activity.CreatedAt)
		if err != nil {
			return nil, 0, fmt.Errorf("error scanning study activity: %v", err)
		}
		activities = append(activities, activity)
	}

	return activities, total, nil
}

// GetStudyActivityStats retrieves statistics for a study activity
func (db *DB) GetStudyActivityStats(activityID int) (int, float64, error) {
	var totalSessions int
	var avgCorrect float64

	query := `
		SELECT 
			COUNT(DISTINCT s.id) as total_sessions,
			COALESCE(AVG(CASE WHEN r.correct = 1 THEN 1.0 ELSE 0.0 END), 0.0) as avg_correct
		FROM study_sessions s
		LEFT JOIN word_review_items r ON s.id = r.study_session_id
		WHERE s.study_activity_id = ?
		GROUP BY s.study_activity_id
	`
	err := db.QueryRow(query, activityID).Scan(&totalSessions, &avgCorrect)
	if err != nil && err != sql.ErrNoRows {
		return 0, 0, fmt.Errorf("error getting study activity stats: %v", err)
	}

	return totalSessions, avgCorrect, nil
}

// GetStudyActivitySessions returns all study sessions for a specific activity
func (db *DB) GetStudyActivitySessions(activityID, page, pageSize int) ([]StudySession, int, error) {
	var sessions []StudySession
	var total int

	// Get total count
	err := db.QueryRow(`
		SELECT COUNT(*) 
		FROM study_sessions 
		WHERE study_activity_id = ?`, activityID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("error getting total count: %v", err)
	}

	// Get paginated sessions
	query := `
		SELECT s.id, s.group_id, s.created_at, s.study_activity_id, g.name as group_name
		FROM study_sessions s
		LEFT JOIN groups g ON s.group_id = g.id
		WHERE s.study_activity_id = ?
		ORDER BY s.created_at DESC
		LIMIT ? OFFSET ?
	`

	offset := (page - 1) * pageSize
	rows, err := db.Query(query, activityID, pageSize, offset)
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
