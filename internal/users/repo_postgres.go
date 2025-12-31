package users

import (
	"context"
	"database/sql"
	"log/slog"
	"strings"

	"github.com/lib/pq"
)

type PostgresRepo struct {
	db  *sql.DB
	log *slog.Logger
}

func NewPostgresRepo(db *sql.DB, log *slog.Logger) *PostgresRepo {
	return &PostgresRepo{
		db:  db,
		log: log.With(slog.String("component", "users_postgres_repo")),
	}
}

func (r *PostgresRepo) Create(ctx context.Context, user User) (User, error) {
	query := `
		INSERT INTO users (email, password_hash) 
		VALUES ($1, $2) 
		RETURNING id, created_at, updated_at`

	args := []any{user.Email, user.PasswordHash}

	r.log.Debug("executing query", slog.String("query", query))
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		// Check for unique violation
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code.Name() == "unique_violation" {
			if strings.Contains(pqErr.Message, "users_email_key") {
				return User{}, ErrDuplicateEmail
			}
		}
		r.log.Error("failed to create user", "error", err)
		return User{}, err
	}

	return user, nil
}

func (r *PostgresRepo) GetByEmail(ctx context.Context, email string) (User, error) {
	var u User
	query := `
		SELECT id, email, password_hash, created_at, updated_at 
		FROM users WHERE email = $1`
	
	r.log.Debug("executing query", slog.String("query", query))
	err := r.db.QueryRowContext(ctx, query, email).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return User{}, ErrNotFound
		}
		r.log.Error("failed to get user by email", "error", err)
		return User{}, err
	}

	return u, nil
}

func (r *PostgresRepo) GetByID(ctx context.Context, id int64) (User, error) {
	var u User
	query := `
		SELECT id, email, password_hash, created_at, updated_at 
		FROM users WHERE id = $1`

	r.log.Debug("executing query", slog.String("query", query))
	err := r.db.QueryRowContext(ctx, query, id).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return User{}, ErrNotFound
		}
		r.log.Error("failed to get user by id", "error", err)
		return User{}, err
	}

	return u, nil
}
