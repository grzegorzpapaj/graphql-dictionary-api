package repository

import (
	"context"
	"database/sql"

	"github.com/grzegorzpapaj/graphql-dictionary-api/internal/graph/model"
)

type PolishWordRepository struct {
	DB *sql.DB
}

func (pwr *PolishWordRepository) AddPolishWord(ctx context.Context, polishWord model.AddPolishWordInput) (*model.PolishWord, error) {

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

func (pwr *PolishWordRepository) GetAllPolishWords(ctx context.Context) ([]*model.PolishWord, error) {
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
		polishWords = append(polishWords, &pw)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return polishWords, nil
}
