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

	id, version, err := esr.insertExampleSentence(ctx, translationID, newExampleSentence.SentencePl, newExampleSentence.SentenceEn)
	if err != nil {
		return nil, err
	}
	newExampleSentence.ID = id
	newExampleSentence.Version = version

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
	err := esr.DB.QueryRowContext(ctx, "DELETE FROM example_sentences WHERE id = $1 RETURNING sentence_pl, sentence_en, translation_id, version", id).
		Scan(&deletedEs.SentencePl, &deletedEs.SentenceEn, &deletedEs.Translation.ID, &deletedEs.Version)

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

func (esr *ExampleSentenceRepositoryDB) UpdateExampleSentence(ctx context.Context, id string, edits model.EditExampleSentenceInput) (*model.ExampleSentence, error) {

	es := &model.ExampleSentence{
		ID:          id,
		Translation: &model.Translation{},
	}

	var translationID string

	err := esr.DB.QueryRowContext(ctx, "SELECT sentence_pl, sentence_en, translation_id, version FROM example_sentences WHERE id = $1", id).
		Scan(&es.SentencePl, &es.SentenceEn, &translationID, &es.Version)

	if err != nil {
		return nil, err
	}

	if err := UpdateSingleExampleSentence(ctx, esr.DB, es, &edits); err != nil {
		return nil, err
	}

	translation, err := esr.fetchTranslationAndPolishWord(ctx, translationID)
	if err != nil {
		return nil, err
	}

	es.Translation = translation

	return es, nil
}

func (esr *ExampleSentenceRepositoryDB) GetSingleExampleSentence(ctx context.Context, id string) (*model.ExampleSentence, error) {

	es := &model.ExampleSentence{
		ID:          id,
		Translation: &model.Translation{},
	}

	var translationID string

	err := esr.DB.QueryRowContext(ctx, "SELECT sentence_pl, sentence_en, translation_id, version FROM example_sentences WHERE id = $1", id).
		Scan(&es.SentencePl, &es.SentenceEn, &translationID, &es.Version)

	if err != nil {
		return nil, err
	}

	translation, err := esr.fetchTranslationAndPolishWord(ctx, translationID)
	if err != nil {
		return nil, err
	}

	es.Translation = translation

	return es, nil
}

func (esr *ExampleSentenceRepositoryDB) GetExampleSentencesByTranslationId(ctx context.Context, translationID string) ([]*model.ExampleSentence, error) {

	exampleSentences, err := GetCurrentExampleSentencesFromDB(ctx, esr.DB, translationID)

	if err != nil {
		return nil, err
	}

	for _, es := range exampleSentences {

		translation, err := esr.fetchTranslationAndPolishWord(ctx, translationID)

		if err != nil {
			return nil, err
		}

		es.Translation = translation
	}

	return exampleSentences, nil
}
