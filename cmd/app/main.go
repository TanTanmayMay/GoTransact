package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"rest1/internal/repository"

	"github.com/jackc/pgx/v4"
)

func main() {
	// Replace these values with your PostgreSQL details
	username := "nishant"
	password := "nishant"
	host := "localhost"
	port := "5432"
	database := "nishant"
	fmt.Println("Connecting to .....")
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
	 _ , err = conn.Exec(context.Background() , "CREATE TABLE users (id serial PRIMARY KEY,name VARCHAR ( 50 ) UNIQUE NOT NULL,password VARCHAR ( 50 ) NOT NULL,accountNo INT);")
	 err = conn.QueryRow(context.Background(), "INSERT INTO users(id , name, accountNo, password) VALUES($1, $2, $3 , $4)", 124, "Om", 123 , "123").Scan("123")
	if(err != nil) {
		fmt.Println(err)

	}
	fmt.Println("PostgreSQL version:", version)
	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		users, err := repository.NewUserRepo(conn).GetAll()
		if err != nil {
			http.Error(w, fmt.Sprintf("Error getting users: %s", err), http.StatusInternalServerError)
			return
		}

		// Print the users to the response.
		for _, user := range users {
			fmt.Fprintf(w, "ID: %s, Name: %s, AccountNo: %s, Password: %s\n", user.ID, user.Name, user.AccountNo, user.Password)
		}
	})
}