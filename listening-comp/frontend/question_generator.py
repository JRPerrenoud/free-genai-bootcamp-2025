from typing import Dict, Optional
import boto3
import json
import sys
import os
import random

# Add backend to path so we can import our modules
sys.path.append(os.path.join(os.path.dirname(__file__), '..', 'backend'))
from chat import BedrockChat

class QuestionGenerator:
    def __init__(self):
        """Initialize Bedrock chat client"""
        self.chat = BedrockChat()

    def generate_question(self, practice_type: str) -> Optional[Dict]:
        """
        Generate a new question based on practice type
        
        Args:
            practice_type (str): Type of practice (Conversación, Vocabulario, etc.)
            
        Returns:
            Dict containing:
            - introduction: Brief intro to the scenario
            - conversation: The dialogue
            - question: The question to answer
            - options: List of 4 options
            - correct_answer: Index of correct answer
            - feedback: Dictionary of feedback for each option
        """
        # Create prompt based on practice type
        prompt = f"""
        You are a Spanish language teacher. Create an interactive Spanish practice question for {practice_type}.

        Follow these rules:
        1. Write a brief introduction in Spanish explaining the context
        2. Write a natural conversation in Spanish between 2-3 people
        3. Write a question in Spanish about the conversation
        4. Create 4 multiple choice options in Spanish
        5. Mark which option is correct by adding "[CORRECT]" at the start of that option
        6. Provide brief feedback in Spanish for each option explaining why it's correct/incorrect
        
        Additional rules based on practice type:
        - For "Conversación": Focus on natural dialogue and comprehension
        - For "Vocabulario": Include useful vocabulary in context
        - For "Comprensión Auditiva": Focus on listening comprehension details
        - For "Situaciones Cotidianas": Create realistic daily life scenarios
        - For "Gramática en Contexto": Naturally incorporate grammar concepts
        
        Format your response as JSON with these exact keys:
        {{
            "introduction": "brief intro in Spanish",
            "conversation": "the dialogue in Spanish",
            "question": "the question in Spanish",
            "options": ["[CORRECT]correct answer", "wrong1", "wrong2", "wrong3"],
            "feedback": {{
                "0": "why first option is correct/incorrect",
                "1": "why second option is correct/incorrect",
                "2": "why third option is correct/incorrect",
                "3": "why fourth option is correct/incorrect"
            }}
        }}
        """

        try:
            # Generate response
            completion = self.chat.generate_response(prompt)
            if not completion:
                return None
                
            try:
                question_data = json.loads(completion)
            except json.JSONDecodeError as e:
                print(f"Error parsing JSON response: {completion}")
                print(f"JSON error: {str(e)}")
                return None
            
            # Validate required fields
            required_fields = ["introduction", "conversation", "question", "options", "feedback"]
            missing_fields = [field for field in required_fields if field not in question_data]
            if missing_fields:
                print(f"Missing required fields in response: {missing_fields}")
                return None

            # Find correct answer and clean up options
            options = question_data["options"]
            correct_answer = None
            clean_options = []
            
            for i, option in enumerate(options):
                if option.startswith("[CORRECT]"):
                    correct_answer = i
                    clean_options.append(option.replace("[CORRECT]", "").strip())
                else:
                    clean_options.append(option.strip())
            
            if correct_answer is None:
                print("No correct answer marked in options")
                return None
                
            # Shuffle options and adjust feedback
            combined = list(zip(clean_options, range(len(clean_options)), question_data["feedback"].values()))
            random.shuffle(combined)
            shuffled_options, original_indices, feedback_values = zip(*combined)
            
            # Find new index of correct answer
            new_correct_index = original_indices.index(correct_answer)
            
            # Create new feedback dictionary with shuffled indices
            new_feedback = {str(i): fb for i, fb in enumerate(feedback_values)}
            
            return {
                "introduction": question_data["introduction"],
                "conversation": question_data["conversation"],
                "question": question_data["question"],
                "options": list(shuffled_options),
                "correct_answer": new_correct_index,
                "feedback": new_feedback
            }
            
        except Exception as e:
            print(f"Error generating question: {str(e)}")
            import traceback
            print(f"Traceback: {traceback.format_exc()}")
            return None

    def get_feedback(self, question_data: Dict, selected_index: int) -> str:
        """Get feedback for the selected answer"""
        return question_data["feedback"][str(selected_index)]
