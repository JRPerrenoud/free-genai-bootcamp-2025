from typing import List, Dict, Optional, Any
import boto3
import json
from dataclasses import dataclass
import os

# Model ID
MODEL_ID = "amazon.nova-micro-v1:0"

@dataclass
class DELEQuestion:
    introduction: str
    conversation: str
    question: str

class TranscriptStructurer:
    def __init__(self, model_id: str = MODEL_ID):
        """Initialize Bedrock client"""
        self.bedrock_client = boto3.client('bedrock-runtime', region_name="us-east-1")
        self.model_id = model_id

    def structure_transcript(self, transcript_path: str, inference_config: Optional[Dict[str, Any]] = None) -> List[DELEQuestion]:
        """
        Structure the transcript into DELE questions using Amazon Bedrock
        
        Args:
            transcript_path (str): Path to the transcript file
            inference_config (Optional[Dict[str, Any]]): Configuration for inference
            
        Returns:
            List[DELEQuestion]: List of structured DELE questions
        """
        if not os.path.exists(transcript_path):
            print(f"Error: Transcript file not found at {transcript_path}")
            return []

        try:
            with open(transcript_path, 'r', encoding='utf-8') as f:
                transcript_text = f.read()
        except Exception as e:
            print(f"Error reading transcript file: {str(e)}")
            return []

        if inference_config is None:
            inference_config = {"temperature": 0.3}  # Lower temperature for more focused responses

        prompt = f"""
        Analyze this transcript and extract DELE listening practice questions.
        For each question, identify these three parts:
        1. introduction: The introductory context
        2. conversation: The dialogue or audio content
        3. question: The specific question being asked

        Return the result as a valid JSON array where each object has these exact fields:
        [
            {{
                "introduction": "introduction text here",
                "conversation": "conversation text here",
                "question": "question text here"
            }}
        ]

        Here's the transcript:

        {transcript_text}
        """

        messages = [{
            "role": "user",
            "content": [{"text": prompt}]
        }]
        
        try:
            response = self.bedrock_client.converse(
                modelId=self.model_id,
                messages=messages,
                inferenceConfig=inference_config
            )
            
            response_text = response['output']['message']['content'][0]['text']
            
            # Try to extract JSON from the response if it's wrapped in other text
            try:
                # Find the first [ and last ] to extract just the JSON array
                start = response_text.find('[')
                end = response_text.rfind(']') + 1
                if start != -1 and end != 0:
                    response_text = response_text[start:end]
                
                # Parse the response and convert to DELEQuestion objects
                structured_data = json.loads(response_text)
                questions = []
                
                for item in structured_data:
                    question = DELEQuestion(
                        introduction=item.get('introduction', ''),
                        conversation=item.get('conversation', ''),
                        question=item.get('question', '')
                    )
                    questions.append(question)
                
                return questions
            except json.JSONDecodeError as e:
                print(f"Error parsing JSON: {str(e)}")
                print("Raw response:", response_text)
                return []
            
        except Exception as e:
            print(f"Error processing structured data: {str(e)}")
            return []

    def save_structured_data(self, questions: List[DELEQuestion], filename: str) -> bool:
        """
        Save structured questions to a JSON file
        
        Args:
            questions (List[DELEQuestion]): List of structured questions
            filename (str): Output filename
            
        Returns:
            bool: True if successful, False otherwise
        """
        try:
            data = [
                {
                    'introduction': q.introduction,
                    'conversation': q.conversation,
                    'question': q.question
                }
                for q in questions
            ]
            
            with open(filename, 'w', encoding='utf-8') as f:
                json.dump(data, f, indent=2, ensure_ascii=False)
            return True
            
        except Exception as e:
            print(f"Error saving structured data: {str(e)}")
            return False


if __name__ == "__main__":
    structurer = TranscriptStructurer()
    transcript_path = "transcripts/O2_ROLywXrM.txt"
    structured_text = structurer.structure_transcript(transcript_path)
    structurer.save_structured_data(structured_text, "questions/questions.json")
    