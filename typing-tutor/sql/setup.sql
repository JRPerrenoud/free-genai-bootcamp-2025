-- Words table for Spanish vocabulary
CREATE TABLE words (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  english TEXT NOT NULL,
  spanish TEXT NOT NULL
);

-- Groups for organizing words
CREATE TABLE groups (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT NOT NULL
);

-- Junction table for word-group relationships
CREATE TABLE word_groups (
  word_id INTEGER NOT NULL,
  group_id INTEGER NOT NULL,
  FOREIGN KEY (word_id) REFERENCES words(id),
  FOREIGN KEY (group_id) REFERENCES groups(id)
);

-- Table for tracking user performance on words
CREATE TABLE word_reviews (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  word_id INTEGER NOT NULL,
  correct_count INTEGER DEFAULT 0,
  wrong_count INTEGER DEFAULT 0,
  last_reviewed TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (word_id) REFERENCES words(id)
);