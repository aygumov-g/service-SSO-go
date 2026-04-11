package ready

import (
	"context"
)

type Usecase struct {
	checker HealthChecker
}

func NewUsecase(checker HealthChecker) *Usecase {
	return &Usecase{
		checker: checker,
	}
}

func (uc *Usecase) Execute(ctx context.Context) bool {
	return uc.checker.Check(ctx)
}
