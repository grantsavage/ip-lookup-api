package graph

import "database/sql"

//go:generate go run github.com/99designs/gqlgen

type Resolver struct {
	// Database holds a pointer to the database connection
	Database *sql.DB
}
