import streamlit as st
from typing import Dict
import json
from collections import Counter
import re

import sys
import os
sys.path.append(os.path.dirname(os.path.dirname(os.path.abspath(__file__))))

from backend.get_transcript import YouTubeTranscriptDownloader

from backend.chat import BedrockChat


# Page config
st.set_page_config(
    page_title="Spanish Learning Assistant",
    page_icon="🇪🇸",
    layout="wide"
)

# Initialize session state
if 'transcript' not in st.session_state:
    st.session_state.transcript = None
if 'messages' not in st.session_state:
    st.session_state.messages = []

def render_header():
    """Render the header section"""    
    #st.markdown("<h1>🇪🇸 Spanish Learning Assistant</h1>", unsafe_allow_html=True)
    #st.title("🇪🇸 Spanish Learning Assistant")
    st.title("Spanish Learning Assistant")
    st.markdown("""
    Transform YouTube transcripts into interactive Spanish learning experiences.
    
    This tool demonstrates:
    - Base LLM Capabilities
    - RAG (Retrieval Augmented Generation)
    - Amazon Bedrock Integration
    - Agent-based Learning Systems
    """)

def render_sidebar():
    """Render the sidebar with component selection"""
    with st.sidebar:
        st.header("Development Stages")
        
        # Main component selection
        selected_stage = st.radio(
            "Select Stage:",
            [
                "1. Chat with Nova",
                "2. Raw Transcript",
                "3. Structured Data",
                "4. RAG Implementation",
                "5. Interactive Learning"
            ]
        )
        
        # Stage descriptions
        stage_info = {
            "1. Chat with Nova": """
            **Current Focus:**
            - Basic Spanish learning
            - Understanding LLM capabilities
            - Identifying limitations
            """,
            
            "2. Raw Transcript": """
            **Current Focus:**
            - YouTube transcript download
            - Raw text visualization
            - Initial data examination
            """,
            
            "3. Structured Data": """
            **Current Focus:**
            - Text cleaning
            - Dialogue extraction
            - Data structuring
            """,
            
            "4. RAG Implementation": """
            **Current Focus:**
            - Bedrock embeddings
            - Vector storage
            - Context retrieval
            """,
            
            "5. Interactive Learning": """
            **Current Focus:**
            - Scenario generation
            - Audio synthesis
            - Interactive practice
            """
        }
        
        st.markdown("---")
        st.markdown(stage_info[selected_stage])
        
        return selected_stage

def render_chat_stage():
    """Render an improved chat interface"""
    st.header("Chat with Nova")

    # Initialize BedrockChat instance if not in session state
    if 'bedrock_chat' not in st.session_state:
        st.session_state.bedrock_chat = BedrockChat()

    # Introduction text
    st.markdown("""
    Start by exploring Nova's base Spanish language capabilities. Try asking questions about Spanish grammar, 
    vocabulary, or cultural aspects.
    """)

    # Initialize chat history if not exists
    if "messages" not in st.session_state:
        st.session_state.messages = []

    # Display chat messages
    for message in st.session_state.messages:
        with st.chat_message(message["role"], avatar="🧑‍💻" if message["role"] == "user" else "🤖"):
            st.markdown(message["content"])

    # Chat input area
    if prompt := st.chat_input("Ask about Spanish language..."):
        # Process the user input
        process_message(prompt)

    # Example questions in sidebar
    with st.sidebar:
        st.markdown("### Try These Examples")
        example_questions = [
            "How do I say 'Where is the train station?' in Spanish?",
             "What's the polite form of eat?",
            "How do I count objects in Spanish?",
            "What's the difference between ¡Buenos días! and ¡Buenas tardes! and ¡Buenas noches!?",
            "How do I ask for directions politely?"
        ]
        
        for q in example_questions:
            if st.button(q, use_container_width=True, type="secondary"):
                # Process the example question
                process_message(q)
                st.rerun()

    # Add a clear chat button
    if st.session_state.messages:
        if st.button("Clear Chat", type="primary"):
            st.session_state.messages = []
            st.rerun()

def process_message(message: str):
    """Process a message and generate a response"""
    # Add user message to state and display
    st.session_state.messages.append({"role": "user", "content": message})
    with st.chat_message("user", avatar="🧑‍💻"):
        st.markdown(message)

    # Generate and display assistant's response
    with st.chat_message("assistant", avatar="🤖"):
        response = st.session_state.bedrock_chat.generate_response(message)
        if response:
            st.markdown(response)
            st.session_state.messages.append({"role": "assistant", "content": response})



def count_characters(text):
    """Count Spanish and total characters in text"""
    if not text:
        return 0, 0
        
    def is_spanish(char):
        return any([
            '\u0041' <= char <= '\u007A',  # Basic Latin
            '\u00C0' <= char <= '\u00FF',  # Latin-1 Supplement            
            '\u0100' <= char <= '\u017F',  # Latin Extended-A
        ])
   
    sp_chars = sum(1 for char in text if is_spanish(char))
    return sp_chars, len(text)

def render_transcript_stage():
    """Render the raw transcript stage"""
    st.header("Raw Transcript Processing")
    
    # URL input
    url = st.text_input(
        "YouTube URL",
        placeholder="Enter a Spanish lesson YouTube URL"
    )
    
    # Download button and processing
    if url:
        if st.button("Download Transcript"):
            progress_text = "Downloading transcript..."
            progress_bar = st.progress(0, text=progress_text)
            
            try:
                # Update progress
                progress_bar.progress(25, text="Initializing downloader...")
                downloader = YouTubeTranscriptDownloader()
                
                progress_bar.progress(50, text="Fetching transcript...")
                transcript = downloader.get_transcript(url)
                
                if transcript:
                    progress_bar.progress(75, text="Processing transcript...")
                    # Store the raw transcript text in session state
                    transcript_text = "\n".join([entry['text'] for entry in transcript])
                    st.session_state.transcript = transcript_text
                    
                    # Save the transcript file
                    video_id = downloader.extract_video_id(url)
                    if video_id and downloader.save_transcript(transcript, video_id):
                        progress_bar.progress(100, text="Complete!")
                        st.success(f"Transcript downloaded and saved to transcripts/{video_id}.txt")
                    else:
                        progress_bar.progress(100, text="Partial completion")
                        st.warning("Transcript downloaded but could not be saved to file")
                else:
                    progress_bar.empty()
                    st.error("No transcript found for this video.")
            except Exception as e:
                progress_bar.empty()
                st.error(f"Error downloading transcript: {str(e)}")
                st.error("If the process seems stuck, try refreshing the page and using a different video.")
    
    col1, col2 = st.columns(2)
    
    with col1:
        st.subheader("Raw Transcript")
        if st.session_state.transcript:
            st.text_area(
                label="Raw text",
                value=st.session_state.transcript,
                height=400,
                disabled=True
            )
    
        else:
            st.info("No transcript loaded yet")
    
    with col2:
        st.subheader("Transcript Stats")
        if st.session_state.transcript:
            # Calculate stats
            sp_chars, total_chars = count_characters(st.session_state.transcript)
            total_lines = len(st.session_state.transcript.split('\n'))
            
            # Display stats
            st.metric("Total Characters", total_chars)
            st.metric("Spanish Characters", sp_chars)
            st.metric("Total Lines", total_lines)
        else:
            st.info("Load a transcript to see statistics")

def render_structured_stage():
    """Render the structured data stage"""
    st.header("Structured Data Processing")
    
    col1, col2 = st.columns(2)
    
    with col1:
        st.subheader("Dialogue Extraction")
        # Placeholder for dialogue processing
        st.info("Dialogue extraction will be implemented here")
        
    with col2:
        st.subheader("Data Structure")
        # Placeholder for structured data view
        st.info("Structured data view will be implemented here")

def render_rag_stage():
    """Render the RAG implementation stage"""
    st.header("RAG System")
    
    # Query input
    query = st.text_input(
        "Test Query",
        placeholder="Enter a question about Spanish..."
    )
    
    col1, col2 = st.columns(2)
    
    with col1:
        st.subheader("Retrieved Context")
        # Placeholder for retrieved contexts
        st.info("Retrieved contexts will appear here")
        
    with col2:
        st.subheader("Generated Response")
        # Placeholder for LLM response
        st.info("Generated response will appear here")

def render_interactive_stage():
    """Render the interactive learning stage"""
    st.header("Interactive Learning")
    
    # Initialize question generator in session state if not exists
    if 'question_generator' not in st.session_state:
        from question_generator import QuestionGenerator
        st.session_state.question_generator = QuestionGenerator()
    
    if 'current_question' not in st.session_state:
        st.session_state.current_question = None
    
    # Practice type selection
    practice_type = st.selectbox(
        "Select Practice Type",
        [
            "Conversación (Dialogue Comprehension)", 
            "Vocabulario (Vocabulary Practice)", 
            "Comprensión Auditiva (Listening Skills)",
            "Situaciones Cotidianas (Daily Situations)",
            "Gramática en Contexto (Grammar in Context)"
        ]
    )
    
    # Generate new question button
    if st.button("Generate New Question"):
        with st.spinner("Generating question..."):
            st.session_state.current_question = st.session_state.question_generator.generate_question(practice_type)
            st.session_state.selected_answer = None
    
    if st.session_state.current_question:
        col1, col2 = st.columns([2, 1])
        
        with col1:
            st.subheader("Practice Scenario")
            # Show introduction and conversation
            st.info(st.session_state.current_question["introduction"])
            st.text_area("Conversation", st.session_state.current_question["conversation"], height=200)
            
            # Show question and options
            st.write("**" + st.session_state.current_question["question"] + "**")
            
            # Initialize answer selection in session state if not exists
            if 'selected_answer' not in st.session_state:
                st.session_state.selected_answer = None
            
            # Create radio buttons with no default selection
            options = st.session_state.current_question["options"]
            selected = st.radio(
                "Choose your answer:",
                options,
                index=None,  # No default selection
                key=f"answer_radio_{hash(str(options))}"  # Unique key to force refresh
            )
            
            # Get selected index
            if selected:
                selected_index = options.index(selected)
                st.session_state.selected_answer = selected_index
        
        with col2:
            st.subheader("Audio")
            # TODO: Implement text-to-speech for the conversation
            st.info("Audio feature coming soon!")
            
            st.subheader("Feedback")
            if st.session_state.selected_answer is not None:
                feedback = st.session_state.question_generator.get_feedback(
                    st.session_state.current_question,
                    st.session_state.selected_answer
                )
                
                # Show if answer is correct
                is_correct = st.session_state.selected_answer == st.session_state.current_question["correct_answer"]
                if is_correct:
                    st.success("¡Correcto! 🎉")
                else:
                    st.error("Incorrecto")
                    correct_option = st.session_state.current_question["options"][st.session_state.current_question["correct_answer"]]
                    st.warning(f"La respuesta correcta es: {correct_option}")
                
                # Show feedback
                st.info(feedback)
    else:
        st.info("Click 'Generate New Question' to start practicing!")

def main():
    render_header()
    selected_stage = render_sidebar()
    
    # Render appropriate stage
    if selected_stage == "1. Chat with Nova":
        render_chat_stage()
    elif selected_stage == "2. Raw Transcript":
        render_transcript_stage()
    elif selected_stage == "3. Structured Data":
        render_structured_stage()
    elif selected_stage == "4. RAG Implementation":
        render_rag_stage()
    elif selected_stage == "5. Interactive Learning":
        render_interactive_stage()
    
    # Debug section at the bottom
    with st.expander("Debug Information"):
        st.json({
            "selected_stage": selected_stage,
            "transcript_loaded": st.session_state.transcript is not None,
            "chat_messages": len(st.session_state.messages)
        })

if __name__ == "__main__":
    main()