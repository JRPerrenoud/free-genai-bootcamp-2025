package seeder

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"lang_portal_go/internal/models"
)

// SeedData represents the structure of our seed JSON file
type SeedData struct {
	Groups          []GroupSeed          `json:"groups"`
	Words           []WordSeed           `json:"words"`
	StudyActivities []StudyActivitySeed  `json:"study_activities"`
	StudySessions   []StudySessionSeed   `json:"study_sessions"`
	WordReviewItems []WordReviewItemSeed `json:"word_review_items"`
}

type GroupSeed struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type WordSeed struct {
	Spanish    string   `json:"spanish"`
	English    string   `json:"english"`
	GroupNames []string `json:"group_names"`
}

type StudyActivitySeed struct {
	Name         string `json:"name"`
	Description  string `json:"description"`
	ThumbnailURL string `json:"thumbnail_url"`
	LaunchURL    string `json:"launch_url"`
}

type StudySessionSeed struct {
	GroupID         int    `json:"group_id"`
	StudyActivityID int    `json:"study_activity_id"`
	CreatedAt       string `json:"created_at"`
	StartTime       string `json:"start_time"`
	EndTime         string `json:"end_time"`
}

type WordReviewItemSeed struct {
	WordID         int    `json:"word_id"`
	StudySessionID int    `json:"study_session_id"`
	Correct        bool   `json:"correct"`
	CreatedAt      string `json:"created_at"`
}

// LoadSeedData reads and parses the seed JSON file
func LoadSeedData(filePath string) (*SeedData, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error opening seed file: %v", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("error reading seed file: %v", err)
	}

	var seedData SeedData
	if err := json.Unmarshal(data, &seedData); err != nil {
		return nil, fmt.Errorf("error parsing seed data: %v", err)
	}

	return &seedData, nil
}

// SeedDatabase populates the database with initial data
func SeedDatabase(db *models.DB, seedData *SeedData) error {
	// Create groups first
	groupNameToID := make(map[string]int)
	for _, groupSeed := range seedData.Groups {
		group := &models.Group{
			Name:        groupSeed.Name,
			Description: groupSeed.Description,
		}
		if err := db.CreateGroup(group); err != nil {
			return fmt.Errorf("error creating group %s: %v", group.Name, err)
		}
		groupNameToID[group.Name] = group.ID
	}

	// Create words and associate them with groups
	for _, wordSeed := range seedData.Words {
		word := &models.Word{
			Spanish: wordSeed.Spanish,
			English: wordSeed.English,
		}
		if err := db.CreateWord(word); err != nil {
			return fmt.Errorf("error creating word %s: %v", word.Spanish, err)
		}

		// Associate word with groups
		for _, groupName := range wordSeed.GroupNames {
			groupID, exists := groupNameToID[groupName]
			if !exists {
				return fmt.Errorf("group %s not found for word %s", groupName, word.Spanish)
			}
			if err := db.AddWordToGroup(word.ID, groupID); err != nil {
				return fmt.Errorf("error adding word %s to group %s: %v", word.Spanish, groupName, err)
			}
		}
	}

	// Create study activities
	for _, activitySeed := range seedData.StudyActivities {
		activity := &models.StudyActivity{
			Name:        activitySeed.Name,
			Description: activitySeed.Description,
		}
		if err := db.CreateStudyActivity(activity); err != nil {
			return fmt.Errorf("error creating study activity %s: %v", activity.Name, err)
		}
	}

	// Create study sessions
	for _, sessionSeed := range seedData.StudySessions {
		_, err := db.CreateStudySession(sessionSeed.GroupID, sessionSeed.StudyActivityID)
		if err != nil {
			return fmt.Errorf("error creating study session: %v", err)
		}
	}

	// Create word review items
	for _, reviewSeed := range seedData.WordReviewItems {
		reviewItem, err := db.CreateWordReviewItem(reviewSeed.WordID, reviewSeed.StudySessionID, reviewSeed.Correct)
		if err != nil {
			return fmt.Errorf("error creating word review item: %v", err)
		}
		_ = reviewItem
	}

	return nil
}
