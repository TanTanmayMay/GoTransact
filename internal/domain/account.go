package domain

type Account struct {
	AccountNo int32 `json:"accountNo"`
	Balance float32 `json:"balance"`
	MinBalance float32 `json:"minbalance"`

}


type AccountMethods interface {
	GetByNo(id int) (*User, error)
	Create(use *User) error

}