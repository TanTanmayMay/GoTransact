package domain

type User struct {
	ID			string 		`json:"id"`
	Name 		string 		`json:"name"`
	AccountNo	string 		`json:"accountNo"`
	Password 	string 		`json:"password"`
}


type UserMethods interface {
	GetAll() ([]User, error)
	GetByID(id int) (*User, error)
	Create(user *User) error
	Update(user *User) error
}