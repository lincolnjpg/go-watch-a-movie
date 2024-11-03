package repository

import "backend/internal/models"

type GenreRepository interface {
	GetAllGenres() ([]*models.Genre, error)
}
