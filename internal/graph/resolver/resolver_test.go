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

func TestAddPolishWord(t *testing.T) {

	mockRepo := new(mocks.MockPolishWordRepository)

	mutation := &mutationResolver{
		Resolver: &Resolver{
			PolishWordRepo: mockRepo,
		},
	}

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
	mockRepo := new(mocks.MockPolishWordRepository)

	mutation := &mutationResolver{
		Resolver: &Resolver{
			PolishWordRepo: mockRepo,
		},
	}

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
	mockRepo := new(mocks.MockPolishWordRepository)

	mutation := &mutationResolver{
		Resolver: &Resolver{
			PolishWordRepo: mockRepo,
		},
	}

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
	mockRepo := new(mocks.MockPolishWordRepository)

	mutation := &mutationResolver{
		Resolver: &Resolver{
			PolishWordRepo: mockRepo,
		},
	}

	mockRepo.On("DeletePolishWord", mock.Anything, (*string)(nil), (*string)(nil)).
		Return(nil, errors.New("either id or word must be provided")).Once()

	result, err := mutation.DeletePolishWord(context.Background(), nil, nil)

	require.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}
