# graphql-dictionary-api
System collecting translations of Polish words into English in a relational database with GraphQL API. The API allows end users to manage translations - create new ones, receive, modify and remove existing ones. A user is able to send a Polish word and its translation to English together with exemplary sentences presenting usage of the word. Multiple translations of a single word are possible.

# Technologies
- **Programming Language**: Go
- **GraphQL Framework**: gqlgen
- **Database**: PostgreSQL
- **Containerization**: Docker, Docker Compose

# Data Model
[ERD diagram here]

Additionally, version columns were added to every table to provide optimistic concurrency control.

# Installation

## Prerequisites
- **Go**
- **Docker and Docker Compose**

## Clone the repository

```
git clone https://github.com/grzegorzpapaj/graphql-dictionary-api.git
cd graphql-dictionary-api
```

## Configuration
Store all sensitive configuration in an .env file. Below is an example configuration


```
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=dictionary-db
```

## Running the Application

```
docker-compose up -d
go run main.go
```

The API will be accessible at http://localhost:8080.

# GraphQL API Usage

The API exposes GraphQL endpoints for performing CRUD operations on database entries.

## Example Mutations and Queries

### Polish Words
Adding a polish word:

```graph
mutation addPolishWordMutation{
  addPolishWord(polishWord: {
    word:"przykład"
    translations:[
      {
        englishWord: "example"
        exampleSentences: [
          {
            sentencePl: "To zdanie podane jest jako przykład."
            sentenceEn: "This sentence is given as an example."
          }
          {
            sentencePl: "To jest bardzo dobry przykład."
            sentenceEn: "It is a very good example."
          }
        ]
      }
    ]
  }) {
    id
    word
  	translations {
      id
      englishWord
      exampleSentences {
        id
        sentencePl
        sentenceEn
      }
    }
  }
}
```

Retrieving all polish words:
```graph
query retrieveAllPolishWordsQuery {
	polishWords {
    id
    word
    version
    translations {
      id
      englishWord
      version
      exampleSentences {
        id
        sentencePl
        sentenceEn
        version
      }
    }
  }
}
```

Retrieving single polish word by word:
```graph
query retrieveSinglePolishWordByWord {
  polishWord(word:"przykład"){
    id
    word
    version
    translations {
      id
      englishWord
      version
      exampleSentences {
        id
        sentencePl
        sentenceEn
        version
      }
    }
  }
}
```
Alternatively, the user can retrieve a single polish word by ID.


Updating a polish word by word:
```graph
mutation updatePolishWordByWord {
  updatePolishWord(
    word: "przykład"
    edits: {
      version: 1
      translations: [
        {
          version: 1
          englishWord: "updated_example"
          exampleSentences: [
            {
              version: 1
              sentencePl: "Zaktualizowane zdanie PL"
              sentenceEn: "Updated sentence EN"
            }
          ]
        }
      ]
    }
  ) {
    id
    word
      	translations {
      id
      englishWord
      exampleSentences {
        id
        sentencePl
        sentenceEn
      }
    }
  }
}
```

Alternatively, the user can update polish words by ID.

Deleting a polish word:
```
mutation deletePolishWordByWord{
  deletePolishWord(word:"przykład") {
  	 id
    word
      	translations {
      id
      englishWord
      exampleSentences {
        id
        sentencePl
        sentenceEn
      }
    }
  }
}
```

Alternatively, the user can delete polish words by ID.


### Translations
Adding a translation by the word field of Polish Word:
```
mutation addTranslationByPolishwordWord{
  addTranslation(
    polishWord: "przykład"
    translation: {
      englishWord: "AddedTranslation"
      exampleSentences: {
        sentencePl: "Przykładowe zdanie z AddedTranslation"
        sentenceEn: "Example sentence with AddedTranslation"
      }
    }
  ) {
    id
    englishWord
    exampleSentences {
      sentencePl
      sentenceEn
    }
    polishWord {
      id
      word
    }
  }
}
```

Alternatively, the user can add translations by ID field of Polish Word.

Retrieving single translation by its ID:
```graph
query retrieveSingleTranslationByID{
  translation(id:"16") {
    id
    englishWord
    exampleSentences {
      sentencePl
      sentenceEn
    }
    polishWord {
      id
      word
    }
  }
}
```

Updating translation by ID:
```graph
query retrieveSingleTranslationByID{
  translation(id:"1") {
    id
    englishWord
    exampleSentences {
      sentencePl
      sentenceEn
    }
    polishWord {
      id
      word
    }
  }
}
```