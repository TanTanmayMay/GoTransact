package usecases

import (
	"log"
	"math/rand"
	"rest1/internal/domain"
	"rest1/internal/repository"
	"time"

	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

type AccountUserCaseMethods interface {
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
func (a *AccountUsecase) CreateAccount(userID int , conn *pgx.Conn) (int, error) {
	var newAccount domain.Account
	rand.Seed(time.Now().UnixNano())
	randomNumber := rand.Intn(500)

	newAccount.AccountNo = userID + randomNumber
	newAccount.Balance = 0
	newAccount.MinBalance = 500
	//err := repository.NewAccountRepo(conn , a.logger).CreateAccount(&newAccount)
	id, err := a.repo.CreateAccount(&newAccount)
	if err != nil {
		log.Fatal(err)
		return 0, err
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

func (a * AccountUsecase) GetAllAccounts(conn *pgx.Conn) ([] domain.Account , error){
	//accounts , err := repository.NewAccountRepo(conn , a.logger).GetAll()
	accounts, err := a.repo.GetAll()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return accounts, err
}