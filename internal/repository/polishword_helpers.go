package repository

import (
	"context"
	"fmt"

	"github.com/grzegorzpapaj/graphql-dictionary-api/internal/graph/model"
)

func (pwr *PolishWordRepositoryDB) fetchPolishWords(ctx context.Context, id *string, word *string) (*model.PolishWord, error) {

	var fetchedPolishWord model.PolishWord
	if id != nil {
		err := pwr.DB.QueryRowContext(ctx, "SELECT id, word FROM polish_words WHERE id = $1",
			*id).Scan(&fetchedPolishWord.ID, &fetchedPolishWord.Word)
		if err != nil {
			return nil, err
		}
	} else if word != nil {
		err := pwr.DB.QueryRowContext(ctx, "SELECT id, word FROM polish_words WHERE word = $1",
			*word).Scan(&fetchedPolishWord.ID, &fetchedPolishWord.Word)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("either id or word must be provided")
	}

	return &fetchedPolishWord, nil
}

func (pwr *PolishWordRepositoryDB) updateTranslations(
	ctx context.Context,
	polishWordID string,
	editTranslations []*model.EditTranslationInput,
) ([]*model.Translation, error) {

	currentTranslationsFromDB, err := pwr.getCurrentTranslationsFromDB(ctx, polishWordID)
	if err != nil {
		return nil, err
	}

	for i, editTr := range editTranslations {

		if i < len(currentTranslationsFromDB) {

			err := pwr.updateSingleTranslation(ctx, currentTranslationsFromDB[i], editTr)

			if err != nil {
				return nil, err
			}

		} else {

			if editTr.EnglishWord != nil {
				var newTranslationID string
				err := pwr.DB.QueryRowContext(ctx,
					"INSERT INTO translations (english_word, polish_word_id) VALUES ($1, $2) RETURNING id",
					*editTr.EnglishWord, polishWordID).Scan(&newTranslationID)
				if err != nil {
					return nil, err
				}

				newTranslation := &model.Translation{
					ID:               newTranslationID,
					EnglishWord:      *editTr.EnglishWord,
					ExampleSentences: []*model.ExampleSentence{},
				}

				if editTr.ExampleSentences != nil {
					exampleSentencesForTranslation, err := pwr.insertExampleSentences(ctx, newTranslationID, editTr.ExampleSentences)

					if err != nil {
						return nil, err
					}

					newTranslation.ExampleSentences = exampleSentencesForTranslation
				}

				currentTranslationsFromDB = append(currentTranslationsFromDB, newTranslation)
			}

		}
	}

	return currentTranslationsFromDB, nil

}

func (pwr *PolishWordRepositoryDB) updateExampleSentences(
	ctx context.Context,
	translationID string,
	editExamples []*model.EditExampleSentenceInput,
) ([]*model.ExampleSentence, error) {

	rows, err := pwr.DB.QueryContext(ctx,
		"SELECT id, sentence_pl, sentence_en FROM example_sentences WHERE translation_id = $1 ORDER BY id", translationID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var currentExampleSentencesFromDB []*model.ExampleSentence
	for rows.Next() {
		var es model.ExampleSentence
		if err := rows.Scan(&es.ID, &es.SentencePl, &es.SentenceEn); err != nil {
			return nil, err
		}

		currentExampleSentencesFromDB = append(currentExampleSentencesFromDB, &es)
	}

	for i, editEs := range editExamples {

		if i < len(currentExampleSentencesFromDB) {

			sentencePl := currentExampleSentencesFromDB[i].SentencePl
			sentenceEn := currentExampleSentencesFromDB[i].SentenceEn

			if editEs.SentencePl != nil {
				sentencePl = *editEs.SentencePl
			}

			if editEs.SentenceEn != nil {
				sentenceEn = *editEs.SentenceEn
			}

			_, err := pwr.DB.ExecContext(ctx,
				"UPDATE example_sentences SET sentence_pl = $1, sentence_en = $2 WHERE id = $3",
				sentencePl, sentenceEn, currentExampleSentencesFromDB[i].ID)

			if err != nil {
				return nil, err
			}

			currentExampleSentencesFromDB[i].SentencePl = sentencePl
			currentExampleSentencesFromDB[i].SentenceEn = sentenceEn
		} else {

			var newExampleSentenceID string
			sentencePl := ""
			sentenceEn := ""

			if editEs.SentencePl != nil {
				sentencePl = *editEs.SentencePl
			}

			if editEs.SentenceEn != nil {
				sentenceEn = *editEs.SentenceEn
			}

			err := pwr.DB.QueryRowContext(ctx,
				"INSERT INTO example_sentences (sentence_pl, sentence_en, translation_id) VALUES ($1, $2, $3) RETURNING id",
				sentencePl, sentenceEn, translationID).Scan(&newExampleSentenceID)

			if err != nil {
				return nil, err
			}

			currentExampleSentencesFromDB = append(currentExampleSentencesFromDB, &model.ExampleSentence{
				ID:         newExampleSentenceID,
				SentencePl: sentencePl,
				SentenceEn: sentenceEn,
			})
		}
	}

	return currentExampleSentencesFromDB, nil
}

func (pwr *PolishWordRepositoryDB) insertExampleSentences(
	ctx context.Context,
	translationID string,
	editExamples []*model.EditExampleSentenceInput,
) ([]*model.ExampleSentence, error) {

	var exampleSentences []*model.ExampleSentence
	for _, editEs := range editExamples {
		var newExampleSentenceID string
		err := pwr.DB.QueryRowContext(ctx,
			"INSERT INTO example_sentences (sentence_pl, sentence_en, translation_id) VALUES ($1, $2, $3) RETURNING id",
			editEs.SentencePl, editEs.SentenceEn, translationID).Scan(&newExampleSentenceID)

		if err != nil {
			return nil, err
		}

		exampleSentences = append(exampleSentences, &model.ExampleSentence{
			ID:         newExampleSentenceID,
			SentencePl: *editEs.SentencePl,
			SentenceEn: *editEs.SentenceEn,
		})
	}

	return exampleSentences, nil
}

func (pwr *PolishWordRepositoryDB) getCurrentTranslationsFromDB(
	ctx context.Context,
	polishWordID string,
) ([]*model.Translation, error) {
	rows, err := pwr.DB.QueryContext(ctx,
		"SELECT id, english_word FROM translations WHERE polish_word_id = $1 ORDER BY id", polishWordID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var currentTranslationsFromDB []*model.Translation
	for rows.Next() {
		var t model.Translation
		if err := rows.Scan(&t.ID, &t.EnglishWord); err != nil {
			return nil, err
		}
		currentTranslationsFromDB = append(currentTranslationsFromDB, &t)
	}

	return currentTranslationsFromDB, nil
}

func (pwr *PolishWordRepositoryDB) updateSingleTranslation(ctx context.Context, translation *model.Translation, editTr *model.EditTranslationInput) error {

	if editTr.EnglishWord != nil {
		_, err := pwr.DB.ExecContext(ctx, "UPDATE translations SET english_word = $1 WHERE id = $2",
			*editTr.EnglishWord, translation.ID)

		if err != nil {
			return err
		}

		translation.EnglishWord = *editTr.EnglishWord
	}

	if editTr.ExampleSentences != nil {
		exampleSentencesForTranslation, err := pwr.updateExampleSentences(ctx, translation.ID, editTr.ExampleSentences)

		if err != nil {
			return err
		}
		translation.ExampleSentences = exampleSentencesForTranslation
	}

	return nil

}
