package config

import "fmt"

type Postgres struct {
	dBHost     string
	dBUser     string
	dBPassword string
	dBName     string
}

func (p Postgres) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		p.dBUser,
		p.dBPassword,
		p.dBHost,
		p.dBName)
}
