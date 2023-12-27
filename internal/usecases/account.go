package usecases

import (
	"fmt"
	"log"
	"rest1/internal/domain"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type AccountRepository interface {
	DropAccountsTable() error
	CreateTable() error
	GetByNo(accountNo uuid.UUID) (*domain.Account, error)
	CreateAccount(account *domain.Account) (uuid.UUID, error)
	GetAll() ([]domain.Account, error)
	GetAccByUserId(userid uuid.UUID) (*domain.Account, error)
}

type AccountUsecase struct {
	repo   AccountRepository
	logger *zap.Logger
}

func NewAccountUseCase(reposi AccountRepository, logger *zap.Logger) *AccountUsecase {
	return &AccountUsecase{
		repo:   reposi,
		logger: logger,
	}
}

func (a *AccountUsecase) DropAccountsTable() error {
	err := a.repo.DropAccountsTable()
	if err != nil {
		fmt.Println("Error while deleting account table")
		return err
	}

	return nil
}

func (a *AccountUsecase) CreateAccountTable() error {
	err := a.repo.CreateTable()
	if err != nil {
		fmt.Println("Error while creating account table")
		return err
	}

	return nil
}
func (a *AccountUsecase) CreateAccount(userID uuid.UUID) (uuid.UUID, error) {
	var newAccount domain.Account

	newAccount.UserID = userID
	newAccount.Balance = float64(0.0)
	newAccount.MinBalance = float64(500.0)
	newAccount.AccountNo = uuid.New()
	fmt.Println("usecase account id", newAccount.AccountNo)

	//err := repository.NewAccountRepo(conn , a.logger).CreateAccount(&newAccount)
	accid, err := a.repo.CreateAccount(&newAccount)
	if err != nil {
		log.Fatal(err)
		return uuid.MustParse("00000000-0000-0000-0000-000000000000"), err
	}
	return accid, nil
}

func (a *AccountUsecase) GetByAccountNo(accountNo uuid.UUID) (*domain.Account, error) {
	//account , err := repository.NewAccountRepo(conn , a.logger).GetByNo(accountNo)
	account, err := a.repo.GetByNo(accountNo)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return account, nil
}

func (a *AccountUsecase) GetAccountByUserID(userid uuid.UUID) (*domain.Account, error) {
	account, err := a.repo.GetAccByUserId(userid)
	if err != nil {
		a.logger.Error("Could not Ftech account")
		return nil, err
	}
	return account, err
}

func (a *AccountUsecase) GetAllAccounts() ([]domain.Account, error) {
	//accounts , err := repository.NewAccountRepo(conn , a.logger).GetAll()
	accounts, err := a.repo.GetAll()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return accounts, err
}
