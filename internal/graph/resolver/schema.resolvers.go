package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.66

import (
	"context"
	"fmt"

	"github.com/grzegorzpapaj/graphql-dictionary-api/internal/graph/generated"
	"github.com/grzegorzpapaj/graphql-dictionary-api/internal/graph/model"
)

// AddPolishWord is the resolver for the addPolishWord field.
func (r *mutationResolver) AddPolishWord(ctx context.Context, polishWord model.AddPolishWordInput) (*model.PolishWord, error) {
	return r.PolishWordRepo.AddPolishWord(ctx, polishWord)
}

// DeletePolishWord is the resolver for the deletePolishWord field.
func (r *mutationResolver) DeletePolishWord(ctx context.Context, id *string, word *string) (*model.PolishWord, error) {
	return r.PolishWordRepo.DeletePolishWord(ctx, id, word)
}

// UpdatePolishWord is the resolver for the updatePolishWord field.
func (r *mutationResolver) UpdatePolishWord(ctx context.Context, id *string, word *string, edits *model.EditPolishWordInput) (*model.PolishWord, error) {
	return r.PolishWordRepo.UpdatePolishWord(ctx, id, word, edits)
}

// AddTranslation is the resolver for the addTranslation field.
func (r *mutationResolver) AddTranslation(ctx context.Context, polishWordID *string, polishWord *string, translation *model.AddTranslationInput) (*model.Translation, error) {
	return r.TranslationRepo.AddTranslation(ctx, polishWordID, polishWord, translation)
}

// DeleteTranslation is the resolver for the deleteTranslation field.
func (r *mutationResolver) DeleteTranslation(ctx context.Context, id string) (*model.Translation, error) {
	return r.TranslationRepo.DeleteTranslation(ctx, id)
}

// UpdateTranslation is the resolver for the updateTranslation field.
func (r *mutationResolver) UpdateTranslation(ctx context.Context, id string, edits model.EditTranslationInput) (*model.Translation, error) {
	return r.TranslationRepo.UpdateTranslation(ctx, id, edits)
}

// AddExampleSentence is the resolver for the addExampleSentence field.
func (r *mutationResolver) AddExampleSentence(ctx context.Context, translationID string, exampleSentence model.AddExampleSentenceInput) (*model.ExampleSentence, error) {
	panic(fmt.Errorf("not implemented: AddExampleSentence - addExampleSentence"))
}

// DeleteExampleSentence is the resolver for the deleteExampleSentence field.
func (r *mutationResolver) DeleteExampleSentence(ctx context.Context, id string) (*model.ExampleSentence, error) {
	panic(fmt.Errorf("not implemented: DeleteExampleSentence - deleteExampleSentence"))
}

// UpdateExampleSentence is the resolver for the updateExampleSentence field.
func (r *mutationResolver) UpdateExampleSentence(ctx context.Context, id string, edits model.EditExampleSentenceInput) (*model.ExampleSentence, error) {
	panic(fmt.Errorf("not implemented: UpdateExampleSentence - updateExampleSentence"))
}

// PolishWord is the resolver for the polishWord field.
func (r *queryResolver) PolishWord(ctx context.Context, id *string, word *string) (*model.PolishWord, error) {
	panic(fmt.Errorf("not implemented: PolishWord - polishWord"))
}

// PolishWords is the resolver for the polishWords field.
func (r *queryResolver) PolishWords(ctx context.Context) ([]*model.PolishWord, error) {
	return r.PolishWordRepo.GetAllPolishWords(ctx)
}

// Translation is the resolver for the translation field.
func (r *queryResolver) Translation(ctx context.Context, id string) (*model.Translation, error) {
	return r.TranslationRepo.GetSingleTranslationByID(ctx, id)
}

// ExampleSentence is the resolver for the exampleSentence field.
func (r *queryResolver) ExampleSentence(ctx context.Context, id string) (*model.ExampleSentence, error) {
	panic(fmt.Errorf("not implemented: ExampleSentence - exampleSentence"))
}

// ExampleSentences is the resolver for the exampleSentences field.
func (r *queryResolver) ExampleSentences(ctx context.Context, translationID string) ([]*model.ExampleSentence, error) {
	panic(fmt.Errorf("not implemented: ExampleSentences - exampleSentences"))
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
