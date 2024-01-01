package repository

import (
	"context"
	"rest1/internal/domain"
	"rest1/internal/usecases"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

// AtomicAccountRepository satisfies class.AtomicRepository and is capable
// of beginning transactions.
type AtomicAccountRepository struct {
	conn   *pgxpool.Pool
	logger *zap.Logger
}

var _ usecases.AtomicAccountRepository = (*AtomicAccountRepository)(nil)

// NewAtomicAccountRepo instantiates a new AtomicAccountRepository using the pgxpool provided.
func NewAtomicAccountRepo(conn *pgxpool.Pool, logger *zap.Logger) *AtomicAccountRepository {
	return &AtomicAccountRepository{
		conn:   conn,
		logger: logger,
	}
}

// AccountRepository satisfies class.Repository and uses *pgxpool.Pool directly.
type AccountRepository struct {
	conn   *pgxpool.Pool
	logger *zap.Logger
}

var _ usecases.AccountRepository = (*AccountRepository)(nil)

func (ar *AccountRepository) GetByNo(accountNo uuid.UUID) (*domain.Account, error) {
	var account domain.Account
	err := ar.conn.QueryRow(
		context.Background(),
		"SELECT accountno, userid, balance, minBalance FROM accounts WHERE accountno = $1",
		accountNo,
	).Scan(&account.AccountNo, &account.UserID, &account.Balance, &account.MinBalance)

	if err != nil {
		ar.logger.Error("Failed to get account by ID from Database", zap.Error(err))
		return nil, err
	}

	return &account, nil
}

func (ar *AccountRepository) DropAccountsTable() error {
	_, err := ar.conn.Exec(context.Background(), "DROP TABLE accounts")

	if err != nil {
		ar.logger.Error("Failed to drop account table", zap.Error(err))
		return err
	}

	return nil
}

func (ar *AccountRepository) CreateTable() error {
	_, err := ar.conn.Exec(
		context.Background(),
		"CREATE TABLE accounts (accountno UUID PRIMARY KEY, userid UUID, balance FLOAT, minBalance FLOAT, CONSTRAINT constrain_fk FOREIGN KEY (userid) REFERENCES users(userid) );",
	)
	if err != nil {
		ar.logger.Error("Failed to create table in Database", zap.Error(err))
		return err
	}
	return nil
}

func (ar *AccountRepository) CreateAccount(account *domain.Account) (uuid.UUID, error) {
	_, err := ar.conn.Exec(
		context.Background(),
		"INSERT INTO accounts (accountno, userid, balance, minBalance) VALUES($1, $2, $3, $4)",
		account.AccountNo, account.UserID, account.Balance, account.MinBalance,
	)
	if err != nil {
		ar.logger.Error("Failed to create account in Database", zap.Error(err))
		return uuid.MustParse("00000000-0000-0000-0000-000000000000"), err
	}
	return account.AccountNo, nil
}

func (ar *AccountRepository) GetAll() ([]domain.Account, error) {
	rows, err := ar.conn.Query(context.Background(), "SELECT accountno, userid, balance, minBalance FROM accounts")
	if err != nil {
		ar.logger.Error("Failed to get all accounts from Database", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var accounts []domain.Account
	for rows.Next() {
		var account domain.Account
		if err := rows.Scan(&account.AccountNo, &account.UserID, &account.Balance, &account.MinBalance); err != nil {
			ar.logger.Error("Failed to get account by ID from Database and append it to ds", zap.Error(err))
			return nil, err
		}
		accounts = append(accounts, account)
	}

	return accounts, nil
}

func (ar *AccountRepository) GetAccByUserId(userid uuid.UUID) (*domain.Account, error) {
	var account domain.Account
	err := ar.conn.QueryRow(
		context.Background(),
		"SELECT accountno, userid, balance, minBalance FROM accounts WHERE userid = $1",
		userid.String(),
	).Scan(&account.AccountNo, &account.UserID, &account.Balance, &account.MinBalance)
	if err != nil {
		ar.logger.Error("Failed to get account by ID from Database", zap.Error(err))
		return nil, err
	}
	return &account, nil
}

func (ar *AtomicAccountRepository) Execute(op usecases.AtomicAccountOperation) error {
	tx, err := ar.conn.Begin(context.Background())
	if err != nil {
		return err
	}

	defer func() { _ = tx.Rollback(context.Background()) }()

	// Create a new single-use AccountRepository backed by the transaction.
	accountRepoWithTransaction := AccountRepository{
		conn:   ar.conn,
		logger: ar.logger,
	}

	// Perform the AtomicAccountOperation using the AccountRepository.
	if err := op(&accountRepoWithTransaction); err != nil {
		return err
	}

	if err := tx.Commit(context.Background()); err != nil {
		return err
	}

	return nil
}
