package session

import (
	"context"
	"errors"
	"time"

	account_d "github.com/aygumov-g/service-SSO-go/internal/domain/account"
	session_d "github.com/aygumov-g/service-SSO-go/internal/domain/session"
)

type Service struct {
	sessionsRepo SessionRepository
	tokens       TokenProvider
	clk          Clock
	ttl          time.Duration
}

func NewService(
	sessionsRepo SessionRepository,
	tokens TokenProvider,
	clk Clock,
	ttl time.Duration,
) *Service {
	return &Service{
		sessionsRepo: sessionsRepo,
		tokens:       tokens,
		clk:          clk,
		ttl:          ttl,
	}
}

func (s *Service) Create(ctx context.Context, account *account_d.Account) (
	string,
	string,
	error,
) {
	now := s.clk.Now()

	refreshToken, err := s.tokens.GenerateRefreshToken()
	if err != nil {
		return "", "", err
	}

	session := &session_d.Session{
		AccountID: account.ID,
		TokenHash: s.tokens.HashRefreshToken(refreshToken),
		ExpiresAt: now.Add(s.ttl),
		CreatedAt: now,
	}

	if err := s.sessionsRepo.Create(ctx, session); err != nil {
		return "", "", err
	}

	accessToken, err := s.tokens.IssueAccessToken(account.ID, account.TokenVersion)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *Service) Rotate(ctx context.Context, refreshToken string) (
	int64,
	string,
	string,
	error,
) {
	now := s.clk.Now()

	oldHash := s.tokens.HashRefreshToken(refreshToken)

	newRefreshToken, err := s.tokens.GenerateRefreshToken()
	if err != nil {
		return 0, "", "", err
	}

	newSession := &session_d.Session{
		TokenHash: s.tokens.HashRefreshToken(newRefreshToken),
		ExpiresAt: now.Add(s.ttl),
		CreatedAt: now,
	}

	accountID, tokenVersion, err := s.sessionsRepo.RotateByTokenHash(
		ctx,
		oldHash,
		newSession,
		now,
	)
	if err != nil {
		if errors.Is(err, session_d.ErrNotFound) ||
			errors.Is(err, session_d.ErrExpired) ||
			errors.Is(err, session_d.ErrRevoked) {

			return 0, "", "", ErrInvalidRefreshToken
		}
	}

	accessToken, err := s.tokens.IssueAccessToken(accountID, tokenVersion)
	if err != nil {
		return 0, "", "", err
	}

	return accountID, accessToken, newRefreshToken, nil
}

func (s *Service) RevokeAllByAccountID(ctx context.Context, id int64, now time.Time) error {
	return s.sessionsRepo.RevokeAllByAccountID(ctx, id, now)
}
