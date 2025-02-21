package main

import (
	"os"
	"fmt"
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/joho/godotenv"

	"github.com/grzegorzpapaj/graphql-dictionary-api/internal/database"
	"github.com/grzegorzpapaj/graphql-dictionary-api/internal/graph/generated"
	"github.com/grzegorzpapaj/graphql-dictionary-api/internal/graph/resolver"
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

	word := "test"
	insertQuery := "INSERT INTO polish_words (word) VALUES ($1) RETURNING id"

	var id int
	err = db.QueryRow(insertQuery, word).Scan(&id)
	if err != nil {
		log.Fatalf("Error inserting into database: %v", err)
	}

	fmt.Printf("Inserted word '%s' with ID %d\n", word, id)

	startServer()
}

func startServer() {
	port := os.Getenv("PORT")

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &resolver.Resolver{}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)


	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
