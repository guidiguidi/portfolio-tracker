package config

import (
    "github.com/spf13/viper"
)

type Config struct {
    Logger    LoggerConfig    `mapstructure:"logger"`
    Database  DatabaseConfig  `mapstructure:"database"`
    JWT       JWTConfig       `mapstructure:"jwt"`
    Coingecko CoingeckoConfig `mapstructure:"coingecko"`
    Redis     RedisConfig     `mapstructure:"redis"`
    App       AppConfig       `mapstructure:"app"`
}

type LoggerConfig struct {
    Level string `mapstructure:"level"`
    JSON  bool   `mapstructure:"json"`
}

type DatabaseConfig struct {
    Driver          string `mapstructure:"driver"`
    DSN             string `mapstructure:"dsn"`
    MaxOpenConns    int    `mapstructure:"max_open_conns"`
    MaxIdleConns    int    `mapstructure:"max_idle_conns"`
    ConnMaxLifetime string `mapstructure:"conn_max_lifetime"`
}


type JWTConfig struct {
    Secret         string `mapstructure:"secret"`
    AccessTokenTTL string `mapstructure:"access_token_ttl"`
    RefreshTokenTTL string `mapstructure:"refresh_token_ttl"`
}

type CoingeckoConfig struct {
    BaseURL        string `mapstructure:"base_url"`
    RequestTimeout string `mapstructure:"request_timeout"`
}

type RedisConfig struct {
    Enabled bool `mapstructure:"enabled"`
}

type AppConfig struct {
    Version         string `mapstructure:"version"`
    ShutdownTimeout string `mapstructure:"shutdown_timeout"`
    Port            string `mapstructure:"port"`
}


func LoadConfig(path string) (*Config, error) {
    v := viper.New()
    v.SetConfigFile(path)

    if err := v.ReadInConfig(); err != nil {
        return nil, err
    }

    var cfg Config
    if err := v.Unmarshal(&cfg); err != nil { 
        return nil, err
    }

    return &cfg, nil
}