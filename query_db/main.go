package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

func writeToFile(filePath string, data []string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	for _, line := range data {
		_, err := file.WriteString(line + "\n")
		if err != nil {
			return fmt.Errorf("failed to write to file: %v", err)
		}
	}

	return nil
}

func queryTopMoviesByGenre(db *sql.DB, genre string, outputFile string) {
	query := `
		SELECT m.movie_name, m.movie_year, m.movie_rank
		FROM movies m
		JOIN movies_genres mg ON m.movie_id = mg.movie_id
		WHERE mg.genre = ?
		ORDER BY m.movie_rank DESC
		LIMIT 10;
	`

	rows, err := db.Query(query, genre)
	if err != nil {
		fmt.Println("Failed to execute query:", err)
		return
	}
	defer rows.Close()

	var results []string
	results = append(results, fmt.Sprintf("Top 10 Movies in Genre '%s':", genre))

	for rows.Next() {
		var name string
		var year int
		var rank float64
		if err := rows.Scan(&name, &year, &rank); err != nil {
			fmt.Println("Failed to scan row:", err)
			continue
		}
		line := fmt.Sprintf("Name: %s, Year: %d, Rank: %.2f", name, year, rank)
		results = append(results, line)
	}

	if err := writeToFile(outputFile, results); err != nil {
		fmt.Println("Failed to write query results to file:", err)
	} else {
		fmt.Printf("Query results saved to %s\n", outputFile)
	}
}

func queryDirectorsWithGenres(db *sql.DB, outputFile string) {
	query := `
		SELECT d.director_name, dg.genre
		FROM directors d
		JOIN directors_genres dg ON d.director_id = dg.director_id
		ORDER BY d.director_name;
	`

	rows, err := db.Query(query)
	if err != nil {
		fmt.Println("Failed to execute query:", err)
		return
	}
	defer rows.Close()

	var results []string
	results = append(results, "Directors and their Genres:")

	for rows.Next() {
		var name, genre string
		if err := rows.Scan(&name, &genre); err != nil {
			fmt.Println("Failed to scan row:", err)
			continue
		}
		line := fmt.Sprintf("Director: %s, Genre: %s", name, genre)
		results = append(results, line)
	}

	if err := writeToFile(outputFile, results); err != nil {
		fmt.Println("Failed to write query results to file:", err)
	} else {
		fmt.Printf("Query results saved to %s\n", outputFile)
	}
}

func queryActorsWithRoles(db *sql.DB, outputFile string) {
	query := `
		SELECT a.actor_name, r.role_name, m.movie_name
		FROM actors a
		JOIN roles r ON a.actor_id = r.actor_id
		JOIN movies m ON r.movie_id = m.movie_id
		ORDER BY a.actor_name;
	`

	rows, err := db.Query(query)
	if err != nil {
		fmt.Println("Failed to execute query:", err)
		return
	}
	defer rows.Close()

	var results []string
	results = append(results, "Actors and their Roles:")

	for rows.Next() {
		var actor, role, movie string
		if err := rows.Scan(&actor, &role, &movie); err != nil {
			fmt.Println("Failed to scan row:", err)
			continue
		}
		line := fmt.Sprintf("Actor: %s, Role: %s, Movie: %s", actor, role, movie)
		results = append(results, line)
	}

	if err := writeToFile(outputFile, results); err != nil {
		fmt.Println("Failed to write query results to file:", err)
	} else {
		fmt.Printf("Query results saved to %s\n", outputFile)
	}
}

func main() {
	// Open database connection
	db, err := sql.Open("sqlite", "../movie.db")
	if err != nil {
		fmt.Println("Failed to connect to database:", err)
		return
	}
	defer db.Close()

	// Perform queries and save results to files
	queryTopMoviesByGenre(db, "Action", "top_movies_action.txt")
	queryDirectorsWithGenres(db, "directors_genres.txt")
	queryActorsWithRoles(db, "actors_roles.txt")
}
