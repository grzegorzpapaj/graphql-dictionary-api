ALTER TABLE translations 
ADD CONSTRAINT uq_translation_pwid_englishword 
UNIQUE (polish_word_id, english_word);