import os
import sys
import json
sys.path.append(os.path.dirname(os.path.dirname(os.path.abspath(__file__))))
from backend.audio_generator import AudioGenerator

def test_audio():
    """Test audio generation with a sample conversation"""
    print("\nTesting audio generation...")
    
    # Load the first question from saved questions
    questions_file = os.path.join(os.path.dirname(os.path.dirname(os.path.abspath(__file__))), 'data', 'saved_questions.json')
    with open(questions_file, "r", encoding='utf-8') as f:
        questions = json.load(f)
        question = questions[0]
    
    audio_generator = AudioGenerator()
    audio_file = audio_generator.generate_audio(question)
    
    if audio_file and os.path.exists(audio_file):
        print(f"Successfully generated audio file: {audio_file}")
    else:
        print("Failed to generate audio file")

if __name__ == "__main__":
    test_audio()
