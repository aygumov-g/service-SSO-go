package config

import "os"

type Config struct {
	AppPort string
	MainDB  Postgres
}

func Load() *Config {
	return &Config{
		AppPort: os.Getenv("APP_PORT"),
		MainDB: Postgres{
			dBHost:     os.Getenv("POSTGRES_HOST"),
			dBUser:     os.Getenv("POSTGRES_USER"),
			dBPassword: os.Getenv("POSTGRES_PASSWORD"),
			dBName:     os.Getenv("POSTGRES_DB"),
		},
	}
}
