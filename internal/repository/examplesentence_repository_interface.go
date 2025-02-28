package repository

import (
	"context"

	"github.com/grzegorzpapaj/graphql-dictionary-api/internal/graph/model"
)

type ExampleSentenceRepositoryInterface interface {
	AddExampleSentence(ctx context.Context, translationID string, exampleSentence model.AddExampleSentenceInput) (*model.ExampleSentence, error)
}
