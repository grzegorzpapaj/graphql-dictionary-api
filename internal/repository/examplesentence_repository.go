package repository

import (
	"context"
	"database/sql"

	"github.com/grzegorzpapaj/graphql-dictionary-api/internal/graph/model"
)

type ExampleSentenceRepositoryDB struct {
	DB *sql.DB
}

func (esr *ExampleSentenceRepositoryDB) AddExampleSentence(ctx context.Context, translationID string, exampleSentence model.AddExampleSentenceInput) (*model.ExampleSentence, error) {

	newExampleSentence := &model.ExampleSentence{
		SentencePl: exampleSentence.SentencePl,
		SentenceEn: exampleSentence.SentenceEn,
	}

	id, err := esr.insertExampleSentence(ctx, translationID, newExampleSentence.SentencePl, newExampleSentence.SentenceEn)
	if err != nil {
		return nil, err
	}
	newExampleSentence.ID = id

	translation, err := esr.fetchTranslationAndPolishWord(ctx, translationID)
	if err != nil {
		return nil, err
	}
	newExampleSentence.Translation = translation

	return newExampleSentence, nil

}

func (esr *ExampleSentenceRepositoryDB) DeleteExampleSentence(ctx context.Context, id string) (*model.ExampleSentence, error) {

	deletedEs := &model.ExampleSentence{
		ID:          id,
		Translation: &model.Translation{},
	}
	err := esr.DB.QueryRowContext(ctx, "DELETE FROM example_sentences WHERE id = $1 RETURNING sentence_pl, sentence_en, translation_id", id).
		Scan(&deletedEs.SentencePl, &deletedEs.SentenceEn, &deletedEs.Translation.ID)

	if err != nil {
		return nil, err
	}

	translation, err := esr.fetchTranslationAndPolishWord(ctx, deletedEs.Translation.ID)
	if err != nil {
		return nil, err
	}

	deletedEs.Translation = translation

	return deletedEs, nil
}
