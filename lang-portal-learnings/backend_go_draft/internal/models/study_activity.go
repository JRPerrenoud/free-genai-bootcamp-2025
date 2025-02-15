package models

import (
	"database/sql"
	"fmt"
	"time"
)

// CreateStudyActivity creates a new study activity
func (db *DB) CreateStudyActivity(activity *StudyActivity) error {
	query := `
		INSERT INTO study_activities (name, description, thumbnail_url, launch_url)
		VALUES (?, ?, ?, ?)
		RETURNING id, created_at
	`
	return db.QueryRow(query, activity.Name, activity.Description, activity.ThumbnailURL, activity.LaunchURL).Scan(&activity.ID, &activity.CreatedAt)
}

// GetStudyActivity retrieves a study activity by ID
func (db *DB) GetStudyActivity(id int) (*StudyActivity, error) {
	activity := &StudyActivity{}
	err := db.QueryRow(`
		SELECT id, name, description, thumbnail_url, launch_url, created_at
		FROM study_activities
		WHERE id = ?
	`, id).Scan(
		&activity.ID,
		&activity.Name,
		&activity.Description,
		&activity.ThumbnailURL,
		&activity.LaunchURL,
		&activity.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("error querying study activity: %v", err)
	}

	// Get associated study sessions
	query := `
		SELECT s.id, s.group_id, s.study_activity_id, 
			   strftime('%Y-%m-%d %H:%M:%f', s.start_time) as start_time,
			   CASE 
				   WHEN s.end_time IS NULL THEN strftime('%Y-%m-%d %H:%M:%f', '0001-01-01 00:00:00')
				   ELSE strftime('%Y-%m-%d %H:%M:%f', s.end_time)
			   END as end_time,
			   g.name as group_name, a.name as activity_name,
			   COUNT(r.id) as review_items_count
		FROM study_sessions s
		JOIN groups g ON s.group_id = g.id
		JOIN study_activities a ON s.study_activity_id = a.id
		LEFT JOIN word_review_items r ON s.id = r.study_session_id
		WHERE s.study_activity_id = ?
		GROUP BY s.id
		ORDER BY s.start_time DESC
	`
	rows, err := db.Query(query, id)
	if err != nil {
		return nil, fmt.Errorf("error querying study sessions: %v", err)
	}
	defer rows.Close()

	var sessions []StudySession
	for rows.Next() {
		var session StudySession
		var startTimeStr, endTimeStr string
		err := rows.Scan(
			&session.ID,
			&session.GroupID,
			&session.StudyActivityID,
			&startTimeStr,
			&endTimeStr,
			&session.GroupName,
			&session.ActivityName,
			&session.ReviewItemsCount,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning study session: %v", err)
		}

		startTime, err := time.Parse("2006-01-02 15:04:05", startTimeStr[:19])
		if err != nil {
			return nil, fmt.Errorf("error parsing start time: %v", err)
		}
		session.StartTime = startTime

		endTime, err := time.Parse("2006-01-02 15:04:05", endTimeStr[:19])
		if err != nil {
			return nil, fmt.Errorf("error parsing end time: %v", err)
		}
		session.EndTime = endTime

		sessions = append(sessions, session)
	}

	activity.StudySessions = sessions
	return activity, nil
}

// UpdateStudyActivity updates an existing study activity
func (db *DB) UpdateStudyActivity(activity *StudyActivity) error {
	query := `
		UPDATE study_activities
		SET name = ?, description = ?, thumbnail_url = ?, launch_url = ?
		WHERE id = ?
	`
	result, err := db.Exec(query, activity.Name, activity.Description, activity.ThumbnailURL, activity.LaunchURL, activity.ID)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return fmt.Errorf("study activity with id %d not found", activity.ID)
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
	offset := (page - 1) * pageSize

	// Get total count
	var total int
	err := db.QueryRow("SELECT COUNT(*) FROM study_activities").Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Get paginated activities
	query := `
		SELECT id, name, description, thumbnail_url, launch_url, created_at
		FROM study_activities
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`
	rows, err := db.Query(query, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var activities []StudyActivity
	for rows.Next() {
		var activity StudyActivity
		err := rows.Scan(
			&activity.ID,
			&activity.Name,
			&activity.Description,
			&activity.ThumbnailURL,
			&activity.LaunchURL,
			&activity.CreatedAt,
		)
		if err != nil {
			return nil, 0, err
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
	offset := (page - 1) * pageSize

	// Get total count
	var total int
	err := db.QueryRow(`
		SELECT COUNT(*) 
		FROM study_sessions 
		WHERE study_activity_id = ?
	`, activityID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("error getting total count: %v", err)
	}

	// Get paginated sessions
	query := `
		SELECT s.id, s.group_id, s.study_activity_id, 
			   strftime('%Y-%m-%d %H:%M:%f', s.start_time) as start_time,
			   CASE 
				   WHEN s.end_time IS NULL THEN strftime('%Y-%m-%d %H:%M:%f', '0001-01-01 00:00:00')
				   ELSE strftime('%Y-%m-%d %H:%M:%f', s.end_time)
			   END as end_time,
			   g.name as group_name, 
			   a.name as activity_name,
			   COUNT(r.id) as review_items_count
		FROM study_sessions s
		JOIN groups g ON s.group_id = g.id
		JOIN study_activities a ON s.study_activity_id = a.id
		LEFT JOIN word_review_items r ON s.id = r.study_session_id
		WHERE s.study_activity_id = ?
		GROUP BY s.id
		ORDER BY s.start_time DESC
		LIMIT ? OFFSET ?
	`
	rows, err := db.Query(query, activityID, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("error querying study sessions: %v", err)
	}
	defer rows.Close()

	var sessions []StudySession
	for rows.Next() {
		var session StudySession
		var startTimeStr, endTimeStr string
		err := rows.Scan(
			&session.ID,
			&session.GroupID,
			&session.StudyActivityID,
			&startTimeStr,
			&endTimeStr,
			&session.GroupName,
			&session.ActivityName,
			&session.ReviewItemsCount,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("error scanning study session: %v", err)
		}

		startTime, err := time.Parse("2006-01-02 15:04:05", startTimeStr[:19])
		if err != nil {
			return nil, 0, fmt.Errorf("error parsing start time: %v", err)
		}
		session.StartTime = startTime

		endTime, err := time.Parse("2006-01-02 15:04:05", endTimeStr[:19])
		if err != nil {
			return nil, 0, fmt.Errorf("error parsing end time: %v", err)
		}
		session.EndTime = endTime

		sessions = append(sessions, session)
	}

	return sessions, total, nil
}
