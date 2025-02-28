package repository

import (
	"context"

	"github.com/grzegorzpapaj/graphql-dictionary-api/internal/graph/model"
)

func (esr *ExampleSentenceRepositoryDB) insertExampleSentence(ctx context.Context, translationID, sentencePl, sentenceEn string) (string, error) {
	var id string
	err := esr.DB.QueryRowContext(ctx,
		"INSERT INTO example_sentences (sentence_pl, sentence_en, translation_id) VALUES ($1, $2, $3) RETURNING id",
		sentencePl, sentenceEn, translationID,
	).Scan(&id)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (esr *ExampleSentenceRepositoryDB) fetchTranslationAndPolishWord(ctx context.Context, translationID string) (*model.Translation, error) {
	var translation model.Translation
	var polishWord model.PolishWord

	err := esr.DB.QueryRowContext(ctx, `
		SELECT t.id, t.english_word, p.id, p.word
		FROM translations t
		JOIN polish_words p ON t.polish_word_id = p.id
		WHERE t.id = $1`, translationID,
	).Scan(&translation.ID, &translation.EnglishWord, &polishWord.ID, &polishWord.Word)
	if err != nil {
		return nil, err
	}

	translation.PolishWord = &polishWord
	return &translation, nil
}
