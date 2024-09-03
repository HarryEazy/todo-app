package db

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3" // Import the SQLite3 driver
)

// Define the schema for the tasks table
var schema = `
CREATE TABLE IF NOT EXISTS tasks (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT,
    description TEXT,
    status TEXT,
    due_date TEXT 
);`

// DB is the global database connection pool
var DB *sqlx.DB

// InitDB initializes the database and creates the tasks table if it doesn't exist
func InitDB() {
    var err error
    // Use an in-memory SQLite database for testing purposes
    DB, err = sqlx.Connect("sqlite3", ":memory:")
    if err != nil {
        log.Fatalln(err)
    }

    DB.MustExec(schema)
}
