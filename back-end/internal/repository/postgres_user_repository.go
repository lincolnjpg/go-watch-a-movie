package repository

import (
	"backend/internal/models"
	"context"
	"database/sql"
)

type PostgresUserRepository struct {
	Db *sql.DB
}

func (u *PostgresUserRepository) GetUserByEmail(email string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		SELECT
			id,
			email,
			first_name,
			last_name,
			password,
			created_at,
			updated_at
		FROM
			users
		WHERE
			email = $1;
	`

	var user models.User
	row := u.Db.QueryRowContext(ctx, query, email)

	err := row.Scan(
		&user.Id,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return &models.User{}, err
	}

	return &user, nil
}

func (u *PostgresUserRepository) GetUserById(id int) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		SELECT
			id,
			email,
			first_name,
			last_name,
			password,
			created_at,
			updated_at
		FROM
			users
		WHERE
			id = $1;
	`

	var user models.User
	row := u.Db.QueryRowContext(ctx, query, id)

	err := row.Scan(
		&user.Id,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return &models.User{}, err
	}

	return &user, nil
}
