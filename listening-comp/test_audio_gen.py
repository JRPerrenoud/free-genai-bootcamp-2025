from backend.audio_generator import AudioGenerator
from pathlib import Path
import json

def test_audio():
    ag = AudioGenerator()
    
    # Load the first question from saved questions
    with open(Path('data/saved_questions.json'), 'r') as f:
        questions = json.load(f)
        first_question = questions[0]
    
    # Generate audio
    audio_path = ag.generate_audio(first_question)
    print(f"Generated audio at: {audio_path}")

if __name__ == "__main__":
    test_audio()
