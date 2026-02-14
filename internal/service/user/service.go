package user

import (
	"context"
	"errors"

	d_identity "github.com/aygumov-g/service-SSO-go/internal/domain/identity"
	d_user "github.com/aygumov-g/service-SSO-go/internal/domain/user"

	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	users UserRepository
	jwt   JWTService
}

func NewService(users UserRepository, jwt JWTService) *Service {
	return &Service{
		users: users,
		jwt:   jwt,
	}
}

func (s *Service) Login(ctx context.Context, login, password string) (string, error) {
	user, err := s.users.GetByLogin(ctx, login)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return "", ErrInvalidCredentials
		}

		return "", err
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(password),
	); err != nil {
		return "", ErrInvalidCredentials
	}

	identity := d_identity.Identity{
		ID: user.ID,
	}

	return s.jwt.Issue(&identity)
}

func (s *Service) GetByID(ctx context.Context, id int64) (*d_user.User, error) {
	return s.users.GetByID(ctx, id)
}

func (s *Service) Register(ctx context.Context, login, password string) error {
	_, err := s.users.GetByLogin(ctx, login)
	if err == nil {
		return ErrUserAlreadyExists
	}

	if !errors.Is(err, ErrUserNotFound) {
		return err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &d_user.User{
		Login:        login,
		PasswordHash: string(hash),
	}

	return s.users.Create(ctx, user)
}
