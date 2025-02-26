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
	mockRepo := new(mocks.MockPolishWordRepository)
	resolver := &mutationResolver{
		Resolver: &Resolver{
			PolishWordRepo: mockRepo,
		},
	}
	return mockRepo, resolver
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
