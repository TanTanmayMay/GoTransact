package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"rest1/internal/handler"
	"rest1/internal/repository"
	"rest1/internal/usecases"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

/*
create a map (int -> uuid)
and user id to API will be int but internally we'll pass uuid

Update:
Done directly by using UUID and storing it in Database as VARCHAR(255)
*/
func main() {

	// initialize zap
	var logger *zap.Logger
	var err error
	logger, err = zap.NewProduction()
	if err != nil {
		log.Fatalf("Failed to initialize Zap logger: %v", err)
	}

	defer logger.Sync() //buffer
	err = godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		return
	}
	username := os.Getenv("DB_USER") //"nishant"
	password := os.Getenv("DB_PASSWORD")
	host := "db"
	port := "5432"
	database := os.Getenv("DB_NAME")
	// Connection string
	connString := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable", username, password, host, port, database)

	// Establish a connection to the PostgreSQL database

	/*
		Before
		conn, err := pgx.Connect(context.Background(), "your-database-connection-string")

		After
		pool, err := pgxpool.Connect(context.Background(), "your-database-connection-string")

	*/
	// conn, err := pgxpool.Connect(context.Background(), connString)
	conn, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		logger.Fatal("Error connecting to PostgreSQL", zap.Error(err))
	}
	defer conn.Close()
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

	// user routes
	r.Post("/user/register", userHandler.Register)
	r.Get("/user/get/{userid}", userHandler.GetUserById)
	r.Get("/user/getall", userHandler.GetAllUsers)
	// functionality routes
	r.Put("/withdraw/{userid}/amount/{amount}", userHandler.WithdrawHandler) // TODO
	r.Put("/deposit/{userid}/amount/{amount}", userHandler.DepositHandler)   //TODO
	// account routes
	r.Post("/account/create/{userid}", accountHandler.CreateAccountHandler)
	r.Get("/account/get/{accoundId}", accountHandler.GetByAccountNoHandler)
	// Utility Routes
	r.Get("/drop/account/table", accountHandler.DropAccountsTableHandler)
	r.Get("/create/account/table", accountHandler.CreateAccountTableHandler)
	r.Get("/create/users/table", userHandler.CreateUsersTableHandler)
	r.Get("/drop/users/table", userHandler.DropUserTableHandler)

	logger.Info("Server listening on :8000")
	http.ListenAndServe(":8000", r)

}
