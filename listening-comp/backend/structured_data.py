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

    def structure_transcript(self, transcript_path: str, inference_config: Optional[Dict[str, Any]] = None) -> str:
        """
        Structure the transcript into DELE listening practice questions
        
        Args:
            transcript_path (str): Path to the transcript file
            inference_config (Optional[Dict[str, Any]]): Configuration for inference
            
        Returns:
            str: Formatted string with questions in XML-like structure
        """
        if not os.path.exists(transcript_path):
            print(f"Error: Transcript file not found at {transcript_path}")
            return ""

        try:
            with open(transcript_path, 'r', encoding='utf-8') as f:
                transcript_text = f.read()
        except Exception as e:
            print(f"Error reading transcript file: {str(e)}")
            return ""

        if inference_config is None:
            inference_config = {"temperature": 0.3}  # Lower temperature for more focused responses

        prompt = f"""
        Analyze this transcript and extract DELE listening practice questions.
        For each question, identify these three parts:
        1. introduction: A brief setup in Spanish that ONLY describes:
           - The general setting/location (e.g., "en una tienda", "en un restaurante")
           - The number and type of speakers (e.g., "dos amigos", "un cliente y un vendedor")
           - The general topic if relevant (e.g., "hablando sobre planes", "discutiendo sobre compras")
           Always start with "Vas a escuchar..." or "Escucharás..."
           NEVER include specific actions or details from the actual conversation.

        2. conversation: The exact dialogue or audio content in Spanish
        3. question: The specific question being asked in Spanish

        Format each question exactly like this, maintaining the XML-like tags and newlines:
        <question 1>
        Introduction:
        [Brief context-setting introduction following the rules above]

        Conversation:
        [The exact conversation text in Spanish]

        Question:
        [The question in Spanish]
        </question>

        Examples of GOOD introductions:
        - "Vas a escuchar una conversación entre dos amigos en una cafetería hablando sobre sus planes."
        - "Escucharás a un cliente y un vendedor en una tienda de ropa."
        - "Vas a escuchar a dos personas conversando en una plaza del centro."

        Examples of BAD introductions (DO NOT USE - too specific):
        - "Vas a escuchar a una mujer preguntando por un café con leche." (reveals action)
        - "Escucharás a dos amigos donde uno compra un libro." (reveals specific event)
        - "Vas a escuchar una conversación sobre un regalo de zapatos." (reveals specific item)

        Here's the transcript to analyze:

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
            
            return response['output']['message']['content'][0]['text']
            
        except Exception as e:
            print(f"Error processing structured data: {str(e)}")
            return ""

    def save_structured_data(self, formatted_text: str, filename: str) -> bool:
        """
        Save structured questions to a file
        
        Args:
            formatted_text (str): Formatted text with questions
            filename (str): Output filename
            
        Returns:
            bool: True if successful, False otherwise
        """
        try:
            with open(filename, 'w', encoding='utf-8') as f:
                f.write(formatted_text)
            return True
            
        except Exception as e:
            print(f"Error saving structured data: {str(e)}")
            return False


if __name__ == "__main__":
    structurer = TranscriptStructurer()
    transcript_path = "transcripts/O2_ROLywXrM.txt"
    structured_text = structurer.structure_transcript(transcript_path)
    structurer.save_structured_data(structured_text, "questions/questions.txt")