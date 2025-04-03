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

func TestAddWithExampleSentence_Concurrent(t *testing.T) {

	err := godotenv.Load("../../.env")
	if err != nil {
		t.Fatal("Error loading .env file")
	}

	db, err := database.Connect()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	exampleSentenceRepo := &ExampleSentenceRepositoryDB{DB: db}

	translationRepo := &TranslationRepositoryDB{
		ExampleSentenceRepo: exampleSentenceRepo,
		DB:                  db,
	}

	polishRepo := &PolishWordRepositoryDB{
		DB:              db,
		TranslationRepo: translationRepo,
	}

	word := "pisać"

	_, err = polishRepo.AddPolishWord(context.Background(), model.AddPolishWordInput{
		Word:         word,
		Translations: nil,
	})
	if err != nil {
		t.Fatalf("Error setting up Polish word: %v", err)
	}

	translationInput := model.AddTranslationInput{
		EnglishWord: "write",
	}
	trans, err := translationRepo.AddTranslation(context.Background(), nil, &word, &translationInput)
	if err != nil {
		t.Fatalf("Error adding translation: %v", err)
	}

	exampleSentences := []model.AddExampleSentenceInput{
		{
			SentencePl: "Pisze list",
			SentenceEn: "I am writing a letter",
		},
		{
			SentencePl: "Pisze e-mail",
			SentenceEn: "I am writing an email",
		},
	}

	var wg sync.WaitGroup
	wg.Add(len(exampleSentences))

	for _, esInput := range exampleSentences {
		go func(input model.AddExampleSentenceInput) {
			defer wg.Done()
			_, err := exampleSentenceRepo.AddExampleSentence(context.Background(), trans.ID, input)
			if err != nil {
				t.Errorf("Unexpected error adding example sentence: %v", err)
			}
		}(esInput)
	}

	wg.Wait()

	pw, err := polishRepo.GetSinglePolishWord(context.Background(), nil, &word)
	if err != nil {
		t.Fatal(err)
	}

	var foundTrans *model.Translation
	for _, t := range pw.Translations {
		if t.EnglishWord == "write" {
			foundTrans = t
			break
		}
	}
	if foundTrans == nil {
		t.Fatalf("Translation 'write' not found")
	}

	if len(foundTrans.ExampleSentences) != len(exampleSentences) {
		t.Fatalf("Expected %d example sentences, got %d", len(exampleSentences), len(foundTrans.ExampleSentences))
	}

	expectedSentences := map[string]string{
		"Pisze list":   "I am writing a letter",
		"Pisze e-mail": "I am writing an email",
	}

	for _, es := range foundTrans.ExampleSentences {
		expectedEn, ok := expectedSentences[es.SentencePl]
		if !ok {
			t.Errorf("Unexpected Polish sentence: %q", es.SentencePl)
		} else if es.SentenceEn != expectedEn {
			t.Errorf("For Polish sentence %q, expected English %q, got %q", es.SentencePl, expectedEn, es.SentenceEn)
		}
	}
}
