package health

import (
	"context"
	"time"

	"github.com/alexliesenfeld/health"
)

type Checker struct {
	checker health.Checker
}

func New(db Postgres) *Checker {
	return &Checker{
		checker: health.NewChecker(
			health.WithCacheDuration(1*time.Second),
			health.WithTimeout(10*time.Second),
			health.WithChecks(
				health.Check{
					Name:    "postgres",
					Timeout: 2 * time.Second,
					Check: func(ctx context.Context) error {
						return db.Ping(ctx)
					},
				},
			),
		),
	}
}

func (c *Checker) Check(ctx context.Context) bool {
	return c.checker.Check(ctx).Status == health.StatusUp
}
