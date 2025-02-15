# Backend Server Technical Spec

## Business Goal

A language learning school wants to build a prototype learning portal that serves three main purposes:
1. Maintain an inventory of vocabulary words that can be learned
2. Act as a Learning Record Store (LRS), tracking correct and incorrect responses during vocabulary practice
3. Provide a unified platform for launching different learning activities

## Technical Requirements

- Backend Framework: Go with Gin web framework
- Database: SQLite3 (file: `words.db`)
- Task Runner: Mage for database management and server operations
- API Format: JSON responses with standardized success/error format
- Authentication: None (single-user application)
- Response Format:
  ```json
  {
    "success": true,
    "data": {
      // Response data here
    }
  }
  ```
  or for errors:
  ```json
  {
    "success": false,
    "error": "Error message here"
  }
  ```
- Pagination: Default 20 items per page where applicable

## Directory Structure

```text
lang_portal_go/
├── cmd/
│   ├── server/        # Main web server application
│   ├── init_db/       # Database initialization with CLI flags
│   ├── initdb/        # Simple database initialization
│   └── seed/          # Database seeding utility
├── internal/
│   ├── models/        # Data structures and database operations
│   ├── handlers/      # HTTP handlers organized by feature
│   ├── seeder/        # Seeding logic and data loading
│   └── services/      # Business logic
├── db/
│   ├── migrations/    # Database schema and migrations
│   └── seeds/         # Initial data for seeding
├── magefile.go        # Task runner configuration
├── go.mod            # Go module definition
├── go.sum            # Go module checksums
└── words.db         # SQLite database file
```

## Database Schema

The application uses SQLite as its database. The schema is defined in `db/migrations/001_initial_schema.sql`.

### Tables

- words - stored vocabulary words
    - id integer
    - spanish string
    - english string
    - parts json

- words_groups - join table for words and groups
many-to-many
    - id integer
    - word_id integer
    - group_id integer
    
- groups - thematic groups of words
    -id integer
    - name string

- study_sessions - records of study sessions grouping word_review_items
    - id integer
    - group_id integer
    - created_at datetime
    - study_activity_id integer
    
- study_activities - a specific study activity, linking study session to group
    - id integer
    - study_session_id integar
    - group_id integer
    - created_at datetime

- word_review_items - a record of word practice, determining if the word was correct or not
    - word_id integer
    - study_session_id integer
    - correct boolean
    - created_at datetime

## API Endpoints

All API responses follow this standard format:
```json
{
    "success": boolean,
    "data": object (optional),
    "error": string (optional)
}
```

For endpoints that return paginated data, the response data will follow this format:
```json
{
    "items": array,
    "current_page": integer,
    "total_pages": integer,
    "total_items": integer,
    "items_per_page": integer
}
```

### Dashboard Endpoints

#### GET /api/dashboard/last_study_session
Returns information about the most recent study session.

#### JSON Response
```json
{
    "success": true,
    "data": {
        "id": 123,
        "group_id": 456,
        "created_at": "2025-02-12T11:10:14-07:00",
        "study_activity_id": 789,
        "group_name": "Verb"
    }
}
```

#### GET /api/dashboard/study_progress
Returns study progress information.
Please note that the frontend will determine progress by based on total words studied and total available words.


#### JSON Response
```json
{
    "success": true,
    "data": {
        "total_words_studied": 500,
        "total_available_words": 1000
    }
}
```

#### GET /api/dashboard/quick_stats
Returns quick overview statistics.

#### JSON Response
```json
{
    "success": true,
    "data": {
        "total_words_studied": 1000,
        "total_words_correct": 800,
        "total_groups": 10,
        "total_study_sessions": 25,
        "overall_accuracy": 0.8,
        "study_streak": 5
    }
}
```

### GET /api/study_activities/:id
Returns details of a study activity.

#### JSON Response
```json
{
    "id": 123,
    "name": "Vocabulary Quiz",
    "thumbnail_url": "https://example.com/thumbnail.jpg",
    "description": "This is a quiz about vocabulary"
}
```



### GET /api/study_activities/:id/study_sessions
Pagination of study sessions (20 per page)

#### JSON Response
```json
{
    "items": [
     {
        "id": 789,
        "activity_name": "Vocabulary Quiz",
        "group_name": "Basic Vocabulary",
        "start_time": "2025-02-12T11:10:14-07:00",
        "end_time": "2025-02-12T11:25:14-07:00",
        "review_items_count": 10
     }
    ],
    "pagination": {
        "current_page": 1,
        "total_pages": 5,
        "total_items": 100,
        "items_per_page": 20
    }
}
```

### POST /api/study_activities/

#### Request Params
- group_id integer
- study_activity_id integer

#### JSON Response
```json
{
   "id": 123,
   "group_id": 456   
}
```

### Word Endpoints

#### GET /api/words
Returns a paginated list of words. Default page size is 20 items.

Query Parameters:
- `page`: Page number (default: 1)
- `page_size`: Number of items per page (default: 20)
- `group_id`: Optional, filter words by group ID

#### JSON Response
```json
{
    "success": true,
    "data": {
        "items": [
            {
                "id": 1,
                "spanish": "hola",
                "english": "hello",
                "part_of_speech": "interjection"
            }
        ],
        "current_page": 1,
        "total_pages": 10,
        "total_items": 100,
        "items_per_page": 20
    }
}
```

#### GET /api/words/:id
Returns details for a specific word including study statistics.

#### JSON Response
```json
{
    "success": true,
    "data": {
        "word": {
            "id": 1,
            "spanish": "hola",
            "english": "hello",
            "part_of_speech": "interjection"
        },
        "study_stats": {
            "total_reviews": 10,
            "correct_reviews": 8
        }
    }
}
```

#### POST /api/words
Creates a new word.

#### Request Body
```json
{
    "spanish": "hola",
    "english": "hello",
    "part_of_speech": "interjection"
}
```

#### JSON Response
```json
{
    "success": true,
    "data": {
        "id": 1,
        "spanish": "hola",
        "english": "hello",
        "part_of_speech": "interjection"
    }
}
```

#### PUT /api/words/:id
Updates an existing word.

#### Request Body
```json
{
    "spanish": "hola",
    "english": "hello",
    "part_of_speech": "interjection"
}
```

#### JSON Response
```json
{
    "success": true,
    "data": {
        "id": 1,
        "spanish": "hola",
        "english": "hello",
        "part_of_speech": "interjection"
    }
}
```

#### DELETE /api/words/:id
Deletes a word.

#### JSON Response
```json
{
    "success": true,
    "data": {
        "message": "Word deleted successfully"
    }
}
```

### Group Endpoints

#### GET /api/groups/:id
Returns details for a specific group including its words.

#### JSON Response
```json
{
    "success": true,
    "data": {
        "group": {
            "id": 1,
            "name": "interjection"            
        },
        "words": [
            {
                "id": 1,
                "spanish": "hola",
                "english": "hello",
                "part_of_speech": "interjection"
            }
        ],
        "total_words": 25
    }
}
```

#### POST /api/groups
Creates a new group.

#### Request Body
```json
{
    "name": "noun"    
}
```

#### JSON Response
```json
{
    "success": true,
    "data": {
        "id": 1,
        "name": "Basic Vocabulary"        
    }
}
```

#### PUT /api/groups/:id
Updates an existing group.

#### Request Body
```json
{
    "name": "verb"
}
```

#### JSON Response
```json
    "success": true,
    "data": {
        "id": 1,
        "name": "verb"        
    }
}
```

#### DELETE /api/groups/:id
Deletes a group.

#### JSON Response
```json
{
    "success": true,
    "data": {
        "message": "Group deleted successfully"
    }
}
```

#### POST /api/groups/:id/words
Adds a word to a group.

#### Request Body
```json
{
    "word_id": 1
}
```

#### JSON Response
```json
{
    "success": true,
    "data": {
        "message": "Word added to group successfully"
    }
}
```

#### DELETE /api/groups/:id/words/:word_id
Removes a word from a group.

#### JSON Response
```json
{
    "success": true,
    "data": {
        "message": "Word removed from group successfully"
    }
}
```

### GET /api/groups/:id
This endpoint will return a single group

#### JSON Response
```json
{
    "id": 1,
    "name": "noun",
    "statistics": {
        "total_word_count": 100
    }
}
```


### GET /api/groups/:id/words
This endpoint will return a list of words in a group

#### JSON Response
```json
{
    "items": [
     {
        "spanish": "hola",
        "english": "hello",
        "correct_count": 10,
        "wrong_count": 5,
        "parts": ["Noun"]
     }
    ],
    "pagination": {
        "current_page": 1,
        "total_pages": 10,
        "total_items": 100,
        "items_per_page": 100
    }
}
```


### GET /api/groups/:id/study_sessions
This endpoint will return a list of study sessions in a group

#### JSON Response
```json
{
    "items": [
     {
        "id": 789,
        "activity_name": "Vocabulary Quiz",
        "group_name": "Basic Vocabulary",
        "start_time": "2025-02-12T11:10:14-07:00",
        "end_time": "2025-02-12T11:25:14-07:00",
        "review_items_count": 20
     }
    ],
    "pagination": {
        "current_page": 1,
        "total_pages": 1,
        "total_items": 5,
        "items_per_page": 100
    }
}
```


### GET /api/study_sessions/
This endpoint will return a list of study sessions

#### JSON Response
```json
{
    "items": [
     {
        "id": 789,
        "activity_name": "Vocabulary Quiz",
        "group_name": "Basic Vocabulary",
        "start_time": "2025-02-12T11:10:14-07:00",
        "end_time": "2025-02-12T11:25:14-07:00",
        "review_items_count": 20        
     }
    ],
    "pagination": {
        "current_page": 1,
        "total_pages": 1,
        "total_items": 100,
        "items_per_page": 100
    }
}
```


### GET /api/study_sessions/:id
This endpoint will return a single study session

#### JSON Response
```json
{
    "id": 789,
    "activity_name": "Vocabulary Quiz",
    "group_name": "Basic Vocabulary",
    "start_time": "2025-02-12T11:10:14-07:00",
    "end_time": "2025-02-12T11:25:14-07:00",
    "review_items_count": 10
}
```


### GET /api/study_sessions/:id/words
Pagination with 100 items per page

#### JSON Response
```json
{
    "items": [
     {
        "spanish": "hola",
        "english": "hello",
        "correct_count": 10,
        "wrong_count": 5,
        "parts": ["Noun"]
     }
    ],
    "pagination": {
        "current_page": 1,
        "total_pages": 10,
        "total_items": 100,
        "items_per_page": 100
    }
}
```


### POST /api/reset_history
This endpoint will reset the history of a study session

#### JSON Response
```json
{
    "success": true,
    "message": "Study history has been reset"
}
```

### POST /api/full_reset
This endpoint will reset the history of all study sessions

#### JSON Response
```json
{
    "success": true,
    "message": "Study history has been reset"
}
```


### POST /api/study_sessions/:id/words:word_id/review
This endpoint will update the review status of a word in a study session

#### Request Params
- id (study_session_id) integer
- word_id (word_id) integer
- correct boolean


#### Request Payload
```json
{
    "correct": true
}
```
#### JSON Response
```json
{
    "success": true,
    "word_id": 123,
    "study_session_id": 789,
    "correct": true,
    "created_at": "2025-02-12T11:10:14-07:00"    
}
```


### Study Activity Endpoints

#### GET /api/study-activities
Returns a paginated list of study activities. Default page size is 20 items.

Query Parameters:
- `page`: Page number (default: 1)
- `page_size`: Number of items per page (default: 20)

#### JSON Response
```json
{
    "success": true,
    "data": {
        "items": [
            {
                "id": 1,
                "name": "Flashcards",
                "description": "Basic flashcard review",
                "created_at": "2025-02-13T21:20:34-07:00"
            }
        ],
        "current_page": 1,
        "total_pages": 5,
        "total_items": 100,
        "items_per_page": 20
    }
}
```

#### GET /api/study-activities/:id
Returns details for a specific study activity.

#### JSON Response
```json
{
    "success": true,
    "data": {
        "id": 1,
        "name": "Flashcards",
        "description": "Basic flashcard review",
        "created_at": "2025-02-13T21:20:34-07:00"
    }
}
```

#### POST /api/study-activities
Creates a new study activity.

#### Request Body
```json
{
    "name": "Flashcards",
    "description": "Basic flashcard review"
}
```

#### JSON Response
```json
{
    "success": true,
    "data": {
        "id": 1,
        "name": "Flashcards",
        "description": "Basic flashcard review",
        "created_at": "2025-02-13T21:20:34-07:00"
    }
}
```

#### PUT /api/study-activities/:id
Updates an existing study activity.

#### Request Body
```json
{
    "name": "Flashcards",
    "description": "Basic flashcard review"
}
```

#### JSON Response
```json
{
    "success": true,
    "data": {
        "id": 1,
        "name": "Flashcards",
        "description": "Basic flashcard review",
        "created_at": "2025-02-13T21:20:34-07:00"
    }
}
```

#### DELETE /api/study-activities/:id
Deletes a study activity.

#### JSON Response
```json
{
    "success": true,
    "data": {
        "message": "Study activity deleted successfully"
    }
}
```

#### GET /api/study_activities/:id/sessions
Returns a paginated list of study sessions for a specific activity.

Query Parameters:
- `page`: Page number (default: 1)
- `page_size`: Number of items per page (default: 20)

#### JSON Response
```json
{
    "success": true,
    "data": {
        "items": [
            {
                "id": 1,
                "group_id": 1,
                "group_name": "Typing Tutor",
                "study_activity_id": 1,
                "created_at": "2025-02-13T21:20:34-07:00"
            }
        ],
        "current_page": 1,
        "total_pages": 5,
        "total_items": 100,
        "items_per_page": 20
    }
}
```

#### POST /api/study_activities/:id/sessions
Starts a new study session for an activity.

#### Request Body
```json
{
    "group_id": 1
}
```

#### JSON Response
```json
{
    "success": true,
    "data": {
        "id": 1,
        "group_id": 1,
        "group_name": "Typing Tutor",
        "study_activity_id": 1,
        "created_at": "2025-02-13T21:20:34-07:00"
    }
}
```

#### POST /api/study_activities/:id/sessions/:session_id/results
Records a study result for a session.

#### Request Body
```json
{
    "word_id": 1,
    "correct": true
}
```

#### JSON Response
```json
{
    "success": true,
    "data": {
        "id": 1,
        "word_id": 1,
        "study_session_id": 1,
        "correct": true,
        "created_at": "2025-02-13T21:20:34-07:00"
    }
}
```

## Task Runner Tasks

The application uses [Mage](https://magefile.org/) as its task runner. Here are the available tasks:

### Database Tasks (namespace: `db`)

#### `mage db:init`
Initializes a new SQLite database using the schema defined in `db/migrations/001_initial_schema.sql`.

#### `mage db:clean`
Removes the existing database file (`words.db`).

#### `mage db:seed`
Populates the database with initial data using the seeding utility in `cmd/seed/main.go`.

#### `mage db:reset`
Performs a complete database reset by running the following tasks in sequence:
1. `db:clean` - Removes existing database
2. `db:init` - Initializes new database
3. `db:seed` - Seeds with initial data

### Server Tasks (namespace: `server`)

#### `mage server:start`
Starts the application server by running `cmd/server/main.go`.

### Default Task

When running `mage` without any target, it defaults to `db:init`.

### Seed Data Format

The seed data is stored in `db/seeds/initial_data.json` using the following format:

```json
{
    "groups": [
        {
            "name": "Basic Greetings",
            "description": "Common greetings and introductions in Spanish"
        },
        {
            "name": "Numbers 1-20",
            "description": "Basic numbers in Spanish"
        }
    ],
    "words": [
        {
            "spanish": "hola",
            "english": "hello",
            "part_of_speech": "interjection",
            "group_names": ["Basic Greetings"]
        },
        {
            "spanish": "gracias",
            "english": "thank you",
            "part_of_speech": "interjection",
            "group_names": ["Basic Greetings"]
        }
    ]
}
```

The seeder will:
1. Create the groups first
2. Create the words
3. Establish the word-group relationships based on the `group_names` array
