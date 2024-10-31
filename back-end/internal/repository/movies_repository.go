package repository

import (
	"backend/internal/models"
	"database/sql"
)

type MoviesRepository interface {
	Connection() *sql.DB
	GetAllMovies() ([]*models.Movie, error)
}
