package repository

import (
	"context"
	"fmt"
	"rest1/internal/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)


type AccountRepo struct {
	conn *pgxpool.Pool
	logger *zap.Logger
}

func NewAccountRepo(conn *pgxpool.Pool, logger *zap.Logger) *AccountRepo {
	return &AccountRepo{
		conn: conn,
		logger: logger,
	}
}

func (a *AccountRepo) DropAccountsTable() error {
	_, err := a.conn.Exec(context.Background(), "DROP TABLE accounts")

	if err != nil {
		a.logger.Error("Failed to drop account table", zap.Error(err))
		return err
	}

	return nil
}

func (a *AccountRepo) CreateTable() error {
	_, err := a.conn.Exec(context.Background(), "CREATE TABLE accounts (accountno VARCHAR(255) PRIMARY KEY, userid VARCHAR (255) , balance FLOAT, minBalance FLOAT , CONSTRAINT constrain_fk FOREIGN KEY (userid) REFERENCES users(userid) );")
	if err != nil {
		a.logger.Error("Failed to create table in Database", zap.Error(err))
		return err
	}
	return nil
	
}

func (a *AccountRepo) GetByNo(accountNo uuid.UUID) (*domain.Account, error) {
	var account domain.Account
	err := a.conn.QueryRow(context.Background(), "SELECT accountno, userid, balance, minBalance  FROM accounts WHERE accountno = $1", accountNo).Scan(&account.AccountNo, &account.UserID , &account.Balance, &account.MinBalance)

	if err != nil {
		a.logger.Error("Failed to get account by ID from Database", zap.Error(err))
		return nil, err
	}

	return &account, nil
}

func (a *AccountRepo) CreateAccount(account *domain.Account) (uuid.UUID, error) {
	fmt.Println("repository acc id", account.AccountNo)
	fmt.Println("repository userid ", account.UserID)
	fmt.Println("repository ", account.Balance)
	fmt.Println("repository ", account.MinBalance)
	_, err := a.conn.Exec(context.Background(), "INSERT INTO accounts (accountno, userid, balance, minBalance) VALUES($1, $2, $3 , $4)", account.AccountNo , account.UserID, account.Balance , account.MinBalance )
	if err != nil {
		a.logger.Error("Failed to create account in Database")
		return uuid.MustParse("00000000-0000-0000-0000-000000000000"), err // Return the error
	}
	return account.AccountNo, nil
}

func (a *AccountRepo) GetAll() ([]domain.Account, error) {
	rows, err := a.conn.Query(context.Background(), "SELECT accountno, userid, balance, minBalance FROM accounts")
	if err != nil {
		a.logger.Error("Failed to get all accounts from Database", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var accounts []domain.Account
	for rows.Next() {
		var account domain.Account
		if err := rows.Scan(&account.AccountNo, &account.Balance, &account.MinBalance , &account.UserID); err != nil {
			a.logger.Error("Failed to get account by ID from Database and append it to ds", zap.Error(err))
			return nil, err
		}
		accounts = append(accounts, account)
	}

	return accounts, nil
}


func (a* AccountRepo) GetAccByUserId(userid uuid.UUID) (*domain.Account , error) {
	var account = domain.Account{}
	fmt.Println("Repo userrrrid", userid)

	// idstr := userid.String()
	err := a.conn.QueryRow(context.Background(), "SELECT accountno, userid, balance, minBalance FROM accounts WHERE userid = $1", userid.String()).Scan(&account.AccountNo, &account.UserID, &account.Balance, &account.MinBalance )
	fmt.Println("Repo AccountNumber ", account.AccountNo)
	fmt.Println("Repo AccountUID ", account.UserID)
	fmt.Println("Repo AccountBalan ", account.Balance)
	fmt.Println("Repo AccMinBal ", account.MinBalance)
	if err != nil {
		a.logger.Error("Failed to get account by ID from Database", zap.Error(err))
		return nil, err
	}
	return &account, nil
}

