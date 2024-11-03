package repository

import (
	"backend/internal/models"
	"context"
	"database/sql"
)

type PostgresGenresRepository struct {
	Db *sql.DB
}

func (g *PostgresGenresRepository) GetAllGenres() ([]*models.Genre, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		SELECT
			id,
			genre,
			created_at,
			updated_at
		FROM
			genres
		ORDER BY
			genre
	`

	rows, err := g.Db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var genres []*models.Genre

	for rows.Next() {
		var genre models.Genre
		err := rows.Scan(
			&genre.Id,
			&genre.Genre,
			&genre.CreatedAt,
			&genre.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		genres = append(genres, &genre)
	}

	return genres, nil
}
