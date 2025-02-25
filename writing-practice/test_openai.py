from openai import OpenAI
import os
from dotenv import load_dotenv

# Load environment variables
load_dotenv()

# Initialize client
client = OpenAI()

try:
    # Try a simple completion
    response = client.chat.completions.create(
        model="gpt-3.5-turbo",
        messages=[
            {"role": "user", "content": "Say hello"}
        ]
    )
    print("API call successful!")
    print("Response:", response.choices[0].message.content)
    
except Exception as e:
    print(f"Error occurred: {str(e)}")
