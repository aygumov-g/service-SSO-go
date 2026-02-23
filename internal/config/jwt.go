package config

import "time"

type JWT struct {
	Secret []byte
	TTL    time.Duration
}
