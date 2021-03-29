package db

import (
	"database/sql"
	"errors"
	"net"

	"github.com/grantsavage/ip-lookup-api/graph/model"
	_ "github.com/mattn/go-sqlite3"
)

var ErrorNotFound error = errors.New("could not find a result for IP")

// Connect opens the connection to the database
func Connect(datastore string) *sql.DB {
	// Initialize SQLite database
	db, err := sql.Open("sqlite3", datastore)
	if err != nil {
		panic(err.Error())
	}
	return db
}

// SetupDatabase creates the required tables for the application
func SetupDatabase(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS address_results 
	(
		uuid TEXT UNIQUE, 
		response_code TEXT, 
		ip_address TEXT UNIQUE PRIMARY KEY,
		created_at TEXT, 
		updated_at TEXT
	)
	`
	sqlStatement, err := db.Prepare(query)
	if err != nil {
		return err
	}

	_, err = sqlStatement.Exec()
	return err
}

// GeIPLookupResult gets an IP lookup result
func GetIPLookupResult(db *sql.DB, ip net.IP) (*model.IPLookupResult, error) {
	query := `
	SELECT uuid, ip_address, response_code, created_at, updated_at 
	FROM address_results
	WHERE ip_address = $1
	LIMIT 1
	`
	rows, err := db.Query(query, ip.String())
	if err != nil {
		return nil, err
	}

	// Normalize returned row data intro IPLookupResult
	result := &model.IPLookupResult{}
	rows.Next()
	err = rows.Scan(&result.UUID, &result.IPAddress, &result.ResponseCode, &result.CreatedAt, &result.UpdatedAt)

	// Check if no IPLookupResult was found
	if result == nil {
		return nil, ErrorNotFound
	}

	return result, err
}

// UpsertIPLookupResult upserts an IPLookupResult
func UpsertIPLookupResult(db *sql.DB, result model.IPLookupResult) error {
	// This will first try to insert a result, but if a conflict occurs, this is most likely
	// because a record for the IP already exists, so instead we update the response_code and
	// updated_at time
	query := `
	INSERT INTO address_results (uuid, ip_address, response_code, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5)
	ON CONFLICT(ip_address) DO UPDATE SET response_code = $3, updated_at = $5
	WHERE ip_address = $2;
	`
	upsertStatement, err := db.Prepare(query)
	if err != nil {
		return err
	}

	_, err = upsertStatement.Exec(result.UUID, result.IPAddress, result.ResponseCode, result.CreatedAt, result.UpdatedAt)
	return err
}
