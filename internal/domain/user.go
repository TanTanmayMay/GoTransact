package domain

import "github.com/jackc/pgx/v4"

type User struct {
	ID			int 		`json:"id"` //pri
	Name 		string 		`json:"name"`
	AccountNo	int 		`json:"accountNo"` //foreign key
	Password 	string 		`json:"password"`
}


type UserMethods interface {
	GetAll() ([]User, error)
	GetByID(id int) (*User, error)
	Create(user *User) (*Account , error)
	Withdraw(user* User) (error)
	Deposit(user *User) error
	CreateTable() error
}