package resolver

import (
	"context"
	"errors"
	"testing"

	"github.com/grzegorzpapaj/graphql-dictionary-api/internal/graph/model"
	"github.com/grzegorzpapaj/graphql-dictionary-api/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func setupTestMutationResolver() (*mocks.MockPolishWordRepository, *mutationResolver) {

	mockPolishWordRepo := new(mocks.MockPolishWordRepository)
	resolver := &mutationResolver{
		Resolver: &Resolver{
			PolishWordRepo: mockPolishWordRepo,
		},
	}
	return mockPolishWordRepo, resolver
}

func setupTestTranslationMutationResolver() (*mocks.MockTranslationRepository, *mutationResolver) {
	mockTranslationRepo := new(mocks.MockTranslationRepository)
	resolver := &mutationResolver{
		Resolver: &Resolver{
			TranslationRepo: mockTranslationRepo,
		},
	}

	return mockTranslationRepo, resolver
}

func setupTestExampleSentenceMutationResolver() (*mocks.MockExampleSentenceRepository, *mutationResolver) {
	mockExampleSentenceRepo := new(mocks.MockExampleSentenceRepository)
	resolver := &mutationResolver{
		Resolver: &Resolver{
			ExampleSentenceRepo: mockExampleSentenceRepo,
		},
	}

	return mockExampleSentenceRepo, resolver
}

func TestAddPolishWord(t *testing.T) {

	mockRepo, mutation := setupTestMutationResolver()

	input := model.AddPolishWordInput{
		Word: "test_pierwszy",
		Translations: []*model.AddTranslationInput{
			{
				EnglishWord: "test_first",
				ExampleSentences: []*model.AddExampleSentenceInput{
					{
						SentencePl: "Test pierwszy",
						SentenceEn: "Test first",
					},
				},
			},
		},
	}

	expected := &model.PolishWord{
		ID:   "1",
		Word: "test_pierwszy",
		Translations: []*model.Translation{
			{
				ID:          "1",
				EnglishWord: "test_first",
				ExampleSentences: []*model.ExampleSentence{
					{
						ID:         "1",
						SentencePl: "Test pierwszy",
						SentenceEn: "Test first",
					},
				},
			},
		},
	}

	mockRepo.On("AddPolishWord", mock.Anything, input).Return(expected, nil).Once()

	result, err := mutation.AddPolishWord(context.Background(), input)

	require.NoError(t, err)
	assert.Equal(t, expected, result)

	mockRepo.AssertExpectations(t)
}

func TestDeletePolishWord_ByID_Success(t *testing.T) {

	mockRepo, mutation := setupTestMutationResolver()

	id := "1"

	expected := &model.PolishWord{
		ID:   "1",
		Word: "słowo",
	}

	mockRepo.On("DeletePolishWord", mock.Anything, &id, (*string)(nil)).Return(expected, nil).Once()

	result, err := mutation.DeletePolishWord(context.Background(), &id, nil)

	require.NoError(t, err)
	assert.Equal(t, expected, result)
	mockRepo.AssertExpectations(t)
}

func TestDeletePolishWord_ByWord_Success(t *testing.T) {

	mockRepo, mutation := setupTestMutationResolver()

	word := "słowo"

	expected := &model.PolishWord{
		ID:   "1",
		Word: "słowo",
	}

	mockRepo.On("DeletePolishWord", mock.Anything, (*string)(nil), &word).Return(expected, nil).Once()

	result, err := mutation.DeletePolishWord(context.Background(), nil, &word)

	require.NoError(t, err)
	assert.Equal(t, expected, result)
	mockRepo.AssertExpectations(t)
}

func TestDeletePolishWord_MissingParameters(t *testing.T) {

	mockRepo, mutation := setupTestMutationResolver()

	mockRepo.On("DeletePolishWord", mock.Anything, (*string)(nil), (*string)(nil)).
		Return(nil, errors.New("either id or word must be provided")).Once()

	result, err := mutation.DeletePolishWord(context.Background(), nil, nil)

	require.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestUpdatePolishWord(t *testing.T) {

	mockRepo, mutation := setupTestMutationResolver()

	id := "1"
	newWord := "newWord"
	newEnglishWord := "newEnglishWord"
	newSentencePL := "New sentence PL"
	newSentenceEN := "New sentence EN"
	newSentencePL2 := "New sentence PL2"
	newSentenceEN2 := "New sentence EN2"

	editTranslation := &model.EditTranslationInput{
		EnglishWord: &newEnglishWord,
		ExampleSentences: []*model.EditExampleSentenceInput{
			{
				SentencePl: &newSentencePL,
				SentenceEn: &newSentenceEN,
			},
			{
				SentencePl: &newSentencePL2,
				SentenceEn: &newSentenceEN2,
			},
		},
	}

	editPolishWord := &model.EditPolishWordInput{
		Word:         &newWord,
		Translations: []*model.EditTranslationInput{editTranslation},
	}

	expected := &model.PolishWord{
		ID:   id,
		Word: newWord,
		Translations: []*model.Translation{
			{
				ID:          "1",
				EnglishWord: newEnglishWord,
				ExampleSentences: []*model.ExampleSentence{
					{
						ID:         "1",
						SentencePl: newSentencePL,
						SentenceEn: newSentenceEN,
					},
					{
						ID:         "2",
						SentencePl: newSentencePL2,
						SentenceEn: newSentencePL2,
					},
				},
			},
		},
	}

	mockRepo.On("UpdatePolishWord", mock.Anything, &id, (*string)(nil), editPolishWord).Return(expected, nil).Once()

	result, err := mutation.UpdatePolishWord(context.Background(), &id, nil, editPolishWord)

	require.NoError(t, err)
	assert.Equal(t, expected, result)

	mockRepo.AssertExpectations(t)

}

// Translation tests

func TestAddTranslation(t *testing.T) {
	mockRepo, mutation := setupTestTranslationMutationResolver()

	input := model.AddTranslationInput{
		EnglishWord: "test_translation",
		ExampleSentences: []*model.AddExampleSentenceInput{
			{
				SentencePl: "Test tłumaczenia",
				SentenceEn: "Test translation",
			},
		},
	}

	expected := &model.Translation{
		ID:          "1",
		EnglishWord: "test_first",
		PolishWord: &model.PolishWord{
			ID:   "1",
			Word: "test_translation",
		},
		ExampleSentences: []*model.ExampleSentence{
			{
				ID:         "1",
				SentencePl: "Test tłumaczenia",
				SentenceEn: "Test translation",
			},
		},
	}

	mockRepo.On("AddTranslation", mock.Anything, mock.Anything, mock.Anything, &input).Return(expected, nil).Once()

	polishWordId := "1"

	result, err := mutation.TranslationRepo.AddTranslation(context.Background(), &polishWordId, nil, &input)
	require.NoError(t, err)
	assert.Equal(t, expected, result)
	mockRepo.AssertExpectations(t)

}

func TestDeleteTranslation(t *testing.T) {
	mockRepo, mutation := setupTestTranslationMutationResolver()

	idToDelete := "1"

	expected := &model.Translation{
		ID:          "1",
		EnglishWord: "delete_translation",
		PolishWord: &model.PolishWord{
			ID:   "1",
			Word: "test_polish_word",
		},

		ExampleSentences: []*model.ExampleSentence{},
	}

	mockRepo.On("DeleteTranslation", mock.Anything, idToDelete).Return(expected, nil).Once()

	result, err := mutation.DeleteTranslation(context.Background(), idToDelete)
	require.NoError(t, err)
	assert.Equal(t, expected, result)

	mockRepo.AssertExpectations(t)
}

func TestUpdateTranslation(t *testing.T) {
	mockRepo, mutation := setupTestTranslationMutationResolver()

	id := "1"
	newEnglishWord := "updated_translation"
	newSentencePL := "Nowe zdanie PL"
	newSentenceEN := "New sentence EN"

	editTranslation := model.EditTranslationInput{
		EnglishWord: &newEnglishWord,
		ExampleSentences: []*model.EditExampleSentenceInput{
			{
				SentencePl: &newSentencePL,
				SentenceEn: &newSentenceEN,
			},
		},
	}

	expected := &model.Translation{
		ID:          id,
		EnglishWord: newEnglishWord,
		PolishWord: &model.PolishWord{
			ID:   "1",
			Word: "test_polish_word",
		},
		ExampleSentences: []*model.ExampleSentence{
			{
				ID:         "1",
				SentencePl: newSentencePL,
				SentenceEn: newSentenceEN,
			},
		},
	}

	mockRepo.On("UpdateTranslation", mock.Anything, id, editTranslation).Return(expected, nil).Once()

	result, err := mutation.TranslationRepo.UpdateTranslation(context.Background(), id, editTranslation)

	require.NoError(t, err)
	assert.Equal(t, expected, result)

	mockRepo.AssertExpectations(t)
}

func TestGetSingleTranslationByID(t *testing.T) {
	mockRepo, query := setupTestTranslationMutationResolver()

	translationID := "1"

	expected := &model.Translation{
		ID:          translationID,
		EnglishWord: "translation",
		PolishWord: &model.PolishWord{
			ID:   "1",
			Word: "test_polishword",
		},
		ExampleSentences: []*model.ExampleSentence{
			{
				ID:         "1",
				SentencePl: "Testowe zdanie",
				SentenceEn: "Test sentence",
			},
			{
				ID:         "2",
				SentencePl: "Kolejne testowe zdanie",
				SentenceEn: "Another test sentence",
			},
		},
	}

	mockRepo.On("GetSingleTranslationByID", mock.Anything, translationID).Return(expected, nil).Once()

	result, err := query.TranslationRepo.GetSingleTranslationByID(context.Background(), translationID)

	require.NoError(t, err)
	assert.Equal(t, expected, result)

	mockRepo.AssertExpectations(t)
}

func TestAddExample(t *testing.T) {
	mockRepo, mutation := setupTestExampleSentenceMutationResolver()

	translationID := "1"

	input := model.AddExampleSentenceInput{
		SentencePl: "Zdanie testowe PL",
		SentenceEn: "Test sentence EN",
	}

	expected := &model.ExampleSentence{
		ID:         "1",
		SentencePl: "Zdanie testowe PL",
		SentenceEn: "Test sentence EN",
	}

	mockRepo.On("AddExampleSentence", mock.Anything, translationID, input).Return(expected, nil).Once()

	result, err := mutation.ExampleSentenceRepo.AddExampleSentence(context.Background(), translationID, input)

	require.NoError(t, err)
	assert.Equal(t, expected, result)

	mockRepo.AssertExpectations(t)
}

func TestDeleteExample(t *testing.T) {
	mockRepo, mutation := setupTestExampleSentenceMutationResolver()

	exampleID := "1"

	expectedExample := &model.ExampleSentence{
		ID:         exampleID,
		SentencePl: "Zdanie testowe PL",
		SentenceEn: "Test sentence EN",
		Translation: &model.Translation{
			ID: "1",
		},
	}

	mockRepo.On("DeleteExampleSentence", mock.Anything, exampleID).Return(expectedExample, nil).Once()
	result, err := mutation.ExampleSentenceRepo.DeleteExampleSentence(context.Background(), exampleID)

	require.NoError(t, err)
	assert.Equal(t, expectedExample, result)

	mockRepo.AssertExpectations(t)
}

func TestUpdateExample(t *testing.T) {

	mockRepo, mutation := setupTestExampleSentenceMutationResolver()

	exampleID := "1"
	newSentencePL := "Zaktualizowane zdanie PL"
	newSentenceEN := "Updated sentence EN"

	edits := model.EditExampleSentenceInput{
		SentencePl: &newSentencePL,
		SentenceEn: &newSentenceEN,
	}

	expectedExample := &model.ExampleSentence{
		ID:         exampleID,
		SentencePl: newSentencePL,
		SentenceEn: newSentenceEN,
		Translation: &model.Translation{
			ID: "1",
		},
	}

	mockRepo.On("UpdateExampleSentence", mock.Anything, exampleID, edits).Return(expectedExample, nil).Once()

	result, err := mutation.ExampleSentenceRepo.UpdateExampleSentence(context.Background(), exampleID, edits)

	require.NoError(t, err)
	assert.Equal(t, expectedExample, result)

	mockRepo.AssertExpectations(t)
}

func TestGetSinglePolishWordByID(t *testing.T) {
	mockRepo, query := setupTestMutationResolver()

	id := "1"
	var word *string = nil

	expected := &model.PolishWord{
		ID:   "1",
		Word: "test_word",
		Translations: []*model.Translation{
			{
				ID:          "1",
				EnglishWord: "test_english_word",
			},
		},
	}

	mockRepo.On("GetSinglePolishWord", mock.Anything, &id, word).Return(expected, nil).Once()

	result, err := query.PolishWordRepo.GetSinglePolishWord(context.Background(), &id, word)

	require.NoError(t, err)
	assert.Equal(t, expected, result)

	mockRepo.AssertExpectations(t)
}

func TestGetSinglePolishWordByWord(t *testing.T) {
	mockRepo, query := setupTestMutationResolver()

	var id *string = nil
	word := "test_word"

	expected := &model.PolishWord{
		ID:   "1",
		Word: "test_word",
		Translations: []*model.Translation{
			{
				ID:          "1",
				EnglishWord: "test_english_word",
			},
		},
	}

	mockRepo.On("GetSinglePolishWord", mock.Anything, id, &word).Return(expected, nil).Once()

	result, err := query.PolishWordRepo.GetSinglePolishWord(context.Background(), id, &word)

	require.NoError(t, err)
	assert.Equal(t, expected, result)

	mockRepo.AssertExpectations(t)
}

func TestGetSinglePolishWordNoParametersError(t *testing.T) {
	mockRepo, query := setupTestMutationResolver()

	var id *string = nil
	var word *string = nil

	mockRepo.On("GetSinglePolishWord", mock.Anything, id, word).Return(nil, errors.New("either id or word must be provided")).Once()

	result, err := query.PolishWordRepo.GetSinglePolishWord(context.Background(), id, word)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "either id or word must be provided")

	mockRepo.AssertExpectations(t)
}

func TestGetAllPolishWords(t *testing.T) {
	mockRepo, query := setupTestMutationResolver()

	expected := []*model.PolishWord{
		{
			ID:   "1",
			Word: "test_word",
			Translations: []*model.Translation{
				{
					ID:          "1",
					EnglishWord: "test_english_word",
				},
			},
		},
		{
			ID:   "2",
			Word: "another_word",
			Translations: []*model.Translation{
				{
					ID:          "2",
					EnglishWord: "another_english_word",
				},
			},
		},
	}

	mockRepo.On("GetAllPolishWords", mock.Anything).Return(expected, nil).Once()

	result, err := query.PolishWordRepo.GetAllPolishWords(context.Background())
	require.NoError(t, err)
	assert.Equal(t, expected, result)

	mockRepo.AssertExpectations(t)
}

func TestGetSingleExampleSentence(t *testing.T) {
	mockRepo, query := setupTestExampleSentenceMutationResolver()

	exampleID := "1"
	expectedExample := &model.ExampleSentence{
		ID:         "exampleID",
		SentencePl: "Zdanie testowe PL",
		SentenceEn: "Test sentence EN",
		Translation: &model.Translation{
			ID:          "1",
			EnglishWord: "test_translation",
		},
	}

	mockRepo.On("GetSingleExampleSentence", mock.Anything, exampleID).Return(expectedExample, nil).Once()

	result, err := query.ExampleSentenceRepo.GetSingleExampleSentence(context.Background(), exampleID)

	require.NoError(t, err)
	assert.Equal(t, expectedExample, result)

	mockRepo.AssertExpectations(t)
}
