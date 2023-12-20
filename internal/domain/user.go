package domain

type User struct {
	ID			string 		`json:"id"` //pri
	Name 		string 		`json:"name"`
	AccountNo	string 		`json:"accountNo"` //foreign key
	Password 	string 		`json:"password"`
}


type UserMethods interface {
	GetAll() ([]User, error)
	GetByID(id int) (*User, error)
	Create(user *User) (*Account , error)
	Withdraw(user* User) (error)
	Deposit(user *User) error
}