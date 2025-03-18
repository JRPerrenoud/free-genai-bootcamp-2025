import gradio as gr
import requests
import json
import random
import logging
from openai import OpenAI
import os
import dotenv
import yaml
import shutil
import pytesseract
import cv2
from PIL import Image
import numpy as np

dotenv.load_dotenv()

def load_prompts():
    """Load prompts from YAML file"""
    with open('prompts.yaml', 'r', encoding='utf-8') as f:
        return yaml.safe_load(f)

# Setup logging
logger = logging.getLogger('spanish_app')
logger.setLevel(logging.DEBUG)
fh = logging.FileHandler('gradio_app.log')
fh.setLevel(logging.DEBUG)
formatter = logging.Formatter('%(asctime)s - %(levelname)s - %(message)s')
fh.setFormatter(formatter)
logger.addHandler(fh)

class SpanishWritingApp:
    def __init__(self):
        self.client = OpenAI()
        self.vocabulary = None
        self.current_word = None
        self.current_sentence = None
        self.current_translation = None
        self.load_vocabulary()
        
        # Configure Tesseract for Spanish language
        # Note: This requires the Spanish language data files to be installed
        # on the system where the app is running
        self.tesseract_config = r'--oem 1 --psm 6 -l spa'

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

    def generate_sentence(self, word):
        """Generate a sentence using OpenAI API"""
        logger.debug(f"Generating sentence for word: {word.get('english', '')}")
        
        try:
            prompts = load_prompts()
            messages = [
                {"role": "system", "content": prompts['sentence_generation']['system']},
                {"role": "user", "content": prompts['sentence_generation']['user'].format(word=word.get('english', ''))}
            ]
            logger.debug(f"Messages for API call: {messages}")
            
            response = self.client.chat.completions.create(
                model="gpt-3.5-turbo",
                messages=messages,
                temperature=0.7,
                max_tokens=100
            )
            sentence = response.choices[0].message.content.strip()
            logger.info(f"Generated sentence: {sentence}")
            return sentence
        except Exception as e:
            logger.error(f"Error generating sentence: {str(e)}")
            return "Error generating sentence. Please try again."
            
    def translate_sentence(self, sentence):
        """Translate a Spanish sentence to English"""
        logger.debug(f"Translating sentence: {sentence}")
        
        try:
            prompts = load_prompts()
            messages = [
                {"role": "system", "content": prompts['translation']['system']},
                {"role": "user", "content": prompts['translation']['user'].format(text=sentence)}
            ]
            
            response = self.client.chat.completions.create(
                model="gpt-3.5-turbo",
                messages=messages,
                temperature=0.3,
                max_tokens=100
            )
            translation = response.choices[0].message.content.strip()
            logger.info(f"Generated translation: {translation}")
            return translation
        except Exception as e:
            logger.error(f"Error translating sentence: {str(e)}")
            return "Error translating sentence. Please try again."

    def get_random_word_and_sentence(self):
        """Get a random word and generate a sentence"""
        logger.debug("Getting random word and generating sentence")
        
        if not self.vocabulary or not self.vocabulary.get('words'):
            logger.error("No vocabulary loaded")
            return "No vocabulary loaded", "No translation available", "", ""
            
        self.current_word = random.choice(self.vocabulary['words'])
        logger.debug(f"Selected word: {self.current_word}")
        self.current_sentence = self.generate_sentence(self.current_word)
        self.current_translation = self.translate_sentence(self.current_sentence)
        
        return (
            self.current_sentence,
            self.current_translation,
            self.current_word.get('english', ''),
            self.current_word.get('spanish', ''),          
        )

    def preprocess_image(self, image_path):
        """Preprocess the image to improve OCR accuracy"""
        try:
            # Read the image
            img = cv2.imread(image_path)
            
            # Convert to grayscale
            gray = cv2.cvtColor(img, cv2.COLOR_BGR2GRAY)
            
            # Apply thresholding to get a binary image
            _, binary = cv2.threshold(gray, 150, 255, cv2.THRESH_BINARY_INV)
            
            # Apply dilation to make text more visible
            kernel = np.ones((1, 1), np.uint8)
            dilated = cv2.dilate(binary, kernel, iterations=1)
            
            # Invert back
            preprocessed = cv2.bitwise_not(dilated)
            
            # Save the preprocessed image temporarily
            temp_path = "temp_preprocessed.png"
            cv2.imwrite(temp_path, preprocessed)
            
            return temp_path
        except Exception as e:
            logger.error(f"Error preprocessing image: {str(e)}")
            return image_path

    def grade_submission(self, image):
        """Process image submission and grade it using Tesseract OCR and LLM"""
        try:
            # Preprocess the image to improve OCR accuracy
            logger.info("Preprocessing image for OCR")
            preprocessed_image_path = self.preprocess_image(image)
            
            # Use Tesseract to extract text from the image
            logger.info("Transcribing image with Tesseract OCR")
            transcription = pytesseract.image_to_string(
                Image.open(preprocessed_image_path), 
                config=self.tesseract_config
            )
            transcription = transcription.strip()
            logger.debug(f"Transcription result: {transcription}")
            
            # Clean up temporary file if it was created
            if preprocessed_image_path != image and os.path.exists(preprocessed_image_path):
                os.remove(preprocessed_image_path)
            
            # Load prompts
            prompts = load_prompts()
            
            # Get literal translation
            logger.info("Getting literal translation")
            translation_response = self.client.chat.completions.create(
                model="gpt-3.5-turbo",
                messages=[
                    {"role": "system", "content": prompts['translation']['system']},
                    {"role": "user", "content": prompts['translation']['user'].format(text=transcription)}
                ],
                temperature=0.3
            )
            translation = translation_response.choices[0].message.content.strip()
            logger.debug(f"Translation: {translation}")
            
            # Get grading and feedback
            logger.info("Getting grade and feedback")
            grading_response = self.client.chat.completions.create(
                model="gpt-3.5-turbo",
                messages=[
                    {"role": "system", "content": prompts['grading']['system']},
                    {"role": "user", "content": prompts['grading']['user'].format(
                        target_sentence=self.current_sentence,
                        submission=transcription,
                        translation=translation
                    )}
                ],
                temperature=0.3
            )
            
            feedback = grading_response.choices[0].message.content.strip()
            # Parse grade and feedback from response
            grade = 'C'  # Default grade
            if 'Grade: A' in feedback:
                grade = 'A'
            elif 'Grade: B' in feedback:
                grade = 'B'
            elif 'Grade: C' in feedback:
                grade = 'C'
            
            # Extract just the feedback part
            feedback = feedback.split('Feedback:')[-1].strip()
            
            logger.info(f"Grading complete: {grade}")
            logger.debug(f"Feedback: {feedback}")
            
            return transcription, translation, grade, feedback
            
        except Exception as e:
            logger.error(f"Error in grade_submission: {str(e)}")
            return "Error processing submission", "Error processing submission", "C", f"An error occurred: {str(e)}"

def create_ui():
    app = SpanishWritingApp()
    
    # Custom CSS for larger text and font fixes
    custom_css = """
    /* Fix font issues by using system fonts */
    * {
        font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif, "Apple Color Emoji", "Segoe UI Emoji", "Segoe UI Symbol" !important;
    }
    
    /* Reduce Spanish sentence font size to 2/3 of original */
    .spanish-text-output textarea {
        font-size: 26px !important;
        line-height: 1.5 !important;
    }
    
    /* Style for English translation to match */
    .english-text-output textarea {
        font-size: 26px !important;
        line-height: 1.5 !important;
    }
    
    /* Transcription output styling */
    .transcription-output textarea {
        font-size: 26px !important;
        line-height: 1.5 !important;
    }
    
    /* Hide manifest error in console by disabling the request */
    @media (display-mode: browser) {
        html {
            --pwa-manifest: none;
        }
    }
    
    /* Style for hint button */
    .hint-button {
        margin-left: 10px;
    }
    """
    
    with gr.Blocks(
        title="Spanish Writing Practice",
        css=custom_css,
        head="""
        <link rel="manifest" href="data:application/json;base64,ewogICJuYW1lIjogIlNwYW5pc2ggV3JpdGluZyBQcmFjdGljZSIsCiAgInNob3J0X25hbWUiOiAiU3BhbmlzaEFwcCIsCiAgInN0YXJ0X3VybCI6ICIvIiwKICAiZGlzcGxheSI6ICJzdGFuZGFsb25lIiwKICAiYmFja2dyb3VuZF9jb2xvciI6ICIjZmZmZmZmIiwKICAidGhlbWVfY29sb3IiOiAiIzRhOTBlMiIsCiAgImljb25zIjogW10KfQ==">
        <meta name="viewport" content="width=device-width, initial-scale=1.0">
        <style>
        @font-face {
            font-family: 'ui-sans-serif';
            src: local('Arial'), local('Helvetica'), local('San Francisco'), local('Segoe UI');
            font-weight: normal;
            font-style: normal;
        }
        </style>
        """
    ) as interface:
        gr.Markdown("# Spanish Writing Practice")
        
        # Store the current translation and word
        current_translation = gr.State("")
        current_english_word = gr.State("")
        
        with gr.Row():
            with gr.Column():
                generate_btn = gr.Button("Generate New Sentence", variant="primary")
                
                # Spanish sentence output
                gr.Markdown("### Sentence to Practice")
                spanish_sentence_output = gr.Textbox(
                    label="Spanish",
                    lines=3,
                    scale=2,
                    show_label=True,
                    container=True,
                    elem_classes=["spanish-text-output"]
                )
                
                # English translation section with hint button
                with gr.Row():
                    gr.Markdown("### English")
                    show_english_btn = gr.Button("Show", size="sm", elem_classes=["hint-button"])
                
                # Initially hidden English translation
                english_translation_output = gr.Textbox(
                    label="",
                    lines=3,
                    scale=2,
                    show_label=False,
                    container=True,
                    elem_classes=["english-text-output"],
                    visible=False
                )
                
                # Word information section
                gr.Markdown("### Word Information")
                
                # Spanish word output
                spanish_word_output = gr.Textbox(label="Spanish", interactive=False)
                
                # English word section with hint button
                with gr.Row():
                    gr.Markdown("#### English")
                    show_english_word_btn = gr.Button("Show", size="sm", elem_classes=["hint-button"])
                
                # Initially hidden English word
                english_word_output = gr.Textbox(
                    label="",
                    interactive=False,
                    visible=False
                )
            
            with gr.Column():
                image_input = gr.Image(label="Upload your handwritten sentence", type="filepath")
                submit_btn = gr.Button("Submit", variant="secondary")
                
                with gr.Group():
                    gr.Markdown("### Feedback")
                    transcription_output = gr.Textbox(
                        label="Transcription",
                        lines=3,
                        scale=2,
                        show_label=True,
                        container=True,
                        elem_classes=["transcription-output"]
                    )
                    translation_output = gr.Textbox(label="Translation", lines=2)
                    grade_output = gr.Textbox(label="Grade")
                    feedback_output = gr.Textbox(label="Feedback", lines=3)

        # Event handlers
        def handle_generate_click():
            logger.debug("Generate button clicked")
            spanish, english_trans, english_word, spanish_word = app.get_random_word_and_sentence()
            logger.debug("Finished processing generate button click")
            # Hide the English text when generating a new sentence
            return [spanish, spanish_word, gr.update(visible=False), gr.update(visible=False), english_trans, english_word]

        generate_btn.click(
            fn=handle_generate_click,
            outputs=[
                spanish_sentence_output, 
                spanish_word_output,
                english_translation_output,
                english_word_output,
                current_translation,
                current_english_word
            ]
        )
        
        def handle_submission(image):
            return app.grade_submission(image)
            
        submit_btn.click(
            fn=handle_submission,
            inputs=[image_input],
            outputs=[transcription_output, translation_output, grade_output, feedback_output]
        )
        
        # Show/hide English translation
        def show_english_translation(translation):
            return gr.update(value=translation, visible=True)
        
        show_english_btn.click(
            fn=show_english_translation,
            inputs=[current_translation],
            outputs=[english_translation_output]
        )
        
        # Show/hide English word
        def show_english_word(word):
            return gr.update(value=word, visible=True)
        
        show_english_word_btn.click(
            fn=show_english_word,
            inputs=[current_english_word],
            outputs=[english_word_output]
        )

    return interface

if __name__ == "__main__":
    interface = create_ui()
    # Launch the app
    interface.launch(server_name="0.0.0.0", server_port=8501)
