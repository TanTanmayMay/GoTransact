package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
	"rest1/internal/domain"
)

type AccountRepo struct {
	conn *pgx.Conn
	logger *zap.Logger
}

func NewAccountRepo(conn *pgx.Conn, logger *zap.Logger) *AccountRepo {
	return &AccountRepo{
		conn: conn,
		logger: logger,
	}
}

func (a *AccountRepo) CreateTable() error {
	_, err := a.conn.Exec(context.Background(), "CREATE TABLE accounts (accountno INT PRIMARY KEY, balance FLOAT, minBalance FLOAT);")
	if err != nil {
		a.logger.Error("Failed to create table in Database", zap.Error(err))
		fmt.Println(err)
	}
	return nil
}

func (a *AccountRepo) GetByNo(accountNo int) (*domain.Account, error) {
	var account domain.Account
	err := a.conn.QueryRow(context.Background(), "SELECT accountNo, balance, minBalance FROM accounts WHERE accountNo = $1", accountNo).
		Scan(&account.AccountNo, &account.Balance, &account.MinBalance)

	if err != nil {
		a.logger.Error("Failed to get account by ID from Database", zap.Error(err))
		return nil, err
	}

	return &account, nil
}

func (a *AccountRepo) CreateAccount(account *domain.Account) (int ,error) {
	var accId int
	err := a.conn.QueryRow(context.Background(), "INSERT INTO accounts(accountNo, balance, minBalance) VALUES($1, $2, $3)", account.AccountNo, account.Balance, account.MinBalance).Scan(&accId)
	if err != nil {
		a.logger.Error("Failed to create account in Database", zap.Error(err))
		fmt.Println(err)
	}
	return accId, nil
}

func (a *AccountRepo) GetAll() ([]domain.Account, error) {
	rows, err := a.conn.Query(context.Background(), "SELECT accountNo, balance, minBalance FROM accounts")
	if err != nil {
		a.logger.Error("Failed to get all accounts from Database", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var accounts []domain.Account
	for rows.Next() {
		var account domain.Account
		if err := rows.Scan(&account.AccountNo, &account.Balance, &account.MinBalance); err != nil {
			a.logger.Error("Failed to get account by ID from Database and append it to ds", zap.Error(err))
			return nil, err
		}
		accounts = append(accounts, account)
	}

	return accounts, nil
}


