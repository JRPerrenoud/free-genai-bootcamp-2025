import json
from backend.audio_generator import AudioGenerator
import os

def test_audio_generation():
    # Initialize the audio generator
    generator = AudioGenerator()
    
    # Print the audio directory location
    print(f"\nAudio files will be stored in: {generator.audio_dir}")
    
    # Disable cleanup to keep temporary files
    generator.enable_cleanup = False
    
    # Load test data
    with open('data/saved_questions.json', 'r', encoding='utf-8') as f:
        questions = json.load(f)
    
    # Take the first question
    question = questions[0]
    
    # Generate audio
    audio_path = generator.generate_audio(question)
    print(f"\nFinal audio file: {audio_path}")
    
    # List all generated files
    print("\nAll generated audio files:")
    for file in sorted(os.listdir(generator.audio_dir)):
        if file.endswith('.mp3'):
            full_path = os.path.join(generator.audio_dir, file)
            print(f"- {full_path}")

if __name__ == "__main__":
    test_audio_generation()
