package usecases

import (
	"fmt"
	"rest1/internal/domain"
	"sync"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type AtomicUserOperation func(UserRepository) error

type AtomicUserRepository interface {
	Execute(AtomicUserOperation) error
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
	atomicRepo AtomicUserRepository
	logger     *zap.Logger
}

func NewUserUseCase(reposi AtomicUserRepository, logger *zap.Logger) *UserUsecase {
	return &UserUsecase{
		atomicRepo: reposi,
		logger:     logger,
	}
}

func (a *UserUsecase) DropUserTable() error {
	dropTableAtomicOp := func(repo UserRepository) error {
		return repo.DropUserTable()
	}

	if err := a.atomicRepo.Execute(dropTableAtomicOp); err != nil {
		fmt.Println("Error while deleting user table")
		return err
	}

	return nil
}

func (a *UserUsecase) CreateUserTable() error {
	createTableAtomicOp := func(repo UserRepository) error {
		return repo.CreateUserTable()
	}

	if err := a.atomicRepo.Execute(createTableAtomicOp); err != nil {
		fmt.Println("Error while creating account table")
		return err
	}

	return nil
}

func (a *UserUsecase) CreateUser(user *domain.User) error {
	//err := repository.NewUserRepo(conn , a.logger).CreateUser(user)
	createAccountAtomicOp := func(repo UserRepository) error {
		err := repo.CreateUser(user)
		return err
	}

	if err := a.atomicRepo.Execute(createAccountAtomicOp); err != nil {
		fmt.Println("Error while creating User")

		return err
	}

	return nil
}

func (a *UserUsecase) GetUserById(id uuid.UUID) (*domain.User, error) {
	var user *domain.User
	//user, err := repository.NewUserRepo(conn , a.logger).GetByID(id)
	getByiD := func(repo UserRepository) error {
		var err error
		user, err = repo.GetByID(id)
		if err != nil {
			return err
		}
		// Perform additional business logic or validations if needed
		return nil
	}

	if err := a.atomicRepo.Execute(getByiD); err != nil {
		// log.Fatal(err)
		return nil, err
	}
	return user, nil
}

func (a *UserUsecase) GetAll() ([]domain.User, error) {
	var userList []domain.User
	//userList, err := repository.NewUserRepo(conn , a.logger).GetAll()
	getAll := func(repo UserRepository) error {
		var err error
		userList, err = repo.GetAll()
		if err != nil {
			return err
		}
		// Perform additional business logic or validations if needed
		return nil
	}
	if err := a.atomicRepo.Execute(getAll); err != nil {
		// log.Fatal(err)
		return nil, err
	}
	return userList, nil
}

/*
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
*/
