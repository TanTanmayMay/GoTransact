package repository

import (
	"rest1/internal/repository"
)

type UserRepository struct {
	Storage *[]domain.User
}

