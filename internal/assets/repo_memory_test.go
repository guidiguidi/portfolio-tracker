package assets

import (
    "context"
    "io"
    "log/slog"
    "testing"
)

func TestMemoryRepo_CreateAndGetByID(t *testing.T) {
    log := slog.New(slog.NewTextHandler(io.Discard, nil))
    repo := NewMemoryRepo(log)
    ctx := context.Background()

    created, err := repo.Create(ctx, Asset{
        Symbol: "BTC",
        Name:   "Bitcoin",
    })
    if err != nil {
        t.Fatalf("Create returned error: %v", err)
    }
    if created.ID == 0 {
        t.Fatalf("expected non-zero ID")
    }

    got, err := repo.GetByID(ctx, created.ID)
    if err != nil {
        t.Fatalf("GetByID returned error: %v", err)
    }

    if got.Symbol != "BTC" || got.Name != "Bitcoin" {
        t.Fatalf("unexpected asset: %+v", got)
    }
}
