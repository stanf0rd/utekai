package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq" // Postgres DB driver
)

var db *sql.DB

func init() {
	host := os.Getenv("PG_APP_HOST")
	user := os.Getenv("PG_APP_USER")
	password := os.Getenv("PG_APP_PASSWORD")
	dbname := os.Getenv("PG_APP_DB")

	connStr := fmt.Sprintf(`
		host=%v
		user=%v
		password=%v
		dbname=%v
		sslmode=disable
	`, host, user, password, dbname)

	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Successfully connected to database")
	}

	err = initQuestions()
	if err != nil {
		log.Fatalf("Cannot add initial questions: %v", err)
	}
}
