package models

import (
	"database/sql"
	"fmt"
)

// CreateWord creates a new word in the database
func (db *DB) CreateWord(word *Word) error {
	query := `
		INSERT INTO words (spanish, english, part_of_speech)
		VALUES (?, ?, ?)
	`

	result, err := db.Exec(query, word.Spanish, word.English, word.PartOfSpeech)
	if err != nil {
		return fmt.Errorf("error creating word: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("error getting last insert id: %v", err)
	}

	word.ID = int(id)
	return nil
}

// GetWord retrieves a word by ID
func (db *DB) GetWord(id int) (*Word, error) {
	word := &Word{}
	query := `SELECT id, spanish, english, part_of_speech FROM words WHERE id = ?`

	err := db.QueryRow(query, id).Scan(&word.ID, &word.Spanish, &word.English, &word.PartOfSpeech)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error getting word: %v", err)
	}

	return word, nil
}

// UpdateWord updates an existing word
func (db *DB) UpdateWord(word *Word) error {
	query := `
		UPDATE words
		SET spanish = ?, english = ?, part_of_speech = ?
		WHERE id = ?
	`

	result, err := db.Exec(query, word.Spanish, word.English, word.PartOfSpeech, word.ID)
	if err != nil {
		return fmt.Errorf("error updating word: %v", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %v", err)
	}
	if rows == 0 {
		return fmt.Errorf("word not found")
	}

	return nil
}

// DeleteWord deletes a word by ID
func (db *DB) DeleteWord(id int) error {
	query := `DELETE FROM words WHERE id = ?`

	result, err := db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error deleting word: %v", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %v", err)
	}
	if rows == 0 {
		return fmt.Errorf("word not found")
	}

	return nil
}

// ListWords returns a paginated list of words
func (db *DB) ListWords(page, pageSize int, groupID *int) ([]Word, int, error) {
	var words []Word
	var total int

	// Get total count
	countQuery := `SELECT COUNT(*) FROM words`
	if groupID != nil {
		countQuery = `
			SELECT COUNT(DISTINCT w.id)
			FROM words w
			JOIN word_groups wg ON w.id = wg.word_id
			WHERE wg.group_id = ?
		`
	}

	var err error
	if groupID != nil {
		err = db.QueryRow(countQuery, *groupID).Scan(&total)
	} else {
		err = db.QueryRow(countQuery).Scan(&total)
	}
	if err != nil {
		return nil, 0, fmt.Errorf("error getting total count: %v", err)
	}

	// Get paginated words
	query := `
		SELECT id, spanish, english, part_of_speech
		FROM words
	`
	if groupID != nil {
		query = `
			SELECT DISTINCT w.id, w.spanish, w.english, w.part_of_speech
			FROM words w
			JOIN word_groups wg ON w.id = wg.word_id
			WHERE wg.group_id = ?
		`
	}
	query += ` LIMIT ? OFFSET ?`

	offset := (page - 1) * pageSize
	var rows *sql.Rows
	if groupID != nil {
		rows, err = db.Query(query, *groupID, pageSize, offset)
	} else {
		rows, err = db.Query(query, pageSize, offset)
	}
	if err != nil {
		return nil, 0, fmt.Errorf("error querying words: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var word Word
		err := rows.Scan(&word.ID, &word.Spanish, &word.English, &word.PartOfSpeech)
		if err != nil {
			return nil, 0, fmt.Errorf("error scanning word: %v", err)
		}
		words = append(words, word)
	}

	return words, total, nil
}

// GetWordStats retrieves statistics for a word
func (db *DB) GetWordStats(wordID int) (int, int, error) {
	// First verify the word exists
	var exists bool
	err := db.QueryRow(`SELECT EXISTS(SELECT 1 FROM words WHERE id = ?)`, wordID).Scan(&exists)
	if err != nil {
		return 0, 0, fmt.Errorf("error checking word existence: %v", err)
	}
	if !exists {
		return 0, 0, fmt.Errorf("word not found")
	}

	// Get review statistics
	var totalReviews, correctReviews int
	query := `
		SELECT 
			COUNT(*) as total_reviews,
			COALESCE(SUM(CASE WHEN correct = 1 THEN 1 ELSE 0 END), 0) as correct_reviews
		FROM word_review_items
		WHERE word_id = ?
	`
	err = db.QueryRow(query, wordID).Scan(&totalReviews, &correctReviews)
	if err != nil {
		return 0, 0, fmt.Errorf("error getting word stats: %v", err)
	}

	return totalReviews, correctReviews, nil
}
