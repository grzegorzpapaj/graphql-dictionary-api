package repository

import (
	"context"
	"fmt"

	"github.com/grzegorzpapaj/graphql-dictionary-api/internal/graph/model"
)

func (pwr *PolishWordRepositoryDB) fetchPolishWords(ctx context.Context, id *string, word *string) (*model.PolishWord, error) {

	var fetchedPolishWord model.PolishWord
	if id != nil {
		err := pwr.DB.QueryRowContext(ctx, "SELECT id, word FROM polish_words WHERE id = $1",
			*id).Scan(&fetchedPolishWord.ID, &fetchedPolishWord.Word)
		if err != nil {
			return nil, err
		}
	} else if word != nil {
		err := pwr.DB.QueryRowContext(ctx, "SELECT id, word FROM polish_words WHERE word = $1",
			*word).Scan(&fetchedPolishWord.ID, &fetchedPolishWord.Word)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("either id or word must be provided")
	}

	return &fetchedPolishWord, nil
}
