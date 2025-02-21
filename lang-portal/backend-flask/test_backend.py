import os
from app import create_app
import requests

def test_backend():
    # Create test config to use a test database
    test_config = {
        'DATABASE': 'test_words.db',
        'TESTING': True
    }
    
    # Create and configure the test app
    app = create_app(test_config)
    
    # Initialize the database
    with app.app_context():
        app.db.init(app)
    
    # Start the Flask server
    app.run(port=5001, debug=False)

if __name__ == '__main__':
    # Remove test database if it exists
    if os.path.exists('test_words.db'):
        os.remove('test_words.db')
    
    test_backend()
