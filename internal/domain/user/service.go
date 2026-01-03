package user

import (
	"context"
	"math/rand"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateRandom(ctx context.Context) (User, error) {
	u := User{
		ID:   rand.Int63n(1_000_000_000),
		Name: "Test User",
	}

	return s.repo.Create(ctx, u)
}
