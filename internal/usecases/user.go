package usecases

import (
	"fmt"
	"rest1/internal/domain"
	"rest1/internal/repository"
	"github.com/jackc/pgx/v5"
)


type UserUsercasesMethods interface {
	// user *domain.User
	CreateUser(user *domain.User, conn *pgx.Conn) error
    Deposit(user *domain.User, amount int, conn *pgx.Conn) error
    Withdraw(user *domain.User, amount int, conn *pgx.Conn) error
	GetAll(conn *pgx.Conn) ([]domain.User, error)
	GetAccountByID(user *domain.User, conn* pgx.Conn) (*domain.User ,error)
}


type UserUsecase struct {
	repo UserUsercasesMethods
	conn *pgx.Conn
}

/* 
func NewAccountHandler(useCase *usecases.AccountUsecase , conn *pgx.Conn) *AccountHandler{
	return &AccountHandler{
		UseCase: useCase,
		conn: conn,
	}
}
*/

func NewUserUseCase (repo *repository.UserRepo, conn *pgx.Conn) {
	return &UserUsecase{
		repo: 
	}
}
func (a *UserUsecase) CreateUser(user* domain.User, conn *pgx.Conn) error {
    err := repository.NewUserRepo(conn).CreateUser(user)
	if(err != nil){
		fmt.Println(err)
		return err
	}
	
	return nil
}

func (a *UserUsecase) GetAccountByID(id int, conn* pgx.Conn ) (*domain.User ,error) {
	user, err := repository.NewUserRepo(conn).GetByID(id)
	if(err != nil){
		fmt.Println(err)
		return nil, err
	}
	return user, err
}

func (a *UserUsecase)  GetAll(conn *pgx.Conn) ([]domain.User, error) {
	// var userList []domain.User
	userList, err := repository.NewUserRepo(conn).GetAll()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return userList, nil
}


func (a *UserUsecase) Withdraw(user *domain.User, amount int, conn *pgx.Conn) error {
	// check if minBalance violated
	account , err:= AccountUserCaseMethods.GetByAccountNo(user.AccountNo , conn)
	if account.Balance - amount < account.MinBalance {
		fmt.Println("Cannot Withdraw the given amount as your minimum Balance has to be : ", account.MinBalance)
		return nil //Custom Error possible ??

	} 
	err = repository.NewUserRepo(conn).Withdraw(user, amount)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil	
}

func (a *UserUsecase) Deposit(user *domain.User, amount int, conn *pgx.Conn) error {
	err := repository.NewUserRepo(conn).Deposit(user , amount);
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}