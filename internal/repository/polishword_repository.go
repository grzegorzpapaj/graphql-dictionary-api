package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/grzegorzpapaj/graphql-dictionary-api/internal/graph/model"
)

type PolishWordRepositoryDB struct {
	DB *sql.DB
}

func (pwr *PolishWordRepositoryDB) AddPolishWord(ctx context.Context, polishWord model.AddPolishWordInput) (*model.PolishWord, error) {

	newPolishWord := &model.PolishWord{
		Word:         polishWord.Word,
		Translations: []*model.Translation{},
	}

	err := pwr.DB.QueryRowContext(ctx, "INSERT INTO polish_words (word) VALUES ($1) RETURNING id", newPolishWord.Word).Scan(&newPolishWord.ID)
	if err != nil {
		return nil, err
	}

	for _, t := range polishWord.Translations {
		newTranslation := &model.Translation{
			EnglishWord:      t.EnglishWord,
			ExampleSentences: []*model.ExampleSentence{},
		}

		err = pwr.DB.QueryRowContext(ctx, "INSERT INTO translations(english_word, polish_word_id) VALUES ($1, $2) RETURNING id",
			newTranslation.EnglishWord, newPolishWord.ID).Scan(&newTranslation.ID)

		if err != nil {
			return nil, err
		}

		for _, es := range t.ExampleSentences {
			newExampleSentence := &model.ExampleSentence{
				SentencePl: es.SentencePl,
				SentenceEn: es.SentenceEn,
			}

			err = pwr.DB.QueryRowContext(ctx, "INSERT INTO example_sentences(sentence_pl, sentence_en, translation_id) VALUES ($1, $2, $3) RETURNING id",
				newExampleSentence.SentencePl, newExampleSentence.SentenceEn, newTranslation.ID).Scan(&newExampleSentence.ID)

			if err != nil {
				return nil, err
			}

			newTranslation.ExampleSentences = append(newTranslation.ExampleSentences, newExampleSentence)
		}

		newPolishWord.Translations = append(newPolishWord.Translations, newTranslation)
	}

	return newPolishWord, nil
}

func (pwr *PolishWordRepositoryDB) DeletePolishWord(ctx context.Context, id *string, word *string) (*model.PolishWord, error) {
	var deletedPolishWord model.PolishWord

	if id != nil {
		err := pwr.DB.QueryRowContext(ctx, "DELETE FROM polish_words WHERE id = $1 RETURNING id, word",
			*id).Scan(&deletedPolishWord.ID, &deletedPolishWord.Word)

		if err != nil {
			return nil, err
		}
	} else if word != nil {
		err := pwr.DB.QueryRowContext(ctx, "DELETE FROM polish_words WHERE word = $1 RETUNING id, word",
			*word).Scan(&deletedPolishWord.ID, &deletedPolishWord.Word)

		if err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("either id or word must be provided")
	}

	return &deletedPolishWord, nil
}

func (pwr *PolishWordRepositoryDB) UpdatePolishWord(ctx context.Context, id *string, word *string, edits *model.EditPolishWordInput) (*model.PolishWord, error) {

	polishWordToEdit, err := pwr.fetchPolishWords(ctx, id, word)
	if err != nil {
		return nil, err
	}

	if word == nil && edits.Word != nil {
		_, err := pwr.DB.ExecContext(ctx,
			"UPDATE polish_words SET word = $1 WHERE id = $2",
			*edits.Word, polishWordToEdit.ID)

		if err != nil {
			return nil, err
		}
		polishWordToEdit.Word = *edits.Word
	}

	if edits.Translations != nil {

		rows, err := pwr.DB.QueryContext(ctx,
			"SELECT id, english_word FROM translations WHERE polish_word_id = $1 ORDER BY id",
			polishWordToEdit.ID)

		if err != nil {
			return nil, err
		}
		defer rows.Close()

		var existingTranslations []*model.Translation

		for rows.Next() {
			var t model.Translation
			if err := rows.Scan(&t.ID, &t.EnglishWord); err != nil {
				return nil, err
			}
			existingTranslations = append(existingTranslations, &t)
		}

		for i, editTr := range edits.Translations {
			if i < len(existingTranslations) {
				if editTr.EnglishWord != nil {
					_, err := pwr.DB.ExecContext(ctx,
						"UPDATE translations SET english_word = $1 WHERE id = $2",
						*editTr.EnglishWord, existingTranslations[i].ID)

					if err != nil {
						return nil, err
					}
					existingTranslations[i].EnglishWord = *editTr.EnglishWord
				}

				if editTr.ExampleSentences != nil {
					esRows, err := pwr.DB.QueryContext(ctx,
						"SELECT id, sentence_pl, sentence_en FROM example_sentences WHERE translation_id = $1 ORDER BY id",
						existingTranslations[i].ID)

					if err != nil {
						return nil, err
					}

					var existingExamples []*model.ExampleSentence
					for esRows.Next() {
						var es model.ExampleSentence
						if err := esRows.Scan(&es.ID, &es.SentencePl, &es.SentenceEn); err != nil {
							esRows.Close()
							return nil, err
						}

						existingExamples = append(existingExamples, &es)
					}
					esRows.Close()

					for j, editES := range editTr.ExampleSentences {
						if j < len(existingExamples) {

							sentencePl := existingExamples[j].SentencePl
							sentenceEn := existingExamples[j].SentenceEn
							if editES.SentencePl != nil {
								sentencePl = *editES.SentencePl
							}
							if editES.SentenceEn != nil {
								sentenceEn = *editES.SentenceEn
							}

							_, err := pwr.DB.ExecContext(ctx,
								"UPDATE example_sentences SET sentence_pl = $1, sentence_en = $2 WHERE id = $3",
								sentencePl, sentenceEn, existingExamples[j].ID)
							if err != nil {
								return nil, err
							}

							existingExamples[j].SentencePl = sentencePl
							existingExamples[j].SentenceEn = sentenceEn
						} else {

							var newExampleSentenceID string
							sentencePl := ""
							sentenceEn := ""

							if editES.SentencePl != nil {
								sentencePl = *editES.SentencePl
							}

							if editES.SentenceEn != nil {
								sentenceEn = *editES.SentenceEn
							}

							err := pwr.DB.QueryRowContext(ctx, "INSERT INTO example_sentences(sentence_pl, sentence_en, translation_id) VALUES ($1, $2, $3) RETURNING id",
								sentencePl, sentenceEn, existingTranslations[i].ID).Scan(&newExampleSentenceID)

							if err != nil {
								return nil, err
							}

							newExample := &model.ExampleSentence{
								ID:         newExampleSentenceID,
								SentencePl: sentencePl,
								SentenceEn: sentenceEn,
							}

							existingExamples = append(existingExamples, newExample)

						}
					}

					existingTranslations[i].ExampleSentences = existingExamples
				}
			}
		}

		if len(edits.Translations) > len(existingTranslations) {

			for i := len(existingTranslations); i < len(edits.Translations); i++ {

				editTr := edits.Translations[i]
				if editTr.EnglishWord != nil {
					var newTranslationID string
					err := pwr.DB.QueryRowContext(ctx, "INSERT INTO translations(english_word, polish_word_id) VALUES ($1, $2) RETURNING id",
						*editTr.EnglishWord, polishWordToEdit.ID).Scan(&newTranslationID)

					if err != nil {
						return nil, err
					}

					newTranslation := &model.Translation{
						ID:               newTranslationID,
						EnglishWord:      *editTr.EnglishWord,
						ExampleSentences: []*model.ExampleSentence{},
					}

					if editTr.ExampleSentences != nil {

						for _, editEs := range editTr.ExampleSentences {
							var newExampleSentenceID string
							sentencePl := ""
							sentenceEn := ""

							if editEs.SentencePl != nil {
								sentencePl = *editEs.SentencePl
							}

							if editEs.SentenceEn != nil {
								sentenceEn = *editEs.SentenceEn
							}

							err := pwr.DB.QueryRowContext(ctx, "INSERT INTO example_sentences(sentence_pl, sentence_en, translation_id) VALUES ($1, $2, $3) RETURNING id",
								sentencePl, sentenceEn, newTranslationID).Scan(&newExampleSentenceID)

							if err != nil {
								return nil, err
							}

							newExample := &model.ExampleSentence{
								ID:         newExampleSentenceID,
								SentencePl: sentencePl,
								SentenceEn: sentenceEn,
							}

							newTranslation.ExampleSentences = append(newTranslation.ExampleSentences, newExample)
						}
					}
					existingTranslations = append(existingTranslations, newTranslation)
				}
			}
		}

		polishWordToEdit.Translations = existingTranslations
	}

	return polishWordToEdit, nil

}

func (pwr *PolishWordRepositoryDB) GetAllPolishWords(ctx context.Context) ([]*model.PolishWord, error) {
	rows, err := pwr.DB.QueryContext(ctx, "SELECT id, word FROM polish_words")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var polishWords []*model.PolishWord
	for rows.Next() {
		var pw model.PolishWord

		if err := rows.Scan(&pw.ID, &pw.Word); err != nil {
			return nil, err
		}
		polishWords = append(polishWords, &pw)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return polishWords, nil
}
