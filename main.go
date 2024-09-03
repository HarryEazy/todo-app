package main

import (
	"log"
	"net/http"
	"todo_app/db"
	"todo_app/graph"
	"todo_app/graph/generated"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

// go get -u github.com/99designs/gqlgen
// go get -u github.com/mattn/go-sqlite3
// go get -u github.com/jmoiron/sqlx
func main() {
    // Initialize the database
    db.InitDB()

    port := "8080"
    srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))

    http.Handle("/", playground.Handler("GraphQL playground", "/query"))
    http.Handle("/query", srv)

    log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
    log.Fatal(http.ListenAndServe(":"+port, nil))
}
