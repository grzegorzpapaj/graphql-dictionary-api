ALTER TABLE polish_words ADD COLUMN version INTEGER NOT NULL DEFAULT 1;
ALTER TABLE translations ADD COLUMN version INTEGER NOT NULL DEFAULT 1;
ALTER TABLE example_sentences ADD COLUMN version INTEGER NOT NULL DEFAULT 1;