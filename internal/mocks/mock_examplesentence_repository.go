package mocks

import (
	"context"

	"github.com/grzegorzpapaj/graphql-dictionary-api/internal/graph/model"
	"github.com/stretchr/testify/mock"
)

type MockExampleSentenceRepository struct {
	mock.Mock
}

func (m *MockExampleSentenceRepository) AddExampleSentence(ctx context.Context, translationID string, exampleSentence model.AddExampleSentenceInput) (*model.ExampleSentence, error) {

	return GetMockResult[*model.ExampleSentence](m.Called(ctx, translationID, exampleSentence))
}

func (m *MockExampleSentenceRepository) DeleteExampleSentence(ctx context.Context, id string) (*model.ExampleSentence, error) {

	return GetMockResult[*model.ExampleSentence](m.Called(ctx, id))
}

func (m *MockExampleSentenceRepository) UpdateExampleSentence(ctx context.Context, id string, edits model.EditExampleSentenceInput) (*model.ExampleSentence, error) {

	return GetMockResult[*model.ExampleSentence](m.Called(ctx, id, edits))
}

func (m *MockExampleSentenceRepository) GetSingleExampleSentence(ctx context.Context, id string) (*model.ExampleSentence, error) {

	return GetMockResult[*model.ExampleSentence](m.Called(ctx, id))
}

func (m *MockExampleSentenceRepository) GetExampleSentencesByTranslationId(ctx context.Context, translationID string) ([]*model.ExampleSentence, error) {

	return GetMockResult[[]*model.ExampleSentence](m.Called(ctx, translationID))
}
