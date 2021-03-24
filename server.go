package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/grantsavage/ip-lookup-api/db"
	"github.com/grantsavage/ip-lookup-api/graph"
	"github.com/grantsavage/ip-lookup-api/graph/generated"
)

const defaultPort = "8080"

func main() {
	// Get and setup app configuration
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	// Open connection to the database
	db.Connect("./database.db")
	defer db.Database.Close()

	// Setup the database
	err := db.SetupDatabase()
	if err != nil {
		log.Fatal("error setting up the database: " + err.Error())
	}

	// Create and setup new GraphQL server
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	// Start listening for requests
	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
