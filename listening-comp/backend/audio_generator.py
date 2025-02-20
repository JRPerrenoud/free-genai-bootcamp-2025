import boto3
import json
import os
import re
from typing import Dict, List, Tuple, Optional, Union
import tempfile
import subprocess
from datetime import datetime
import streamlit as st

class AudioGenerator:
    def __init__(self):
        """Initialize Bedrock and Polly clients"""
        self.bedrock_client = boto3.client('bedrock-runtime', region_name="us-east-1")
        self.polly_client = boto3.client('polly', region_name="us-east-1")  # Polly standard engine is supported in us-east-1
        self.model_id = "amazon.nova-micro-v1:0"
        
        # Spanish voices by gender
        self.male_voices = ["Miguel", "Enrique"]
        self.female_voices = ["Lucia", "Lupe", "Conchita"]
        
        # Create audio directory if it doesn't exist
        self.audio_dir = os.path.join(os.path.dirname(os.path.dirname(os.path.abspath(__file__))), 'data', 'audio')
        os.makedirs(self.audio_dir, exist_ok=True)
        
        # Flag to control cleanup behavior
        self.enable_cleanup = True
        
        # Spanish voices by gender
        self.male_voices = ["Miguel", "Enrique"]
        self.female_voices = ["Lucia", "Lupe", "Conchita"]
        
        # Create audio directory if it doesn't exist
        self.audio_dir = os.path.join(os.path.dirname(os.path.dirname(os.path.abspath(__file__))), 'data', 'audio')
        os.makedirs(self.audio_dir, exist_ok=True)

    def _parse_conversation(self, conversation: str) -> List[Dict]:
        """
        Parse conversation and identify speakers
        Returns list of {speaker: str, text: str, gender: str}
        """
        text_parts = []
        current_speaker = "Announcer"
        current_gender = "male"
        
        # Split by newlines and process each line
        for line in conversation.split('\n'):
            line = line.strip()
            if not line:
                continue
                
            # Check for special markers
            if line.startswith('Announcer:'):
                current_speaker = "Announcer"
                current_gender = "male"
                text = line.split(':', 1)[1].strip()
            elif line.startswith('Question:'):
                current_speaker = "Announcer"
                current_gender = "male"
                text = line.split(':', 1)[1].strip()
            else:
                # For conversation lines, alternate between male and female voices
                # but preserve the exact order
                text = line
                if text.startswith('—') or text.startswith('-'):
                    text = text[1:].strip()  # Remove the dash
                current_gender = "female" if len(text_parts) % 2 == 0 else "male"
                current_speaker = f"Speaker {len(text_parts) + 1}"
            
            if text:
                text_parts.append({
                    "speaker": current_speaker,
                    "text": text,
                    "gender": current_gender
                })
        
        return text_parts

    def _simple_parse_conversation(self, conversation: str) -> List[Dict]:
        """Fallback method for simple text-based parsing"""
        text_parts = []
        current_speaker = "Announcer"
        current_gender = "male"
        
        # Split by newlines and process each line
        for line in conversation.split('\n'):
            line = line.strip()
            if not line:
                continue
                
            # Check for speaker patterns
            if line.startswith('Question:'):
                current_speaker = "Announcer"
                current_gender = "male"
                text = line.split(':', 1)[1].strip()
            elif line.startswith('Announcer:'):
                current_speaker = "Announcer"
                current_gender = "male"
                text = line.split(':', 1)[1].strip()
            elif 'Speaker' in line:
                # Try to extract speaker number
                match = re.search(r'Speaker\s*(\d+)', line)
                if match:
                    speaker_num = match.group(1)
                    current_speaker = f"Speaker {speaker_num}"
                    # Even numbered speakers are male, odd are female
                    current_gender = "male" if int(speaker_num) % 2 == 0 else "female"
                    # Try to get text after colon
                    if ':' in line:
                        text = line.split(':', 1)[1].strip()
                    else:
                        text = line.strip()
                else:
                    # If we can't parse speaker number, continue with current speaker
                    text = line.strip()
            else:
                text = line.strip()
            
            if text:
                text_parts.append({
                    "speaker": current_speaker,
                    "text": text,
                    "gender": current_gender
                })
        
        return text_parts

    def _generate_speech(self, text: str, voice_id: str) -> Optional[str]:
        """Generate speech from text using Amazon Polly"""
        try:
            # Create a unique filename
            timestamp = datetime.now().strftime("%Y%m%d_%H%M%S_%f")
            output_path = os.path.join(self.audio_dir, f"speech_{timestamp}.mp3")
            
            print(f"\nGenerating speech for text: '{text[:50]}...' with voice: {voice_id}")
            
            response = self.polly_client.synthesize_speech(
                Engine='standard',
                LanguageCode='es-ES',
                OutputFormat='mp3',
                Text=text,
                VoiceId=voice_id
            )
            
            if "AudioStream" in response:
                with open(output_path, 'wb') as file:
                    file.write(response['AudioStream'].read())
                print(f"Created audio file: {output_path}")
                return output_path
            return None
            
        except Exception as e:
            print(f"Error generating speech: {str(e)}")
            return None

    def _generate_silence(self, duration_ms: int = 1000) -> Optional[str]:
        """Generate a silent audio file of specified duration"""
        try:
            timestamp = datetime.now().strftime("%Y%m%d_%H%M%S_%f")  # Add microseconds for uniqueness
            output_path = os.path.join(self.audio_dir, f"silence_{timestamp}.mp3")
            
            # Generate silence using ffmpeg with optimized settings
            subprocess.run([
                'ffmpeg', '-f', 'lavfi', '-i', f'anullsrc=r=22050:cl=mono',
                '-t', str(duration_ms/1000),
                '-c:a', 'libmp3lame', '-b:a', '32k',  # Lower bitrate for silence
                '-y',  # Overwrite output file
                output_path
            ], check=True, stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL)
            
            return output_path
        except Exception as e:
            print(f"Error generating silence: {str(e)}")
            return None

    def _cleanup_files(self, files: List[str]) -> None:
        """Clean up temporary audio files"""
        if not self.enable_cleanup:
            print("Cleanup disabled - keeping temporary files")
            return
            
        for file in files:
            try:
                if os.path.exists(file):
                    os.remove(file)
                    print(f"Removed: {file}")
            except Exception as e:
                print(f"Error cleaning up file {file}: {str(e)}")

    def _format_with_nova(self, question_data: Dict) -> Dict:
        """Use Bedrock Nova to format the conversation with clear speaker roles"""
        prompt = f"""Format this Spanish language learning scenario into a clear conversation with distinct speaker roles.
The output should have three parts:
1. An introduction (narrated)
2. A conversation between speakers (with clear speaker indicators)
3. A final question (narrated)

Here's the scenario:
Introduction: {question_data.get('introduction', '')}
Conversation: {question_data.get('conversation', '')}
Question: {question_data.get('question', '')}

Format the output as JSON with this structure:
{{
    "introduction": "narrated text",
    "conversation": [
        {{"speaker": "Speaker 1", "gender": "male/female", "text": "..."}},
        {{"speaker": "Speaker 2", "gender": "male/female", "text": "..."}}
    ],
    "question": "narrated text"
}}"""

        try:
            body = json.dumps({
                "messages": [
                    {
                        "role": "user",
                        "content": prompt
                    }
                ]
            })

            print("Calling Nova with prompt:", prompt)
            response = self.bedrock_client.invoke_model(
                modelId=self.model_id,
                body=body.encode()
            )

            response_body = json.loads(response.get('body').read())
            print("Nova response:", response_body)
            formatted_response = json.loads(response_body.get('completion'))
            print("Formatted response:", formatted_response)
            return formatted_response

        except Exception as e:
            print(f"Error in Nova formatting: {str(e)}")
            st.error(f"Error formatting with Nova: {str(e)}")
            return question_data

    def get_cached_audio_path(self, question_data: Dict) -> Optional[str]:
        """Check if audio already exists for this question"""
        try:
            # Create a unique identifier for this question based on its content
            content_hash = hash(f"{question_data.get('introduction', '')}{question_data.get('conversation', '')}{question_data.get('question', '')}")
            
            # Search in audio directory for matching file
            if not os.path.exists(self.audio_dir):
                return None
                
            for filename in os.listdir(self.audio_dir):
                if filename.endswith(".mp3"):
                    # Use content hash in filename
                    if str(content_hash) in filename:
                        audio_path = os.path.join(self.audio_dir, filename)
                        if os.path.exists(audio_path):
                            return audio_path
            
            return None
            
        except Exception as e:
            print(f"Error checking cached audio: {str(e)}")
            return None

    def generate_audio(self, question_data: Union[str, Dict]) -> Optional[str]:
        """Generate audio for the conversation and return the file path"""
        print("\n=== Starting Audio Generation ===\n")
        
        try:
            # If input is a string, try to parse as JSON
            if isinstance(question_data, str):
                question_data = json.loads(question_data)
            
            # Create unique output filename
            timestamp = datetime.now().strftime("%Y%m%d_%H%M%S")
            unique_id = str(hash(json.dumps(question_data, sort_keys=True)))
            output_filename = f"conversation_{timestamp}_{unique_id}.mp3"
            output_path = os.path.join(self.audio_dir, output_filename)
            
            # Check if audio already exists
            cached_path = self.get_cached_audio_path(question_data)
            if cached_path:
                print(f"Using cached audio file: {cached_path}")
                return cached_path
            
            # List to store temporary audio files
            temp_files = []
            
            # Process introduction
            print("--- Processing Introduction ---")
            intro_text = question_data.get('introduction', '')
            print(f"Introduction text: '{intro_text[:50]}...'\n")
            
            intro_audio = self._generate_speech(intro_text, "Miguel")
            if intro_audio:
                temp_files.append(intro_audio)
                print(f"Added introduction audio: {intro_audio}\n")
            
            # Add pause after introduction
            pause_audio = self._generate_silence(1000)  # 1 second
            if pause_audio:
                temp_files.append(pause_audio)
            
            # Process conversation
            print("--- Processing Conversation ---")
            conv = question_data.get('conversation', '').strip()
            lines = [line.strip() for line in conv.split('\n') if line.strip()]
            print(f"Found {len(lines)} conversation lines\n")
            
            # Generate audio for each conversation line
            for i, line in enumerate(lines):
                print(f"Processing line {i + 1}/{len(lines)}")
                if line.startswith('—') or line.startswith('-'):
                    line = line[1:].strip()
                print(f"Line text: '{line[:50]}...'")
                
                voice_id = self.male_voices[0] if i % 2 == 0 else self.female_voices[0]
                print(f"Using voice: {voice_id}\n")
                
                audio_path = self._generate_speech(line, voice_id)
                if audio_path:
                    temp_files.append(audio_path)
                    print(f"Added conversation audio {i + 1}: {audio_path}\n")
                
                # Add small pause between lines
                pause_audio = self._generate_silence(500)  # 0.5 seconds
                if pause_audio:
                    temp_files.append(pause_audio)
            
            # Add pause after conversation (1.2 seconds)
            long_pause = self._generate_silence(1200)
            if long_pause:
                temp_files.append(long_pause)
            
            # Process question
            print("--- Processing Question ---")
            question_text = question_data.get('question', '')
            print(f"Question text: '{question_text}'\n")
            
            question_audio = self._generate_speech(question_text, "Miguel")
            if question_audio:
                temp_files.append(question_audio)
                print(f"Added question audio: {question_audio}\n")
            
            # Create output file
            print("\n--- Concatenating Audio Files ---")
            print(f"Output file will be: {output_path}")
            
            # Create concatenation file
            print("\nCreating ffmpeg concat file...")
            with tempfile.NamedTemporaryFile(mode='w', suffix='.txt', delete=False) as f:
                print("Writing files to concat in this order:")
                for i, file in enumerate(temp_files, 1):
                    print(f"{i}. {file}")
                    f.write(f"file '{file}'\n")
                f.flush()
                print(f"Concat file created: {f.name}")
                
                print("\nRunning ffmpeg concat command...")
                result = subprocess.run([
                    'ffmpeg', '-f', 'concat', '-safe', '0',
                    '-i', f.name,
                    '-c', 'copy',
                    '-y',
                    output_path
                ], capture_output=True, text=True)
                
                print("\nffmpeg output:")
                print(result.stdout)
                print(result.stderr)
            
            print("\n--- Cleaning Up ---")
            print("Removing temporary audio files...")
            for file in temp_files:
                try:
                    if os.path.exists(file):
                        os.remove(file)
                        print(f"Removed: {file}")
                except Exception as e:
                    print(f"Error removing {file}: {str(e)}")
            
            if os.path.exists(f.name):
                os.remove(f.name)
                print(f"Removed concat file: {f.name}")
            
            print(f"\n=== Audio Generation Complete ===")
            print(f"Final output file: {output_path}")
            return output_path
            
        except Exception as e:
            print(f"\nError during audio generation: {str(e)}")
            print("Cleaning up temporary files...")
            for file in temp_files:
                try:
                    if os.path.exists(file):
                        os.remove(file)
                        print(f"Removed: {file}")
                except:
                    pass
            raise e
            
        except Exception as e:
            print(f"\nFatal error in generate_audio: {str(e)}")
            st.error(f"Error generating audio: {str(e)}")
            return None

    def _assign_voices(self, speakers: List[Dict]) -> Dict[str, str]:
        """Assign consistent Polly voices to speakers based on gender"""
        voice_assignments = {}
        male_count = 0
        female_count = 0
        
        for speaker in speakers:
            if speaker["speaker"] not in voice_assignments:
                if speaker["gender"].lower() == "male":
                    voice_assignments[speaker["speaker"]] = self.male_voices[male_count % len(self.male_voices)]
                    male_count += 1
                else:
                    voice_assignments[speaker["speaker"]] = self.female_voices[female_count % len(self.female_voices)]
                    female_count += 1
        
        return voice_assignments
