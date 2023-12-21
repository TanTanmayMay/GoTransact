package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"rest1/internal/handler"
	"rest1/internal/repository"
	"rest1/internal/usecases"
	"go.uber.org/zap"
	// "rest1/internal/domain"
	// "rest1/internal/repository"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"os"
	"github.com/joho/godotenv"
)

func main() {

	// initialize zap
	var logger *zap.Logger
	var err error
    logger, err = zap.NewProduction()
    if err != nil {
        log.Fatalf("Failed to initialize Zap logger: %v", err)
    }

    defer logger.Sync()  //buffer
	err = godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		return
	}
	username :=  os.Getenv("DB_USER")//"nishant"
	password := os.Getenv("DB_PASSWORD")
	host := "db"
	port := "5432"
	database :=  os.Getenv("DB_NAME")
	// Connection string
	connString := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable", username, password, host, port, database)

	// Establish a connection to the PostgreSQL database
	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		logger.Fatal("Error connecting to PostgreSQL", zap.Error(err))
	}
	defer conn.Close(context.Background())
	// Check the connection
	err = conn.Ping(context.Background())
	if err != nil {
		logger.Panic("Connection not established!")
	}
	logger.Info("Connected to PostgreSQL")

	r := chi.NewRouter()

	// Initialize UseCase and Handler
	// func NewAccountRepo(conn *pgx.Conn)
	userRepo := repository.NewUserRepo(conn, logger)
	accountRepo := repository.NewAccountRepo(conn, logger)
	userUseCase := usecases.NewUserUseCase(userRepo, conn, logger)
	accountUseCase := usecases.NewAccountUseCase(accountRepo, conn, logger)
	userHandler := handler.NewUserHandler(userUseCase, conn, logger)
	accountHandler := handler.NewAccountHandler(accountUseCase, conn, logger)

	// Routes
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})
	r.Post("/user/register", userHandler.Register)
	r.Get("/user/{userid}" , userHandler.GetAccountById)
	r.Get("/users", userHandler.GetAllUsers)
	
	r.Post("/{userid}/account/create" , accountHandler.CreateAccountHandler)
	r.Get("/account/{accoundId}", accountHandler.GetByAccountNoHandler)


	logger.Info("Server listening on :8000")
	http.ListenAndServe(":8000", r)
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
