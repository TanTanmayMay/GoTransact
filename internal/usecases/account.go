package usecases

import (
	"fmt"
	"log"
	"rest1/internal/domain"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// Define AtomicAccountOperation interface
type AtomicAccountOperation func(AccountRepository) error

// Define AtomicAccountRepository interface
type AtomicAccountRepository interface {
	Execute(AtomicAccountOperation) error
}

// Define AccountRepository interface
type AccountRepository interface {
	DropAccountsTable() error
	CreateTable() error
	GetByNo(accountNo uuid.UUID) (*domain.Account, error)
	CreateAccount(account *domain.Account) (uuid.UUID, error)
	GetAll() ([]domain.Account, error)
	GetAccByUserId(userid uuid.UUID) (*domain.Account, error)
}

// Define AccountUsecase struct
type AccountUsecase struct {
	atomicRepo AtomicAccountRepository
	logger     *zap.Logger
}

// Define constructor
func NewAccountUsecase(atomicRepo AtomicAccountRepository, logger *zap.Logger) *AccountUsecase {
	return &AccountUsecase{
		atomicRepo: atomicRepo,
		logger:     logger,
	}
}

func (a *AccountUsecase) DropAccountsTable() error {
	dropTableAtomicOp := func(repo AccountRepository) error {
		return repo.DropAccountsTable()
	}

	if err := a.atomicRepo.Execute(dropTableAtomicOp); err != nil {
		fmt.Println("Error while deleting account table")
		return err
	}

	return nil
}

// Implement CreateAccountTable method
func (a *AccountUsecase) CreateAccountTable() error {
	createTableAtomicOp := func(repo AccountRepository) error {
		return repo.CreateTable()
	}

	if err := a.atomicRepo.Execute(createTableAtomicOp); err != nil {
		fmt.Println("Error while creating account table")
		return err
	}

	return nil
}

// Implement CreateAccount method
func (a *AccountUsecase) CreateAccount(userID uuid.UUID) (uuid.UUID, error) {
	var newAccount domain.Account

	newAccount.UserID = userID
	newAccount.Balance = float64(0.0)
	newAccount.MinBalance = float64(500.0)
	newAccount.AccountNo = uuid.New()
	fmt.Println("usecase account id", newAccount.AccountNo)

	createAccountAtomicOp := func(repo AccountRepository) error {
		_, err := repo.CreateAccount(&newAccount)
		return err
	}

	if err := a.atomicRepo.Execute(createAccountAtomicOp); err != nil {
		log.Fatal(err)
		return uuid.MustParse("00000000-0000-0000-0000-000000000000"), err
	}

	return newAccount.AccountNo, nil
}

// Implement GetByAccountNo method
func (a *AccountUsecase) GetByAccountNo(accountNo uuid.UUID) (*domain.Account, error) {
	var account *domain.Account
	getByAccountNoAtomicOp := func(repo AccountRepository) error {
		var err error
		account, err = repo.GetByNo(accountNo)
		if err != nil {
			return err
		}
		// Perform additional business logic or validations if needed
		return nil
	}

	if err := a.atomicRepo.Execute(getByAccountNoAtomicOp); err != nil {
		log.Fatal(err)
		return nil, err
	}

	return account, nil
}

// Implement GetAccountByUserID method
func (a *AccountUsecase) GetAccountByUserID(userid uuid.UUID) (*domain.Account, error) {
	var account *domain.Account
	getAccountByUserIDAtomicOp := func(repo AccountRepository) error {
		var err error
		account, err = repo.GetAccByUserId(userid)
		return err
	}

	if err := a.atomicRepo.Execute(getAccountByUserIDAtomicOp); err != nil {
		a.logger.Error("Could not fetch account")
		return nil, err
	}

	return account, nil
}

// Implement GetAllAccounts method
func (a *AccountUsecase) GetAllAccounts() ([]domain.Account, error) {
	var accounts []domain.Account
	getAllAccountsAtomicOp := func(repo AccountRepository) error {
		var err error
		accounts, err = repo.GetAll()
		return err
	}

	if err := a.atomicRepo.Execute(getAllAccountsAtomicOp); err != nil {
		log.Fatal(err)
		return nil, err
	}

	return accounts, nil
}
