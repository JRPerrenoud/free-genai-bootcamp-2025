import sqlite3
import json

def init_db():
    # Connect to SQLite database (creates it if it doesn't exist)
    conn = sqlite3.connect('words.db')
    cursor = conn.cursor()

    # Create words table with new schema
    cursor.execute('''
    CREATE TABLE IF NOT EXISTS words (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        english TEXT NOT NULL,
        spanish TEXT NOT NULL
    )
    ''')

    # Create word_reviews table
    cursor.execute('''
    CREATE TABLE IF NOT EXISTS word_reviews (
        word_id INTEGER PRIMARY KEY,
        correct_count INTEGER DEFAULT 0,
        wrong_count INTEGER DEFAULT 0,
        FOREIGN KEY (word_id) REFERENCES words (id)
    )
    ''')

    # Create groups table
    cursor.execute('''
    CREATE TABLE IF NOT EXISTS groups (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL
    )
    ''')

    # Create word_groups table
    cursor.execute('''
    CREATE TABLE IF NOT EXISTS word_groups (
        word_id INTEGER,
        group_id INTEGER,
        PRIMARY KEY (word_id, group_id),
        FOREIGN KEY (word_id) REFERENCES words (id),
        FOREIGN KEY (group_id) REFERENCES groups (id)
    )
    ''')

    # Load and insert adjectives
    with open('seed/data_adjectives.json', 'r') as f:
        adjectives = json.load(f)
        for adj in adjectives:
            cursor.execute('INSERT INTO words (english, spanish) VALUES (?, ?)',
                         (adj['english'], adj['spanish']))
    
    # Load and insert verbs
    with open('seed/data_verbs.json', 'r') as f:
        verbs = json.load(f)
        for verb in verbs:
            cursor.execute('INSERT INTO words (english, spanish) VALUES (?, ?)',
                         (verb['english'], verb['spanish']))

    # Create a default group for each type
    cursor.execute('INSERT INTO groups (name) VALUES (?)', ('Adjectives',))
    adj_group_id = cursor.lastrowid
    cursor.execute('INSERT INTO groups (name) VALUES (?)', ('Verbs',))
    verb_group_id = cursor.lastrowid

    # Add words to their respective groups
    cursor.execute('SELECT id FROM words WHERE english IN (SELECT english FROM json_each(?))',
                  (json.dumps([adj['english'] for adj in adjectives]),))
    adj_ids = cursor.fetchall()
    for adj_id in adj_ids:
        cursor.execute('INSERT INTO word_groups (word_id, group_id) VALUES (?, ?)',
                      (adj_id[0], adj_group_id))

    cursor.execute('SELECT id FROM words WHERE english IN (SELECT english FROM json_each(?))',
                  (json.dumps([verb['english'] for verb in verbs]),))
    verb_ids = cursor.fetchall()
    for verb_id in verb_ids:
        cursor.execute('INSERT INTO word_groups (word_id, group_id) VALUES (?, ?)',
                      (verb_id[0], verb_group_id))

    conn.commit()
    conn.close()

if __name__ == '__main__':
    init_db()
