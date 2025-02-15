# Endpoints
Some of these can be consolidated with query parameters and request bodies but this is where we are going to start to test out AI.



## Dashboard Endpoints
### GET /api/dashboard/last_study_session
### GET /api/dashboard/study_progress
### GET /api/dashboard/quick_stats

## Study Activity Endpoints
### GET /api/study_activities
### GET /api/study_activities/:id
### POST /api/study_activities
### POST /api/study_activities/:id/sessions (alias: study_activities/:id/launch)
### POST /api/study_activities/:id/sessions/:session_id/results

## Session Endpoints
### GET /api/study_sessions/:id
### GET /api/study_sessions
### GET /api/study_sessions/:id/words
### POST /api/study_sessions/:id/words/:word_id/review

## Word Endpoints
### GET /api/words
### GET /api/words/:id
### POST /api/words
### PUT /api/words/:id
### DELETE /api/words/:id

## Group Endpoints
### GET /api/groups/:id
### GET /api/groups/:id/study_sessions
### POST /api/groups
### PUT /api/groups/:id
### DELETE /api/groups/:id
### POST /api/groups/:id/words
### DELETE /api/groups/:id/words/:word_id

## Settings Endpoints
### POST /api/reset_history
### POST /api/full_reset


Y