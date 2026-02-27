package config

import "time"

type JWT struct {
	Secret string
	TTL    time.Duration
}
