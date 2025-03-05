package repository

import (
	"context"
	"database/sql"

	"github.com/grzegorzpapaj/graphql-dictionary-api/internal/graph/model"
)

type TranslationRepositoryDB struct {
	DB                  *sql.DB
	ExampleSentenceRepo *ExampleSentenceRepositoryDB
}

func (tr *TranslationRepositoryDB) AddTranslation(ctx context.Context, polishWordID *string, polishWord *string, translation *model.AddTranslationInput) (*model.Translation, error) {

	targetPolishWordID, err := tr.getTargetPolishWordID(ctx, polishWordID, polishWord)

	if err != nil {
		return nil, err
	}

	newTranslation := &model.Translation{
		EnglishWord:      translation.EnglishWord,
		ExampleSentences: []*model.ExampleSentence{},
	}

	err = tr.DB.QueryRowContext(ctx, "INSERT INTO translations(english_word, polish_word_id) VALUES ($1, $2) RETURNING id, version",
		newTranslation.EnglishWord, *targetPolishWordID).Scan(&newTranslation.ID, &newTranslation.Version)

	if err != nil {
		return nil, err
	}

	if polishWord != nil {
		newTranslation.PolishWord = &model.PolishWord{
			ID:   *targetPolishWordID,
			Word: *polishWord,
		}
	} else {
		newTranslation.PolishWord, err = tr.prepareWordWithId(ctx, targetPolishWordID)

		if err != nil {
			return nil, err
		}
	}

	for _, es := range translation.ExampleSentences {
		newExampleSentence, err := tr.ExampleSentenceRepo.AddExampleSentence(ctx, newTranslation.ID, *es)

		if err != nil {
			return nil, err
		}
		newTranslation.ExampleSentences = append(newTranslation.ExampleSentences, newExampleSentence)
	}
	return newTranslation, nil
}

func (tr *TranslationRepositoryDB) DeleteTranslation(ctx context.Context, id string) (*model.Translation, error) {
	var deletedTranslation model.Translation
	deletedTranslation.PolishWord = &model.PolishWord{}

	exampleSentences, err := GetCurrentExampleSentencesFromDB(ctx, tr.DB, id)

	if err != nil {
		return nil, err
	}

	deletedTranslation.ExampleSentences = exampleSentences

	err = tr.DB.QueryRowContext(ctx, "DELETE FROM translations WHERE id = $1 RETURNING id, english_word, polish_word_id, version", id).
		Scan(&deletedTranslation.ID, &deletedTranslation.EnglishWord, &deletedTranslation.PolishWord.ID, &deletedTranslation.Version)

	if err != nil {
		return nil, err
	}

	var fetchedPolishWord string
	var fetchedPolishWordVersion int
	err = tr.DB.QueryRowContext(ctx, "SELECT word, version FROM polish_words WHERE id = $1", deletedTranslation.PolishWord.ID).Scan(&fetchedPolishWord, &fetchedPolishWordVersion)

	if err != nil {
		return nil, err
	}
	deletedTranslation.PolishWord.Word = fetchedPolishWord
	deletedTranslation.PolishWord.Version = fetchedPolishWordVersion

	return &deletedTranslation, nil
}

func (tr *TranslationRepositoryDB) UpdateTranslation(ctx context.Context, id string, edits model.EditTranslationInput) (*model.Translation, error) {
	var translation model.Translation
	translation.PolishWord = &model.PolishWord{}

	err := tr.DB.QueryRowContext(ctx, "SELECT id, english_word, polish_word_id, version FROM translations WHERE id = $1", id).
		Scan(&translation.ID, &translation.EnglishWord, &translation.PolishWord.ID, &translation.Version)

	if err != nil {
		return nil, err
	}

	if err := UpdateSingleTranslation(ctx, tr.DB, &translation, &edits); err != nil {
		return nil, err
	}

	err = tr.DB.QueryRowContext(ctx, "SELECT word, version FROM polish_words WHERE id = $1", translation.PolishWord.ID).
		Scan(&translation.PolishWord.Word, &translation.PolishWord.Version)

	if err != nil {
		return nil, err
	}

	return &translation, nil
}

func (tr *TranslationRepositoryDB) GetSingleTranslationByID(ctx context.Context, id string) (*model.Translation, error) {
	var translation model.Translation
	translation.PolishWord = &model.PolishWord{}

	err := tr.DB.QueryRowContext(ctx, "SELECT id, english_word, polish_word_id, version FROM translations WHERE id = $1", id).
		Scan(&translation.ID, &translation.EnglishWord, &translation.PolishWord.ID, &translation.Version)

	if err != nil {
		return nil, err
	}

	err = tr.DB.QueryRowContext(ctx, "SELECT word, version FROM polish_words WHERE id = $1", translation.PolishWord.ID).
		Scan(&translation.PolishWord.Word, &translation.PolishWord.Version)

	if err != nil {
		return nil, err
	}

	examples, err := GetCurrentExampleSentencesFromDB(ctx, tr.DB, translation.ID)

	if err != nil {
		return nil, err
	}

	translation.ExampleSentences = examples

	return &translation, nil
}
