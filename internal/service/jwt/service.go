package jwt

import (
	"strconv"
	"time"

	d_identity "github.com/aygumov-g/service-SSO-go/internal/domain/identity"

	"github.com/golang-jwt/jwt/v5"
)

type JWTService struct {
	secret []byte
	ttl    time.Duration
	clk    Clock
}

func NewJWTService(secret []byte, ttl time.Duration, clk Clock) *JWTService {
	return &JWTService{
		secret: secret,
		ttl:    ttl,
		clk:    clk,
	}
}

func (s *JWTService) Issue(accountID int64) (string, error) {
	now := s.clk.Now()

	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.FormatInt(accountID, 10),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.ttl)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secret)
}

func (s *JWTService) Parse(tokenStr string) (*d_identity.Identity, error) {
	token, err := jwt.ParseWithClaims(
		tokenStr,
		&Claims{},
		func(token *jwt.Token) (any, error) {
			if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
				return nil, ErrInvalidSigningMethod
			}

			return s.secret, nil
		},
		jwt.WithTimeFunc(s.clk.Now),
	)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	accountID, err := strconv.ParseInt(claims.Subject, 10, 64)
	if err != nil {
		return nil, ErrInvalidToken
	}

	identity := d_identity.Identity{
		ID: accountID,
	}

	return &identity, nil
}
