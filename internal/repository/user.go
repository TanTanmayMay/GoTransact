package repository

import (
	"context"
	"rest1/internal/domain"
	"fmt"
	"github.com/jackc/pgx/v4"
	// type UserRepository struct {
	// 	Storage *[]domain.User
	// }
)

type UserRepo struct {
	conn *pgx.Conn
}

func NewUserRepo(conn *pgx.Conn) *UserRepo {
	return &UserRepo{conn: conn}
}

func (u *UserRepo) CreateTable() error {
	_ , err := u.conn.Exec(context.Background() , "CREATE TABLE users (id INT PRIMARY KEY,name VARCHAR ( 50 )  NOT NULL,password VARCHAR ( 50 ) NOT NULL,accountNo INT);")
	if err != nil{
		fmt.Println(err)
	}
	return nil
}
func (u *UserRepo) GetAll(conn *pgx.Conn) ([]domain.User, error) {
	rows, err := conn.Query(context.Background(), "SELECT id, name, accountNo, password FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var user domain.User
		if err := rows.Scan(&user.ID, &user.Name, &user.AccountNo, &user.Password); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (u *UserRepo) GetByID(id string) (*domain.User, error) {
	var user domain.User
	err := u.conn.QueryRow(context.Background(), "SELECT id, name, accountNo, password FROM users WHERE id = $1", id).
		Scan(&user.ID, &user.Name, &user.AccountNo, &user.Password)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *UserRepo) Create(user *domain.User) (*domain.User, error) {
	var id int
	err := u.conn.QueryRow(context.Background(), "INSERT INTO users(id , name, accountNo, password) VALUES($1, $2, $3 , $4) RETURNING id", 125, "Nishant", 123 , "123").Scan(&id)
	if(err != nil) {
		fmt.Println(err)
	}
	return user , nil
}

func (u *UserRepo) Withdraw(user *domain.User) error {
	// Implement withdrawal logic here
	return nil
}

func (u *UserRepo) Deposit(user *domain.User) error {
	// Implement deposit logic here
	return nil
}