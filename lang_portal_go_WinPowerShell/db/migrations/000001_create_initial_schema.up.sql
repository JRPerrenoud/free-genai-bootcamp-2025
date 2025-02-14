-- Drop existing tables if they exist
DROP TABLE IF EXISTS word_reviews;
DROP TABLE IF EXISTS study_sessions;
DROP TABLE IF EXISTS study_activities;
DROP TABLE IF EXISTS group_words;
DROP TABLE IF EXISTS groups;
DROP TABLE IF EXISTS words;

-- Create words table
CREATE TABLE IF NOT EXISTS words (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    spanish TEXT NOT NULL,
    english TEXT NOT NULL,
    parts TEXT NOT NULL,
    correct_count INTEGER DEFAULT 0,
    wrong_count INTEGER DEFAULT 0
);

-- Create groups table
CREATE TABLE IF NOT EXISTS groups (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL
);

-- Create group_words table
CREATE TABLE IF NOT EXISTS group_words (
    group_id INTEGER NOT NULL,
    word_id INTEGER NOT NULL,
    PRIMARY KEY (group_id, word_id),
    FOREIGN KEY (group_id) REFERENCES groups(id),
    FOREIGN KEY (word_id) REFERENCES words(id)
);

-- Create study_activities table
CREATE TABLE IF NOT EXISTS study_activities (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    description TEXT NOT NULL
);

-- Create study_sessions table
CREATE TABLE IF NOT EXISTS study_sessions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    activity_id INTEGER NOT NULL,
    start_time DATETIME NOT NULL,
    end_time DATETIME,
    FOREIGN KEY (activity_id) REFERENCES study_activities(id)
);

-- Create word_reviews table
CREATE TABLE IF NOT EXISTS word_reviews (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    study_session_id INTEGER NOT NULL,
    word_id INTEGER NOT NULL,
    is_correct BOOLEAN NOT NULL,
    FOREIGN KEY (study_session_id) REFERENCES study_sessions(id),
    FOREIGN KEY (word_id) REFERENCES words(id)
);

-- Create indexes
CREATE INDEX idx_words_spanish ON words(spanish);
CREATE INDEX idx_words_english ON words(english);
CREATE INDEX idx_group_words_group_id ON group_words(group_id);
CREATE INDEX idx_group_words_word_id ON group_words(word_id);
CREATE INDEX idx_word_reviews_session_id ON word_reviews(study_session_id);
CREATE INDEX idx_word_reviews_word_id ON word_reviews(word_id);
