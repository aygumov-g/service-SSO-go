package account

import (
	"context"
	"errors"

	account_d "github.com/aygumov-g/service-SSO-go/internal/domain/account"
	identity_d "github.com/aygumov-g/service-SSO-go/internal/domain/identity"
	session_srv "github.com/aygumov-g/service-SSO-go/internal/service/session"

	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	accountsRepo AccountRepository
	sessions     SessionService
	clk          Clock
}

func NewService(
	accountsRepo AccountRepository,
	sessions SessionService,
	clk Clock,
) *Service {
	return &Service{
		accountsRepo: accountsRepo,
		sessions:     sessions,
		clk:          clk,
	}
}

func (s *Service) Register(ctx context.Context, login, password string) error {
	hash, err := getHash(password)
	if err != nil {
		return err
	}

	now := s.clk.Now()

	account := &account_d.Account{
		Login:        login,
		PasswordHash: string(hash),
		TokenVersion: 0,
		Role:         "user",
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	return s.accountsRepo.Create(ctx, account)
}

func (s *Service) Login(ctx context.Context, login, password string) (*session_srv.TokenPair, error) {
	account, err := s.accountsRepo.GetByLogin(ctx, login)
	if err != nil {
		if errors.Is(err, ErrAccountNotFound) {
			_ = bcrypt.CompareHashAndPassword(
				[]byte("$2a$12$C6UzMDM.H6dfI/f/IKcEeO5Gk6tYd6L9uP0u9e6QnqK9lE9iFQ6eK"),
				[]byte(password),
			)

			return nil, ErrInvalidCredentials
		}

		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(account.PasswordHash),
		[]byte(password),
	); err != nil {
		return nil, ErrInvalidCredentials
	}

	return s.sessions.Create(ctx, account)
}

func (s *Service) ChangePassword(ctx context.Context, id int64, oldPassword, newPassword string) error {
	account, err := s.GetByID(ctx, id)
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

	hash, err := getHash(newPassword)
	if err != nil {
		return err
	}

	now := s.clk.Now()

	account.TokenVersion++
	account.PasswordHash = string(hash)
	account.UpdatedAt = now

	if err := s.accountsRepo.Update(ctx, account); err != nil {
		return err
	}

	return s.sessions.RevokeAllByAccountID(ctx, account.ID, now)
}

func (s *Service) ValidateTokenVersion(ctx context.Context, identity *identity_d.Identity) error {
	account, err := s.GetByID(ctx, identity.ID)
	if err != nil {
		return ErrInvalidCredentials
	}

	if account.TokenVersion != identity.TokenVersion {
		return ErrInvalidCredentials
	}

	return nil
}

func (s *Service) GetByID(ctx context.Context, id int64) (*account_d.Account, error) {
	return s.accountsRepo.GetByID(ctx, id)
}
