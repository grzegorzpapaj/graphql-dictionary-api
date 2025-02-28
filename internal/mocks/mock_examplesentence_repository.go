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
