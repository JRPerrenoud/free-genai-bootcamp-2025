# End Points

## Here's a complete list of all endpoints in the backend server:

### Words Endpoints
GET /api/words - Get a paginated list of words (supports sorting and filtering by group)
GET /api/words/<id> - Get details of a specific word by ID

### Groups Endpoints
GET /api/groups - Get a paginated list of word groups
GET /api/groups/<id> - Get details of a specific group by ID
GET /api/groups/<id>/words - Get paginated list of words in a specific group
GET /api/groups/<id>/words/raw - Get all words in a group without pagination
POST /api/groups - Create a new group
PUT /api/groups/<id> - Update an existing group
DELETE /api/groups/<id> - Delete a group

### Study Activities Endpoints
GET /api/study_activities - Get all study activities
GET /api/study_activities/<id> - Get details of a specific study activity
GET /api/study_activities/<id>/sessions - Get paginated list of study sessions for a specific activity
GET /api/study_activities/<id>/launch - Get launch data for a study activity (includes available groups)

### Study Sessions Endpoints
GET /api/study_sessions - Get a paginated list of all study sessions
POST /api/study_sessions - Create a new study session
GET /api/study_sessions/<id> - Get details of a specific study session including reviewed words
POST /api/study_sessions/<id>/review - Submit word reviews for a study session
POST /api/study_sessions/reset - Reset all study session data

### Dashboard Endpoints
GET /api/dashboard/recent_session - Get the most recent study session with results
GET /api/dashboard/stats - Get overall study statistics (vocabulary count, mastered words, success rate, etc.)


#### NOTES 
All endpoints support CORS (Cross-Origin Resource Sharing) and return JSON responses. Most endpoints include:
- Pagination support (page number and items per page)
- Sorting options
- Error handling
- Proper database connection management

The API follows RESTful conventions and includes proper error handling, returning appropriate HTTP status codes (200 for success, 404 for not found, 500 for server errors).