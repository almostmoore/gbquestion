package main

import (
	"log"
	"os"

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

	qs := NewQuestionStorage(db)

	server := NewServer(os.Getenv("HTTP_BIND"), qs)
	server.Run()
}
