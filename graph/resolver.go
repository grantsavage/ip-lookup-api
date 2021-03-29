package graph

import "database/sql"

//go:generate go run github.com/99designs/gqlgen

type Resolver struct {
	Database *sql.DB
}
