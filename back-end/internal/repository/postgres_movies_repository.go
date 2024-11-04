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
			title;
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

func (m *PostgresMoviesRepository) GetMovieById(id int) (*models.Movie, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

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
		WHERE
			id = $1
	`

	row := m.Db.QueryRowContext(ctx, query, id)

	var movie models.Movie

	err := row.Scan(
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

	query = `
		SELECT
			g.id,
			g.genre
		FROM
			movies_genres mg
		LEFT JOIN
			genres g ON mg.genre_id = g.id
		WHERE
			mg.movie_id = $1
		ORDER BY
			g.genre
	`

	rows, err := m.Db.QueryContext(ctx, query, id)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	defer rows.Close()

	var genres []*models.Genre

	for rows.Next() {
		var g models.Genre

		err := rows.Scan(
			&g.Id,
			&g.Genre,
		)
		if err != nil {
			return nil, err
		}

		genres = append(genres, &g)
	}

	movie.Genres = genres

	return &movie, nil
}

// for admin
func (m *PostgresMoviesRepository) GetMovieByIdForEdit(id int) (*models.Movie, []*models.Genre, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

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
		WHERE
			id = $1
	`

	row := m.Db.QueryRowContext(ctx, query, id)

	var movie models.Movie

	err := row.Scan(
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
		return nil, nil, err
	}

	query = `
		SELECT
			g.id,
			g.genre
		FROM
			movies_genres mg
		LEFT JOIN
			genres g ON mg.genre_id = g.id
		WHERE
			mg.movie_id = $1
		ORDER BY
			g.genre
	`

	rows, err := m.Db.QueryContext(ctx, query, id)
	if err != nil && err != sql.ErrNoRows {
		return nil, nil, err
	}
	defer rows.Close()

	var genres []*models.Genre
	var genresArray []int

	for rows.Next() {
		var g models.Genre

		err := rows.Scan(
			&g.Id,
			&g.Genre,
		)
		if err != nil {
			return nil, nil, err
		}

		genres = append(genres, &g)
		genresArray = append(genresArray, g.Id)
	}

	movie.Genres = genres
	movie.GenresArray = genresArray

	var allGenres []*models.Genre

	query = `
		SELECT
			id,
			genre
		FROM
			genres
		ORDER BY
			genre
	`

	genreRows, err := m.Db.QueryContext(ctx, query)
	if err != nil {
		return nil, nil, err
	}
	defer genreRows.Close()

	for genreRows.Next() {
		var g models.Genre

		err := genreRows.Scan(
			&g.Id,
			&g.Genre,
		)
		if err != nil {
			return nil, nil, err
		}

		allGenres = append(allGenres, &g)
	}

	return &movie, allGenres, nil
}

func (m *PostgresMoviesRepository) InsertMovie(movie models.Movie) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `
		INSERT INTO
			movies(title, description, release_date, runtime, mpaa_rating, created_at, updated_at, image)
		VALUES
			($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING
			id

	`

	var newId int

	err := m.Db.QueryRowContext(
		ctx,
		stmt,
		movie.Title,
		movie.Description,
		movie.ReleaseDate,
		movie.RunTime,
		movie.MpaaRating,
		movie.CreatedAt,
		movie.UpdatedAt,
		movie.Image,
	).Scan(&newId)
	if err != nil {
		return 0, err
	}

	return newId, nil
}

func (m *PostgresMoviesRepository) UpdateMovieGenres(id int, genreIds []int) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `
		DELETE FROM
			movies_genres
		WHERE
			movie_id = $1
	`

	_, err := m.Db.ExecContext(ctx, stmt, id)
	if err != nil {
		return err
	}

	for _, n := range genreIds {
		stmt := `
			INSERT INTO
				movies_genres (movie_id, genre_id)
			VALUES ($1, $2)
		`

		_, err := m.Db.ExecContext(ctx, stmt, id, n)
		if err != nil {
			return err
		}
	}

	return nil
}
