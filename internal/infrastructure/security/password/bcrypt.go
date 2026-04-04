package password

import "golang.org/x/crypto/bcrypt"

type bcryptHasher struct {
	Cost int
}

func NewBcryptHasher(cost int) *bcryptHasher {
	return &bcryptHasher{Cost: cost}
}

func (b *bcryptHasher) Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword(
		[]byte(password),
		b.Cost,
	)
}

func (b *bcryptHasher) CompareHash(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword(
		[]byte(hashedPassword),
		[]byte(password),
	)
}

func (b *bcryptHasher) FakeCompareHash(password string) error {
	return bcrypt.CompareHashAndPassword(
		[]byte("$2a$12$C6UzMDM.H6dfI/f/IKcEeO5Gk6tYd6L9uP0u9e6QnqK9lE9iFQ6eK"),
		[]byte(password),
	)
}
