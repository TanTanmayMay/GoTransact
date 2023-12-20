package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"rest1/internal/domain"
)

type AccountRepo struct {
	conn *pgx.Conn
}

func NewAccountRepo(conn *pgx.Conn) *AccountRepo {
	return &AccountRepo{conn: conn}
}

func (a *AccountRepo) CreateTable() error {
	_, err := a.conn.Exec(context.Background(), "CREATE TABLE accounts (accountno INT PRIMARY KEY, balance FLOAT, minBalance FLOAT);")
	if err != nil {
		fmt.Println(err)
	}
	return nil
}

func (a *AccountRepo) GetByNo(accountNo int) (*domain.Account, error) {
	var account domain.Account
	err := a.conn.QueryRow(context.Background(), "SELECT accountNo, balance, minBalance FROM accounts WHERE accountNo = $1", accountNo).
		Scan(&account.AccountNo, &account.Balance, &account.MinBalance)

	if err != nil {
		return nil, err
	}

	return &account, nil
}

func (a *AccountRepo) CreateAccount(account *domain.Account) error {
	_, err := a.conn.Exec(context.Background(), "INSERT INTO accounts(accountNo, balance, minBalance) VALUES($1, $2, $3)", account.AccountNo, account.Balance, account.MinBalance)
	if err != nil {
		fmt.Println(err)
	}
	return nil
}

func (a *AccountRepo) GetAll() ([]domain.Account, error) {
	rows, err := a.conn.Query(context.Background(), "SELECT accountNo, balance, minBalance FROM accounts")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []domain.Account
	for rows.Next() {
		var account domain.Account
		if err := rows.Scan(&account.AccountNo, &account.Balance, &account.MinBalance); err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}

	return accounts, nil
}


