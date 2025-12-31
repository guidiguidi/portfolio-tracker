package assets

import "context"

type Repository interface {
    Create(ctx context.Context, a Asset) (Asset, error)
    GetByID(ctx context.Context, id int64) (Asset, error)
    List(ctx context.Context) ([]Asset, error)
}

