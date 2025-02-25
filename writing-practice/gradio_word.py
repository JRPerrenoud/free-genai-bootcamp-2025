import gradio as gr
import requests
import json
import random
import logging
import cv2
from openai import OpenAI
import os
import dotenv
import yaml
import paddleocr

dotenv.load_dotenv()

def load_prompts():
    """Load prompts from YAML file"""
    with open('prompts.yaml', 'r', encoding='utf-8') as f:
        return yaml.safe_load(f)

# Setup logging
logger = logging.getLogger('spanish_app')
logger.setLevel(logging.DEBUG)
fh = logging.FileHandler('gradio_word.log')
fh.setLevel(logging.DEBUG)
formatter = logging.Formatter('%(asctime)s - %(levelname)s - %(message)s')
fh.setFormatter(formatter)
logger.addHandler(fh)

class SpanishWritingApp:
    def __init__(self):
        self.client = OpenAI()
        self.vocabulary = None
        self.current_word = None
        self.reader = None
        self.load_vocabulary()

    def load_vocabulary(self):
        """Fetch vocabulary from API using group_id"""
        try:
            # Get group_id from environment variable or use default
            group_id = os.getenv('GROUP_ID', '1')
            url = f"http://localhost:5000/api/groups/{group_id}/words/raw"
            logger.debug(f"Fetching vocabulary from: {url}")
            
            response = requests.get(url)
            if response.status_code == 200:
                self.vocabulary = response.json()
                logger.info(f"Loaded {len(self.vocabulary.get('words', []))} words")
            else:
                logger.error(f"Failed to load vocabulary. Status code: {response.status_code}")
                self.vocabulary = {"words": []}
        except Exception as e:
            logger.error(f"Error loading vocabulary: {str(e)}")
            self.vocabulary = {"words": []}

    def get_random_word(self):
        """Get a random word from vocabulary"""
        logger.debug("Getting random word")
        
        if not self.vocabulary or not self.vocabulary.get('words'):
            logger.error("No vocabulary loaded")
            return "No vocabulary loaded", "", ""
        
        self.current_word = random.choice(self.vocabulary['words'])
        logger.debug(f"Selected word: {self.current_word}")
        
        return (
            f"{self.current_word.get('spanish', '')}",
            f"English: {self.current_word.get('english', '')}",
            f"Spanish: {self.current_word.get('spanish', '')}",          
        )

    def grade_submission(self, image):
        """Grade the user's submission"""
        try:
            if not self.reader:
                logger.info("Initializing PaddleOCR")
                self.reader = paddleocr.PaddleOCR(use_angle_cls=True, lang='es')
            
            # Log the type and content of the image
            logger.debug(f"Image type: {type(image)}")
            logger.debug(f"Image content: {image}")

            # Read the image from the file path
            image_path = image
            image = cv2.imread(image_path)

            # Verify image reading
            if image is None:
                logger.error("Failed to read image")
                return "Image read error", "Error", "Please check the image file."

            # Log image dimensions and type
            logger.debug(f"Image shape: {image.shape}")
            logger.debug(f"Image dtype: {image.dtype}")

            # Convert image to grayscale
            image_gray = cv2.cvtColor(image, cv2.COLOR_BGR2GRAY)

            # Adjust contrast using histogram equalization
            image_enhanced = cv2.equalizeHist(image_gray)

            # Transcribe the image
            logger.info("Transcribing enhanced image with PaddleOCR")
            result = self.reader.ocr(image_enhanced, cls=True)
            
            # Print OCR results for debugging
            logger.debug(f"OCR results: {result}")

            # Extract recognized text from OCR results
            if result and isinstance(result, list):
                transcription = ' '.join([line[1][0] for block in result for line in block if len(line) > 1 and isinstance(line[1], tuple)])
            else:
                logger.error("Unexpected OCR result structure")
                return "OCR result error", "Error", "Unexpected OCR result structure."

            logger.debug(f"Transcription result: {transcription}")
            
            # Check if current_word is initialized
            if not self.current_word:
                logger.error("Current word is not initialized")
                return "No word selected", "Error", "Please select a word first."
            
            # Compare transcription with target word
            is_correct = transcription.strip().lower() == self.current_word.get('spanish', '').strip().lower()
            result = " Correct!" if is_correct else " Incorrect"
            
            logger.debug(f"Current word: {self.current_word}")
            logger.debug(f"Transcription: {transcription}, Target: {self.current_word.get('spanish', '')}, Is correct: {is_correct}")
            
            # Submit grading result to backend
            word_id = self.current_word.get('id')
            review_data = {
                'words': [
                    {
                        'word_id': word_id,
                        'correct': is_correct
                    }
                ]
            }

            session_id = os.getenv('SESSION_ID', '1')  # Use a default session_id for testing

            try:
                response = requests.post(f"http://localhost:5000/api/study_sessions/{session_id}/review", json=review_data)
                if response.status_code == 200:
                    logger.info("Successfully submitted review to backend")
                else:
                    logger.error(f"Failed to submit review. Status code: {response.status_code}")
            except Exception as e:
                logger.error(f"Error submitting review: {str(e)}")
            
            # Feedback
            feedback = "Well done!" if is_correct else "Please try again."
            logger.info(f"Grading complete: {result}")
            
            return transcription, result, feedback
            
        except Exception as e:
            logger.error(f"Error in grade_submission: {str(e)}")
            return "Error processing submission", "Error processing submission", f"An error occurred: {str(e)}"

    def create_study_session(self, group_id, study_activity_id):
        """Create a new study session with the given group_id and study_activity_id"""
        try:
            session_data = {
                'group_id': group_id,
                'study_activity_id': study_activity_id
            }
            response = requests.post("http://localhost:5000/api/study_sessions", json=session_data)
            if response.status_code == 201:
                session_id = response.json().get('id')
                logger.info(f"New study session created with ID: {session_id}")
                return session_id
            else:
                logger.error(f"Failed to create study session. Status code: {response.status_code}")
                return None
        except Exception as e:
            logger.error(f"Error creating study session: {str(e)}")
            return None

def create_ui():
    app = SpanishWritingApp()
    
    # Custom CSS for larger text
    custom_css = """
    .large-text-output textarea {
        font-size: 40px !important;
        line-height: 1.5 !important;
        font-family: 'Noto Sans JP', sans-serif !important;
    }
    """
    
    with gr.Blocks(
        title="Spanish Writing Practice",
        css=custom_css
    ) as interface:
        gr.Markdown("# Spanish Writing Practice")
        
        with gr.Row():
            with gr.Column():
                generate_btn = gr.Button("Get New Word", variant="primary")
                # Make word output more prominent with larger text and more lines
                word_output = gr.Textbox(
                    label="Generated Word",
                    lines=3,
                    scale=2,  # Make the component larger
                    show_label=True,
                    container=True,
                    # Add custom CSS for larger text
                    elem_classes=["large-text-output"]
                )
                word_info = gr.Markdown("### Word Information")
                english_output = gr.Textbox(label="English", interactive=False)
                spanish_output = gr.Textbox(label="Spanish", interactive=False)                
            
            with gr.Column():
                image_input = gr.Image(label="Upload your handwritten translation", type="filepath")
                submit_btn = gr.Button("Submit", variant="secondary")
                
                with gr.Group():
                    gr.Markdown("### Feedback")
                    transcription_output = gr.Textbox(
                        label="Transcription",
                        lines=3,
                        scale=2,
                        show_label=True,
                        container=True,
                        elem_classes=["large-text-output"]
                    )
                    grade_output = gr.Textbox(label="Grade")
                    feedback_output = gr.Textbox(label="Feedback", lines=3)

        # Event handlers
        def handle_generate_click():
            logger.debug("Generate button clicked")
            result = app.get_random_word()
            logger.debug("Finished processing generate button click")
            return result
        
        generate_btn.click(
            fn=handle_generate_click,
            outputs=[word_output, english_output, spanish_output]
        )
        
        def handle_submission(image):
            return app.grade_submission(image)
            
        submit_btn.click(
            fn=handle_submission,
            inputs=[image_input],
            outputs=[transcription_output, grade_output, feedback_output]
        )

    return interface

if __name__ == "__main__":
    interface = create_ui()
    interface.launch(server_name="0.0.0.0", server_port=8501)
