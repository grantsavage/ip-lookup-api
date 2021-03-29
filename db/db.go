package db

import (
	"database/sql"
	"errors"
	"net"

	"github.com/grantsavage/ip-lookup-api/graph/model"
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
		uuid TEXT UNIQUE, 
		response_code TEXT, 
		ip_address TEXT UNIQUE PRIMARY KEY,
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

// GeIPLookupResult gets an IP lookup result
func GetIPLookupResult(ip net.IP) (*model.IPLookupResult, error) {
	query := `
	SELECT uuid, ip_address, response_code, created_at, updated_at 
	FROM address_results
	WHERE ip_address = $1
	LIMIT 1
	`
	rows, err := Database.Query(query, ip.String())
	if err != nil {
		return nil, err
	}

	// Normalize returned row data intro IPLookupResult
	result := &model.IPLookupResult{}
	rows.Next()
	err = rows.Scan(&result.UUID, &result.IPAddress, &result.ResponseCode, &result.CreatedAt, &result.UpdatedAt)

	if result == nil {
		return nil, errors.New("could not find a result for IP " + ip.String())
	}

	return result, err
}

// StoreIPLookupResult stores a new IP lookup result
func UpsertIPLookupResult(result model.IPLookupResult) error {
	query := `
	INSERT INTO address_results (uuid, ip_address, response_code, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5)
	ON CONFLICT(ip_address) DO UPDATE SET response_code = $3, updated_at = $5
	WHERE ip_address = $2;
	`
	insertStatement, err := Database.Prepare(query)
	if err != nil {
		return err
	}

	_, err = insertStatement.Exec(result.UUID, result.IPAddress, result.ResponseCode, result.CreatedAt, result.UpdatedAt)
	return err
}
