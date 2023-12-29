package usecases

import (
	"fmt"
	"log"
	"rest1/internal/domain"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type AtomicAccountRepository interface {
	Transactor

	DropAccountsTable() error
	CreateTable() error
	GetByNo(accountNo uuid.UUID) (*domain.Account, error)
	CreateAccount(account *domain.Account) (uuid.UUID, error)
	GetAll() ([]domain.Account, error)
	GetAccByUserId(userid uuid.UUID) (*domain.Account, error)
}

type AtomicAccountRepositoryFactory func() AtomicAccountRepository

type AccountUsecase struct {
	newARepo AtomicAccountRepositoryFactory
	logger   *zap.Logger
}

func NewAccountUseCase(accreposiFactory AtomicAccountRepositoryFactory, logger *zap.Logger) *AccountUsecase {
	return &AccountUsecase{
		newARepo: accreposiFactory,
		logger:   logger,
	}
}

func (a *AccountUsecase) DropAccountsTable() error {
	repo := a.newARepo()

	if err := repo.Begin(); err != nil {
		a.logger.Error("Failed to Begin drop account table", zap.Error(err))
		return err
	}

	defer func() {
		_ = repo.Rollback()
	}()

	err := repo.DropAccountsTable()
	if err != nil {
		fmt.Println("Error while deleting account table")
		return err
	}

	if err := repo.Commit(); err != nil {
		a.logger.Error("Failed to Commit delete account table", zap.Error(err))
		return err
	}

	return nil
}

func (a *AccountUsecase) CreateAccountTable() error {
	repo := a.newARepo()

	if err := repo.Begin(); err != nil {
		a.logger.Error("Failed to Begin drop account table", zap.Error(err))
		return err
	}

	defer func() {
		_ = repo.Rollback()
	}()

	err := repo.CreateTable()
	if err != nil {
		fmt.Println("Error while creating account table")
		return err
	}

	if err := repo.Commit(); err != nil {
		a.logger.Error("Failed to Commit delete account table", zap.Error(err))
		return err
	}

	return nil
}
func (a *AccountUsecase) CreateAccount(userID uuid.UUID) (uuid.UUID, error) {
	repo := a.newARepo()

	if err := repo.Begin(); err != nil {
		a.logger.Error("Failed to Begin drop account table", zap.Error(err))
		return uuid.MustParse("00000000-0000-0000-0000-000000000000"), err
	}

	defer func() {
		_ = repo.Rollback()
	}()

	var newAccount domain.Account

	newAccount.UserID = userID
	newAccount.Balance = float64(0.0)
	newAccount.MinBalance = float64(500.0)
	newAccount.AccountNo = uuid.New()
	fmt.Println("usecase account id", newAccount.AccountNo)

	//err := repository.NewAccountRepo(conn , a.logger).CreateAccount(&newAccount)
	accid, err := repo.CreateAccount(&newAccount)
	if err != nil {
		log.Fatal(err)
		return uuid.MustParse("00000000-0000-0000-0000-000000000000"), err
	}

	if err := repo.Commit(); err != nil {
		a.logger.Error("Failed to Commit delete account table", zap.Error(err))
		return uuid.MustParse("00000000-0000-0000-0000-000000000000"), err
	}

	return accid, nil
}

func (a *AccountUsecase) GetByAccountNo(accountNo uuid.UUID) (*domain.Account, error) {
	repo := a.newARepo()

	if err := repo.Begin(); err != nil {
		a.logger.Error("Failed to Begin drop account table", zap.Error(err))
		return nil, err
	}

	defer func() {
		_ = repo.Rollback()
	}()
	//account , err := repository.NewAccountRepo(conn , a.logger).GetByNo(accountNo)
	account, err := repo.GetByNo(accountNo)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	if err := repo.Commit(); err != nil {
		a.logger.Error("Failed to Commit delete account table", zap.Error(err))
		return nil, err
	}
	return account, nil
}

func (a *AccountUsecase) GetAccountByUserID(userid uuid.UUID) (*domain.Account, error) {
	repo := a.newARepo()

	if err := repo.Begin(); err != nil {
		a.logger.Error("Failed to Begin drop account table", zap.Error(err))
		return nil, err
	}

	defer func() {
		_ = repo.Rollback()
	}()

	account, err := repo.GetAccByUserId(userid)
	if err != nil {
		a.logger.Error("Could not Ftech account")
		return nil, err
	}

	if err := repo.Commit(); err != nil {
		a.logger.Error("Failed to Commit delete account table", zap.Error(err))
		return nil, err
	}

	return account, nil
}

func (a *AccountUsecase) GetAllAccounts() ([]domain.Account, error) {
	repo := a.newARepo()

	if err := repo.Begin(); err != nil {
		a.logger.Error("Failed to Begin drop account table", zap.Error(err))
		return nil, err
	}

	defer func() {
		_ = repo.Rollback()
	}()

	//accounts , err := repository.NewAccountRepo(conn , a.logger).GetAll()
	accounts, err := repo.GetAll()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	if err := repo.Commit(); err != nil {
		a.logger.Error("Failed to Commit delete account table", zap.Error(err))
		return nil, err
	}

	return accounts, nil
}
