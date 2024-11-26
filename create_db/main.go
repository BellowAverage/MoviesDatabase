package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

	_ "modernc.org/sqlite"
)

func createTables(db *sql.DB) error {
	// SQL statements to create tables
	createActorsTable := `
		CREATE TABLE IF NOT EXISTS actors (
			actor_id INTEGER PRIMARY KEY,
			actor_name TEXT NOT NULL
		);`
	createDirectorsTable := `
		CREATE TABLE IF NOT EXISTS directors (
			director_id INTEGER PRIMARY KEY,
			director_name TEXT NOT NULL
		);`
	createDirectorsGenresTable := `
		CREATE TABLE IF NOT EXISTS directors_genres (
			director_id INTEGER,
			genre TEXT NOT NULL,
			PRIMARY KEY (director_id, genre)
		);`
	createMoviesTable := `
		CREATE TABLE IF NOT EXISTS movies (
			movie_id INTEGER PRIMARY KEY,
			movie_name TEXT NOT NULL,
			movie_year INTEGER,
			movie_rank REAL
		);`
	createMoviesGenresTable := `
		CREATE TABLE IF NOT EXISTS movies_genres (
			movie_id INTEGER,
			genre TEXT NOT NULL,
			PRIMARY KEY (movie_id, genre)
		);`
	createRolesTable := `
		CREATE TABLE IF NOT EXISTS roles (
			role_id INTEGER PRIMARY KEY,
			actor_id INTEGER,
			movie_id INTEGER,
			role_name TEXT NOT NULL
		);`

	// Execute table creation queries
	tables := []string{createActorsTable, createDirectorsTable, createDirectorsGenresTable, createMoviesTable, createMoviesGenresTable, createRolesTable}
	for _, query := range tables {
		_, err := db.Exec(query)
		if err != nil {
			return fmt.Errorf("failed to create table: %v", err)
		}
	}
	return nil
}

func importCSV(db *sql.DB, filePath string, query string, columns int) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %v", filePath, err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to read file %s: %v", filePath, err)
	}

	// Skip header row
	for _, record := range records[1:] {
		if len(record) < columns {
			// Ignore rows with fewer fields than expected
			continue
		}

		values := make([]interface{}, columns)
		for i := 0; i < columns; i++ {
			values[i] = record[i]
			if i == 0 { // Convert first column to integer
				values[i], _ = strconv.Atoi(record[i])
			}
		}

		_, err := db.Exec(query, values...)
		if err != nil {
			// Ignore insert errors and move to the next row
			continue
		}
	}

	return nil
}

func main() {
	// Open database connection
	db, err := sql.Open("sqlite", "../movie.db")
	if err != nil {
		fmt.Println("Failed to connect to database:", err)
		return
	}
	defer db.Close()

	// Create tables
	if err := createTables(db); err != nil {
		fmt.Println("Failed to create tables:", err)
		return
	}

	// Import CSV data
	files := []struct {
		path    string
		query   string
		columns int
	}{
		{"../IMDB-actors.csv", "INSERT INTO actors (actor_id, actor_name) VALUES (?, ?)", 2},
		{"../IMDB-directors.csv", "INSERT INTO directors (director_id, director_name) VALUES (?, ?)", 2},
		{"../IMDB-directors_genres.csv", "INSERT INTO directors_genres (director_id, genre) VALUES (?, ?)", 2},
		{"../IMDB-movies.csv", "INSERT INTO movies (movie_id, movie_name, movie_year, movie_rank) VALUES (?, ?, ?, ?)", 4},
		{"../IMDB-movies_genres.csv", "INSERT INTO movies_genres (movie_id, genre) VALUES (?, ?)", 2},
		{"../IMDB-roles.csv", "INSERT INTO roles (role_id, actor_id, movie_id, role_name) VALUES (?, ?, ?, ?)", 4},
	}

	for _, file := range files {
		if err := importCSV(db, file.path, file.query, file.columns); err != nil {
			fmt.Printf("Failed to import %s: %v\n", file.path, err)
		}
	}

	fmt.Println("Database and data import completed successfully.")
}
