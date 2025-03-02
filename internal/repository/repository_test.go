package repository

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/grzegorzpapaj/graphql-dictionary-api/internal/graph/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdatePolishWordVersionConflict(t *testing.T) {

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := &PolishWordRepositoryDB{
		DB: db,
	}

	ctx := context.Background()
	id := "1"

	mock.ExpectQuery("SELECT id, word, version FROM polish_words WHERE id = \\$1").
		WithArgs(id).
		WillReturnRows(sqlmock.NewRows([]string{"id", "word", "version"}).
			AddRow(id, "old_word", 1))

	newWord := "new_word"
	mock.ExpectExec("UPDATE polish_words SET word = \\$1, version = version \\+ 1 WHERE id = \\$2 AND version = \\$3").
		WithArgs(newWord, id, 1).
		WillReturnResult(sqlmock.NewResult(0, 0))

	edits := &model.EditPolishWordInput{
		Word:    &newWord,
		Version: 1,
	}

	_, err = repo.UpdatePolishWord(ctx, &id, nil, edits)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "this polish word has been modified by a different process")

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateTranslationVersionConflict(t *testing.T) {

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := &TranslationRepositoryDB{
		DB: db,
	}

	ctx := context.Background()
	id := "1"

	mock.ExpectQuery("SELECT id, english_word, polish_word_id, version FROM translations WHERE id = \\$1").
		WithArgs(id).
		WillReturnRows(sqlmock.NewRows([]string{"id", "english_word", "polish_word_id", "version"}).
			AddRow(id, "old_translation", "1", 1))

	newTranslation := "new_translation"
	mock.ExpectExec("UPDATE translations SET english_word = \\$1, version = version \\+ 1 WHERE id = \\$2 AND version = \\$3").
		WithArgs(newTranslation, id, 1).
		WillReturnResult(sqlmock.NewResult(0, 0))

	edits := model.EditTranslationInput{
		EnglishWord: &newTranslation,
		Version:     1,
	}

	_, err = repo.UpdateTranslation(ctx, id, edits)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "this translation has been modified by a different process")

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateExampleSentenceVersionConflict(t *testing.T) {

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := &ExampleSentenceRepositoryDB{
		DB: db,
	}

	ctx := context.Background()
	id := "1"

	mock.ExpectQuery("SELECT sentence_pl, sentence_en, translation_id, version FROM example_sentences WHERE id = \\$1").
		WithArgs(id).
		WillReturnRows(sqlmock.NewRows([]string{"sentence_pl", "sentence_en", "translation_id", "version"}).
			AddRow("old_sentence_pl", "old_sentence_en", "1", 1))

	newSentencePl := "new_sentence_pl"
	newSentenceEn := "new_sentence_en"
	mock.ExpectExec("UPDATE example_sentences SET sentence_pl = \\$1, sentence_en = \\$2, version = version \\+ 1 WHERE id = \\$3 AND version = \\$4").
		WithArgs(newSentencePl, newSentenceEn, id, 1).
		WillReturnResult(sqlmock.NewResult(0, 0))

	edits := model.EditExampleSentenceInput{
		SentencePl: &newSentencePl,
		SentenceEn: &newSentenceEn,
		Version:    1,
	}

	_, err = repo.UpdateExampleSentence(ctx, id, edits)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "this example sentence has been modified by a different process")

	require.NoError(t, mock.ExpectationsWereMet())
}
