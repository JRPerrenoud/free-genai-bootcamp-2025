# free-genai-bootcamp-2025
repo for language learning AI project

Not all of the components of lang-portal are currently working together -- however they should work stand alone.

Below is how to run the various components of lang-portal that are currently integrated and working with lang-portal backend and frontend.

# Lang-Portal/backend-flask
    init_db.py 	
        1) Deletes words.db database if it exists
        2) Creates new database
        3) Creates all tables from scratch
        4) Populates the database with data from seeds/ directory
    
    Python3 -m venv venv
        Need to create a venv and install requirements.txt before starting the server
    
    python3  -m flask run
        Starts the backend (uses flask environmental variables) on port 5000
    
    Python3 app.py
        Also starts the backend (runs app in debug mode as configured in the script) on port 5000
    
# Lang-Portal/frontend-react
    nmp run dev
        Start development server for the frontend on port 5173 


# Typing-Tudor Game
    npx http-server ./public -p 8080
        Starts the typing tudor game on port 8080
    
# Writing-Practice Application
    python3 -m venv venv
        Need to create a venv and install requirements.txt before starting the server

    python3 gradio_app.py
        Starts the app on port 8501
        NOTE: The OCR Sucks - need to either send in text typed or all Capital letters helps
    
    
# Listening-comp
    wget https://repo.anaconda.com/miniconda/Miniconda3-latest-Linux-x86_64.sh
        Need to install Conda if it isn't installed
    
    Run the .sh script
        	Run the .sh Script
    
    Make sure to add to path
        echo 'export PATH="$HOME/miniconda3/bin:$PATH"' >> ~/.bashrc
        source ~/.bashrc
        which conda  (to make sure it's there)
        
    Create directory
        conda create -n llapp python=3.12.9  # or whatever Python version is appropriate
    
    Activate conda
        conda activate llapp
        NOTE: NEED to run this again if you close the app down)	conda activate llapp
    
    Install requirements
        pip install -r backend/requirements.txt

    Configure AWS Credentials	

    Verify AWS Credentials are good	python3
        backend/test_aws.py

