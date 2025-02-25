import gradio as gr
import requests
import json
import random
import logging
from openai import OpenAI
import os
import dotenv
import yaml

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
        self.mocr = None
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
            if not self.mocr:
                logger.info("Initializing MangaOCR")
                from manga_ocr import MangaOcr
                self.mocr = MangaOcr()
            
            # Transcribe the image
            logger.info("Transcribing image with MangaOCR")
            transcription = self.mocr(image)
            logger.debug(f"Transcription result: {transcription}")
            
            # Compare transcription with target word
            is_correct = transcription.strip().lower() == self.current_word.get('spanish', '').strip().lower()
            result = " Correct!" if is_correct else " Incorrect"
            
            logger.debug(f"Current word: {self.current_word}")
            logger.debug(f"Transcription: {transcription}, Target: {self.current_word.get('spanish', '')}, Is correct: {is_correct}")
            
            # Feedback
            feedback = "Well done!" if is_correct else "Please try again."
            logger.info(f"Grading complete: {result}")
            
            return transcription, result, feedback
            
        except Exception as e:
            logger.error(f"Error in grade_submission: {str(e)}")
            return "Error processing submission", "Error processing submission", f"An error occurred: {str(e)}"

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
