# End Points

## Here's a complete list of all endpoints in the backend server:

### Words Routes (/routes/words.py):
GET /api/words - Get paginated list of words with optional group filtering and sorting
GET /api/words/<id> - Get a single word with its details and groups

### Groups Routes (/routes/groups.py):
GET /api/groups - Get paginated list of groups with sorting
GET /api/groups/<id> - Get a single group's details
GET /api/groups/<id>/words - Get paginated list of words in a group
GET /api/groups/<id>/study_sessions - Get study sessions for a group

### Study Activities Routes (/routes/study_activities.py):
GET /api/study_activities - Get list of all study activities
GET /api/study_activities/<id> - Get details of a single study activity
GET /api/study_activities/<id>/sessions - Get study sessions for an activity
GET /api/study_activities/<id>/launch - Get launch data for an activity


### Study Sessions Routes (/routes/study_sessions.py):
GET /api/study_sessions - Get paginated list of all study sessions
GET /api/study_sessions/<id> - Get details of a single study session with reviewed words
POST /api/study_sessions/reset - Reset all study session data


There are also some TODO endpoints noted in the code:

POST /study_sessions/:id/review (planned)
GET /groups/:id/words/raw (planned)



#### NOTES 
All endpoints support CORS (Cross-Origin Resource Sharing) and return JSON responses. Most endpoints include:
- Pagination support (page number and items per page)
- Sorting options
- Error handling
- Proper database connection management

The API follows RESTful conventions and includes proper error handling, returning appropriate HTTP status codes (200 for success, 404 for not found, 500 for server errors).