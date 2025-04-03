package repository

import (
	"context"
	"sync"
	"testing"

	"github.com/grzegorzpapaj/graphql-dictionary-api/internal/database"
	"github.com/grzegorzpapaj/graphql-dictionary-api/internal/graph/model"
	"github.com/joho/godotenv"
)

func TestAddPolishWord_ConcurrentTranslations(t *testing.T) {

	err := godotenv.Load("../../.env")
	if err != nil {
		t.Fatal("Error loading .env file")
	}

	db, err := database.Connect()
	if err != nil {
		t.Fatal(err)
	}

	defer db.Close()

	translationRepo := &TranslationRepositoryDB{
		DB: db,
	}
	polishRepo := &PolishWordRepositoryDB{
		DB:              db,
		TranslationRepo: translationRepo,
	}

	word := "pisać"
	translations := []model.AddTranslationInput{
		{EnglishWord: "write"},
		{EnglishWord: "type"},
	}

	var wg sync.WaitGroup
	wg.Add(len(translations))

	for _, translation := range translations {
		go func(tInput model.AddTranslationInput) {
			defer wg.Done()
			input := model.AddPolishWordInput{
				Word:         word,
				Translations: []*model.AddTranslationInput{&tInput},
			}
			_, err := polishRepo.AddPolishWord(context.Background(), input)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		}(translation)
	}

	wg.Wait()

	pw, err := polishRepo.GetSinglePolishWord(context.Background(), nil, &word)
	if err != nil {
		t.Fatal(err)
	}

	if len(pw.Translations) != 2 {
		t.Fatalf("Expected 2 translations, got %d", len(pw.Translations))
	}

	found := map[string]bool{
		"write": false,
		"type":  false,
	}

	for _, trans := range pw.Translations {
		found[trans.EnglishWord] = true
	}

	for k, v := range found {
		if !v {
			t.Errorf("Translation %q not found", k)
		}
	}
}
