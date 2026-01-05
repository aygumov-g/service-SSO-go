package auth

import (
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           int64
	Login        string
	PasswordHash string
}

type Service struct {
	users  UserReaderByLogin
	tokens TokenManager
}

func NewService(users UserReaderByLogin, tokens TokenManager) *Service {
	return &Service{
		users:  users,
		tokens: tokens,
	}
}

func (s *Service) Login(ctx context.Context, login, password string) (string, error) {
	user, err := s.users.GetByLogin(ctx, login)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return "", ErrInvalidCredentials
		}

		return "", ErrInternal
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(password),
	); err != nil {
		return "", ErrInvalidCredentials
	}

	return s.tokens.Issue(user.ID)
}

func (s *Service) Register(ctx context.Context, login, password string) error {
	_, err := s.users.GetByLogin(ctx, login)
	if err == nil {
		return ErrUserAlreadyExists
	}

	if !errors.Is(err, ErrUserNotFound) {
		return ErrInternal
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return ErrInternal
	}

	user := User{
		Login:        login,
		PasswordHash: string(hash),
	}

	if err := s.users.(UserCreator).Create(ctx, user); err != nil {
		return ErrInternal
	}

	return nil
}
