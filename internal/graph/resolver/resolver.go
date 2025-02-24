package resolver

import (
	"database/sql"

	"github.com/grzegorzpapaj/graphql-dictionary-api/internal/repository"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	DB *sql.DB
	PolishWordRepo repository.PolishWordRepository
}
