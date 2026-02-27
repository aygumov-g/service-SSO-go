package token

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"strconv"
	"time"

	identity_d "github.com/aygumov-g/service-SSO-go/internal/domain/identity"

	"github.com/golang-jwt/jwt/v5"
)

type JWTProvider struct {
	secret []byte
	ttl    time.Duration
	clk    Clock
}

func NewJWTProvider(secret string, ttl time.Duration, clk Clock) *JWTProvider {
	return &JWTProvider{
		secret: []byte(secret),
		ttl:    ttl,
		clk:    clk,
	}
}

func (j *JWTProvider) IssueAccessToken(accountID int64, tokenVersion int) (string, error) {
	now := j.clk.Now()

	claims := Claims{
		TokenVersion: tokenVersion,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.FormatInt(accountID, 10),
			ExpiresAt: jwt.NewNumericDate(now.Add(j.ttl)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secret)
}

func (j *JWTProvider) Parse(tokenStr string) (*identity_d.Identity, error) {
	token, err := jwt.ParseWithClaims(
		tokenStr,
		&Claims{},
		func(token *jwt.Token) (any, error) {
			if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
				return nil, ErrInvalidSigningMethod
			}

			return j.secret, nil
		},
		jwt.WithTimeFunc(j.clk.Now),
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

	identity := identity_d.Identity{
		ID:           accountID,
		TokenVersion: claims.TokenVersion,
	}

	return &identity, nil
}

func (j *JWTProvider) GenerateRefreshToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(b), nil
}

func (j *JWTProvider) HashRefreshToken(t string) string {
	hash := sha256.Sum256([]byte(t))
	return base64.RawURLEncoding.EncodeToString(hash[:])
}
