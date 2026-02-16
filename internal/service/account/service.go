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

func (s *Service) Login(ctx context.Context, login, password string) (string, error) {
	account, err := s.accounts.GetByLogin(ctx, login)
	if err != nil {
		if errors.Is(err, ErrAccountNotFound) {
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

func (s *Service) GetByID(ctx context.Context, id int64) (*d_account.Account, error) {
	return s.accounts.GetByID(ctx, id)
}

func (s *Service) Register(ctx context.Context, login, password string) error {
	_, err := s.accounts.GetByLogin(ctx, login)
	if err == nil {
		return ErrAccountAlreadyExists
	}

	if !errors.Is(err, ErrAccountNotFound) {
		return err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	account := &d_account.Account{
		Login:        login,
		PasswordHash: string(hash),
		Role:         "user",
		CreatedAt:    s.clk.Now(),
		UpdatedAt:    s.clk.Now(),
	}

	return s.accounts.Create(ctx, account)
}
