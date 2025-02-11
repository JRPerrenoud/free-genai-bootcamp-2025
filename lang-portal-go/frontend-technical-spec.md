# Frontend Technical Spec

## Pages

### Dashboard  '/dashboard'

#### Purpose
The purpse of this page is to provide a summary of learning and act as the summary page when user visits webapp

#### Components
- Last Study Session
    - shows last activity used
    - shows when last activity used
    - summarizes wrong vs correct from last activity
    - has a link to the group

- Study Progress
    - total words study eg. 3/124
        - across all study sessions show the total words studied out of all possible words in our database
        - display a master progress 

- Quick Stats
    - success rate eg. 80%
    - total study sessions eg. 3
    - total active groups eg. 4
    - study streak eg. 4 days

- Start Studying Button
    - goes to study activities page


#### Needed API Endpoints
- GET /api/dashboard/last_study_session
- GET /api/dashboard/study_progress
- GET /api/dashboard/quick-states


### Study Activities '/study_activities'

### Purpose
The purpose of this page is to show a collection of study activities with a thumbnail and its name to either launch or view the study activity.

#### Components
- Study Activity Card
    - show a thumbnail of the study activity
    - the name of the study activity
    - a launch button to take to launch page
    - a view button to view more information about past study session for this study activity

#### Needed API Endpoints
-GET /api/study_activities

### Study Activity Show '/study_activities/:id'

###Purpose
The purpose of this page is to show the details of a study activity in its past study sessions.

#### Components
- Name of study activity
- Thumbnail of study activity
- Description of study activity
- Launch button
- Study Activities Paginated List
    - id
    - activity name
    - group name
    - start time
    - end time (inferred by the last word_review_item submitted) 
    - number of review items


#### Needed API Endpoints
- GET /api/study_activities/:id
- GET /api/study_activities/:id/study_sessions


### Study Activities Launch '/study_activities/:id/launch'




-------
### Purpose
#### Components
#### Needed API Endpoints


