package mocks

import (
	"context"

	"github.com/grzegorzpapaj/graphql-dictionary-api/internal/graph/model"
	"github.com/stretchr/testify/mock"
)

type MockTranslationRepository struct {
	mock.Mock
}

func (m *MockTranslationRepository) AddTranslation(ctx context.Context, polishWordID *string, polishWord *string, translation *model.AddTranslationInput) (*model.Translation, error) {

	args := m.Called(ctx, polishWordID, polishWord, translation)

	if result, ok := args.Get(0).(*model.Translation); ok {
		return result, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockTranslationRepository) DeleteTranslation(ctx context.Context, id string) (*model.Translation, error) {

	args := m.Called(ctx, id)

	if result, ok := args.Get(0).(*model.Translation); ok {
		return result, args.Error(1)
	}
	return nil, args.Error(1)
}
