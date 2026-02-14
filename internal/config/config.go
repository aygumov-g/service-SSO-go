package config

import (
	"os"
)

type Config struct {
	AppPort string
	DB      Postgres
	JWT     JWT
}

func Load() *Config {
	return &Config{
		AppPort: os.Getenv("APP_PORT"),
		DB: Postgres{
			dBHost:     os.Getenv("POSTGRES_HOST"),
			dBUser:     os.Getenv("POSTGRES_USER"),
			dBPassword: os.Getenv("POSTGRES_PASSWORD"),
			dBName:     os.Getenv("POSTGRES_DB"),
		},
		JWT: JWT{
			Secret: []byte(os.Getenv("JWT_SECRET")),
			TTL:    TTL(os.Getenv("JWT_TTL")),
		},
	}
}
