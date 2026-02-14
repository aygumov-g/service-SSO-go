package jwt

import (
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

func (j *JWTService) Issue(identity *d_identity.Identity) (string, error) {
	claims := jwt.MapClaims{
		"sub": identity.ID,
		"exp": j.clk.Now().Add(j.ttl).Unix(),
		"iat": j.clk.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secret)
}

func (j *JWTService) Parse(tokenStr string) (*d_identity.Identity, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
		return j.secret, nil
	})
	if err != nil {
		return nil, err
	}

	claims := token.Claims.(jwt.MapClaims)
	identity := d_identity.Identity{
		ID: int64(claims["sub"].(float64)),
	}

	return &identity, nil
}
