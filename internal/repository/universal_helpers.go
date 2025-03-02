package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/grzegorzpapaj/graphql-dictionary-api/internal/graph/model"
)

func UpdateSingleTranslation(ctx context.Context, db *sql.DB, translation *model.Translation, editTr *model.EditTranslationInput) error {

	if editTr.EnglishWord != nil {
		result, err := db.ExecContext(ctx, "UPDATE translations SET english_word = $1, version = version + 1 WHERE id = $2 AND version = $3",
			*editTr.EnglishWord, translation.ID, translation.Version)

		if err != nil {
			return err
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return err
		}

		if rowsAffected == 0 {
			return fmt.Errorf("this translation has been modified by a different process")
		}

		translation.EnglishWord = *editTr.EnglishWord
		translation.Version++
	}

	if editTr.ExampleSentences != nil {
		exampleSentencesForTranslation, err := UpdateExampleSentences(ctx, db, translation.ID, editTr.ExampleSentences)

		if err != nil {
			return err
		}
		translation.ExampleSentences = exampleSentencesForTranslation
	}

	return nil

}

func UpdateExampleSentences(
	ctx context.Context,
	db *sql.DB,
	translationID string,
	editExamples []*model.EditExampleSentenceInput,
) ([]*model.ExampleSentence, error) {

	currentExampleSentencesFromDB, err := GetCurrentExampleSentencesFromDB(ctx, db, translationID)
	if err != nil {
		return nil, err
	}

	for i, editEs := range editExamples {

		if i < len(currentExampleSentencesFromDB) {

			err := UpdateSingleExampleSentence(ctx, db, currentExampleSentencesFromDB[i], editEs)
			if err != nil {
				return nil, err
			}

		} else {

			newExampleSentence, err := InsertExampleSentence(ctx, db, translationID, editEs)
			if err != nil {
				return nil, err
			}

			currentExampleSentencesFromDB = append(currentExampleSentencesFromDB, newExampleSentence)
		}
	}

	return currentExampleSentencesFromDB, nil
}

func GetCurrentExampleSentencesFromDB(
	ctx context.Context,
	db *sql.DB,
	translationID string,
) ([]*model.ExampleSentence, error) {
	rows, err := db.QueryContext(ctx,
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

	return currentExampleSentencesFromDB, nil
}

func UpdateSingleExampleSentence(
	ctx context.Context,
	db *sql.DB,
	exampleSentence *model.ExampleSentence,
	editEs *model.EditExampleSentenceInput,
) error {
	sentencePl := exampleSentence.SentencePl
	sentenceEn := exampleSentence.SentenceEn

	if editEs.SentencePl != nil {
		sentencePl = *editEs.SentencePl
	}

	if editEs.SentenceEn != nil {
		sentenceEn = *editEs.SentenceEn
	}

	result, err := db.ExecContext(ctx,
		"UPDATE example_sentences SET sentence_pl = $1, sentence_en = $2, version = version + 1 WHERE id = $3 AND version = $4",
		sentencePl, sentenceEn, exampleSentence.ID, exampleSentence.Version)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("this example sentence has been modified by a different process")
	}

	exampleSentence.SentencePl = sentencePl
	exampleSentence.SentenceEn = sentenceEn
	exampleSentence.Version++

	return nil

}

func InsertExampleSentence(
	ctx context.Context,
	db *sql.DB,
	translationID string,
	editEs *model.EditExampleSentenceInput) (*model.ExampleSentence, error) {

	var newExampleSentenceID string
	sentencePl := ""
	sentenceEn := ""

	if editEs.SentencePl != nil {
		sentencePl = *editEs.SentencePl
	}

	if editEs.SentenceEn != nil {
		sentenceEn = *editEs.SentenceEn
	}

	err := db.QueryRowContext(ctx,
		"INSERT INTO example_sentences (sentence_pl, sentence_en, translation_id) VALUES ($1, $2, $3) RETURNING id",
		sentencePl, sentenceEn, translationID).Scan(&newExampleSentenceID)

	if err != nil {
		return nil, err
	}

	return &model.ExampleSentence{
		ID:         newExampleSentenceID,
		SentencePl: sentencePl,
		SentenceEn: sentenceEn,
	}, nil
}
