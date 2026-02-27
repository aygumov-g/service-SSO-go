package password

import "golang.org/x/crypto/bcrypt"

type BcryptHasher struct {
	cost int
}

func NewBcryptHasher(cost int) *BcryptHasher {
	return &BcryptHasher{cost: cost}
}

func (b *BcryptHasher) Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword(
		[]byte(password),
		b.cost,
	)
}

func (b *BcryptHasher) CompareHash(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword(
		[]byte(hashedPassword),
		[]byte(password),
	)
}

func (b *BcryptHasher) FakeCompareHash(password string) error {
	return bcrypt.CompareHashAndPassword(
		[]byte("$2a$12$C6UzMDM.H6dfI/f/IKcEeO5Gk6tYd6L9uP0u9e6QnqK9lE9iFQ6eK"),
		[]byte(password),
	)
}
