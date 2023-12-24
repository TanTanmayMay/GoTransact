package repository

import (
	"context"
	"errors"
	"fmt"
	"rest1/internal/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type UserRepo struct {
	conn   *pgxpool.Pool
	logger *zap.Logger
}

func NewUserRepo(conn *pgxpool.Pool, logger *zap.Logger) *UserRepo {
	return &UserRepo{
		conn:   conn,
		logger: logger,
	}
}

func (u *UserRepo) DropUserTable() error {
	_, err := u.conn.Exec(context.Background(), "DROP TABLE users;")
	if err != nil {
		u.logger.Error("Failed to create table in Database", zap.Error(err))
		fmt.Println(err)
	}
	return nil
}
func (u *UserRepo) CreateUserTable() error {
	_, err := u.conn.Exec(context.Background(), "CREATE TABLE users (userid varchar(255), name VARCHAR ( 50 )  NOT NULL,password VARCHAR ( 50 ) NOT NULL, PRIMARY KEY(userid));")
	if err != nil {
		u.logger.Error("Failed to create table in Database", zap.Error(err))
		fmt.Println(err)
	}
	return nil
}

func (u *UserRepo) GetAll() ([]domain.User, error) {
	rows, err := u.conn.Query(context.Background(), "SELECT userid, name, password FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var user domain.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Password); err != nil {
			u.logger.Error("Failed to get all accounts from Database", zap.Error(err))
			return nil, err

		}
		users = append(users, user)
	}

	return users, nil
}

func (u *UserRepo) GetByID(id uuid.UUID) (*domain.User, error) {
	var user domain.User
	err := u.conn.QueryRow(context.Background(), "SELECT userid, name, password FROM users WHERE userid = $1", id).Scan(&user.ID, &user.Name, &user.Password)

	if err != nil {
		u.logger.Error("Failed to get account by ID from Database", zap.Error(err))

		return nil, err
	}

	return &user, nil
}

func (u *UserRepo) CreateUser(user *domain.User) error {

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

func (u *UserRepo) Withdraw(account *domain.Account, amount int) error {
	qry := "UPDATE accounts SET accounts.balance = (accounts.balance - $1) WHERE accounts.userid = $2"
	_, err := u.conn.Exec(context.Background(), qry, amount, account.UserID)
	if err != nil {
		u.logger.Error("Failed to withdraw from account", zap.Error(err))
		return err
	}
	return nil
}

func (u *UserRepo) Deposit(account *domain.Account, amount int) error {
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
