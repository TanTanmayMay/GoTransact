package domain

import (
	"github.com/google/uuid"
)
type User struct {
	ID			uuid.UUID 		`json:"id"` //pri
	Name 		string 			`json:"name"`
	Password 	string 			`json:"password"`
}


type UserMethods interface {
	GetAll() ([]User, error)
	GetByID(id uuid.UUID) (*User, error)
	CreateUser(user *User) error
	Withdraw(user* User, amount int) (error)
	Deposit(user *User, amount int) error
	CreateTable() error
}