package repository

import (
	"context"

	"github.com/grzegorzpapaj/graphql-dictionary-api/internal/graph/model"
)

type ExampleSentenceRepositoryInterface interface {
	AddExampleSentence(ctx context.Context, translationID string, exampleSentence model.AddExampleSentenceInput) (*model.ExampleSentence, error)
	DeleteExampleSentence(ctx context.Context, id string) (*model.ExampleSentence, error)
	UpdateExampleSentence(ctx context.Context, id string, edits model.EditExampleSentenceInput) (*model.ExampleSentence, error)
	GetSingleExampleSentence(ctx context.Context, id string) (*model.ExampleSentence, error)
}
