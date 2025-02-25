package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/joho/godotenv"

	"github.com/grzegorzpapaj/graphql-dictionary-api/internal/database"
	"github.com/grzegorzpapaj/graphql-dictionary-api/internal/graph/generated"
	"github.com/grzegorzpapaj/graphql-dictionary-api/internal/graph/resolver"
	"github.com/grzegorzpapaj/graphql-dictionary-api/internal/repository"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := database.Connect()
	if err != nil {
		log.Fatalf("Could not connect to PostgreSQL: %v", err)
	}

	fmt.Println("Successfully connected to PostgreSQL")

	startServer(db)
}

func startServer(db *sql.DB) {
	port := os.Getenv("PORT")

	polishWordRepo := &repository.PolishWordRepositoryDB{DB: db}

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &resolver.Resolver{
		PolishWordRepo: polishWordRepo,
	}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
