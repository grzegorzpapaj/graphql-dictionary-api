CREATE TABLE polish_words (
    id SERIAL PRIMARY KEY,
    word VARCHAR(50) NOT NULL UNIQUE
);

CREATE TABLE translations (
    id SERIAL PRIMARY KEY,
    polish_word_id INTEGER NOT NULL,
    english_word VARCHAR(50) NOT NULL,

    CONSTRAINT fk_polish_word FOREIGN KEY (polish_word_id) REFERENCES polish_words (id) ON DELETE CASCADE
);

CREATE TABLE example_sentences (
    id SERIAL PRIMARY KEY,
    translation_id INTEGER NOT NULL,
    sentence_pl TEXT NOT NULL,
    sentence_en TEXT NOT NULL,

    CONSTRAINT fk_translation FOREIGN KEY (translation_id) REFERENCES translations (id) ON DELETE CASCADE
);