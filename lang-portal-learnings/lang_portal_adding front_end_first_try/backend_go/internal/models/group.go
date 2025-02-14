package models

import (
	"database/sql"
	"fmt"
)

// Group represents a collection of words with the same part of speech
type Group struct {
	Name      string `json:"name"`       // The part of speech
	WordCount int    `json:"word_count"` // Number of words in this group
}

// ListGroups returns a list of all parts of speech and their word counts
func (db *DB) ListGroups() ([]Group, error) {
	query := `
		SELECT part_of_speech, COUNT(*) as word_count
		FROM words
		GROUP BY part_of_speech
		ORDER BY part_of_speech ASC
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error listing groups: %v", err)
	}
	defer rows.Close()

	var groups []Group
	for rows.Next() {
		var group Group
		if err := rows.Scan(&group.Name, &group.WordCount); err != nil {
			return nil, fmt.Errorf("error scanning group: %v", err)
		}
		groups = append(groups, group)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating groups: %v", err)
	}

	return groups, nil
}

// GetGroupWords retrieves all words with a specific part of speech
func (db *DB) GetGroupWords(partOfSpeech string, page, pageSize int) ([]Word, int, error) {
	offset := (page - 1) * pageSize

	// Get total count
	var totalCount int
	countQuery := `SELECT COUNT(*) FROM words WHERE part_of_speech = ?`
	err := db.QueryRow(countQuery, partOfSpeech).Scan(&totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("error getting total count: %v", err)
	}

	// Get paginated words
	query := `
		SELECT id, spanish, english, part_of_speech
		FROM words
		WHERE part_of_speech = ?
		ORDER BY spanish ASC
		LIMIT ? OFFSET ?
	`

	rows, err := db.Query(query, partOfSpeech, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("error getting group words: %v", err)
	}
	defer rows.Close()

	var words []Word
	for rows.Next() {
		var word Word
		if err := rows.Scan(&word.ID, &word.Spanish, &word.English, &word.PartOfSpeech); err != nil {
			return nil, 0, fmt.Errorf("error scanning word: %v", err)
		}
		words = append(words, word)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating words: %v", err)
	}

	return words, totalCount, nil
}

// GetGroupStats returns statistics for a specific part of speech group
func (db *DB) GetGroupStats(partOfSpeech string) (*Group, error) {
	var group Group
	query := `
		SELECT part_of_speech, COUNT(*) as word_count
		FROM words
		WHERE part_of_speech = ?
		GROUP BY part_of_speech
	`
	
	err := db.QueryRow(query, partOfSpeech).Scan(&group.Name, &group.WordCount)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error getting group stats: %v", err)
	}

	return &group, nil
}
