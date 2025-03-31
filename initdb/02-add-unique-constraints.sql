ALTER TABLE translations 
ADD CONSTRAINT uq_translation_pwid_englishword 
UNIQUE (polish_word_id, english_word);

ALTER TABLE example_sentences 
ADD CONSTRAINT uq_example_sentence_tid_senpl_senen
UNIQUE (translation_id, sentence_pl, sentence_en);
