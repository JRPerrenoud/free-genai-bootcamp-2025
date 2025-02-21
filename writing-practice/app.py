import streamlit as st
import requests
import json
from PIL import Image
import io
import os
from dotenv import load_dotenv
import openai

# Load environment variables
load_dotenv()

# Initialize OpenAI client
client = openai.OpenAI()

# Application states
STATES = {
    'SETUP': 'setup',
    'PRACTICE': 'practice',
    'REVIEW': 'review'
}

def fetch_word_groups():
    """Fetch word groups from the API"""
    try:
        response = requests.get('http://localhost:8000/api/groups/1/raw')
        return response.json()
    except:
        # For development, return mock data
        return {
            "words": [
                {"spanish": "libro", "english": "book"},
                {"spanish": "comer", "english": "to eat"},
                {"spanish": "beber", "english": "to drink"}
            ]
        }

def generate_sentence(word):
    """Generate a simple sentence using GPT"""
    prompt = f"""Generate a sentence using the following word: {word}
    The grammar should use A1 DELE grammar.
    You can use the following vocabulary to construct a simple sentence:
    - simple object eg. book, car, noodle
    - simple verbs eg. to drink, to eat, to meet
    - simple times eg. tomorrow, today, yesterday
    Return only the sentence, nothing else."""
    
    response = client.chat.completions.create(
        model="gpt-4",
        messages=[{"role": "user", "content": prompt}]
    )
    return response.choices[0].message.content.strip()

def grade_submission(english_sentence, image):
    """Grade the submitted image"""
    # Convert PIL Image to bytes
    img_byte_arr = io.BytesIO()
    image.save(img_byte_arr, format=image.format)
    img_byte_arr = img_byte_arr.getvalue()
    
    # First, get the transcription using GPT-4 Vision
    response = client.chat.completions.create(
        model="gpt-4-vision-preview",
        messages=[
            {
                "role": "user",
                "content": [
                    {"type": "text", "text": "Please transcribe the Spanish text in this image exactly as written:"},
                    {
                        "type": "image_url",
                        "image_url": {
                            "url": f"data:image/{image.format.lower()};base64,{img_byte_arr}",
                        }
                    }
                ]
            }
        ]
    )
    transcription = response.choices[0].message.content
    
    # Get translation and grading
    prompt = f"""Given:
    Original English: {english_sentence}
    Student's Spanish: {transcription}
    
    Please provide:
    1. Literal translation of student's Spanish
    2. Grade (A-F)
    3. Brief feedback on accuracy and suggestions for improvement
    
    Format as JSON:
    {{
        "translation": "...",
        "grade": "...",
        "feedback": "..."
    }}"""
    
    response = client.chat.completions.create(
        model="gpt-4",
        messages=[{"role": "user", "content": prompt}]
    )
    
    grading = json.loads(response.choices[0].message.content)
    return {
        "transcription": transcription,
        **grading
    }

# Initialize session state
if 'state' not in st.session_state:
    st.session_state.state = STATES['SETUP']
if 'current_sentence' not in st.session_state:
    st.session_state.current_sentence = None
if 'word_groups' not in st.session_state:
    st.session_state.word_groups = fetch_word_groups()

# App title
st.title("Language Writing Practice")

# Handle different states
if st.session_state.state == STATES['SETUP']:
    if st.button("Generate Sentence"):
        # Randomly select a word and generate a sentence
        word = st.session_state.word_groups["words"][0]["english"]  # For simplicity, using first word
        st.session_state.current_sentence = generate_sentence(word)
        st.session_state.state = STATES['PRACTICE']
        st.rerun()

elif st.session_state.state == STATES['PRACTICE']:
    st.write("### English Sentence:")
    st.write(st.session_state.current_sentence)
    
    uploaded_file = st.file_uploader("Upload your written Spanish answer (image)", type=['png', 'jpg', 'jpeg'])
    
    if uploaded_file is not None:
        image = Image.open(uploaded_file)
        st.image(image, caption="Your uploaded answer", use_column_width=True)
        
        if st.button("Submit for Review"):
            with st.spinner("Grading your submission..."):
                grading_result = grade_submission(st.session_state.current_sentence, image)
                st.session_state.grading_result = grading_result
                st.session_state.state = STATES['REVIEW']
                st.rerun()

elif st.session_state.state == STATES['REVIEW']:
    st.write("### English Sentence:")
    st.write(st.session_state.current_sentence)
    
    st.write("### Review Results:")
    grading_result = st.session_state.grading_result
    
    st.write("**Transcription of your answer:**")
    st.write(grading_result["transcription"])
    
    st.write("**Translation:**")
    st.write(grading_result["translation"])
    
    st.write("**Grade:**")
    st.write(grading_result["grade"])
    
    st.write("**Feedback:**")
    st.write(grading_result["feedback"])
    
    if st.button("Next Question"):
        st.session_state.state = STATES['SETUP']
        st.rerun()