package repository

import (
	"context"
	"database/sql"

	"github.com/grzegorzpapaj/graphql-dictionary-api/internal/graph/model"
)

type TranslationRepositoryDB struct {
	DB *sql.DB
}

func (tr *TranslationRepositoryDB) AddTranslation(ctx context.Context, polishWordID *string, polishWord *string, translation *model.AddTranslationInput) (*model.Translation, error) {

	newTranslation := &model.Translation{
		EnglishWord:      translation.EnglishWord,
		ExampleSentences: []*model.ExampleSentence{},
	}

	err := tr.DB.QueryRowContext(ctx, "INSERT INTO translations(english_word, polish_word_id) VALUES ($1, $2) RETURNING id",
		newTranslation.EnglishWord, polishWordID).Scan(&newTranslation.ID)

	if err != nil {
		return nil, err
	}

	for _, es := range translation.ExampleSentences {
		newExampleSentence := &model.ExampleSentence{
			SentencePl: es.SentencePl,
			SentenceEn: es.SentenceEn,
		}

		err = tr.DB.QueryRowContext(ctx, "INSERT INTO example_sentences(sentence_pl, sentence_en, translation_id) VALUES ($1, $2, $3) RETURNING id",
			newExampleSentence.SentencePl, newExampleSentence.SentenceEn, newTranslation.ID).Scan(&newExampleSentence.ID)

		if err != nil {
			return nil, err
		}

		newTranslation.ExampleSentences = append(newTranslation.ExampleSentences, newExampleSentence)
	}
	return newTranslation, nil
}
