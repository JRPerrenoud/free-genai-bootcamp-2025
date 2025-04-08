# These are the TODOs to be completed for lang-portal project

## A) Improve Language Learning app in listening-comp

### Here is what the app currently does and how it flows:
This is a Spanish language learning application that transforms YouTube transcripts into interactive learning experiences. It demonstrates a progression from basic language learning with an LLM to a more sophisticated RAG (Retrieval Augmented Generation) implementation.

### Key Features and Flow:
1. Five Development Stages:
    • Chat with Nova (Basic LLM interaction)
    • Raw Transcript (YouTube transcript download and viewing)
    • Structured Data (Text cleaning and dialogue extraction)
    • RAG Implementation (Using Bedrock embeddings and vector storage)
    • Interactive Learning (Generating practice scenarios with audio)

2. Interactive Learning Component:
    • Generates different types of practice questions:
    • Dialogue Comprehension
    • Vocabulary Practice
    • Listening Skills
    • Daily Situations
    • Grammar in Context
    • Allows saving questions for later practice
    • Provides audio synthesis for listening practice

3. Technical Components:
    • Uses Amazon Bedrock for text generation (Nova) and embeddings (Titan)
    • Processes YouTube transcripts as knowledge sources
    • Implements RAG for context-aware responses
    • Includes audio generation capabilities


### How to Launch the App:
To run the app, you need to start both the backend and frontend components:

1. First, start the backend:
sh
CopyInsert
cd /mnt/c/GitHub/free-genai-bootcamp-2025/listening-comp
pip install -r backend/requirements.txt
python backend/get_transcript.py  # This appears to be the main backend component

2. Then, start the frontend:
sh
CopyInsert
cd /mnt/c/GitHub/free-genai-bootcamp-2025/listening-comp
streamlit run frontend/main.py

The app should then be accessible in your web browser at the URL provided by Streamlit (typically http://localhost:8501).
Note: The app requires AWS credentials to be properly configured for accessing Amazon Bedrock services. Make sure your AWS credentials are set up correctly before running the app.

### TODOs:
1. Change the port to unused port in lang-portal
2. Add instructions for how to run the app
3. Update README.md for listening-comp section

## B) Fix typing tudor to actually have to learn some translation - either by typing a spanish work from English prompt or other way around



