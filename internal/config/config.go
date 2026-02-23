package config

import (
	"os"
	"time"
)

type Config struct {
	App     App
	DB      Postgres
	Refresh Refresh
	JWT     JWT
}

func Load() *Config {
	return &Config{
		App: App{
			Port: os.Getenv("APP_PORT"),
		},
		DB: Postgres{
			dBHost:     os.Getenv("POSTGRES_HOST"),
			dBUser:     os.Getenv("POSTGRES_USER"),
			dBPassword: os.Getenv("POSTGRES_PASSWORD"),
			dBName:     os.Getenv("POSTGRES_DB"),
		},
		JWT: JWT{
			Secret: []byte(os.Getenv("JWT_SECRET")),
			TTL:    ttl(os.Getenv("JWT_TTL")),
		},
		Refresh: Refresh{
			TTL: ttl(os.Getenv("REFRESH_TTL")),
		},
	}
}

func ttl(ttlStr string) time.Duration {
	ttl, err := time.ParseDuration(ttlStr)
	if err != nil {
		panic(err)
	}

	return ttl
}
