package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/go-chi/chi"
	"github.com/grantsavage/ip-lookup-api/auth"
	"github.com/grantsavage/ip-lookup-api/db"
	"github.com/grantsavage/ip-lookup-api/graph"
	"github.com/grantsavage/ip-lookup-api/graph/generated"
)

// defaultPort is the default port to bind the server to
const defaultPort = "8080"

// main sets up the database and starts the GraphQL server
func main() {
	// Get and setup app configuration
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	// Open connection to the database
	database, err := db.Connect("./database.db")
	if err != nil {
		log.Fatal("error connecting to the database", err.Error())
	}
	defer database.Close()

	// Setup the database
	err = db.SetupDatabase(database)
	if err != nil {
		log.Fatal("error setting up the database", err.Error())
	}

	// Setup router and middleware
	router := chi.NewRouter()
	router.Use(auth.Middleware)

	// Create and setup new GraphQL server
	config := generated.Config{
		Resolvers: &graph.Resolver{
			Database: database,
		},
	}
	server := handler.NewDefaultServer(generated.NewExecutableSchema(config))

	// Handle panics
	server.SetRecoverFunc(func(ctx context.Context, err interface{}) error {
		log.Printf("internal server error %s", err)
		return errors.New("internal server error")
	})

	// Bind GraphQL server to /graphql route
	router.Handle("/graphql", server)

	// Start listening for requests
	log.Fatal(http.ListenAndServe(":"+port, router))
}
