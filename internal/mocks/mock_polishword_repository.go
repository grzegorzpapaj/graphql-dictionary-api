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
	args := m.Called(ctx, polishWord)

	if result, ok := args.Get(0).(*model.PolishWord); ok {
		return result, args.Error(1)
	}

	return nil, args.Error(1)
}

func (m *MockPolishWordRepository) DeletePolishWord(ctx context.Context, id *string, word *string) (*model.PolishWord, error) {
	args := m.Called(ctx, id, word)

	if result, ok := args.Get(0).(*model.PolishWord); ok {
		return result, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockPolishWordRepository) UpdatePolishWord(ctx context.Context, id *string, word *string, edits *model.EditPolishWordInput) (*model.PolishWord, error) {
	args := m.Called(ctx, id, word, edits)

	if result, ok := args.Get(0).(*model.PolishWord); ok {
		return result, args.Error(1)
	}
	return nil, nil
}

func (m *MockPolishWordRepository) GetAllPolishWords(ctx context.Context) ([]*model.PolishWord, error) {
	args := m.Called(ctx)

	if result, ok := args.Get(0).([]*model.PolishWord); ok {
		return result, args.Error(1)
	}

	return nil, args.Error(1)
}
