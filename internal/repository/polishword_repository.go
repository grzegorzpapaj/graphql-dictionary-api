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

	var pw model.PolishWord

	err := pwr.DB.QueryRowContext(ctx, `
			
				INSERT INTO polish_words (word)
				VALUES ($1)
				ON CONFLICT (word) DO UPDATE SET word = EXCLUDED.word
				RETURNING id, version
			
			
			`, polishWord.Word).Scan(&pw.ID, &pw.Version)

	if err != nil {
		return nil, fmt.Errorf("failed to upsert polish word: %w", err)
	}

	pw.Word = polishWord.Word
	pw.Translations = []*model.Translation{}

	for _, t := range polishWord.Translations {

		newTranslation, err := pwr.TranslationRepo.AddTranslation(ctx, &pw.ID, &pw.Word, t)

		if err != nil {
			return nil, err
		}

		pw.Translations = append(pw.Translations, newTranslation)
	}

	return &pw, nil
}

func (pwr *PolishWordRepositoryDB) DeletePolishWord(ctx context.Context, id *string, word *string) (*model.PolishWord, error) {
	var deletedPolishWord model.PolishWord

	id, err := pwr.TranslationRepo.getTargetPolishWordID(ctx, id, word)

	if err != nil {
		return nil, err
	}

	deletedPolishWord.ID = *id

	translations, err := pwr.getTranslationsWithExampleSentences(ctx, *id)
	if err != nil {
		return nil, err
	}

	deletedPolishWord.Translations = translations

	err = pwr.DB.QueryRowContext(ctx, "DELETE FROM polish_words WHERE id = $1 RETURNING id, word, version",
		*id).Scan(id, &deletedPolishWord.Word, &deletedPolishWord.Version)
	if err != nil {
		return nil, err
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
	rows, err := pwr.DB.QueryContext(ctx, "SELECT id, word, version FROM polish_words")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var polishWords []*model.PolishWord
	for rows.Next() {
		var pw model.PolishWord

		if err := rows.Scan(&pw.ID, &pw.Word, &pw.Version); err != nil {
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
