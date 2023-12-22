package usecases

import (
	"fmt"
	"rest1/internal/domain"
	"rest1/internal/repository"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)


type UserUsercasesMethods interface {
	// user *domain.User
	DropUserTable() error
	CreateUserTable() error
	CreateUser(user *domain.User, conn *pgx.Conn) error
    Deposit(user *domain.User, amount int, conn *pgx.Conn) error
    Withdraw(user *domain.User, amount int, conn *pgx.Conn) error
	GetAll(conn *pgx.Conn) ([]domain.User, error)
	GetAccountByID(user *domain.User, conn* pgx.Conn) (*domain.User ,error)
}


type UserUsecase struct {
	repo *repository.UserRepo
	conn *pgx.Conn
	logger *zap.Logger
	AccountUsecase
}

/* 
func NewAccountHandler(useCase *usecases.AccountUsecase , conn *pgx.Conn) *AccountHandler{
	return &AccountHandler{
		UseCase: useCase,
		conn: conn,
	}
}
*/



func NewUserUseCase (reposi *repository.UserRepo, conn *pgx.Conn , logger *zap.Logger) *UserUsecase{
	return &UserUsecase{
		repo: reposi,
		conn: conn,
		logger: logger,
	}
}

func (a *UserUsecase) DropUserTable() error{
	err := a.repo.DropUserTable()
	if err != nil {
		a.logger.Error("Failed to create user table", zap.Error(err))
        return err
	}

	return nil
}
func (a *UserUsecase) CreateUserTable() error {
	err := a.repo.CreateUserTable()
	if err != nil {
        a.logger.Error("Failed to create user table", zap.Error(err))
        return err
    }

    return nil
}

func (a *UserUsecase) CreateUser(user* domain.User, conn *pgx.Conn) error {
    //err := repository.NewUserRepo(conn , a.logger).CreateUser(user)
	err := a.repo.CreateUser(user)
	if err != nil {
        return err
    }

    return nil
}

func (a *UserUsecase) GetUserById(id uuid.UUID, conn* pgx.Conn ) (*domain.User ,error) {
	//user, err := repository.NewUserRepo(conn , a.logger).GetByID(id)
	user, err := a.repo.GetByID(id)
	if(err != nil){
		a.logger.Error("Error performing user operation get account by id", zap.Error(err))
		return nil, err
	}
	return user, err
}

func (a *UserUsecase)  GetAll(conn *pgx.Conn) ([]domain.User, error) {
	// var userList []domain.User
	//userList, err := repository.NewUserRepo(conn , a.logger).GetAll()
	userList, err := a.repo.GetAll()
	if err != nil {
		a.logger.Error("Error performing user operation get all accounts", zap.Error(err))
		return nil, err
	}
	return userList, nil
}


func (a *UserUsecase) Withdraw(user *domain.User, amount int, conn *pgx.Conn) error {
	// check if minBalance violated
	account , err := a.AccountUsecase.repo.GetAccByUserId(user.ID);
	// account, err := a.repo.GetByAccountNo(user.AccountNo)
	if account.Balance - float64(amount) < account.MinBalance {
		a.logger.Error("Error performing user operation get withdrawal due to min balance violation", zap.Error(err))
		return nil //Custom Error possible ??

	} 
	//err = repository.NewUserRepo(conn , a.logger).Withdraw(user, amount)
	err = a.repo.Withdraw(account, amount)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil	
}

func (a *UserUsecase) Deposit(user *domain.User, amount int, conn *pgx.Conn) error {
	//err := repository.NewUserRepo(conn , a.logger).Deposit(user , amount);
	account , err := a.AccountUsecase.repo.GetAccByUserId(user.ID);
	err = a.repo.Deposit(account, amount)
	if err != nil {
		a.logger.Error("Error performing user operation desposit " , zap.Error(err))
		return err
	}
	return nil
}