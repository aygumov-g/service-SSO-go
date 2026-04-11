package ready

import "context"

type HealthChecker interface {
	Check(ctx context.Context) bool
}
