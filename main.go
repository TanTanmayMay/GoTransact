package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"log"
)

func main() {
	// Replace these values with your PostgreSQL details
	username := "nishant"
	password := "nishant"
	host := "localhost"
	port := "5432"
	database := "nishant"

	// Connection string
	connString := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable", username, password, host, port, database)

	// Establish a connection to the PostgreSQL database
	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close(context.Background())

	// Check the connection
	err = conn.Ping(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to PostgreSQL!")

	// Perform database operations here...

	// Example: Querying the version of the PostgreSQL server
	var version string
	err = conn.QueryRow(context.Background(), "SELECT version()").Scan(&version)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("PostgreSQL version:", version)
}