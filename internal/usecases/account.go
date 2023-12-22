package usecases

import (
	"fmt"
	"log"
	"rest1/internal/domain"
	"rest1/internal/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

type AccountUserCaseMethods interface {
	DropAccountsTable() error
	CreateAccountTable() error
	CreateAccount(userId int , conn *pgx.Conn) (int, error)
    GetByAccountNo(accountNo int , conn *pgx.Conn) (* domain.Account , error)
	GetAllAccounts(conn *pgx.Conn) ([]domain.Account , error)	//[]domain.Account
}

type AccountUsecase struct {
	repo *repository.AccountRepo 
	conn *pgx.Conn
	logger *zap.Logger
}

func NewAccountUseCase (reposi *repository.AccountRepo, conn *pgx.Conn , logger *zap.Logger) *AccountUsecase{
	return &AccountUsecase{
		repo: reposi,
		conn: conn,
		logger: logger,
	}
}

func (a *AccountUsecase) DropAccountsTable() error{
	err := a.repo.DropAccountsTable()
	if err != nil {
		fmt.Println("Error while deleting account table")
		return err
	}

	return nil
}

func (a *AccountUsecase) CreateAccountTable() error{
	err := a.repo.CreateTable()
	if err != nil {
		fmt.Println("Error while creating account table")
		return err
	}

	return nil
}
func (a *AccountUsecase) CreateAccount(userID uuid.UUID , conn *pgx.Conn) (uuid.UUID, error) {
	var newAccount domain.Account

	newAccount.UserID = userID 
	newAccount.Balance = float64(0.0)
	newAccount.MinBalance = float64(500.0)
	newAccount.AccountNo = uuid.New()
	//err := repository.NewAccountRepo(conn , a.logger).CreateAccount(&newAccount)
	id, err := a.repo.CreateAccount(&newAccount)
	if err != nil {
		log.Fatal(err)
		return uuid.MustParse("00000000-0000-0000-0000-000000000000"), err
	}
	return id, nil
}

func (a *AccountUsecase) GetByAccountNo(accountNo int , conn* pgx.Conn) (* domain.Account , error) {
	//account , err := repository.NewAccountRepo(conn , a.logger).GetByNo(accountNo)
	account, err := a.repo.GetByNo(accountNo)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return account, nil
}

func (a *AccountUsecase) GetAccountByUserID(userid uuid.UUID) (* domain.Account, error) {
	account, err := a.repo.GetAccByUserId(userid)
	if err != nil {
		a.logger.Error("Could not Ftech account")
		return nil, err
	}
	return account, err
}

func (a * AccountUsecase) GetAllAccounts(conn *pgx.Conn) ([] domain.Account , error){
	//accounts , err := repository.NewAccountRepo(conn , a.logger).GetAll()
	accounts, err := a.repo.GetAll()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return accounts, err
}