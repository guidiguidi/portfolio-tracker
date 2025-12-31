package assets

import (
    "bytes"
    "io"
    "log/slog"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/gin-gonic/gin"
)

func TestCreateAsset(t *testing.T) {
    gin.SetMode(gin.TestMode)

    // логгер для тестов
    log := slog.New(slog.NewTextHandler(io.Discard, nil))

    repo := NewMemoryRepo(log)
    h := NewHandler(repo, log)

    r := gin.New()
    r.POST("/assets", h.CreateAsset)

    body := []byte(`{"symbol":"BTC","name":"Bitcoin"}`)
    req, _ := http.NewRequest(http.MethodPost, "/assets", bytes.NewReader(body))
    req.Header.Set("Content-Type", "application/json")

    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)

    if w.Code != http.StatusCreated {
        t.Fatalf("expected status %d, got %d", http.StatusCreated, w.Code)
    }
}
