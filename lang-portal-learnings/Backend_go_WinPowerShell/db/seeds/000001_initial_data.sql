-- Insert initial words
INSERT INTO words (spanish, english, parts) VALUES
    ('hola', 'hello', 'interjection'),
    ('gracias', 'thank you', 'interjection'),
    ('por favor', 'please', 'interjection'),
    ('buenos días', 'good morning', 'interjection'),
    ('buenas noches', 'good night', 'interjection');

-- Insert initial groups
INSERT INTO groups (name) VALUES
    ('Greetings'),
    ('Polite Phrases');

-- Link words to groups
INSERT INTO group_words (group_id, word_id) VALUES
    (1, 1),  -- hola -> Greetings
    (1, 4),  -- buenos días -> Greetings
    (1, 5),  -- buenas noches -> Greetings
    (2, 2),  -- gracias -> Polite Phrases
    (2, 3);  -- por favor -> Polite Phrases

-- Insert study activities
INSERT INTO study_activities (name, description) VALUES
    ('Flashcards', 'Practice vocabulary with flashcards'),
    ('Multiple Choice', 'Choose the correct translation'),
    ('Writing Practice', 'Write the translation from memory');
