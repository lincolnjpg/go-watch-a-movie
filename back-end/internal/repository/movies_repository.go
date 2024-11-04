package repository

import (
	"backend/internal/models"
	"database/sql"
)

type MoviesRepository interface {
	Connection() *sql.DB
	GetAllMovies() ([]*models.Movie, error)
	GetMovieById(id int) (*models.Movie, error)
	GetMovieByIdForEdit(id int) (*models.Movie, []*models.Genre, error)
	InsertMovie(movie models.Movie) (int, error)
	UpdateMovieGenres(id int, genreIds []int) error
}
