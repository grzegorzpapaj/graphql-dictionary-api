package repository

import (
	"context"

	"github.com/grzegorzpapaj/graphql-dictionary-api/internal/graph/model"
)

type PolishWordRepositoryInterface interface {
	AddPolishWord(ctx context.Context, polishWord model.AddPolishWordInput) (*model.PolishWord, error)
	GetAllPolishWords(ctx context.Context) ([]*model.PolishWord, error)
}
