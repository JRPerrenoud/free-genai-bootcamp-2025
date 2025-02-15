package models

import (
	"database/sql"
	"fmt"
)

// CreateGroup creates a new group
func (db *DB) CreateGroup(group *Group) error {
	query := `
		INSERT INTO groups (name, description)
		VALUES (?, ?)
	`

	result, err := db.Exec(query, group.Name, group.Description)
	if err != nil {
		return fmt.Errorf("error creating group: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("error getting last insert id: %v", err)
	}

	group.ID = int(id)
	return nil
}

// GetGroup retrieves a group by ID
func (db *DB) GetGroup(id int) (*Group, error) {
	group := &Group{}
	query := `SELECT id, name, description FROM groups WHERE id = ?`
	
	err := db.QueryRow(query, id).Scan(&group.ID, &group.Name, &group.Description)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error getting group: %v", err)
	}

	return group, nil
}

// UpdateGroup updates an existing group
func (db *DB) UpdateGroup(group *Group) error {
	query := `
		UPDATE groups
		SET name = ?, description = ?
		WHERE id = ?
	`

	result, err := db.Exec(query, group.Name, group.Description, group.ID)
	if err != nil {
		return fmt.Errorf("error updating group: %v", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %v", err)
	}
	if rows == 0 {
		return fmt.Errorf("group not found")
	}

	return nil
}

// DeleteGroup deletes a group by ID
func (db *DB) DeleteGroup(id int) error {
	query := `DELETE FROM groups WHERE id = ?`
	
	result, err := db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error deleting group: %v", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %v", err)
	}
	if rows == 0 {
		return fmt.Errorf("group not found")
	}

	return nil
}

// ListGroups returns a paginated list of groups
func (db *DB) ListGroups(page, pageSize int) ([]Group, int, error) {
	var groups []Group
	var total int

	// Get total count
	err := db.QueryRow(`SELECT COUNT(*) FROM groups`).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("error getting total count: %v", err)
	}

	// Get paginated groups
	query := `
		SELECT id, name, description
		FROM groups
		LIMIT ? OFFSET ?
	`

	offset := (page - 1) * pageSize
	rows, err := db.Query(query, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("error querying groups: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var group Group
		err := rows.Scan(&group.ID, &group.Name, &group.Description)
		if err != nil {
			return nil, 0, fmt.Errorf("error scanning group: %v", err)
		}
		groups = append(groups, group)
	}

	return groups, total, nil
}

// AddWordToGroup adds a word to a group
func (db *DB) AddWordToGroup(wordID, groupID int) error {
	query := `INSERT INTO word_groups (word_id, group_id) VALUES (?, ?)`
	
	result, err := db.Exec(query, wordID, groupID)
	if err != nil {
		return fmt.Errorf("error adding word to group: %v", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %v", err)
	}
	if rows == 0 {
		return fmt.Errorf("failed to add word to group")
	}

	return nil
}

// RemoveWordFromGroup removes a word from a group
func (db *DB) RemoveWordFromGroup(wordID, groupID int) error {
	query := `DELETE FROM word_groups WHERE word_id = ? AND group_id = ?`
	
	result, err := db.Exec(query, wordID, groupID)
	if err != nil {
		return fmt.Errorf("error removing word from group: %v", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %v", err)
	}
	if rows == 0 {
		return fmt.Errorf("word not found in group")
	}

	return nil
}

// GetGroupWords retrieves all words in a group
func (db *DB) GetGroupWords(groupID int, page, pageSize int) ([]Word, int, error) {
	var words []Word
	var total int

	// Get total count
	countQuery := `
		SELECT COUNT(*) 
		FROM words w 
		JOIN word_groups wg ON w.id = wg.word_id 
		WHERE wg.group_id = ?
	`
	err := db.QueryRow(countQuery, groupID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("error getting total count: %v", err)
	}

	// Get words
	query := `
		SELECT w.id, w.spanish, w.english
		FROM words w 
		JOIN word_groups wg ON w.id = wg.word_id 
		WHERE wg.group_id = ? 
		ORDER BY w.spanish 
		LIMIT ? OFFSET ?
	`
	offset := (page - 1) * pageSize
	rows, err := db.Query(query, groupID, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("error querying group words: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var word Word
		err := rows.Scan(&word.ID, &word.Spanish, &word.English)
		if err != nil {
			return nil, 0, fmt.Errorf("error scanning word: %v", err)
		}
		words = append(words, word)
	}

	return words, total, nil
}
