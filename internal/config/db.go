package config

import (
    "database/sql"
    "log"
    "os"
    "time"

    _ "github.com/lib/pq"
)

func NewDB(cfg *Config) (*sql.DB, error) {
    dsn := cfg.Database.DSN
    if env := os.Getenv("DB_DSN"); env != "" {
        dsn = env
    }

    log.Printf("DB DSN: %s", dsn)

    db, err := sql.Open("postgres", dsn)
    if err != nil {
        return nil, err
    }

    db.SetMaxOpenConns(cfg.Database.MaxOpenConns)
    db.SetMaxIdleConns(cfg.Database.MaxIdleConns)
    if dur, err := time.ParseDuration(cfg.Database.ConnMaxLifetime); err == nil {
        db.SetConnMaxLifetime(dur)
    }

    if err := db.Ping(); err != nil {
        return nil, err
    }

    return db, nil
}
