package repository

import (
	"context"
	"errors"
	"fmt"
	"rest1/internal/domain"
	"rest1/internal/usecases"
	"sync"

	"github.com/google/uuid"
	// "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type AtomicUserRepository struct {
	conn   *pgxpool.Pool
	logger *zap.Logger
}

var _ usecases.AtomicUserRepository = (*AtomicUserRepository)(nil)

// NewAtomicUserRepo instantiates a new AtomicUserRepository using the pgxpool provided.
func NewAtomicUserRepo(conn *pgxpool.Pool, logger *zap.Logger) *AtomicUserRepository {
	return &AtomicUserRepository{
		conn:   conn,
		logger: logger,
	}
}

// UserRepository satisfies usecase.UserRepository and uses *pgxpool.Pool directly.
type UserRepository struct {
	conn   *pgxpool.Pool
	logger *zap.Logger
}

var _ usecases.UserRepository = (*UserRepository)(nil)

func (ar *UserRepository) DropUserTable() error {
	_, err := ar.conn.Exec(context.Background(), "DROP TABLE users;")

	if err != nil {
		ar.logger.Error("Failed to drop user table", zap.Error(err))
		return err
	}

	return nil
}

func (ar *UserRepository) CreateUserTable() error {
	_, err := ar.conn.Exec(context.Background(), "CREATE TABLE users (userid varchar(255), name VARCHAR ( 50 )  NOT NULL,password VARCHAR ( 50 ) NOT NULL, PRIMARY KEY(userid));")

	if err != nil {
		ar.logger.Error("Failed to create table in Database", zap.Error(err))
		return err
	}

	return nil
}

func (ar *UserRepository) GetIndividual(idd uuid.UUID, wg *sync.WaitGroup, userChan chan<- domain.User) {
	defer wg.Done()
	var user domain.User
	err := ar.conn.QueryRow(context.Background(), "SELECT userid, name, password FROM users WHERE userid = $1", idd).Scan(&user.ID, &user.Name, &user.Password)

	if err != nil {
		ar.logger.Error("Failed to get account by ID from Database", zap.Error(err))

		return
	}

	userChan <- user
}

func (ar *UserRepository) GetAll() ([]domain.User, error) {
	rows, err := ar.conn.Query(context.Background(), "SELECT userid FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var wg sync.WaitGroup
	userChan := make(chan domain.User, 20)

	var users []domain.User
	for rows.Next() {
		var uidd uuid.UUID
		err := rows.Scan(&uidd)
		if err != nil {
			ar.logger.Error("Failed to scan user ID", zap.Error(err))
			return nil, err
		}

		wg.Add(1)
		go ar.GetIndividual(uidd, &wg, userChan)
	}

	go func() {
		wg.Wait()
		close(userChan)
	}()

	for user := range userChan {
		users = append(users, user)
	}

	return users, nil
}

func (ar *UserRepository) GetByID(id uuid.UUID) (*domain.User, error) {
	var user domain.User
	err := ar.conn.QueryRow(context.Background(), "SELECT userid, name, password FROM users WHERE userid = $1", id).Scan(&user.ID, &user.Name, &user.Password)

	if err != nil {
		ar.logger.Error("Failed to get account by ID from Database", zap.Error(err))

		return nil, err
	}

	return &user, nil
}

func (ar *UserRepository) CreateUser(user *domain.User) error {

	// Added a Check on Password before Creating the User
	lenght := len(user.Password)
	if lenght < 5 {
		return errors.New("Length of Password is Short")
	}

	_, err := ar.conn.Exec(context.Background(), "INSERT INTO users(userid , name, password) VALUES($1, $2, $3)", user.ID, user.Name, user.Password)
	if err != nil {
		return err
	}
	fmt.Println("Added User to Database!!")
	return nil
}

func (ar *UserRepository) Withdraw(user *domain.User, amount int) error {
	var account domain.Account
	err := ar.conn.QueryRow(context.Background(), "SELECT accountno, userid, balance, minbalance FROM accounts WHERE userid = $1", user.ID).Scan(&account.AccountNo, &account.UserID, &account.Balance, &account.MinBalance)

	// First check if the given withdraw is possible or not
	if account.Balance-float64(amount) < account.MinBalance {
		return errors.New("Can't Withdraw as terms are violated")
	}
	qry := "UPDATE accounts SET accounts.balance = (accounts.balance - $1) WHERE accounts.userid = $2"
	_, err = ar.conn.Exec(context.Background(), qry, amount, account.UserID)
	if err != nil {
		ar.logger.Error("Failed to withdraw from account", zap.Error(err))
		return err
	}
	return nil
}

func (ar *UserRepository) Deposit(user *domain.User, amount int) error {
	/*
		UPDATE product
		SET net_price = price - price * discount
		FROM product_segment
		WHERE product.segment_id = product_segment.id;
	*/
	var account domain.Account
	err := ar.conn.QueryRow(context.Background(), "SELECT accountno, userid, balance, minbalance FROM accounts WHERE userid = $1", user.ID).Scan(&account.AccountNo, &account.UserID, &account.Balance, &account.MinBalance)

	qry := "UPDATE accounts SET accounts.balance = (accounts.balance + $1) WHERE accounts.userid = $2"
	_, err = ar.conn.Exec(context.Background(), qry, amount, account.UserID)
	if err != nil {
		ar.logger.Error("Failed to deposit into account", zap.Error(err))
		return err
	}
	return nil
}

func (ar *AtomicUserRepository) Execute(op usecases.AtomicUserOperation) error {
	tx, err := ar.conn.Begin(context.Background())
	if err != nil {
		return err
	}

	defer func() { _ = tx.Rollback(context.Background()) }()

	// Create a new single-use AccountRepository backed by the transaction.
	userRepoWithTransaction := UserRepository{
		conn:   ar.conn,
		logger: ar.logger,
	}

	// Perform the AtomicAccountOperation using the AccountRepository.
	if err := op(&userRepoWithTransaction); err != nil {
		return err
	}

	if err := tx.Commit(context.Background()); err != nil {
		return err
	}

	return nil
}
