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

	return GetMockResult[*model.Translation](m.Called(ctx, polishWordID, polishWord, translation))
}

func (m *MockTranslationRepository) DeleteTranslation(ctx context.Context, id string) (*model.Translation, error) {

	return GetMockResult[*model.Translation](m.Called(ctx, id))
}

func (m *MockTranslationRepository) UpdateTranslation(ctx context.Context, id string, edits model.EditTranslationInput) (*model.Translation, error) {

	return GetMockResult[*model.Translation](m.Called(ctx, id, edits))
}

func (m *MockTranslationRepository) GetSingleTranslationByID(ctx context.Context, id string) (*model.Translation, error) {

	return GetMockResult[*model.Translation](m.Called(ctx, id))
}
