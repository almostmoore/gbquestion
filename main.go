package main

import (
	"log"
	"os"

	"github.com/almostmoore/gbquestion/rest"
	"github.com/almostmoore/gbquestion/storage"
	"github.com/boltdb/bolt"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Couldn't load dotenv: %s", err)
	}

	db, err := bolt.Open(os.Getenv("DB_PATH"), 0600, nil)
	if err != nil {
		log.Fatalf("Couldn't load database file (%s): %s", os.Getenv("DB_PATH"), err)
	}
	defer db.Close()

	qs := storage.NewQuestionStorage(db)

	server := rest.NewServer(os.Getenv("HTTP_BIND"), qs)
	server.Run()
}
