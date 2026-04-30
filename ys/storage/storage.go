package storage

import (
	"database/sql"
	"log"

	_ "github.com/glebarez/go-sqlite"
)

var db *sql.DB

// InitDB initializes the database and creates the necessary tables.
func InitDB() {
	var err error
	db, err = sql.Open("sqlite", "./storage/urlshortener.db")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	createTableSQL := `CREATE TABLE IF NOT EXISTS urls (
		"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		"short_url" TEXT NOT NULL UNIQUE,
		"long_url" TEXT NOT NULL
	);`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}
}

// SaveURL saves a short URL and its corresponding long URL to the database.
func SaveURL(shortURL, longURL string) error {
	insertSQL := `INSERT INTO urls(short_url, long_url) VALUES (?, ?)`
	_, err := db.Exec(insertSQL, shortURL, longURL)
	return err
}

// GetURL retrieves the long URL for a given short URL from the database.
func GetURL(shortURL string) (string, error) {
	selectSQL := `SELECT long_url FROM urls WHERE short_url = ?`
	var longURL string
	err := db.QueryRow(selectSQL, shortURL).Scan(&longURL)
	return longURL, err
}
