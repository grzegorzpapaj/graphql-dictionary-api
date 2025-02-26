package repository

import (
	"context"
	"fmt"

	"github.com/grzegorzpapaj/graphql-dictionary-api/internal/graph/model"
)

func (tr *TranslationRepositoryDB) getTargetPolishWordID(ctx context.Context, polishWordID *string, polishWord *string) (*string, error) {

	var targetPolishWordID string

	if polishWordID != nil {
		targetPolishWordID = *polishWordID
	} else if polishWord != nil {
		err := tr.DB.QueryRowContext(ctx, "SELECT id FROM polish_words WHERE word = $1", *polishWord).Scan(&targetPolishWordID)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("either polishWordID or polishWord word must be provided")
	}

	return &targetPolishWordID, nil
}

func (tr *TranslationRepositoryDB) prepareWordWithId(ctx context.Context, targetPolishWordID *string) (*model.PolishWord, error) {

	var word string
	err := tr.DB.QueryRowContext(ctx, "SELECT word FROM polish_words WHERE id = $1", *targetPolishWordID).Scan(&word)

	if err != nil {
		return nil, fmt.Errorf("failed to fetch polish word for id %s: %w", *targetPolishWordID, err)
	}

	return &model.PolishWord{
		ID:   *targetPolishWordID,
		Word: word,
	}, nil
}
