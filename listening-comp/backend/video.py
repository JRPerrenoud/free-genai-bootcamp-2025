import yt_dlp
import os
from typing import Optional

def download_video(video_url: str, output_dir: Optional[str] = None) -> bool:
    """
    Download a YouTube video using yt-dlp
    
    Args:
        video_url (str): URL of the YouTube video
        output_dir (Optional[str]): Directory to save the video. If None, saves in current directory
        
    Returns:
        bool: True if download successful, False otherwise
    """
    try:
        # Create output directory if it doesn't exist
        if output_dir:
            os.makedirs(output_dir, exist_ok=True)
            output_template = os.path.join(output_dir, '%(title)s.%(ext)s')
        else:
            output_template = '%(title)s.%(ext)s'

        # Define options for yt-dlp
        ydl_opts = {
            'format': 'bestvideo[ext=mp4]+bestaudio[ext=m4a]/mp4',  # Prefer MP4 format
            'outtmpl': output_template,
            'quiet': False,
            'no_warnings': False,
            'progress': True,
            'postprocessors': [{
                'key': 'FFmpegVideoConvertor',
                'preferedformat': 'mp4',  # Force MP4 output
            }],
        }

        # Download the video
        with yt_dlp.YoutubeDL(ydl_opts) as ydl:
            print(f"Downloading video from: {video_url}")
            ydl.download([video_url])
            
        print("Download completed successfully!")
        return True

    except Exception as e:
        print(f"An error occurred while downloading: {str(e)}")
        return False

if __name__ == "__main__":
    # Get the path to the videos directory
    backend_dir = os.path.dirname(os.path.abspath(__file__))
    videos_dir = os.path.join(backend_dir, "videos")
    
    # Get video URL from user
    url = input("Enter the YouTube video URL: ")
    
    # Download video
    success = download_video(url, output_dir=videos_dir)
    
    if success:
        print(f"Video saved to: {videos_dir}")
    else:
        print("Failed to download video")