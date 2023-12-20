package domain

type User struct {
	ID			int 		`json:"id"` //pri
	Name 		string 		`json:"name"`
	AccountNo	int 		`json:"accountNo"` //foreign key
	Password 	string 		`json:"password"`
}


type UserMethods interface {
	GetAll() ([]User, error)
	GetByID(id int) (*User, error)
	CreateUser(user *User) error
	Withdraw(user* User, amount int) (error)
	Deposit(user *User, amount int) error
	CreateTable() error
}