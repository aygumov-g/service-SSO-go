package session

import (
	"context"
	"time"

	account_d "github.com/aygumov-g/service-SSO-go/internal/domain/account"
	session_d "github.com/aygumov-g/service-SSO-go/internal/domain/session"
)

type Service struct {
	sessions SessionRepository
	jwt      JWTService
	ttl      time.Duration
	clk      Clock
}

func NewService(sessions SessionRepository, jwt JWTService, ttl time.Duration, clk Clock) *Service {
	return &Service{
		sessions: sessions,
		jwt:      jwt,
		ttl:      ttl,
		clk:      clk,
	}
}

func (s *Service) Create(ctx context.Context, account *account_d.Account) (*TokenPair, error) {
	now := s.clk.Now()

	session, refreshToken, err := s.buildSession(account.ID, now)
	if err != nil {
		return nil, err
	}

	if err := s.sessions.Create(ctx, session); err != nil {
		return nil, err
	}

	accessToken, err := s.issueAccess(account.ID)
	if err != nil {
		return nil, err
	}

	return &TokenPair{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

func (s *Service) Refresh(ctx context.Context, refreshToken string) (*TokenPair, error) {
	now := s.clk.Now()

	newSession, newRefreshToken, err := s.buildSession(0, now)
	if err != nil {
		return nil, err
	}

	accountID, err := s.sessions.RotateByTokenHash(ctx, hashToken(refreshToken), newSession, now)
	if err != nil {
		return nil, err
	}

	accessToken, err := s.issueAccess(accountID)
	if err != nil {
		return nil, err
	}

	return &TokenPair{AccessToken: accessToken, RefreshToken: newRefreshToken}, nil
}

func (s *Service) RevokeAllByAccountID(ctx context.Context, id int64, now time.Time) error {
	return s.sessions.RevokeAllByAccountID(ctx, id, now)
}

func (s *Service) buildSession(accountID int64, now time.Time) (*session_d.Session, string, error) {
	refreshToken, err := generateRefreshToken()
	if err != nil {
		return nil, "", err
	}

	hash := hashToken(refreshToken)

	session := &session_d.Session{
		AccountID: accountID,
		TokenHash: hash,
		ExpiresAt: now.Add(s.ttl),
		CreatedAt: now,
	}

	return session, refreshToken, nil
}

func (s *Service) issueAccess(accountID int64) (string, error) {
	return s.jwt.Issue(accountID)
}
