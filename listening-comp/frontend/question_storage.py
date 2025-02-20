import json
import os
from datetime import datetime
from typing import Dict, List, Optional

class QuestionStorage:
    def __init__(self, storage_file: str = "saved_questions.json"):
        """Initialize question storage with path to JSON file"""
        # Get the directory where this file is located
        current_dir = os.path.dirname(os.path.abspath(__file__))
        # Go up one level to the project root
        project_root = os.path.dirname(current_dir)
        # Create a data directory if it doesn't exist
        data_dir = os.path.join(project_root, "data")
        os.makedirs(data_dir, exist_ok=True)
        # Set the full path for the storage file
        self.storage_file = os.path.join(data_dir, storage_file)
        print(f"Storage file location: {self.storage_file}")  # Debug print
        self.questions = self._load_questions()

    def _load_questions(self) -> List[Dict]:
        """Load questions from JSON file"""
        if os.path.exists(self.storage_file):
            try:
                with open(self.storage_file, 'r', encoding='utf-8') as f:
                    questions = json.load(f)
                    print(f"Loaded {len(questions)} questions")  # Debug print
                    return questions
            except json.JSONDecodeError:
                print(f"Error reading {self.storage_file}. Starting with empty question list.")
                return []
        return []

    def _save_questions(self):
        """Save questions to JSON file"""
        try:
            with open(self.storage_file, 'w', encoding='utf-8') as f:
                json.dump(self.questions, f, ensure_ascii=False, indent=2)
            print(f"Saved {len(self.questions)} questions")  # Debug print
        except Exception as e:
            print(f"Error saving questions: {e}")

    def save_question(self, question: Dict) -> bool:
        """
        Save a new question to storage
        
        Args:
            question (Dict): Question data including introduction, conversation, options, etc.
            
        Returns:
            bool: True if save was successful
        """
        try:
            # Add metadata
            question_to_save = question.copy()
            question_to_save["timestamp"] = datetime.now().isoformat()
            question_to_save["id"] = len(self.questions)  # Simple ID based on position
            
            self.questions.append(question_to_save)
            self._save_questions()
            return True
        except Exception as e:
            print(f"Error saving question: {e}")
            return False

    def get_all_questions(self) -> List[Dict]:
        """Get all saved questions"""
        return self.questions

    def get_question_by_id(self, question_id: int) -> Optional[Dict]:
        """Get a specific question by ID"""
        for question in self.questions:
            if question.get("id") == question_id:
                return question
        return None

    def delete_question(self, question_id: int) -> bool:
        """Delete a question by ID"""
        try:
            self.questions = [q for q in self.questions if q.get("id") != question_id]
            self._save_questions()
            return True
        except Exception as e:
            print(f"Error deleting question: {e}")
            return False

    def clear_all_questions(self) -> bool:
        """Clear all saved questions"""
        try:
            self.questions = []
            self._save_questions()
            return True
        except Exception as e:
            print(f"Error clearing questions: {e}")
            return False
