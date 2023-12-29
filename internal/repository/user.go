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
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

var (
	// ErrTransactionInProgress is returned when more than one concurrent
	// transaction is attempted.
	ErrTransactionInProgress = errors.New("transaction already in progress")

	// ErrTransactionNotStarted is returned when an atomic operation is
	// requested but the repository hasn't begun a transaction.
	ErrTransactionNotStarted = errors.New("transaction not started")
)

type transactor struct {
	conn   *pgxpool.Pool
	tx     *pgx.Tx
	logger *zap.Logger
}

var _ usecases.Transactor = (*transactor)(nil)

func (t *transactor) Begin() error {
	if t.tx != nil {
		return ErrTransactionInProgress
	}

	tx, err := t.conn.Begin(context.Background())
	if err != nil {
		t.logger.Error("Failed to Begin Transaction", zap.Error(err))
		return err
	}
	t.tx = &tx

	return nil
}

func (t *transactor) Commit() error {
	tx := *t.tx
	if err := tx.Commit(context.Background()); err != nil {
		t.logger.Error("Failed to Commit Transaction", zap.Error(err))
		return err
	}
	return nil
}

func (t *transactor) Rollback() error {
	tx := *t.tx
	if err := tx.Rollback(context.Background()); err != nil {
		t.logger.Error("Failed to Rollback Transaction", zap.Error(err))
		return err
	}
	return nil
}

type AtomicUserRepo struct {
	*transactor
}

var _ usecases.AtomicUserRepository = (*AtomicUserRepo)(nil)

func NewAtomicUserRepoFactory(conn *pgxpool.Pool, logger *zap.Logger) usecases.AtomicUserRepositoryFactory {
	return func() usecases.AtomicUserRepository {
		return &AtomicUserRepo{
			transactor: &transactor{
				conn:   conn,
				logger: logger,
			},
		}
	}
}

func (u *AtomicUserRepo) DropUserTable() error {

	if u.tx == nil {
		return ErrTransactionNotStarted
	}

	_, err := u.conn.Exec(context.Background(), "DROP TABLE users;")
	if err != nil {
		u.logger.Error("Failed to create table in Database", zap.Error(err))
		fmt.Println(err)
	}

	return nil
}
func (u *AtomicUserRepo) CreateUserTable() error {
	if u.tx == nil {
		return ErrTransactionNotStarted
	}

	_, err := u.conn.Exec(context.Background(), "CREATE TABLE users (userid varchar(255), name VARCHAR ( 50 )  NOT NULL,password VARCHAR ( 50 ) NOT NULL, PRIMARY KEY(userid));")
	if err != nil {
		u.logger.Error("Failed to create table in Database", zap.Error(err))
		fmt.Println(err)
	}

	return nil
}

func (u *AtomicUserRepo) GetIndividual(idd uuid.UUID, wg *sync.WaitGroup, userChan chan<- domain.User) {

	defer wg.Done()
	var user domain.User
	err := u.conn.QueryRow(context.Background(), "SELECT userid, name, password FROM users WHERE userid = $1", idd).Scan(&user.ID, &user.Name, &user.Password)

	if err != nil {
		u.logger.Error("Failed to get account by ID from Database", zap.Error(err))

		return
	}

	userChan <- user
}

func (u *AtomicUserRepo) GetAll() ([]domain.User, error) {
	if u.tx == nil {
		return nil, ErrTransactionNotStarted
	}

	rows, err := u.conn.Query(context.Background(), "SELECT userid FROM users")
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
			u.logger.Error("Failed to scan user ID", zap.Error(err))
			return nil, err
		}

		wg.Add(1)
		go u.GetIndividual(uidd, &wg, userChan)
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

func (u *AtomicUserRepo) GetByID(id uuid.UUID) (*domain.User, error) {
	if u.tx == nil {
		return nil, ErrTransactionNotStarted
	}

	var user domain.User
	err := u.conn.QueryRow(context.Background(), "SELECT userid, name, password FROM users WHERE userid = $1", id).Scan(&user.ID, &user.Name, &user.Password)

	if err != nil {
		u.logger.Error("Failed to get account by ID from Database", zap.Error(err))

		return nil, err
	}

	return &user, nil
}

func (u *AtomicUserRepo) CreateUser(user *domain.User) error {
	if u.tx == nil {
		return ErrTransactionNotStarted
	}

	// Added a Check on Password before Creating the User
	lenght := len(user.Password)
	if lenght < 5 {
		return errors.New("Length of Password is Short")
	}

	_, err := u.conn.Exec(context.Background(), "INSERT INTO users(userid , name, password) VALUES($1, $2, $3)", user.ID, user.Name, user.Password)
	if err != nil {
		return err
	}
	fmt.Println("Added User to Database!!")
	return nil
}

func (u *AtomicUserRepo) Withdraw(account *domain.Account, amount int) error {
	if u.tx == nil {
		return ErrTransactionNotStarted
	}

	qry := "UPDATE accounts SET accounts.balance = (accounts.balance - $1) WHERE accounts.userid = $2"
	_, err := u.conn.Exec(context.Background(), qry, amount, account.UserID)
	if err != nil {
		u.logger.Error("Failed to withdraw from account", zap.Error(err))
		return err
	}
	return nil
}

func (u *AtomicUserRepo) Deposit(account *domain.Account, amount int) error {
	if u.tx == nil {
		return ErrTransactionNotStarted
	}

	/*
		UPDATE product
		SET net_price = price - price * discount
		FROM product_segment
		WHERE product.segment_id = product_segment.id;
	*/
	qry := "UPDATE accounts SET accounts.balance = (accounts.balance + $1) WHERE userid = $2"
	_, err := u.conn.Exec(context.Background(), qry, amount, account.UserID)
	if err != nil {
		u.logger.Error("Failed to Deposit in Account", zap.Error(err))
		return err
	}
	return nil
}
