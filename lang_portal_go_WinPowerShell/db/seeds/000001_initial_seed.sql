-- Insert sample words
INSERT INTO words (spanish, english, parts) VALUES
    ('hola', 'hello', '["interjection"]'),
    ('gracias', 'thank you', '["interjection"]'),
    ('libro', 'book', '["noun"]'),
    ('casa', 'house', '["noun"]'),
    ('comer', 'to eat', '["verb"]');

-- Insert sample groups
INSERT INTO groups (name) VALUES
    ('Basic Vocabulary'),
    ('Greetings'),
    ('Common Nouns'),
    ('Essential Verbs');

-- Link words to groups
INSERT INTO words_groups (word_id, group_id) VALUES
    (1, 1), -- hola -> Basic Vocabulary
    (1, 2), -- hola -> Greetings
    (2, 1), -- gracias -> Basic Vocabulary
    (2, 2), -- gracias -> Greetings
    (3, 1), -- libro -> Basic Vocabulary
    (3, 3), -- libro -> Common Nouns
    (4, 1), -- casa -> Basic Vocabulary
    (4, 3), -- casa -> Common Nouns
    (5, 1), -- comer -> Basic Vocabulary
    (5, 4); -- comer -> Essential Verbs

-- Insert sample study activities
INSERT INTO study_activities (name, thumbnail_url, description) VALUES
    ('Vocabulary Quiz', 'https://example.com/vocab-quiz.jpg', 'Test your vocabulary knowledge'),
    ('Flashcards', 'https://example.com/flashcards.jpg', 'Learn with flashcards'),
    ('Word Match', 'https://example.com/word-match.jpg', 'Match Spanish words with English translations');

-- Insert sample study sessions
INSERT INTO study_sessions (group_id, study_activity_id, created_at) VALUES
    (1, 1, datetime('now', '-2 days')),
    (2, 2, datetime('now', '-1 days')),
    (3, 3, datetime('now'));

-- Insert sample word review items
INSERT INTO word_review_items (word_id, study_session_id, correct, created_at) VALUES
    (1, 1, true, datetime('now', '-2 days')),
    (2, 1, true, datetime('now', '-2 days')),
    (3, 1, false, datetime('now', '-2 days')),
    (4, 2, true, datetime('now', '-1 days')),
    (5, 2, true, datetime('now', '-1 days'));
