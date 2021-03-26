package db

import (
	"database/sql"
	"net"
	"time"

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

// Check if a lookup result already exists for an IP
func LookupResultFOrIPExists(ip net.IP) (bool, error) {
	query := `
	SELECT COUNT(*)
	FROM address_results
	WHERE ip_address = $1
	`

	rows, err := Database.Query(query, ip.String())
	if err != nil {
		return false, err
	}

	var count int
	rows.Next()
	err = rows.Scan(&count)

	return count >= 1, err
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
func UpdateIPLookupResult(ip net.IP, responseCode net.IP) error {
	query := `
	UPDATE address_results
	SET response_code = $2,
		updated_at = $3
	WHERE ip_address=$1
	`
	insertStatement, err := Database.Prepare(query)
	if err != nil {
		return err
	}

	_, err = insertStatement.Exec(ip, responseCode, time.Now().String())
	return err
}
