# Frontend Technical Spec

## Pages

### Dashboard Index '/dashboard'

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


### Study Activities Index '/study_activities'

#### Purpose
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

#### Purpose
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

#### Purpose
The purpose of this page is to show launch a study activity


#### Components
- Name of study activity
- Launch form
    - select field for group
    - launch now button

#### Behavior
After the form is sumbitted a new tab opens with the study activity based on its URL provided in the database

Also after the form is submitted the page will redirect to the study session show page


#### Needed API Endpoints
- POST /api/study_activities/
     - required params:  group_id, study_activity_id 


### Words Index '/words'

#### Purpose
The purpose of this page is to show all the words in our database

#### Components
- Paginated Word list
    -Columns
        - Fields
        - Spanish
        - English
        - Correct Count
        - Wrong Count
    - Pagination with 100 items per page
    - Clicking the Spanish word will take us to the word show page

#### Needed API Endpoints
- GET /api/words/


### Word Show Page '/words/:id'

#### Purpose
The purpose of this page is to show a single word and its details


#### Components
- Spanish
- English
- Study Statistics
    - Correct Count
    - Wrong Count
- Word Groups
    - show on a series of pills eg. tags
    - when group name is clicked it will take us to the group show page

#### Needed API Endpoints
- GET /api/words/:id


### Word Groups Index '/groups'

#### Purpose
The purpose of this page is to show all the groups in our database

#### Components
- Paginated Group list
    - Columns
        - Group name
        - Word Count
    - Clicking the group name will take us to the group show page

#### Needed API Endpoints
- GET /api/groups/


### Group Show Page '/groups/:id'

#### Purpose
The purpose of this page is to show a single group and its details


#### Components
- Group Name
- Group Statistics
    - Total Word Count
- Words in Group (Paginated list of words)
    - Should use the same component as the word index page
- Study Sessions (Paginated list of study sessions)
    - Should use the same component as the study session index page


#### Needed API Endpoints
- GET /api/groups/:id (the name and group stats)
- GET /api/groups/:id/words
- GET /api/groups/:id/study_sessions


### Study Sessions Index '/study_sessions' 

#### Purpose
The purpose of this page is to show a list of study sessions in our database

#### Components
- Paginated Study Session List
    - Columns
        - Id
        - Activity Name
        - Start Time
        - End Time
        - Number of Review Items
    - Clicking the study session id will take us to the study session show page

#### Needed API Endpoints
- GET /api/study_sessions/


### Study Session Show Page '/study_sessions/:id'

#### Purpose
The purpose of this page is to show a single study session and its details

#### Components
- Study Session Details
    - Activity Name
    - Group Name
    - Start Time
    - End Time
    - Number of Review Items
- Paginated List of Review Items
    - Should use the same component as the word index page

#### Needed API Endpoints
- GET /api/study_sessions/:id/
- GET /api/study_sessions/:id/words


### Settings Page '/settings'

#### Purpose
The purpose of this page is to make configurations to the study portal webapp

#### Components
- Theme Selection eg. Light, Dark, System Default
- Reset History Button
    - this will delete all the study sessions in the database and review items
- Full Reset Button
    - this will delete all tables and re-create with seed data


#### Needed API Endpoints
- POST /api/reset_history
- POST /api/full_reset


