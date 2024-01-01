package usecases

import (
	"fmt"
	"rest1/internal/domain"
	"sync"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/sys/windows/svc"
)

type AtomicOperation func(UserRepository) error

type AtomicUserRepository interface {
	Execute(AtomicOperation) error
}

type UserRepository interface {
	DropUserTable() error
	CreateUserTable() error
	GetIndividual(idd uuid.UUID, wg *sync.WaitGroup, userChan chan<- domain.User)
	GetAll() ([]domain.User, error)
	GetByID(id uuid.UUID) (*domain.User, error)
	CreateUser(user *domain.User) error
	Deposit(account *domain.Account, amount int) error
	Withdraw(account *domain.Account, amount int) error
}

type UserUsecase struct {
	AtomicUserRepo AtomicUserRepository
	logger *zap.Logger
}

/*
func NewAccountHandler(useCase *usecases.AccountUsecase , conn *pgx.Conn) *AccountHandler{
	return &AccountHandler{
		UseCase: useCase,
		conn: conn,
	}
}
*/

func NewUserUseCase(reposi AtomicUserRepository, logger *zap.Logger) *UserUsecase {
	return &UserUsecase{
		AtomicUserRepo:   reposi,
		logger: logger,
	}
}

func (a *UserUsecase) DropUserTable() error {
	drop := func (ur UserRepository) error {
		err := ur.DropUserTable()
		if err != nil {
			return err
		}
		return nil
	}
	
	if err := a.AtomicUserRepo.Execute(drop); err != nil {
		return err
	}

	return nil
}
func (a *UserUsecase) CreateUserTable() error {
	create := func (ur UserRepository) error {
		err := ur.CreateUserTable()
		if err != nil {
			return err
		}
		return nil
	}

	if err := a.AtomicUserRepo.Execute(create); err != nil {
		return err
	}

	return nil
}

func (a *UserUsecase) CreateUser(user *domain.User) error {
	//err := repository.NewUserRepo(conn , a.logger).CreateUser(user)
	createUsr := func (ur UserRepository) error {
		err := ur.CreateUser()
		if err != nil {
			return err
		}
		return nil
	}
	if err := a.AtomicUserRepo.Execute(createUsr); err != nil {
		return err
	}

	return nil
}

func (a *UserUsecase) GetUserById(id uuid.UUID) (*domain.User, error) {
	//user, err := repository.NewUserRepo(conn , a.logger).GetByID(id)
	user, err := a.repo.GetByID(id)
	if err != nil {
		a.logger.Error("Error performing user operation get account by id", zap.Error(err))
		return nil, err
	}
	return user, err
}

func (a *UserUsecase) GetAll() ([]domain.User, error) {
	// var userList []domain.User
	//userList, err := repository.NewUserRepo(conn , a.logger).GetAll()
	userList, err := a.repo.GetAll()
	if err != nil {
		a.logger.Error("Error performing user operation get all accounts", zap.Error(err))
		return nil, err
	}
	return userList, nil
}

func (a *UserUsecase) Withdraw(user *domain.User, amount int) error {
	// check if minBalance violated
	account, err := a.AccountUsecase.repo.GetAccByUserId(user.ID)
	// account, err := a.repo.GetByAccountNo(user.AccountNo)
	if account.Balance-float64(amount) < account.MinBalance {
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

func (a *UserUsecase) Deposit(user *domain.User, amount int) error {
	//err := repository.NewUserRepo(conn , a.logger).Deposit(user , amount);
	account, err := a.AccountUsecase.repo.GetAccByUserId(user.ID)
	fmt.Println("usecases accountbyuserid", account.AccountNo)
	err = a.repo.Deposit(account, amount)
	if err != nil {
		a.logger.Error("Error performing user operation desposit ", zap.Error(err))
		return err
	}
	return nil
}
