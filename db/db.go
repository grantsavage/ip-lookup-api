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

	// Check if not data was returned
	if !rows.Next() {
		return nil, errors.New("Could not find a lookup result for IP " + ip.String())
	}

	// Normalize returned row data intro IPLookupResult
	result := &model.IPLookupResult{}
	err = rows.Scan(&result.UUID, &result.IPAddress, &result.ResponseCode, &result.CreatedAt, &result.UpdatedAt)

	return result, err
}

// StoreIPLookupResult stores a new IP lookup result
func StoreIPLookupResult(result model.IPLookupResult) error {
	query := `
	INSERT INTO address_results 
	(
		uuid, 
		created_at, 
		updated_at, 
		response_code, 
		ip_address
	) 
	VALUES ($1, $2, $3, $4, $5)
	`
	insertStatement, err := Database.Prepare(query)
	if err != nil {
		return err
	}

	_, err = insertStatement.Exec(result.UUID, result.CreatedAt, result.UpdatedAt, result.ResponseCode, result.IPAddress)
	return err
}

// UpdateIPLookupResult updates an IP lookup result
func UpdateIPLookupResult(result model.IPLookupResult) error {
	query := `
	UPDATE address_results 
	(
		response_code,  
		updated_at
	) 
	WHERE ip_address=$1
	VALUES ($2, $3)
	`
	insertStatement, err := Database.Prepare(query)
	if err != nil {
		return err
	}

	_, err = insertStatement.Exec(result.ResponseCode, result.UpdatedAt)
	return err
}
