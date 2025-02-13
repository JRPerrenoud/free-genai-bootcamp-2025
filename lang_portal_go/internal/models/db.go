package models

import (
	"database/sql"
	"fmt"
	"strings"
)

type DB struct {
	*sql.DB
}

func NewDB(db *sql.DB) *DB {
	return &DB{db}
}

// Word methods
func (db *DB) GetWords(page, perPage int) (*PaginatedResponse[Word], error) {
	offset := (page - 1) * perPage
	
	var total int
	err := db.QueryRow("SELECT COUNT(*) FROM words").Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("error counting words: %v", err)
	}

	rows, err := db.Query(`
		SELECT id, spanish, english, parts, correct_count, wrong_count
		FROM words
		LIMIT ? OFFSET ?
	`, perPage, offset)
	if err != nil {
		return nil, fmt.Errorf("error querying words: %v", err)
	}
	defer rows.Close()

	var words []Word
	for rows.Next() {
		var w Word
		var partsStr string
		err := rows.Scan(&w.ID, &w.Spanish, &w.English, &partsStr, &w.CorrectCount, &w.WrongCount)
		if err != nil {
			return nil, fmt.Errorf("error scanning word: %v", err)
		}
		w.Parts = strings.Split(partsStr, ",")
		words = append(words, w)
	}

	return &PaginatedResponse[Word]{
		Items: words,
		Total: total,
		Page:  page,
	}, nil
}

func (db *DB) CreateWord(w *Word) error {
	partsStr := strings.Join(w.Parts, ",")
	result, err := db.Exec(`
		INSERT INTO words (spanish, english, parts, correct_count, wrong_count)
		VALUES (?, ?, ?, ?, ?)
	`, w.Spanish, w.English, partsStr, w.CorrectCount, w.WrongCount)
	if err != nil {
		return fmt.Errorf("error creating word: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("error getting last insert id: %v", err)
	}

	w.ID = id
	return nil
}

func (db *DB) UpdateWord(w *Word) error {
	partsStr := strings.Join(w.Parts, ",")
	_, err := db.Exec(`
		UPDATE words
		SET spanish = ?, english = ?, parts = ?, correct_count = ?, wrong_count = ?
		WHERE id = ?
	`, w.Spanish, w.English, partsStr, w.CorrectCount, w.WrongCount, w.ID)
	if err != nil {
		return fmt.Errorf("error updating word: %v", err)
	}
	return nil
}

func (db *DB) DeleteWord(id int64) error {
	_, err := db.Exec("DELETE FROM words WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("error deleting word: %v", err)
	}
	return nil
}

// Group methods
func (db *DB) GetGroups(page, perPage int) (*PaginatedResponse[Group], error) {
	offset := (page - 1) * perPage
	
	var total int
	err := db.QueryRow("SELECT COUNT(*) FROM groups").Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("error counting groups: %v", err)
	}

	rows, err := db.Query(`
		SELECT g.id, g.name, COUNT(gw.word_id) as word_count
		FROM groups g
		LEFT JOIN group_words gw ON g.id = gw.group_id
		GROUP BY g.id, g.name
		LIMIT ? OFFSET ?
	`, perPage, offset)
	if err != nil {
		return nil, fmt.Errorf("error querying groups: %v", err)
	}
	defer rows.Close()

	var groups []Group
	for rows.Next() {
		var g Group
		err := rows.Scan(&g.ID, &g.Name, &g.WordCount)
		if err != nil {
			return nil, fmt.Errorf("error scanning group: %v", err)
		}
		groups = append(groups, g)
	}

	return &PaginatedResponse[Group]{
		Items: groups,
		Total: total,
		Page:  page,
	}, nil
}

func (db *DB) CreateGroup(g *Group) error {
	result, err := db.Exec("INSERT INTO groups (name) VALUES (?)", g.Name)
	if err != nil {
		return fmt.Errorf("error creating group: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("error getting last insert id: %v", err)
	}

	g.ID = id
	return nil
}

func (db *DB) UpdateGroup(g *Group) error {
	_, err := db.Exec("UPDATE groups SET name = ? WHERE id = ?", g.Name, g.ID)
	if err != nil {
		return fmt.Errorf("error updating group: %v", err)
	}
	return nil
}

func (db *DB) DeleteGroup(id int64) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("error starting transaction: %v", err)
	}
	defer tx.Rollback()

	// Delete group_words entries first
	_, err = tx.Exec("DELETE FROM group_words WHERE group_id = ?", id)
	if err != nil {
		return fmt.Errorf("error deleting group words: %v", err)
	}

	// Delete the group
	_, err = tx.Exec("DELETE FROM groups WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("error deleting group: %v", err)
	}

	return tx.Commit()
}

func (db *DB) GetGroupWords(groupID int64, page, perPage int) (*PaginatedResponse[Word], error) {
	offset := (page - 1) * perPage
	
	var total int
	err := db.QueryRow(`
		SELECT COUNT(*)
		FROM group_words gw
		JOIN words w ON gw.word_id = w.id
		WHERE gw.group_id = ?
	`, groupID).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("error counting group words: %v", err)
	}

	rows, err := db.Query(`
		SELECT w.id, w.spanish, w.english, w.parts, w.correct_count, w.wrong_count
		FROM group_words gw
		JOIN words w ON gw.word_id = w.id
		WHERE gw.group_id = ?
		LIMIT ? OFFSET ?
	`, groupID, perPage, offset)
	if err != nil {
		return nil, fmt.Errorf("error querying group words: %v", err)
	}
	defer rows.Close()

	var words []Word
	for rows.Next() {
		var w Word
		var partsStr string
		err := rows.Scan(&w.ID, &w.Spanish, &w.English, &partsStr, &w.CorrectCount, &w.WrongCount)
		if err != nil {
			return nil, fmt.Errorf("error scanning word: %v", err)
		}
		w.Parts = strings.Split(partsStr, ",")
		words = append(words, w)
	}

	return &PaginatedResponse[Word]{
		Items: words,
		Total: total,
		Page:  page,
	}, nil
}

func (db *DB) AddWordsToGroup(groupID int64, wordIDs []int64) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("error starting transaction: %v", err)
	}
	defer tx.Rollback()

	for _, wordID := range wordIDs {
		_, err := tx.Exec(`
			INSERT INTO group_words (group_id, word_id)
			VALUES (?, ?)
			ON CONFLICT (group_id, word_id) DO NOTHING
		`, groupID, wordID)
		if err != nil {
			return fmt.Errorf("error adding word to group: %v", err)
		}
	}

	return tx.Commit()
}

func (db *DB) RemoveWordsFromGroup(groupID int64, wordIDs []int64) error {
	if len(wordIDs) == 0 {
		return nil
	}

	query := `
		DELETE FROM group_words
		WHERE group_id = ? AND word_id IN (?` + strings.Repeat(",?", len(wordIDs)-1) + ")"

	args := make([]interface{}, len(wordIDs)+1)
	args[0] = groupID
	for i, id := range wordIDs {
		args[i+1] = id
	}

	_, err := db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("error removing words from group: %v", err)
	}
	return nil
}

// Study activity methods
func (db *DB) GetStudyActivities(page, perPage int) (*PaginatedResponse[StudyActivity], error) {
	offset := (page - 1) * perPage
	
	var total int
	err := db.QueryRow("SELECT COUNT(*) FROM study_activities").Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("error counting study activities: %v", err)
	}

	rows, err := db.Query(`
		SELECT id, name, description
		FROM study_activities
		LIMIT ? OFFSET ?
	`, perPage, offset)
	if err != nil {
		return nil, fmt.Errorf("error querying study activities: %v", err)
	}
	defer rows.Close()

	var activities []StudyActivity
	for rows.Next() {
		var a StudyActivity
		err := rows.Scan(&a.ID, &a.Name, &a.Description)
		if err != nil {
			return nil, fmt.Errorf("error scanning study activity: %v", err)
		}
		activities = append(activities, a)
	}

	return &PaginatedResponse[StudyActivity]{
		Items: activities,
		Total: total,
		Page:  page,
	}, nil
}

// Study session methods
func (db *DB) GetStudySessions(page, perPage int) (*PaginatedResponse[StudySession], error) {
	offset := (page - 1) * perPage
	
	var total int
	err := db.QueryRow("SELECT COUNT(*) FROM study_sessions").Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("error counting study sessions: %v", err)
	}

	rows, err := db.Query(`
		SELECT s.id, s.activity_id, s.start_time, s.end_time,
			   COUNT(wr.id) as review_count,
			   SUM(CASE WHEN wr.is_correct THEN 1 ELSE 0 END) as correct_count
		FROM study_sessions s
		LEFT JOIN word_reviews wr ON s.id = wr.study_session_id
		GROUP BY s.id, s.activity_id, s.start_time, s.end_time
		ORDER BY s.start_time DESC
		LIMIT ? OFFSET ?
	`, perPage, offset)
	if err != nil {
		return nil, fmt.Errorf("error querying study sessions: %v", err)
	}
	defer rows.Close()

	var sessions []StudySession
	for rows.Next() {
		var s StudySession
		err := rows.Scan(&s.ID, &s.ActivityID, &s.StartTime, &s.EndTime, &s.ReviewCount, &s.CorrectCount)
		if err != nil {
			return nil, fmt.Errorf("error scanning study session: %v", err)
		}
		sessions = append(sessions, s)
	}

	return &PaginatedResponse[StudySession]{
		Items: sessions,
		Total: total,
		Page:  page,
	}, nil
}

func (db *DB) CreateStudySession(s *StudySession) error {
	result, err := db.Exec(`
		INSERT INTO study_sessions (activity_id, start_time, end_time)
		VALUES (?, ?, ?)
	`, s.ActivityID, s.StartTime, s.EndTime)
	if err != nil {
		return fmt.Errorf("error creating study session: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("error getting last insert id: %v", err)
	}

	s.ID = id
	return nil
}

func (db *DB) CreateWordReviews(reviews []WordReview) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("error starting transaction: %v", err)
	}
	defer tx.Rollback()

	for _, r := range reviews {
		_, err := tx.Exec(`
			INSERT INTO word_reviews (study_session_id, word_id, is_correct)
			VALUES (?, ?, ?)
		`, r.StudySessionID, r.WordID, r.IsCorrect)
		if err != nil {
			return fmt.Errorf("error creating word review: %v", err)
		}

		// Update word statistics
		if r.IsCorrect {
			_, err = tx.Exec("UPDATE words SET correct_count = correct_count + 1 WHERE id = ?", r.WordID)
		} else {
			_, err = tx.Exec("UPDATE words SET wrong_count = wrong_count + 1 WHERE id = ?", r.WordID)
		}
		if err != nil {
			return fmt.Errorf("error updating word statistics: %v", err)
		}
	}

	return tx.Commit()
}
