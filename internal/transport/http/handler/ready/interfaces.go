package ready

import "context"

type ReadyUsecase interface {
	Execute(ctx context.Context) bool
}
