package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"rest1/internal/repository"

	// "rest1/internal/domain"
	// "rest1/internal/repository"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5"
)

func main() {
	username := "nishant" // os.Env
	password := "nishant"
	host := "db"
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

	r := chi.NewRouter()

	// Initialize UseCase and Handler
	// func NewAccountRepo(conn *pgx.Conn)
	userRepo := repository.NewUserRepo(conn)
	accountRepo := repository.NewAccountRepo(conn)
	

	// Routes
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})
	r.Get("/employees", employeeHandler.GetEmployees)
	r.Get("/employees/{id}", employeeHandler.GetEmployeeByID)
	r.Post("/employees", employeeHandler.CreateEmployee)
	r.Put("/employees/{id}", employeeHandler.UpdateEmployee)
	r.Delete("/employees/{id}", employeeHandler.DeleteEmployee)

	http.ListenAndServe(":8080", r)
	// user1  := repository.NewUserRepo(conn)
	// user1.CreateTable()
	// _1 , err1 := conn.Exec(context.Background() , "CREATE TABLE users (id INT GENERATED ALWAYS AS IDENTITY, accountno INT, name VARCHAR ( 50 )  NOT NULL,password VARCHAR ( 50 ) NOT NULL, PRIMARY KEY(id), CONSTRAINT fk_account FOREIGN KEY(accountno) REFERENCES accounts(accountno));")
	// if(err1 != nil){
	// 	fmt.Println(err1)
	// 	fmt.Println(_1)

	// }
	// Perform database operations here...

	// Example: Querying the version of the PostgreSQL server
	// var version string
	// err = conn.QueryRow(context.Background(), "SELECT version()").Scan(&version)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	//  _ , err = conn.Exec(context.Background() , "CREATE TABLE users (id INT PRIMARY KEY,name VARCHAR ( 50 )  NOT NULL,password VARCHAR ( 50 ) NOT NULL,accountNo INT);")
	//  if err != nil{
	// 	fmt.Println(err)
	//  }
	// var user domain.User;
	// user.ID = 2003
	// user.AccountNo = 123
	// user.Name = "NMMMM"
	// user.Password = "abcd"
	// user1  := repository.NewUserRepo(conn)
	// user1.Create(&user)
	// // fmt.Println("PostgreSQL version:", version)
	// users, err := repository.NewUserRepo(conn).GetAll(conn)
	// 	if err != nil {
	// 		fmt.Println(err)
	// 		return
	// 	}
	// 	// Print the users to the response.
	// 	for _, user := range users {
	// 		fmt.Printf("ID: %d, Name: %s, AccountNo: %d, Password: %s\n", user.ID, user.Name, user.AccountNo, user.Password)
	// 	}

	// //8000 port in router

}
