package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var Database *sql.DB

// Connect opens the connection to the database
func Connect(datastore string) {
	// Initialize database
	db, err := sql.Open("sqlite3", datastore)
	if err != nil {
		panic(err.Error())
	}
	Database = db
}

// SetupDatabase sets creates the required tables for the application
func SetupDatabase() error {
	query := `
	CREATE TABLE IF NOT EXISTS address_results 
	(
		uuid TEXT UNIQUE PRIMARY KEY, 
		response_code TEXT, 
		ip_address TEXT,
		created_at TEXT, 
		updated_at TEXT
	)
	`
	sqlStatement, err := Database.Prepare(query)
	if err != nil {
		return err
	}

	_, err = sqlStatement.Exec()
	return err
}
