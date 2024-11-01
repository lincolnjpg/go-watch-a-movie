package repository

import "backend/internal/models"

type UserRepository interface {
	GetUserByEmail(string) (*models.User, error)
	GetUserById(int) (*models.User, error)
}
