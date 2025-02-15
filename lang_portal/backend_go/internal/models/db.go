package models

import (
	"database/sql"
	"fmt"
	"time"
	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	*sql.DB
}

func NewDB(dataSourceName string) (*DB, error) {
	db, err := sql.Open("sqlite3", dataSourceName+"?_foreign_keys=on&_journal_mode=WAL")
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return &DB{db}, nil
}

func (db *DB) GetDashboardStats(groupID int) (*DashboardStats, error) {
	var stats DashboardStats

	// Get total words and groups
	query := `
		SELECT 
			(SELECT COUNT(*) FROM words) as total_words,
			(SELECT COUNT(*) FROM groups) as total_groups,
			(SELECT COUNT(*) FROM study_sessions) as total_sessions,
			(SELECT COUNT(*) FROM word_review_items) as total_reviews,
			(SELECT COUNT(*) FROM word_review_items WHERE correct = 1) as correct_reviews,
			(SELECT COUNT(*) FROM word_review_items WHERE correct = 0) as wrong_reviews`

	if groupID > 0 {
		query = `
			SELECT 
				(SELECT COUNT(*) FROM words_groups WHERE group_id = ?) as total_words,
				1 as total_groups,
				(SELECT COUNT(*) FROM study_sessions WHERE group_id = ?) as total_sessions,
				(SELECT COUNT(*) FROM word_review_items wri 
					JOIN study_sessions ss ON wri.study_session_id = ss.id 
					WHERE ss.group_id = ?) as total_reviews,
				(SELECT COUNT(*) FROM word_review_items wri 
					JOIN study_sessions ss ON wri.study_session_id = ss.id 
					WHERE ss.group_id = ? AND wri.correct = 1) as correct_reviews,
				(SELECT COUNT(*) FROM word_review_items wri 
					JOIN study_sessions ss ON wri.study_session_id = ss.id 
					WHERE ss.group_id = ? AND wri.correct = 0) as wrong_reviews`
	}

	var err error
	if groupID > 0 {
		err = db.QueryRow(query, groupID, groupID, groupID, groupID, groupID).Scan(
			&stats.TotalWords,
			&stats.TotalGroups,
			&stats.TotalSessions,
			&stats.TotalReviews,
			&stats.CorrectReviews,
			&stats.WrongReviews,
		)
	} else {
		err = db.QueryRow(query).Scan(
			&stats.TotalWords,
			&stats.TotalGroups,
			&stats.TotalSessions,
			&stats.TotalReviews,
			&stats.CorrectReviews,
			&stats.WrongReviews,
		)
	}
	if err != nil {
		return nil, fmt.Errorf("error getting dashboard stats: %v", err)
	}

	return &stats, nil
}

func (db *DB) GetStudyProgress(groupID int) (*StudyProgress, error) {
	progress := &StudyProgress{}

	query := `
		SELECT 
			COALESCE(SUM(CASE WHEN correct = 1 THEN 1 ELSE 0 END), 0) as correct_count,
			COALESCE(SUM(CASE WHEN correct = 0 THEN 1 ELSE 0 END), 0) as wrong_count
		FROM word_review_items wri
		JOIN study_sessions ss ON wri.study_session_id = ss.id
		WHERE ss.group_id = ?`

	err := db.QueryRow(query, groupID).Scan(&progress.CorrectCount, &progress.WrongCount)
	if err != nil {
		return nil, fmt.Errorf("error getting study progress: %v", err)
	}

	return progress, nil
}

func (db *DB) GetLastStudySession(groupID int) (*StudySession, error) {
	query := `
		SELECT s.id, s.group_id, s.study_activity_id, s.created_at,
			   g.name, g.description,
			   a.name, a.thumbnail_url, a.launch_url, a.description
		FROM study_sessions s
		LEFT JOIN groups g ON s.group_id = g.id
		LEFT JOIN study_activities a ON s.study_activity_id = a.id`

	if groupID > 0 {
		query += ` WHERE s.group_id = ?`
	}
	query += ` ORDER BY s.created_at DESC LIMIT 1`

	var session StudySession
	var group Group
	var activity StudyActivity
	var thumbURL, activityDesc sql.NullString

	var err error
	if groupID > 0 {
		err = db.QueryRow(query, groupID).Scan(
			&session.ID, &session.GroupID, &session.StudyActivityID, &session.CreatedAt,
			&group.Name, &group.Description,
			&activity.Name, &thumbURL, &activity.LaunchURL, &activityDesc,
		)
	} else {
		err = db.QueryRow(query).Scan(
			&session.ID, &session.GroupID, &session.StudyActivityID, &session.CreatedAt,
			&group.Name, &group.Description,
			&activity.Name, &thumbURL, &activity.LaunchURL, &activityDesc,
		)
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("error getting last study session: %v", err)
	}

	// Handle null strings
	if thumbURL.Valid {
		activity.ThumbnailURL = &thumbURL.String
	}
	if activityDesc.Valid {
		activity.Description = &activityDesc.String
	}

	group.ID = session.GroupID
	activity.ID = session.StudyActivityID
	session.Group = &group
	session.Activity = &activity

	// Get review items
	reviewItems, err := db.GetReviewItemsForSession(session.ID)
	if err != nil {
		return nil, fmt.Errorf("error getting review items: %v", err)
	}
	session.ReviewItems = reviewItems

	return &session, nil
}

func (db *DB) GetStudySession(id int) (*StudySession, error) {
	query := `
		SELECT s.id, s.group_id, s.study_activity_id, s.created_at,
			   g.name, g.description,
			   a.name, a.thumbnail_url, a.launch_url, a.description
		FROM study_sessions s
		LEFT JOIN groups g ON s.group_id = g.id
		LEFT JOIN study_activities a ON s.study_activity_id = a.id
		WHERE s.id = ?`

	var session StudySession
	var group Group
	var activity StudyActivity

	err := db.QueryRow(query, id).Scan(
		&session.ID, &session.GroupID, &session.StudyActivityID, &session.CreatedAt,
		&group.Name, &group.Description,
		&activity.Name, &activity.ThumbnailURL, &activity.LaunchURL, &activity.Description,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("error getting study session: %v", err)
	}

	group.ID = session.GroupID
	activity.ID = session.StudyActivityID
	session.Group = &group
	session.Activity = &activity

	// Get review items
	reviewItems, err := db.GetReviewItemsForSession(session.ID)
	if err != nil {
		return nil, fmt.Errorf("error getting review items: %v", err)
	}
	session.ReviewItems = reviewItems

	return &session, nil
}

func (db *DB) GetStudySessionByID(id int) (*StudySession, error) {
	var session StudySession
	session.Group = &Group{}
	session.Activity = &StudyActivity{}

	var thumbURL, activityDesc sql.NullString

	err := db.QueryRow(`
		SELECT 
			s.id, s.group_id, s.study_activity_id, s.created_at,
			g.name, g.description,
			a.name, a.thumbnail_url, a.launch_url, a.description
		FROM study_sessions s
		JOIN groups g ON s.group_id = g.id
		JOIN study_activities a ON s.study_activity_id = a.id
		WHERE s.id = ?
	`, id).Scan(
		&session.ID,
		&session.GroupID,
		&session.StudyActivityID,
		&session.CreatedAt,
		&session.Group.Name,
		&session.Group.Description,
		&session.Activity.Name,
		&thumbURL,
		&session.Activity.LaunchURL,
		&activityDesc,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error getting study session: %v", err)
	}

	// Handle null strings
	if thumbURL.Valid {
		session.Activity.ThumbnailURL = &thumbURL.String
	}
	if activityDesc.Valid {
		session.Activity.Description = &activityDesc.String
	}

	// Set IDs for nested objects
	session.Group.ID = session.GroupID
	session.Activity.ID = session.StudyActivityID

	// Get review items for this session
	reviewItems, err := db.GetReviewItemsForSession(id)
	if err != nil {
		return nil, fmt.Errorf("error getting review items: %v", err)
	}
	session.ReviewItems = reviewItems

	return &session, nil
}

func (db *DB) GetReviewItemsForSession(sessionID int) ([]WordReviewItem, error) {
	query := `
		SELECT r.id, r.word_id, r.study_session_id, r.correct, r.created_at,
			   w.spanish, w.english, w.type
		FROM word_review_items r
		LEFT JOIN words w ON r.word_id = w.id
		WHERE r.study_session_id = ?`

	rows, err := db.Query(query, sessionID)
	if err != nil {
		return nil, fmt.Errorf("error querying review items: %v", err)
	}
	defer rows.Close()

	var items []WordReviewItem
	for rows.Next() {
		var item WordReviewItem
		var word Word
		err := rows.Scan(
			&item.ID, &item.WordID, &item.StudySessionID, &item.Correct, &item.CreatedAt,
			&word.Spanish, &word.English, &word.Type,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning review item: %v", err)
		}
		word.ID = item.WordID
		item.Word = &word
		items = append(items, item)
	}

	return items, nil
}

func (db *DB) CreateStudySession(session *StudySession) error {
	query := `
		INSERT INTO study_sessions (group_id, study_activity_id, created_at)
		VALUES (?, ?, ?)
		RETURNING id`

	err := db.QueryRow(query, session.GroupID, session.StudyActivityID, time.Now()).Scan(&session.ID)
	if err != nil {
		return fmt.Errorf("error creating study session: %v", err)
	}
	return nil
}

func (db *DB) AddReviewItems(sessionID int, wordIDs []int, correct bool) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("error starting transaction: %v", err)
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`
		INSERT INTO word_review_items (word_id, study_session_id, correct, created_at)
		VALUES (?, ?, ?, ?)`)
	if err != nil {
		return fmt.Errorf("error preparing statement: %v", err)
	}
	defer stmt.Close()

	for _, wordID := range wordIDs {
		_, err = stmt.Exec(wordID, sessionID, correct, time.Now())
		if err != nil {
			return fmt.Errorf("error inserting review item: %v", err)
		}
	}

	return tx.Commit()
}

func (db *DB) GetStudyActivities(groupID int, page, pageSize int) (*PaginatedResponse, error) {
	var total int
	var activities []StudyActivity

	// Get total count
	query := `SELECT COUNT(*) FROM study_activities`
	if groupID > 0 {
		query = `
			SELECT COUNT(DISTINCT sa.id) 
			FROM study_activities sa
			JOIN study_sessions ss ON ss.study_activity_id = sa.id
			WHERE ss.group_id = ?`
	}

	var err error
	if groupID > 0 {
		err = db.QueryRow(query, groupID).Scan(&total)
	} else {
		err = db.QueryRow(query).Scan(&total)
	}
	if err != nil {
		return nil, fmt.Errorf("error getting total study activities: %v", err)
	}

	// Get paginated activities
	query = `
		SELECT id, name, thumbnail_url, launch_url, description
		FROM study_activities
		ORDER BY name
		LIMIT ? OFFSET ?`
	if groupID > 0 {
		query = `
			SELECT DISTINCT sa.id, sa.name, sa.thumbnail_url, sa.launch_url, sa.description
			FROM study_activities sa
			JOIN study_sessions ss ON ss.study_activity_id = sa.id
			WHERE ss.group_id = ?
			ORDER BY sa.name
			LIMIT ? OFFSET ?`
	}

	offset := (page - 1) * pageSize

	var rows *sql.Rows
	if groupID > 0 {
		rows, err = db.Query(query, groupID, pageSize, offset)
	} else {
		rows, err = db.Query(query, pageSize, offset)
	}
	if err != nil {
		return nil, fmt.Errorf("error getting study activities: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var activity StudyActivity
		var thumbURL, desc sql.NullString
		err := rows.Scan(&activity.ID, &activity.Name, &thumbURL, &activity.LaunchURL, &desc)
		if err != nil {
			return nil, fmt.Errorf("error scanning study activity: %v", err)
		}

		if thumbURL.Valid {
			activity.ThumbnailURL = &thumbURL.String
		}
		if desc.Valid {
			activity.Description = &desc.String
		}

		activities = append(activities, activity)
	}

	return &PaginatedResponse{
		TotalItems:   total,
		CurrentPage:  page,
		TotalPages:   (total + pageSize - 1) / pageSize,
		ItemsPerPage: pageSize,
		Items:        activities,
	}, nil
}

func (db *DB) GetWords(groupID int, page, pageSize int) (*PaginatedResponse, error) {
	var total int
	var words []Word

	// Get total count
	query := `SELECT COUNT(*) FROM words`
	if groupID > 0 {
		query = `
			SELECT COUNT(*) 
			FROM words w
			JOIN words_groups wg ON w.id = wg.word_id
			WHERE wg.group_id = ?`
	}

	var err error
	if groupID > 0 {
		err = db.QueryRow(query, groupID).Scan(&total)
	} else {
		err = db.QueryRow(query).Scan(&total)
	}
	if err != nil {
		return nil, fmt.Errorf("error getting total words: %v", err)
	}

	// Get paginated words
	query = `
		SELECT id, spanish, english, type
		FROM words
		ORDER BY spanish
		LIMIT ? OFFSET ?`
	if groupID > 0 {
		query = `
			SELECT w.id, w.spanish, w.english, w.type
			FROM words w
			JOIN words_groups wg ON w.id = wg.word_id
			WHERE wg.group_id = ?
			ORDER BY w.spanish
			LIMIT ? OFFSET ?`
	}

	offset := (page - 1) * pageSize

	var rows *sql.Rows
	if groupID > 0 {
		rows, err = db.Query(query, groupID, pageSize, offset)
	} else {
		rows, err = db.Query(query, pageSize, offset)
	}
	if err != nil {
		return nil, fmt.Errorf("error getting words: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var word Word
		err := rows.Scan(&word.ID, &word.Spanish, &word.English, &word.Type)
		if err != nil {
			return nil, fmt.Errorf("error scanning word: %v", err)
		}
		words = append(words, word)
	}

	return &PaginatedResponse{
		TotalItems:   total,
		CurrentPage:  page,
		TotalPages:   (total + pageSize - 1) / pageSize,
		ItemsPerPage: pageSize,
		Items:        words,
	}, nil
}

func (db *DB) GetWordStats(wordID, groupID int) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Get total reviews and correct/wrong counts
	query := `
		SELECT 
			COUNT(*) as total_reviews,
			SUM(CASE WHEN correct = 1 THEN 1 ELSE 0 END) as correct_reviews,
			SUM(CASE WHEN correct = 0 THEN 1 ELSE 0 END) as wrong_reviews
		FROM word_review_items wri
		WHERE wri.word_id = ?`

	if groupID > 0 {
		query = `
			SELECT 
				COUNT(*) as total_reviews,
				SUM(CASE WHEN wri.correct = 1 THEN 1 ELSE 0 END) as correct_reviews,
				SUM(CASE WHEN wri.correct = 0 THEN 1 ELSE 0 END) as wrong_reviews
			FROM word_review_items wri
			JOIN study_sessions ss ON wri.study_session_id = ss.id
			WHERE wri.word_id = ? AND ss.group_id = ?`
	}

	var totalReviews, correctReviews, wrongReviews int
	var err error

	if groupID > 0 {
		err = db.QueryRow(query, wordID, groupID).Scan(&totalReviews, &correctReviews, &wrongReviews)
	} else {
		err = db.QueryRow(query, wordID).Scan(&totalReviews, &correctReviews, &wrongReviews)
	}
	if err != nil {
		return nil, fmt.Errorf("error getting word stats: %v", err)
	}

	stats["total_reviews"] = totalReviews
	stats["correct_reviews"] = correctReviews
	stats["wrong_reviews"] = wrongReviews

	// Calculate accuracy
	var accuracy float64
	if totalReviews > 0 {
		accuracy = float64(correctReviews) / float64(totalReviews) * 100
	}
	stats["accuracy"] = accuracy

	return stats, nil
}

func (db *DB) GetGroups(page, pageSize int) (*PaginatedResponse, error) {
	offset := (page - 1) * pageSize

	var totalItems int
	err := db.QueryRow("SELECT COUNT(*) FROM groups").Scan(&totalItems)
	if err != nil {
		return nil, err
	}

	rows, err := db.Query(`
		SELECT g.id, g.name, g.description,
			(SELECT COUNT(*) FROM words_groups wg WHERE wg.group_id = g.id) as word_count
		FROM groups g
		ORDER BY g.name
		LIMIT ? OFFSET ?
	`, pageSize, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []Group
	for rows.Next() {
		var group Group
		err := rows.Scan(&group.ID, &group.Name, &group.Description, &group.WordCount)
		if err != nil {
			return nil, err
		}
		groups = append(groups, group)
	}

	totalPages := (totalItems + pageSize - 1) / pageSize

	return &PaginatedResponse{
		TotalItems:   totalItems,
		CurrentPage:  page,
		TotalPages:   totalPages,
		ItemsPerPage: pageSize,
		Items:        groups,
	}, nil
}

func (db *DB) GetGroupByID(id int) (*Group, error) {
	var group Group
	err := db.QueryRow(`
		SELECT g.id, g.name, g.description,
			(SELECT COUNT(*) FROM words_groups wg WHERE wg.group_id = g.id) as word_count
		FROM groups g
		WHERE g.id = ?
	`, id).Scan(&group.ID, &group.Name, &group.Description, &group.WordCount)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &group, nil
}

func (db *DB) GetGroupWords(groupID int) ([]Word, error) {
	rows, err := db.Query(`
		SELECT w.id, w.spanish, w.english, w.type
		FROM words w
		JOIN words_groups wg ON w.id = wg.word_id
		WHERE wg.group_id = ?
		ORDER BY w.spanish
	`, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var words []Word
	for rows.Next() {
		var word Word
		err := rows.Scan(&word.ID, &word.Spanish, &word.English, &word.Type)
		if err != nil {
			return nil, err
		}
		words = append(words, word)
	}

	return words, nil
}

func (db *DB) CreateWord(word *Word) error {
	result, err := db.Exec(`
		INSERT INTO words (spanish, english, type)
		VALUES (?, ?, ?)
	`, word.Spanish, word.English, word.Type)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	word.ID = int(id)
	return nil
}

func (db *DB) UpdateWord(word *Word) error {
	result, err := db.Exec(`
		UPDATE words
		SET spanish = ?, english = ?, type = ?
		WHERE id = ?
	`, word.Spanish, word.English, word.Type, word.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (db *DB) DeleteWord(id int) error {
	// First check if the word exists
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM words WHERE id = ?)", id).Scan(&exists)
	if err != nil {
		return err
	}
	if !exists {
		return sql.ErrNoRows
	}

	// Begin transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Delete from words_groups first (due to foreign key constraint)
	_, err = tx.Exec("DELETE FROM words_groups WHERE word_id = ?", id)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Delete from word_review_items
	_, err = tx.Exec("DELETE FROM word_review_items WHERE word_id = ?", id)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Finally delete the word
	_, err = tx.Exec("DELETE FROM words WHERE id = ?", id)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (db *DB) CreateGroup(group *Group) error {
	result, err := db.Exec(`
		INSERT INTO groups (name, description)
		VALUES (?, ?)
	`, group.Name, group.Description)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	group.ID = int(id)
	return nil
}

func (db *DB) UpdateGroup(group *Group) error {
	result, err := db.Exec(`
		UPDATE groups
		SET name = ?, description = ?
		WHERE id = ?
	`, group.Name, group.Description, group.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (db *DB) DeleteGroup(id int) error {
	// First check if the group exists
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM groups WHERE id = ?)", id).Scan(&exists)
	if err != nil {
		return err
	}
	if !exists {
		return sql.ErrNoRows
	}

	// Begin transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Delete from words_groups first
	_, err = tx.Exec("DELETE FROM words_groups WHERE group_id = ?", id)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Delete the group
	_, err = tx.Exec("DELETE FROM groups WHERE id = ?", id)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (db *DB) AddWordsToGroup(groupID int, wordIDs []int) error {
	// First check if the group exists
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM groups WHERE id = ?)", groupID).Scan(&exists)
	if err != nil {
		return err
	}
	if !exists {
		return sql.ErrNoRows
	}

	// Begin transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Insert each word-group relationship
	for _, wordID := range wordIDs {
		// Check if the word exists
		err := tx.QueryRow("SELECT EXISTS(SELECT 1 FROM words WHERE id = ?)", wordID).Scan(&exists)
		if err != nil {
			tx.Rollback()
			return err
		}
		if !exists {
			tx.Rollback()
			return fmt.Errorf("word with ID %d does not exist", wordID)
		}

		// Check if the relationship already exists
		err = tx.QueryRow(`
			SELECT EXISTS(
				SELECT 1 FROM words_groups 
				WHERE word_id = ? AND group_id = ?
			)`, wordID, groupID).Scan(&exists)
		if err != nil {
			tx.Rollback()
			return err
		}
		if !exists {
			_, err = tx.Exec(`
				INSERT INTO words_groups (word_id, group_id)
				VALUES (?, ?)
			`, wordID, groupID)
			if err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	return tx.Commit()
}

func (db *DB) RemoveWordsFromGroup(groupID int, wordIDs []int) error {
	if len(wordIDs) == 0 {
		return nil
	}

	// Convert wordIDs to string for the IN clause
	wordIDsStr := ""
	for i, id := range wordIDs {
		if i > 0 {
			wordIDsStr += ","
		}
		wordIDsStr += fmt.Sprintf("%d", id)
	}

	// Delete the relationships
	_, err := db.Exec(fmt.Sprintf(`
		DELETE FROM words_groups 
		WHERE group_id = ? AND word_id IN (%s)
	`, wordIDsStr), groupID)

	return err
}

func (db *DB) GetWordByID(id int) (*Word, error) {
	var word Word
	err := db.QueryRow(`
		SELECT id, spanish, english, type
		FROM words
		WHERE id = ?`, id).Scan(&word.ID, &word.Spanish, &word.English, &word.Type)
	
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error getting word: %v", err)
	}
	return &word, nil
}

func (db *DB) GetStudyActivityByID(id int) (*StudyActivity, error) {
	query := `
		SELECT id, name, thumbnail_url, launch_url, description, created_at
		FROM study_activities
		WHERE id = ?`

	var activity StudyActivity
	var thumbURL, desc sql.NullString
	err := db.QueryRow(query, id).Scan(
		&activity.ID,
		&activity.Name,
		&thumbURL,
		&activity.LaunchURL,
		&desc,
		&activity.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("error getting study activity: %v", err)
	}

	if thumbURL.Valid {
		activity.ThumbnailURL = &thumbURL.String
	}
	if desc.Valid {
		activity.Description = &desc.String
	}

	// Get associated study sessions
	sessionsQuery := `
		SELECT s.id, s.group_id, s.created_at,
			   g.name, g.description,
			   COUNT(r.id) as review_items_count
		FROM study_sessions s
		LEFT JOIN groups g ON s.group_id = g.id
		LEFT JOIN word_review_items r ON s.id = r.study_session_id
		WHERE s.study_activity_id = ?
		GROUP BY s.id
		ORDER BY s.created_at DESC`

	rows, err := db.Query(sessionsQuery, id)
	if err != nil {
		return nil, fmt.Errorf("error getting study sessions: %v", err)
	}
	defer rows.Close()

	var sessions []StudySession
	for rows.Next() {
		var session StudySession
		var group Group
		var reviewCount int
		err := rows.Scan(
			&session.ID,
			&session.GroupID,
			&session.CreatedAt,
			&group.Name,
			&group.Description,
			&reviewCount,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning study session: %v", err)
		}
		group.ID = session.GroupID
		session.Group = &group
		session.ReviewItemsCount = reviewCount
		sessions = append(sessions, session)
	}

	activity.Sessions = sessions
	return &activity, nil
}

func (db *DB) GetStudyActivitySessions(activityID, page, pageSize int) (*PaginatedResponse, error) {
	// Get total count
	var totalItems int
	countQuery := `
		SELECT COUNT(DISTINCT s.id)
		FROM study_sessions s
		WHERE s.study_activity_id = ?`
	
	err := db.QueryRow(countQuery, activityID).Scan(&totalItems)
	if err != nil {
		return nil, fmt.Errorf("error getting total sessions count: %v", err)
	}

	// Calculate pagination values
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	// Get paginated sessions
	query := `
		SELECT s.id, s.group_id, s.created_at,
			   g.name, g.description,
			   COUNT(r.id) as review_items_count
		FROM study_sessions s
		LEFT JOIN groups g ON s.group_id = g.id
		LEFT JOIN word_review_items r ON s.id = r.study_session_id
		WHERE s.study_activity_id = ?
		GROUP BY s.id
		ORDER BY s.created_at DESC
		LIMIT ? OFFSET ?`

	rows, err := db.Query(query, activityID, pageSize, offset)
	if err != nil {
		return nil, fmt.Errorf("error getting study sessions: %v", err)
	}
	defer rows.Close()

	var sessions []StudySession
	for rows.Next() {
		var session StudySession
		var group Group
		var reviewCount int
		err := rows.Scan(
			&session.ID,
			&session.GroupID,
			&session.CreatedAt,
			&group.Name,
			&group.Description,
			&reviewCount,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning study session: %v", err)
		}
		group.ID = session.GroupID
		session.Group = &group
		session.ReviewItemsCount = reviewCount
		sessions = append(sessions, session)
	}

	totalPages := (totalItems + pageSize - 1) / pageSize
	return &PaginatedResponse{
		Items:        sessions,
		CurrentPage:  page,
		TotalPages:   totalPages,
		TotalItems:   totalItems,
		ItemsPerPage: pageSize,
	}, nil
}

func (db *DB) CreateStudyActivity(activity *StudyActivity) error {
	query := `
		INSERT INTO study_activities (name, thumbnail_url, launch_url, description)
		VALUES (?, ?, ?, ?)`

	result, err := db.Exec(query, activity.Name, activity.ThumbnailURL, activity.LaunchURL, activity.Description)
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
