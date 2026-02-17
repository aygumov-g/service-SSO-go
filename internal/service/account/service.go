package account

import (
	"context"
	"errors"

	d_account "github.com/aygumov-g/service-SSO-go/internal/domain/account"
	d_identity "github.com/aygumov-g/service-SSO-go/internal/domain/identity"

	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	accounts AccountRepository
	jwt      JWTService
	clk      Clock
}

func NewService(accounts AccountRepository, jwt JWTService, clk Clock) *Service {
	return &Service{
		accounts: accounts,
		jwt:      jwt,
		clk:      clk,
	}
}

func (s *Service) Register(ctx context.Context, login, password string) error {
	hash, err := s.getHash(password)
	if err != nil {
		return err
	}

	now := s.clk.Now()

	account := &d_account.Account{
		Login:        login,
		PasswordHash: string(hash),
		Role:         "user",
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	return s.accounts.Create(ctx, account)
}

func (s *Service) Login(ctx context.Context, login, password string) (string, error) {
	account, err := s.accounts.GetByLogin(ctx, login)
	if err != nil {
		if errors.Is(err, ErrAccountNotFound) {
			_ = bcrypt.CompareHashAndPassword(
				[]byte("$2a$12$C6UzMDM.H6dfI/f/IKcEeO5Gk6tYd6L9uP0u9e6QnqK9lE9iFQ6eK"),
				[]byte(password),
			)

			return "", ErrInvalidCredentials
		}

		return "", err
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(account.PasswordHash),
		[]byte(password),
	); err != nil {
		return "", ErrInvalidCredentials
	}

	identity := d_identity.Identity{
		ID: account.ID,
	}

	return s.jwt.Issue(&identity)
}

func (s *Service) ChangePassword(ctx context.Context, id int64, oldPassword, newPassword string) error {
	account, err := s.accounts.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(account.PasswordHash),
		[]byte(oldPassword),
	); err != nil {
		return ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(account.PasswordHash),
		[]byte(newPassword),
	); err == nil {
		return ErrSamePassword
	}

	hash, err := s.getHash(newPassword)
	if err != nil {
		return err
	}

	account.PasswordHash = string(hash)
	account.UpdatedAt = s.clk.Now()

	return s.accounts.Update(ctx, account)
}

func (s *Service) GetByID(ctx context.Context, id int64) (*d_account.Account, error) {
	return s.accounts.GetByID(ctx, id)
}

func (s *Service) getHash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), 12)
}
