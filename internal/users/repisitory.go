package users

import (
	"context"
	"errors"
	"first/internal/pkg"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserStatus string

const (
	StatusActive UserStatus = "active"
	StatusFired  UserStatus = "fired"
)

type Repository interface {
	Create(ctx context.Context, user *User) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	TerminateEmployee(ctx context.Context, id int64) (*User, error)
}

type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Create(ctx context.Context, user *User) (*User, error) {
	const query = `
		INSERT INTO users (name, last_name, email, role)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at;
	`

	err := r.db.QueryRow(
		ctx,
		query,
		user.Name,
		user.LastName,
		user.Email,
		user.Role,
	).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok {
			if pgErr.Code == "23505" && pgErr.ConstraintName == "users_email_key" {
				return nil, pkg.ErrUserEmailExists
			}
		}
		return nil, fmt.Errorf("create user: %w", err)
	}

	return user, nil
}

func (r *PostgresRepository) GetByEmail(ctx context.Context, email string) (*User, error) {
	var user User

	const query = `
		SELECT id, name, last_name, email, role, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	err := r.db.QueryRow(
		ctx,
		query,
		email,
	).Scan(
		&user.ID,
		&user.Name,
		&user.LastName,
		&user.Email,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, pkg.ErrUserNotFound
		}
		return nil, fmt.Errorf("get user by email: %w", err)
	}

	return &user, nil
}

func (r *PostgresRepository) TerminateEmployee(ctx context.Context, id int64) (*User, error) {
	var user User

	const query = `
		UPDATE users
		SET status = $1,
			updated_at = NOW()
		WHERE id = $2
		RETURNING id, name, last_name, email, role, created_at, updated_at, status;
	`

	err := r.db.QueryRow(
		ctx,
		query,
		StatusFired,
		id,
	).Scan(
		&user.ID,
		&user.Name,
		&user.LastName,
		&user.Email,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, pkg.ErrUserNotFound
		}
		return nil, fmt.Errorf("terminate employee: %w", err)
	}

	return &user, nil
}
