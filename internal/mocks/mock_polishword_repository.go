package mocks

import (
	"context"

	"github.com/grzegorzpapaj/graphql-dictionary-api/internal/graph/model"
	"github.com/stretchr/testify/mock"
)

type MockPolishWordRepository struct {
	mock.Mock
}

func (m *MockPolishWordRepository) AddPolishWord(ctx context.Context, polishWord model.AddPolishWordInput) (*model.PolishWord, error) {

	return GetMockResult[*model.PolishWord](m.Called(ctx, polishWord))
}

func (m *MockPolishWordRepository) DeletePolishWord(ctx context.Context, id *string, word *string) (*model.PolishWord, error) {

	return GetMockResult[*model.PolishWord](m.Called(ctx, id, word))
}

func (m *MockPolishWordRepository) UpdatePolishWord(ctx context.Context, id *string, word *string, edits *model.EditPolishWordInput) (*model.PolishWord, error) {

	return GetMockResult[*model.PolishWord](m.Called(ctx, id, word, edits))
}

func (m *MockPolishWordRepository) GetAllPolishWords(ctx context.Context) ([]*model.PolishWord, error) {

	return GetMockResult[[]*model.PolishWord](m.Called(ctx))
}

func (m *MockPolishWordRepository) GetSinglePolishWord(ctx context.Context, id *string, word *string) (*model.PolishWord, error) {

	return GetMockResult[*model.PolishWord](m.Called(ctx, id, word))
}
