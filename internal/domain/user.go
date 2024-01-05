package domain

import (
	"github.com/google/uuid"
)

type User struct {
	ID       uuid.UUID `json:"id"` //pri
	Name     string    `json:"name"`
	Password string    `json:"password"`
}
