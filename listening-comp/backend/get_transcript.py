from youtube_transcript_api import YouTubeTranscriptApi
from typing import Optional, List, Dict
import requests
from youtube_transcript_api._errors import TranscriptsDisabled, NoTranscriptFound
import os

class YouTubeTranscriptDownloader:
    def __init__(self, languages: List[str] = ["en-US", "en"]):
        self.languages = languages

    def extract_video_id(self, url: str) -> Optional[str]:
        """
        Extract video ID from YouTube URL
        
        Args:
            url (str): YouTube URL
            
        Returns:
            Optional[str]: Video ID if found, None otherwise
        """
        if "v=" in url:
            return url.split("v=")[1][:11]
        elif "youtu.be/" in url:
            return url.split("youtu.be/")[1][:11]
        return None

    def get_transcript(self, video_id: str) -> Optional[List[Dict]]:
        """
        Download YouTube Transcript
        
        Args:
            video_id (str): YouTube video ID or URL
            
        Returns:
            Optional[List[Dict]]: Transcript if successful, None otherwise
        """
        # Extract video ID if full URL is provided
        if "youtube.com" in video_id or "youtu.be" in video_id:
            video_id = self.extract_video_id(video_id)
            
        if not video_id:
            print("Invalid video ID or URL")
            return None

        print(f"Downloading transcript for video ID: {video_id}")
        
        try:
            # First verify the video exists
            response = requests.head(f"https://www.youtube.com/watch?v={video_id}", timeout=5)
            if response.status_code != 200:
                print("Video not found or not accessible")
                return None

            # Get available transcripts
            available = YouTubeTranscriptApi.list_transcripts(video_id)
            
            # Print available languages
            manual_transcripts = list(available._manually_created_transcripts.values())
            auto_transcripts = list(available._generated_transcripts.values())
            
            print("Available transcripts:")
            for t in manual_transcripts:
                print(f"- {t.language_code} (manual)")
            for t in auto_transcripts:
                print(f"- {t.language_code} (auto-generated)")
            
            # Try to get transcript
            try:
                transcript = YouTubeTranscriptApi.get_transcript(video_id, languages=['en-US'])
                print("Using en-US transcript")
                return transcript
            except NoTranscriptFound:
                print("Falling back to general English transcript")
                return YouTubeTranscriptApi.get_transcript(video_id, languages=['en'])
                
        except requests.Timeout:
            print("Error: Request timed out while checking video availability")
            return None
        except TranscriptsDisabled:
            print("Error: Transcripts are disabled for this video")
            return None
        except NoTranscriptFound:
            print("Error: No transcript found in the specified languages:", self.languages)
            return None
        except Exception as e:
            print(f"An error occurred: {str(e)}")
            return None

    def save_transcript(self, transcript: List[Dict], filename: str) -> bool:
        """
        Save transcript to file
        
        Args:
            transcript (List[Dict]): Transcript data
            filename (str): Output filename
            
        Returns:
            bool: True if successful, False otherwise
        """
        if not transcript:
            print("Error: No transcript data to save")
            return False
        
        # Get the path to the backend directory
        backend_dir = os.path.dirname(os.path.abspath(__file__))
        # Create transcripts directory if it doesn't exist
        transcripts_dir = os.path.join(backend_dir, "transcripts")
        os.makedirs(transcripts_dir, exist_ok=True)
        
        # Create full file path
        filepath = os.path.join(transcripts_dir, f"{filename}.txt")
        print(f"Saving transcript to: {filepath}")
        
        try:
            with open(filepath, 'w', encoding='utf-8') as f:
                for entry in transcript:
                    f.write(f"{entry['text']}\n")
            print(f"Successfully saved transcript to {filepath}")
            return True
        except Exception as e:
            print(f"Error saving transcript: {str(e)}")
            return False

def main(video_url: str, print_transcript: bool = False) -> None:
    downloader = YouTubeTranscriptDownloader()
    transcript = downloader.get_transcript(video_url)
    
    if transcript:
        video_id = downloader.extract_video_id(video_url)
        if video_id and downloader.save_transcript(transcript, video_id):
            print(f"Transcript saved successfully to transcripts/{video_id}.txt")
            if print_transcript:
                for entry in transcript:
                    print(f"{entry['text']}")
        else:
            print("Failed to save transcript")
    else:
        print("Failed to get transcript")

if __name__ == "__main__":
    video_id = "https://www.youtube.com/watch?v=O2_ROLywXrM"
    main(video_id, print_transcript=True)