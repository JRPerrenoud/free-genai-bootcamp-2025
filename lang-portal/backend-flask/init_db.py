import sqlite3
import json
import os

def init_db():
    # Remove existing database if it exists
    if os.path.exists('words.db'):
        os.remove('words.db')
        
    # Connect to SQLite database (creates it if it doesn't exist)
    conn = sqlite3.connect('words.db')
    cursor = conn.cursor()

    # Drop all tables if they exist
    cursor.execute('DROP TABLE IF EXISTS word_review_items')
    cursor.execute('DROP TABLE IF EXISTS study_sessions')
    cursor.execute('DROP TABLE IF EXISTS word_groups')
    cursor.execute('DROP TABLE IF EXISTS word_reviews')
    cursor.execute('DROP TABLE IF EXISTS study_activities')
    cursor.execute('DROP TABLE IF EXISTS words')
    cursor.execute('DROP TABLE IF EXISTS groups')

    # Create groups table
    cursor.execute('''
    CREATE TABLE groups (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL
    )
    ''')

    # Create words table with new schema
    cursor.execute('''
    CREATE TABLE words (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        english TEXT NOT NULL,
        spanish TEXT NOT NULL
    )
    ''')

    # Create word_reviews table
    cursor.execute('''
    CREATE TABLE word_reviews (
        word_id INTEGER PRIMARY KEY,
        correct_count INTEGER DEFAULT 0,
        wrong_count INTEGER DEFAULT 0,
        FOREIGN KEY (word_id) REFERENCES words (id)
    )
    ''')

    # Create word_groups table
    cursor.execute('''
    CREATE TABLE word_groups (
        word_id INTEGER,
        group_id INTEGER,
        PRIMARY KEY (word_id, group_id),
        FOREIGN KEY (word_id) REFERENCES words (id),
        FOREIGN KEY (group_id) REFERENCES groups (id)
    )
    ''')

    # Create study_activities table
    cursor.execute('''
    CREATE TABLE study_activities (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL,
        url TEXT NOT NULL,
        preview_url TEXT
    )
    ''')

    # Create study_sessions table
    cursor.execute('''
    CREATE TABLE study_sessions (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        group_id INTEGER NOT NULL,
        study_activity_id INTEGER NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (group_id) REFERENCES groups (id),
        FOREIGN KEY (study_activity_id) REFERENCES study_activities (id)
    )
    ''')

    # Create word_review_items table
    cursor.execute('''
    CREATE TABLE word_review_items (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        study_session_id INTEGER NOT NULL,
        word_id INTEGER NOT NULL,
        correct BOOLEAN NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (study_session_id) REFERENCES study_sessions (id),
        FOREIGN KEY (word_id) REFERENCES words (id)
    )
    ''')

    # Create groups first
    cursor.execute('INSERT INTO groups (name) VALUES (?)', ('Adjectives',))
    adj_group_id = cursor.lastrowid
    cursor.execute('INSERT INTO groups (name) VALUES (?)', ('Verbs',))
    verb_group_id = cursor.lastrowid

    # Load and insert adjectives
    with open('seed/data_adjectives.json', 'r') as f:
        adjectives = json.load(f)
        for adj in adjectives:
            cursor.execute('INSERT INTO words (english, spanish) VALUES (?, ?)',
                         (adj['english'], adj['spanish']))
            word_id = cursor.lastrowid
            # Immediately associate with adjectives group
            cursor.execute('INSERT INTO word_groups (word_id, group_id) VALUES (?, ?)',
                         (word_id, adj_group_id))
    
    # Load and insert verbs
    with open('seed/data_verbs.json', 'r') as f:
        verbs = json.load(f)
        for verb in verbs:
            cursor.execute('INSERT INTO words (english, spanish) VALUES (?, ?)',
                         (verb['english'], verb['spanish']))
            word_id = cursor.lastrowid
            # Immediately associate with verbs group
            cursor.execute('INSERT INTO word_groups (word_id, group_id) VALUES (?, ?)',
                         (word_id, verb_group_id))

    # Load and insert study activities
    with open('seed/study_activities.json', 'r') as f:
        activities = json.load(f)
        for activity in activities:
            cursor.execute('INSERT INTO study_activities (name, url, preview_url) VALUES (?, ?, ?)',
                         (activity['name'], activity['url'], activity.get('preview_url')))

    # Load and insert sample sessions
    with open('seed/sample_sessions.json', 'r') as f:
        sessions_data = json.load(f)
        
        for session in sessions_data['sessions']:
            # Create session with specific timestamp
            cursor.execute('''
            INSERT INTO study_sessions (group_id, study_activity_id, created_at) 
            VALUES (?, ?, ?)
            ''', (session['group_id'], session['study_activity_id'], session['created_at']))
            session_id = cursor.lastrowid
            
            # Add word reviews for this session
            for review in session['word_reviews']:
                # Add review item
                cursor.execute('''
                INSERT INTO word_review_items (study_session_id, word_id, correct, created_at)
                VALUES (?, ?, ?, ?)
                ''', (session_id, review['word_id'], review['correct'], session['created_at']))

    # Update word_reviews table based on the review items
    cursor.execute('''
    INSERT OR REPLACE INTO word_reviews (word_id, correct_count, wrong_count)
    SELECT 
        word_id,
        SUM(CASE WHEN correct THEN 1 ELSE 0 END) as correct_count,
        SUM(CASE WHEN NOT correct THEN 1 ELSE 0 END) as wrong_count
    FROM word_review_items
    GROUP BY word_id
    ''')

    conn.commit()
    conn.close()

if __name__ == '__main__':
    init_db()
