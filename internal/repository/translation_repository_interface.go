package repository

import (
	"context"

	"github.com/grzegorzpapaj/graphql-dictionary-api/internal/graph/model"
)

type TranslationRepositoryInterface interface {
	AddTranslation(ctx context.Context, polishWordID *string, polishWord *string, translation *model.AddTranslationInput) (*model.Translation, error)
	DeleteTranslation(ctx context.Context, id string) (*model.Translation, error)
	UpdateTranslation(ctx context.Context, id string, edits model.EditTranslationInput) (*model.Translation, error)
}
