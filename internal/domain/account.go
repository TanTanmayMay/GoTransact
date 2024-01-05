package domain

import "github.com/google/uuid"

type Account struct {
	AccountNo  uuid.UUID `json:"accountNo"`
	Balance    float64   `json:"balance"`
	MinBalance float64   `json:"minbalance"`
	UserID     uuid.UUID `json:"userID"` //foreign key
}
