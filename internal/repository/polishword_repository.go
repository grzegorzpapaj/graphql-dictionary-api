package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/grzegorzpapaj/graphql-dictionary-api/internal/graph/model"
)

type PolishWordRepositoryDB struct {
	DB *sql.DB
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
		newTranslation := &model.Translation{
			EnglishWord:      t.EnglishWord,
			ExampleSentences: []*model.ExampleSentence{},
		}

		err = pwr.DB.QueryRowContext(ctx, "INSERT INTO translations(english_word, polish_word_id) VALUES ($1, $2) RETURNING id",
			newTranslation.EnglishWord, newPolishWord.ID).Scan(&newTranslation.ID)

		if err != nil {
			return nil, err
		}

		for _, es := range t.ExampleSentences {
			newExampleSentence := &model.ExampleSentence{
				SentencePl: es.SentencePl,
				SentenceEn: es.SentenceEn,
			}

			err = pwr.DB.QueryRowContext(ctx, "INSERT INTO example_sentences(sentence_pl, sentence_en, translation_id) VALUES ($1, $2, $3) RETURNING id",
				newExampleSentence.SentencePl, newExampleSentence.SentenceEn, newTranslation.ID).Scan(&newExampleSentence.ID)

			if err != nil {
				return nil, err
			}

			newTranslation.ExampleSentences = append(newTranslation.ExampleSentences, newExampleSentence)
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
		_, err := pwr.DB.ExecContext(ctx,
			"UPDATE polish_words SET word = $1 WHERE id = $2",
			*edits.Word, polishWordToEdit.ID)

		if err != nil {
			return nil, err
		}
		polishWordToEdit.Word = *edits.Word
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
