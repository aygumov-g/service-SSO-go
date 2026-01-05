package config

import "time"

type JWT struct {
	Secret []byte
	TTL    time.Duration
}

func TTL(ttl_str string) time.Duration {
	ttl, _ := time.ParseDuration(ttl_str)

	return ttl
}
