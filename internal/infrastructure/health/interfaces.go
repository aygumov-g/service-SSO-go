package health

import "context"

type Postgres interface {
	Ping(ctx context.Context) error
}
