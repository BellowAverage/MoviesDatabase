package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"

	_ "modernc.org/sqlite"
)

// ImportCSV imports data from a CSV file into the database using a transformation function.
func importCSV(db *sql.DB, filePath string, query string, columns int, transform func([]string) ([]interface{}, error)) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %v", filePath, err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.LazyQuotes = true       // Enable LazyQuotes
	reader.FieldsPerRecord = -1    // Allow variable number of fields
	reader.TrimLeadingSpace = true // Trim leading space

	// Skip header row
	_, err = reader.Read()
	if err != nil {
		return fmt.Errorf("failed to read header: %v", err)
	}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			// Log error and continue to the next record
			fmt.Printf("Error reading record: %v\n", err)
			continue
		}

		if len(record) < columns {
			// Ignore rows with fewer fields than expected
			continue
		}

		values, err := transform(record)
		if err != nil {
			// Ignore transformation errors and move to the next row
			continue
		}

		_, err = db.Exec(query, values...)
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

	// Create tables (same as your original code)
	if err := createTables(db); err != nil {
		fmt.Println("Failed to create tables:", err)
		return
	}

	// Import CSV data with appropriate transformations
	if err := importCSV(db, "../IMDB-actors.csv", "INSERT INTO actors (actor_id, actor_name) VALUES (?, ?)", 4, transformActors); err != nil {
		fmt.Printf("Failed to import actors: %v\n", err)
	}

	if err := importCSV(db, "../IMDB-directors.csv", "INSERT INTO directors (director_id, director_name) VALUES (?, ?)", 3, transformDirectors); err != nil {
		fmt.Printf("Failed to import directors: %v\n", err)
	}

	if err := importCSV(db, "../IMDB-directors_genres.csv", "INSERT INTO directors_genres (director_id, genre) VALUES (?, ?)", 3, transformDirectorsGenres); err != nil {
		fmt.Printf("Failed to import directors_genres: %v\n", err)
	}

	if err := importCSV(db, "../IMDB-movies.csv", "INSERT INTO movies (movie_id, movie_name, movie_year, movie_rank) VALUES (?, ?, ?, ?)", 4, transformMovies); err != nil {
		fmt.Printf("Failed to import movies: %v\n", err)
	}

	if err := importCSV(db, "../IMDB-movies_genres.csv", "INSERT INTO movies_genres (movie_id, genre) VALUES (?, ?)", 2, transformMoviesGenres); err != nil {
		fmt.Printf("Failed to import movies_genres: %v\n", err)
	}

	if err := importCSV(db, "../IMDB-roles.csv", "INSERT INTO roles (actor_id, movie_id, role_name) VALUES (?, ?, ?)", 3, transformRoles); err != nil {
		fmt.Printf("Failed to import roles: %v\n", err)
	}

	fmt.Println("Database and data import completed successfully.")
}

// Transformation functions for each CSV file

func transformActors(record []string) ([]interface{}, error) {
	id, err := strconv.Atoi(record[0])
	if err != nil {
		return nil, err
	}
	actorName := record[1] + " " + record[2]
	return []interface{}{id, actorName}, nil
}

func transformDirectors(record []string) ([]interface{}, error) {
	id, err := strconv.Atoi(record[0])
	if err != nil {
		return nil, err
	}
	directorName := record[1] + " " + record[2]
	return []interface{}{id, directorName}, nil
}

func transformDirectorsGenres(record []string) ([]interface{}, error) {
	id, err := strconv.Atoi(record[0])
	if err != nil {
		return nil, err
	}
	genre := record[1]
	return []interface{}{id, genre}, nil
}

func transformMovies(record []string) ([]interface{}, error) {
	id, err := strconv.Atoi(record[0])
	if err != nil {
		return nil, err
	}
	name := record[1]

	// Handle 'NULL' and empty year
	var year interface{}
	if record[2] == "NULL" || record[2] == "" {
		year = nil
	} else {
		year, err = strconv.Atoi(record[2])
		if err != nil {
			return nil, err
		}
	}

	// Handle 'NULL' and empty rank
	var rank interface{}
	if record[3] == "NULL" || record[3] == "" {
		rank = nil
	} else {
		rank, err = strconv.ParseFloat(record[3], 64)
		if err != nil {
			return nil, err
		}
	}

	return []interface{}{id, name, year, rank}, nil
}

func transformMoviesGenres(record []string) ([]interface{}, error) {
	id, err := strconv.Atoi(record[0])
	if err != nil {
		return nil, err
	}
	genre := record[1]
	return []interface{}{id, genre}, nil
}

func transformRoles(record []string) ([]interface{}, error) {
	if len(record) < 3 {
		return nil, fmt.Errorf("record too short")
	}
	actorID, err := strconv.Atoi(record[0])
	if err != nil {
		return nil, err
	}
	movieID, err := strconv.Atoi(record[1])
	if err != nil {
		return nil, err
	}
	roleName := record[2]
	return []interface{}{actorID, movieID, roleName}, nil
}

// createTables function remains the same
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
			actor_id INTEGER,
			movie_id INTEGER,
			role_name TEXT NOT NULL,
			PRIMARY KEY (actor_id, movie_id)
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
