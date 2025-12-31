package assets

import (
	"context"
	"errors"
	"log/slog"
	"sync"
)

var ErrNotFound = errors.New("asset not found")

type MemoryRepo struct {
	mu     sync.RWMutex
	lastID int64
	items  map[int64]Asset
	log    *slog.Logger
}

func NewMemoryRepo(log *slog.Logger) *MemoryRepo {
	return &MemoryRepo{
		items: make(map[int64]Asset),
		log:   log.With(slog.String("component", "memory_repo")),
	}
}

func (r *MemoryRepo) Create(ctx context.Context, a Asset) (Asset, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.lastID++
	a.ID = r.lastID
	r.items[a.ID] = a

	r.log.Debug("created asset in memory", slog.Int64("id", a.ID), slog.String("symbol", a.Symbol))

	return a, nil
}

func (r *MemoryRepo) GetByID(ctx context.Context, id int64) (Asset, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	a, ok := r.items[id]
	if !ok {
		return Asset{}, ErrNotFound
	}
	return a, nil
}

func (r *MemoryRepo) List(ctx context.Context) ([]Asset, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	res := make([]Asset, 0, len(r.items))
	for _, a := range r.items {
		res = append(res, a)
	}
	return res, nil
}
