package repository

import (
	"backend/internal/models"
	"context"
	"database/sql"
	"time"
)

type PostgresMoviesRepository struct {
	Db *sql.DB
}

const dbTimeout = time.Second * 3

func (m *PostgresMoviesRepository) Connection() *sql.DB {
	return m.Db
}

func (m *PostgresMoviesRepository) GetAllMovies() ([]*models.Movie, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var movies []*models.Movie

	query := `
		SELECT
			id,
			title,
			release_date,
			runtime,
			mpaa_rating,
			description,
			coalesce(image, ''),
			created_at,
			updated_at
		FROM
			movies
		ORDER BY
			title
	`

	rows, err := m.Db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var movie models.Movie

		err := rows.Scan(
			&movie.Id,
			&movie.Title,
			&movie.ReleaseDate,
			&movie.RunTime,
			&movie.MpaaRating,
			&movie.Description,
			&movie.Image,
			&movie.CreatedAt,
			&movie.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		movies = append(movies, &movie)
	}

	return movies, nil
}
