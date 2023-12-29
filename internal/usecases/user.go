package usecases

import (
	"rest1/internal/domain"
	"sync"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Transactor interface {
	Begin() error
	Commit() error
	Rollback() error
}

type AtomicUserRepository interface {
	Transactor

	DropUserTable() error
	CreateUserTable() error
	GetIndividual(idd uuid.UUID, wg *sync.WaitGroup, userChan chan<- domain.User)
	GetAll() ([]domain.User, error)
	GetByID(id uuid.UUID) (*domain.User, error)
	CreateUser(user *domain.User) error
	Deposit(account *domain.Account, amount int) error
	Withdraw(account *domain.Account, amount int) error
}

type AtomicUserRepositoryFactory func() AtomicUserRepository

type UserUsecase struct {
	newRepo AtomicUserRepositoryFactory
	logger  *zap.Logger
	AccountUsecase
}

func NewUserUseCase(reposiFactory AtomicUserRepositoryFactory, logger *zap.Logger) *UserUsecase {
	return &UserUsecase{
		newRepo: reposiFactory,
		logger:  logger,
	}
}

func (a *UserUsecase) DropUserTable() error {
	repo := a.newRepo()

	if err := repo.Begin(); err != nil {
		a.logger.Error("Failed to Begin create user table", zap.Error(err))
		return err
	}

	defer func() {
		_ = repo.Rollback()
	}()

	err := repo.DropUserTable()
	if err != nil {
		a.logger.Error("Failed to create user table", zap.Error(err))
		return err
	}

	if err := repo.Commit(); err != nil {
		a.logger.Error("Failed to Commit create user table", zap.Error(err))
		return err
	}

	return nil
}

func (a *UserUsecase) CreateUserTable() error {
	repo := a.newRepo()

	if err := repo.Begin(); err != nil {
		a.logger.Error("Failed to Begin create user table", zap.Error(err))
		return err
	}

	defer func() {
		_ = repo.Rollback()
	}()

	err := repo.CreateUserTable()
	if err != nil {
		a.logger.Error("Failed to create user table", zap.Error(err))
		return err
	}

	if err := repo.Commit(); err != nil {
		a.logger.Error("Failed to Commit create user table", zap.Error(err))
		return err
	}

	return nil
}

func (a *UserUsecase) CreateUser(user *domain.User) error {
	//err := repository.NewUserRepo(conn , a.logger).CreateUser(user)
	repo := a.newRepo()

	if err := repo.Begin(); err != nil {
		a.logger.Error("Failed to Begin create user table", zap.Error(err))
		return err
	}

	defer func() {
		_ = repo.Rollback()
	}()

	err := repo.CreateUser(user)
	if err != nil {
		return err
	}

	if err := repo.Commit(); err != nil {
		a.logger.Error("Failed to Commit create user table", zap.Error(err))
		return err
	}

	return nil
}

func (a *UserUsecase) GetUserById(id uuid.UUID) (*domain.User, error) {
	//user, err := repository.NewUserRepo(conn , a.logger).GetByID(id)
	repo := a.newRepo()

	if err := repo.Begin(); err != nil {
		a.logger.Error("Failed to Begin create user table", zap.Error(err))
		return nil, err
	}

	defer func() {
		_ = repo.Rollback()
	}()

	user, err := repo.GetByID(id)
	if err != nil {
		a.logger.Error("Error performing user operation get account by id", zap.Error(err))
		return nil, err
	}

	if err := repo.Commit(); err != nil {
		a.logger.Error("Failed to Commit create user table", zap.Error(err))
		return nil, err
	}

	return user, err
}

func (a *UserUsecase) GetAll() ([]domain.User, error) {
	// var userList []domain.User
	//userList, err := repository.NewUserRepo(conn , a.logger).GetAll()
	repo := a.newRepo()

	if err := repo.Begin(); err != nil {
		a.logger.Error("Failed to Begin create user table", zap.Error(err))
		return nil, err
	}

	defer func() {
		_ = repo.Rollback()
	}()

	userList, err := repo.GetAll()
	if err != nil {
		a.logger.Error("Error performing user operation get all accounts", zap.Error(err))
		return nil, err
	}

	if err := repo.Commit(); err != nil {
		a.logger.Error("Failed to Commit create user table", zap.Error(err))
		return nil, err
	}
	return userList, nil
}

/* 		-> TODO
func (a *UserUsecase) Withdraw(user *domain.User, amount int) error {
	repo := a.newRepo()

	if err := repo.Begin(); err != nil {
		a.logger.Error("Failed to Begin create user table", zap.Error(err))
		return err
	}

	defer func() {
		_ = repo.Rollback()
	}()

	// check if minBalance violated
	account, err := a.AccountUsecase.repo.GetAccByUserId(user.ID)
	// account, err := a.repo.GetByAccountNo(user.AccountNo)
	if account.Balance-float64(amount) < account.MinBalance {
		a.logger.Error("Error performing user operation get withdrawal due to min balance violation", zap.Error(err))
		return nil //Custom Error possible ??

	}
	//err = repository.NewUserRepo(conn , a.logger).Withdraw(user, amount)
	err = repo.Withdraw(account, amount)
	if err != nil {
		fmt.Println(err)
		return err
	}

	if err := repo.Commit(); err != nil {
		a.logger.Error("Failed to Commit create user table", zap.Error(err))
		return err
	}

	return nil
}

func (a *UserUsecase) Deposit(user *domain.User, amount int) error {
	repo := a.newRepo()

	if err := repo.Begin(); err != nil {
		a.logger.Error("Failed to Begin create user table", zap.Error(err))
		return err
	}

	defer func() {
		_ = repo.Rollback()
	}()

	//err := repository.NewUserRepo(conn , a.logger).Deposit(user , amount);
	account, err := a.AccountUsecase.repo.GetAccByUserId(user.ID)
	fmt.Println("usecases accountbyuserid", account.AccountNo)
	err = repo.Deposit(account, amount)
	if err != nil {
		a.logger.Error("Error performing user operation desposit ", zap.Error(err))
		return err
	}

	if err := repo.Commit(); err != nil {
		a.logger.Error("Failed to Commit create user table", zap.Error(err))
		return err
	}

	return nil
}
*/
