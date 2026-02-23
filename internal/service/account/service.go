package account

import (
	"context"
	"errors"

	d_account "github.com/aygumov-g/service-SSO-go/internal/domain/account"
	srv_session "github.com/aygumov-g/service-SSO-go/internal/service/session"

	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	accounts AccountRepository
	sessions SessionService
	clk      Clock
}

func NewService(
	accounts AccountRepository,
	sessions SessionService,
	clk Clock,
) *Service {
	return &Service{
		accounts: accounts,
		sessions: sessions,
		clk:      clk,
	}
}

func (s *Service) Register(ctx context.Context, login, password string) error {
	hash, err := getHash(password)
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

func (s *Service) Login(ctx context.Context, login, password string) (*srv_session.TokenPair, error) {
	account, err := s.accounts.GetByLogin(ctx, login)
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

	hash, err := getHash(newPassword)
	if err != nil {
		return err
	}

	now := s.clk.Now()

	account.PasswordHash = string(hash)
	account.UpdatedAt = now

	if err := s.accounts.Update(ctx, account); err != nil {
		return err
	}

	return s.sessions.RevokeAllByAccountID(ctx, account.ID, now)
}

func (s *Service) GetByID(ctx context.Context, id int64) (*d_account.Account, error) {
	return s.accounts.GetByID(ctx, id)
}
