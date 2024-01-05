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
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
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

// r.Get("/get/{userid}", userHandler.GetUserById)

// middleware for validating uuid we get from URL parameters.
func validateUUIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		uuidFromParam := chi.URLParam(r, "userid")

		_, err := uuid.Parse(uuidFromParam)

		if err != nil {
			http.Error(w, "Invalid UUID", http.StatusBadRequest)
			return
		}

		next.ServeHTTP(w, r)
	})
}

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

	// Initialize UseCase, Handler, Repository(Short Lived)
	atomicUserRepo := repository.NewAtomicUserRepo(conn, logger)
	atomicAccountRepo := repository.NewAtomicAccountRepo(conn, logger)
	userUseCase := usecases.NewUserUseCase(atomicUserRepo, logger)
	accountUseCase := usecases.NewAccountUsecase(atomicAccountRepo, logger)
	userHandler := handler.NewUserHandler(userUseCase, logger)
	accountHandler := handler.NewAccountHandler(accountUseCase, logger)

	// Routes
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	// user routes  -> Grouped Routing
	r.Route("/user", func(r chi.Router) {
		r.Post("/register", userHandler.Register)
		r.With(validateUUIDMiddleware).Get("/get/{userid}", userHandler.GetUserById)
		r.Get("/getall", userHandler.GetAllUsers)
	})

	r.With(validateUUIDMiddleware).Post("/withdraw/{userid}/amount/{amount}", userHandler.WithdrawHandler)
	r.With(validateUUIDMiddleware).Post("/deposit/{userid}/amount/{amount}", userHandler.DepositHandler)

	// account routes -> Grouped Routing
	r.Route("/account", func(r chi.Router) {
		r.With(validateUUIDMiddleware).Post("/create/{userid}", accountHandler.CreateAccountHandler)
		r.Get("/get/{accoundId}", accountHandler.GetByAccountNoHandler)
	})

	// Utility Routes
	r.Get("/drop/account/table", accountHandler.DropAccountsTableHandler)
	r.Get("/create/account/table", accountHandler.CreateAccountTableHandler)
	r.Get("/create/users/table", userHandler.CreateUsersTableHandler)
	r.Get("/drop/users/table", userHandler.DropUserTableHandler)

	port1, port2 := ":8000", ":8001"

	var waitgroup sync.WaitGroup

	waitgroup.Add(1)
	go serverStarter(port1, r, &waitgroup, logger)
	waitgroup.Add(1)
	go serverStarter(port2, r, &waitgroup, logger)

	waitgroup.Wait()
}

func serverStarter(portt string, routerr *chi.Mux, wg *sync.WaitGroup, logger *zap.Logger) {
	defer wg.Done()
	logger.Info("Server listening...") //, portt ????
	http.ListenAndServe(portt, routerr)
}
