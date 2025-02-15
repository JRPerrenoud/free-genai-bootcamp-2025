# Endpoints
Some of these can be consolidated with query parameters and request bodies but this is where we are going to start to test out AI.



## Dashboard Endpoints
### GET /api/dashboard/last_study_session
### GET /api/dashboard/study_progress
### GET /api/dashboard/quick_stats

## Study Activity Endpoints
### GET /api/study_activities
### GET /api/study_activities/:id
### GET /api/study_activities/:id/study_sessions
### POST /api/study_activities

## Word Endpoints
### GET /api/words
### GET /api/words/:id
### POST /api/words
### PUT /api/words/:id
### DELETE /api/words/:id

## Group Endpoints
### GET /api/groups
### GET /api/groups/:id
### POST /api/groups
### PUT /api/groups/:id
### DELETE /api/groups/:id
### POST /api/groups/:id/words
### DELETE /api/groups/:id/words/:word_id

## Session Endpoints
### GET /api/study_sessions/:id
### GET /api/study_sessions
### GET /api/study_sessions/:id/words

## Settings Endpoints
### POST /api/reset_history
### POST /api/full_reset


Y