package repository

import (
	"context"
	"fmt"
	"log"
	"os"
	"rest1/internal/domain"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

/*
	to Run all the UNIT tests of your project do following command

	> go test ./...
*/

func TestGetByID(t *testing.T) {
	var logger *zap.Logger
	var err error
	logger, err = zap.NewProduction()
	if err != nil {
		log.Fatalf("Failed to initialize Zap logger: %v", err)
	}
	defer logger.Sync()
	err = godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		return
	}
	username := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := "db"
	port := "5432"
	database := os.Getenv("DB_NAME")
	connString := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable", username, password, host, port, database)
	conn, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		logger.Fatal("Error connecting to PostgreSQL", zap.Error(err))
	}
	defer conn.Close()
	err = conn.Ping(context.Background())
	if err != nil {
		logger.Panic("Connection not established!")
	}
	logger.Info("Connected to PostgreSQL")

	t.Run("subtest1", func(t *testing.T) {
		// Testing Begins
		r := NewUserRepo(conn, logger)

		dummyuser := domain.User{
			ID:       uuid.New(),
			Name:     "Tanmay",
			Password: "abcdefg",
		}

		err = r.CreateUser(&dummyuser)
		assert.NoError(t, err)
		assert.Equal(t, "Tanmay", dummyuser.Name)
	})

	t.Run("subtest2", func(t *testing.T) {
		// Testing Begins
		r := NewUserRepo(conn, logger)

		dummyuser := domain.User{
			ID:       uuid.New(),
			Name:     "Om",
			Password: "abc",
		}

		err = r.CreateUser(&dummyuser)
		assert.NoError(t, err) //Length of Password is too small Error is caught here
		assert.Equal(t, "Tanmay", dummyuser.Name)
	})
}
