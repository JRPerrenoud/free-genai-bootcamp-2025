import streamlit as st
from typing import Dict, Optional
from question_generator import QuestionGenerator
from question_storage import QuestionStorage
from question_list import render_question_list

class InteractiveLearning:
    def __init__(self):
        """Initialize the interactive learning component"""
        self.question_generator = QuestionGenerator()
        self.storage = QuestionStorage()
        self.current_question: Optional[Dict] = None
        self.selected_answer: Optional[int] = None

    def on_question_load(self, question: Dict):
        """Handle loading a saved question"""
        print("Loading question:", question["introduction"][:50])  # Debug print
        self.current_question = question
        self.selected_answer = None
        st.rerun()

    def render(self):
        """Render the interactive learning interface"""
        st.header("Interactive Learning")
        
        # Create two columns - main content and saved questions
        main_col, sidebar_col = st.columns([3, 1])
        
        with main_col:
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
            
            # Generate/Save question buttons
            col1, col2 = st.columns([2, 1])
            with col1:
                if st.button("Generate New Question"):
                    with st.spinner("Generating question..."):
                        new_question = self.question_generator.generate_question(practice_type)
                        if new_question:
                            new_question["practice_type"] = practice_type
                            self.current_question = new_question
                            self.selected_answer = None
                            print("Generated new question:", new_question["introduction"][:50])  # Debug print
            
            with col2:
                if self.current_question and st.button("Save Current Question"):
                    if self.storage.save_question(self.current_question):
                        st.success("Question saved!")
                        print("Saved question successfully")  # Debug print
                    else:
                        st.error("Failed to save question")
                        print("Failed to save question")  # Debug print
            
            # Display current question
            if self.current_question:
                self._render_question_content()
            else:
                st.info("Click 'Generate New Question' to start practicing!")
        
        # Render saved questions list in sidebar
        with sidebar_col:
            st.subheader("Saved Questions")
            # Debug info about saved questions
            questions = self.storage.get_all_questions()
            print(f"Found {len(questions)} saved questions")  # Debug print
            
            render_question_list(
                storage=self.storage,
                on_question_load=self.on_question_load
            )

    def _render_question_content(self):
        """Render the current question content"""
        col1, col2 = st.columns([2, 1])
        
        with col1:
            st.subheader("Practice Scenario")
            # Show introduction and conversation
            st.info(self.current_question["introduction"])
            st.text_area("Conversation", self.current_question["conversation"], height=200)
            
            # Show question and options
            st.write("**" + self.current_question["question"] + "**")
            
            # Create radio buttons with no default selection
            options = self.current_question["options"]
            selected = st.radio(
                "Choose your answer:",
                options,
                index=None,  # No default selection
                key=f"answer_radio_{hash(str(options))}"  # Unique key to force refresh
            )
            
            # Get selected index
            if selected:
                self.selected_answer = options.index(selected)
        
        with col2:
            st.subheader("Audio")
            # TODO: Implement text-to-speech for the conversation
            st.info("Audio feature coming soon!")
            
            st.subheader("Feedback")
            if self.selected_answer is not None:
                feedback = self.question_generator.get_feedback(
                    self.current_question,
                    self.selected_answer
                )
                
                # Show if answer is correct
                is_correct = self.selected_answer == self.current_question["correct_answer"]
                if is_correct:
                    st.success("¡Correcto! ")
                else:
                    st.error("Incorrecto")
                    correct_option = self.current_question["options"][self.current_question["correct_answer"]]
                    st.warning(f"La respuesta correcta es: {correct_option}")
                
                # Show feedback
                st.info(feedback)
