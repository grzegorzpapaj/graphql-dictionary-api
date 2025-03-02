package repository

import (
	"context"

	"github.com/grzegorzpapaj/graphql-dictionary-api/internal/graph/model"
)

func (esr *ExampleSentenceRepositoryDB) insertExampleSentence(ctx context.Context, translationID, sentencePl, sentenceEn string) (string, int, error) {
	var id string
	var version int

	err := esr.DB.QueryRowContext(ctx,
		"INSERT INTO example_sentences (sentence_pl, sentence_en, translation_id) VALUES ($1, $2, $3) RETURNING id, version",
		sentencePl, sentenceEn, translationID,
	).Scan(&id, &version)
	if err != nil {
		return "", -1, err
	}
	return id, version, nil
}

func (esr *ExampleSentenceRepositoryDB) fetchTranslationAndPolishWord(ctx context.Context, translationID string) (*model.Translation, error) {
	var translation model.Translation
	var polishWord model.PolishWord

	err := esr.DB.QueryRowContext(ctx, `
		SELECT t.id, t.english_word, t.version, p.id, p.word, p.version
		FROM translations t
		JOIN polish_words p ON t.polish_word_id = p.id
		WHERE t.id = $1`, translationID,
	).Scan(&translation.ID, &translation.EnglishWord, &translation.Version, &polishWord.ID, &polishWord.Word, &polishWord.Version)
	if err != nil {
		return nil, err
	}

	translation.PolishWord = &polishWord
	return &translation, nil
}
