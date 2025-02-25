# Technical Specs:

## Initialization Step
When the app first initializes it needs the following:
Fetch from GET localhost:5000/api/groups/:id/raw, this will return a collection of words in json structure. It will have spanish words with theier english translation. We need to store this collection of words in memory.

## Page States
Page states describes the state the single page application should behavior from a user's perspective


## Setup
When a user first starts up the app they will only see a button called generate sentence, it needs to be generated from one of the words in the collection.
When they press the button the app will generate a simple sentence using the Sentence Generator LLM and the state will move to Practice State

## Practice State
When a user is in practice state, they will see the English sentence.
They will also see an upload field under the english sentence.
They will see a button called "submit for review"
When they press the submit for review button an uploaded image will be passed to the Grading System and the app will transition to the Review State


## Review State
When a user is in the review state,
The user will still see the english sentence.
The upload field will be gone.
The user will now see a review of the output from the Grading System.
1 - Transcription of image
2 - Translation of transcription
3 - Grading
    - a letter score
    - a description of whether the attempt was accurate to the english sentence and suggestions.
There will be a button called "Next Question" when clicked
it will generate a new question and place the app into Practice State.


## Sentence Gnerator LLM Prompt
Generate a sentence using the following word: {{word}}
The grammer should use A1 DELE grammer.
You can use the following vocabulary to construct a simple sentence:
- simple object eg. book, car, noodle
- simple verbs eg. to drink, to eat, to meet
- simple times eg. tomorrow, today, yesterday

 
## Grading System
The Grading System will do the following:
    - It will transcribe the image using EasyOCR
    - It will use an LLM to produce a literal translation of the transcription
    - It will use another LLM to product a grade
    - It then returns this data to the frontend app








