package db_connection

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

var db *sql.DB

func OpenDB(dbPath string) (*sql.DB, error) {
	if db != nil {
		return db, nil
	}
	var err error
	db, err = sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to database: %v", err)
	}
	// verify connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("Failed to ping database: %v", err)
	}
	return db, nil
}

// shut down database connection
func CloseDB() error {
	if db != nil {
		err := db.Close()
		db = nil
		return err
	}
	return nil
}
