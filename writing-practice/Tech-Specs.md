# Business Goal 
Students have asked if there could be a learning exercise to practice writing language sentences.
You have been tasked to build a prototyping application which will take a word group, and generate very simple sentences in english, and you must write them in the target lanagueg eg. Spanish.


# Technical Requirements:
Streamlit
BetterOCR (for Spanish) - Not sure if this is going to work -- need to investigate
--i.e. Need to use Managed LLM that has Vision eg. GPT4o (Amazon Bedrock probably also has a vision model)
Be able to upload an image



# Technical Specs:

## Initialization Step
When the app first initializes it needs the following:
Fetch from GET localhost:8000/api/groups/:id/raw, this will return a collection of words in json structure. It will have spanish words with theier english translation. We need to store this collection of words in memory.

## Page States
Page states describes the state the single page application should behavior from a user's perspective


## Setup
When a user first starts up the app they will only see a button called generate sentence.
When they press the button the app will generate a simple sentence using the Sentence Generator LLM and the state will move to Practice State

## Practice State
When a user is in practice state,
they will see an English sentence.
They will also see an upload field under the english sentence.
They will see a button called "submit for review"
When they press the submit for review button an uploaded image
will be passed to the Grading System and the app will transition to the Review State


## Review State
When a user is in the review state,
The user will still see the english sentence.
The upload field will be gone.
The user will now see a review of the output from the Grading System.
1 - Transcription of image
2 - Translation of transcription
3 - Grading
    - a letter score
    - a description of whether the accurate to the english sentence and suggestions.
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
    - It will transcribe the image using BetterOCR
    - It will use an LLM to produce a literal translation of the transcription
    - It will use another LLM to product a grade
    - It then returns this data to the frontend app






