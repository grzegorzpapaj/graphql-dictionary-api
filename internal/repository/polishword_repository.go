package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/grzegorzpapaj/graphql-dictionary-api/internal/graph/model"
)

type PolishWordRepositoryDB struct {
	DB              *sql.DB
	TranslationRepo *TranslationRepositoryDB
}

func (pwr *PolishWordRepositoryDB) AddPolishWord(ctx context.Context, polishWord model.AddPolishWordInput) (*model.PolishWord, error) {

	newPolishWord := &model.PolishWord{
		Word:         polishWord.Word,
		Translations: []*model.Translation{},
	}

	err := pwr.DB.QueryRowContext(ctx, "INSERT INTO polish_words (word) VALUES ($1) RETURNING id", newPolishWord.Word).Scan(&newPolishWord.ID)
	if err != nil {
		return nil, err
	}

	for _, t := range polishWord.Translations {

		newTranslation, err := pwr.TranslationRepo.AddTranslation(ctx, &newPolishWord.ID, &polishWord.Word, t)

		if err != nil {
			return nil, err
		}

		newPolishWord.Translations = append(newPolishWord.Translations, newTranslation)

	}

	return newPolishWord, nil
}

func (pwr *PolishWordRepositoryDB) DeletePolishWord(ctx context.Context, id *string, word *string) (*model.PolishWord, error) {
	var deletedPolishWord model.PolishWord

	if id != nil {
		err := pwr.DB.QueryRowContext(ctx, "DELETE FROM polish_words WHERE id = $1 RETURNING id, word",
			*id).Scan(&deletedPolishWord.ID, &deletedPolishWord.Word)

		if err != nil {
			return nil, err
		}
	} else if word != nil {
		err := pwr.DB.QueryRowContext(ctx, "DELETE FROM polish_words WHERE word = $1 RETUNING id, word",
			*word).Scan(&deletedPolishWord.ID, &deletedPolishWord.Word)

		if err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("either id or word must be provided")
	}

	return &deletedPolishWord, nil
}

func (pwr *PolishWordRepositoryDB) UpdatePolishWord(ctx context.Context, id *string, word *string, edits *model.EditPolishWordInput) (*model.PolishWord, error) {

	polishWordToEdit, err := pwr.fetchPolishWords(ctx, id, word)
	if err != nil {
		return nil, err
	}

	if word == nil && edits.Word != nil {
		result, err := pwr.DB.ExecContext(ctx,
			"UPDATE polish_words SET word = $1, version = version + 1 WHERE id = $2 AND version = $3",
			*edits.Word, polishWordToEdit.ID, polishWordToEdit.Version)

		if err != nil {
			return nil, err
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return nil, err
		}
		if rowsAffected == 0 {
			return nil, fmt.Errorf("this polish word has been modified by a different process")
		}

		polishWordToEdit.Word = *edits.Word
		polishWordToEdit.Version++
	}

	if edits.Translations != nil {
		translations, err := pwr.updateTranslations(ctx, polishWordToEdit.ID, edits.Translations)
		if err != nil {
			return nil, err
		}
		polishWordToEdit.Translations = translations
	}

	return polishWordToEdit, nil

}

func (pwr *PolishWordRepositoryDB) GetAllPolishWords(ctx context.Context) ([]*model.PolishWord, error) {
	rows, err := pwr.DB.QueryContext(ctx, "SELECT id, word FROM polish_words")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var polishWords []*model.PolishWord
	for rows.Next() {
		var pw model.PolishWord

		if err := rows.Scan(&pw.ID, &pw.Word); err != nil {
			return nil, err
		}

		translations, err := pwr.getTranslationsWithExampleSentences(ctx, pw.ID)
		if err != nil {
			return nil, err
		}

		pw.Translations = translations

		polishWords = append(polishWords, &pw)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return polishWords, nil
}

func (pwr *PolishWordRepositoryDB) GetSinglePolishWord(ctx context.Context, id *string, word *string) (*model.PolishWord, error) {
	pw, err := pwr.fetchPolishWords(ctx, id, word)

	if err != nil {
		return nil, err
	}

	translations, err := pwr.getTranslationsWithExampleSentences(ctx, pw.ID)
	if err != nil {
		return nil, err
	}

	pw.Translations = translations

	return pw, nil
}
