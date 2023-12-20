package domain

type Account struct {
	AccountNo int `json:"accountNo"`
	Balance int `json:"balance"`
	MinBalance int `json:"minbalance"`
}

type AccountMethods interface {
	GetByNo(id int) (*User, error)
	CreateAccount(use *User) error
	GetAll() ([]Account, error)
}