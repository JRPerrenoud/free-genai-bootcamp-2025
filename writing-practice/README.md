# Language Writing Practice App

A Streamlit application for practicing language writing skills with AI-powered feedback.

## Features

- Generates simple English sentences for practice
- Allows users to upload handwritten Spanish answers
- Provides AI-powered grading and feedback
- Supports multiple practice sessions

## Setup

1. Install the required dependencies:
```bash
pip install -r requirements.txt
```

2. Set up your environment variables in a `.env` file:
```
OPENAI_API_KEY=your_api_key_here
```

3. Run the Streamlit app:
```bash
streamlit run app.py
```

## How to Use

1. Click "Generate Sentence" to get a new English sentence
2. Write your Spanish translation on paper
3. Take a photo of your written answer
4. Upload the photo using the upload field
5. Click "Submit for Review" to get AI-powered feedback
6. Review your grade and feedback
7. Click "Next Question" to practice with a new sentence

## Technical Details

- Uses OpenAI's GPT-4 for sentence generation and grading
- Uses GPT-4 Vision API for handwriting recognition
- Implements a three-state workflow: Setup → Practice → Review
- Connects to a local API for word groups (configurable)
