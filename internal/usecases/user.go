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
	Deposit(user *domain.User, amount int) error
	Withdraw(user *domain.User, amount int) error
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
	getByiD := func(repo UserRepository) error {
		var err error
		user, err = repo.GetByID(id)
		if err != nil {
			return err
		}
		return nil
	}

	if err := a.atomicRepo.Execute(getByiD); err != nil {
		a.logger.Error("Failed to get user by id from db", zap.Error(err))
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
		return nil
	}
	if err := a.atomicRepo.Execute(getAll); err != nil {
		a.logger.Error("Failed to GetAll users from database", zap.Error(err))
		return nil, err
	}
	return userList, nil
}

func (a *UserUsecase) Withdraw(user *domain.User, amount int) error {
	withdraw := func(repo UserRepository) error {
		err := repo.Withdraw(user, amount)
		if err != nil {
			return err
		}
		return nil
	}
	if err := a.atomicRepo.Execute(withdraw); err != nil {
		a.logger.Error("Failed to Withdraw the given amount", zap.Error(err))
		return err
	}
	return nil
}

func (a *UserUsecase) Deposit(user *domain.User, amount int) error {
	withdraw := func(repo UserRepository) error {
		err := repo.Deposit(user, amount)
		if err != nil {
			return err
		}
		return nil
	}
	if err := a.atomicRepo.Execute(withdraw); err != nil {
		a.logger.Error("Failed to Deposit the given amount", zap.Error(err))
		return err
	}
	return nil
}
