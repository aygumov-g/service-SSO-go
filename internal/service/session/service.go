package session

import (
	"context"
	"time"

	account_d "github.com/aygumov-g/service-SSO-go/internal/domain/account"
	session_d "github.com/aygumov-g/service-SSO-go/internal/domain/session"
)

type Service struct {
	tx           TxManager
	sessionsRepo SessionRepository
	accountsRepo AccountRepository
	tokens       TokenProvider
	clk          Clock
	ttl          time.Duration
}

func NewService(
	tx TxManager,
	sessionsRepo SessionRepository,
	accountsRepo AccountRepository,
	tokens TokenProvider,
	clk Clock,
	ttl time.Duration,
) *Service {
	return &Service{
		tx:           tx,
		sessionsRepo: sessionsRepo,
		accountsRepo: accountsRepo,
		tokens:       tokens,
		clk:          clk,
		ttl:          ttl,
	}
}

func (s *Service) Create(ctx context.Context, account *account_d.Account) (*Output, error) {
	now := s.clk.Now()

	refreshToken, err := s.tokens.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	session := &session_d.Session{
		AccountID: account.ID,
		TokenHash: s.tokens.HashRefreshToken(refreshToken),
		ExpiresAt: now.Add(s.ttl),
		CreatedAt: now,
	}

	if err := s.sessionsRepo.Create(ctx, session); err != nil {
		return nil, err
	}

	accessToken, err := s.tokens.IssueAccessToken(account.ID, account.TokenVersion)
	if err != nil {
		return nil, err
	}

	return &Output{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

func (s *Service) Rotate(ctx context.Context, refreshToken string) (*Output, error) {
	var tokens *Output
	if err := s.tx.Do(ctx, func(txCtx context.Context) error {
		now := s.clk.Now()

		oldHash := s.tokens.HashRefreshToken(refreshToken)
		account_id, err := s.sessionsRepo.GetAccoundIDByHash(txCtx, oldHash)
		if err != nil {
			return err
		}

		account, err := s.accountsRepo.GetByID(txCtx, account_id)
		if err != nil {
			return err
		}

		tokens, err = s.Create(txCtx, account)
		if err != nil {
			return err
		}

		if err := s.sessionsRepo.RevokeByTokenHashStrict(txCtx, oldHash, now); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return tokens, nil
}

func (s *Service) Logout(ctx context.Context, refreshToken string) error {
	now := s.clk.Now()
	hash := s.tokens.HashRefreshToken(refreshToken)

	return s.sessionsRepo.RevokeByTokenHashIfExists(ctx, hash, now)
}

func (s *Service) RevokeAllByAccountID(ctx context.Context, id int64, now time.Time) error {
	return s.sessionsRepo.RevokeAllByAccountID(ctx, id, now)
}
