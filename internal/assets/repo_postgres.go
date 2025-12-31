package assets

import (
	"context"
	"database/sql"
	"log/slog"
)

type PostgresRepo struct {
	db  *sql.DB
	log *slog.Logger
}

func NewPostgresRepo(db *sql.DB, log *slog.Logger) *PostgresRepo {
	return &PostgresRepo{
		db:  db,
		log: log.With(slog.String("component", "postgres_repo")),
	}
}

func (r *PostgresRepo) Create(ctx context.Context, a Asset) (Asset, error) {
	query := `INSERT INTO assets (symbol, name) VALUES ($1, $2) RETURNING id`
	r.log.Debug("executing query", slog.String("query", query))

	err := r.db.QueryRowContext(ctx, query, a.Symbol, a.Name).Scan(&a.ID)
	if err != nil {
		r.log.Error("failed to create asset", "error", err)
		return Asset{}, err
	}
	return a, nil
}

func (r *PostgresRepo) GetByID(ctx context.Context, id int64) (Asset, error) {
	var a Asset
	query := `SELECT id, symbol, name FROM assets WHERE id = $1`
	r.log.Debug("executing query", slog.String("query", query), slog.Int64("id", id))

	err := r.db.QueryRowContext(ctx, query, id).Scan(&a.ID, &a.Symbol, &a.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			r.log.Warn("asset not found", slog.Int64("id", id), "error", err)
			return Asset{}, ErrNotFound
		}
		r.log.Error("failed to get asset by id", slog.Int64("id", id), "error", err)
		return Asset{}, err
	}
	return a, nil
}

func (r *PostgresRepo) List(ctx context.Context) ([]Asset, error) {
    query := `SELECT id, symbol, name FROM assets ORDER BY id`
    r.log.Debug("executing query", slog.String("query", query))

    rows, err := r.db.QueryContext(ctx, query)
    if err != nil {
        r.log.Error("failed to list assets", "error", err)
        return nil, err
    }
    defer rows.Close()

    var res []Asset
    for rows.Next() {
        var a Asset
        if err := rows.Scan(&a.ID, &a.Symbol, &a.Name); err != nil {
            r.log.Error("failed to scan asset row", "error", err)
            return nil, err
        }
        res = append(res, a)
    }
    if err = rows.Err(); err != nil {
        r.log.Error("error iterating asset rows", "error", err)
        return nil, err
    }

    return res, nil
}