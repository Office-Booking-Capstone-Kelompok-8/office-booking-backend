package password

import "golang.org/x/crypto/bcrypt"

type Hash interface {
	GenerateFromPassword(password []byte, cost int) ([]byte, error)
	CompareHashAndPassword(hashedPassword, password []byte) error
}

type HashImpl struct {
}

func NewPasswordFuncImpl() Hash {
	return &HashImpl{}
}

func (HashImpl) GenerateFromPassword(password []byte, cost int) ([]byte, error) {
	return bcrypt.GenerateFromPassword(password, cost)
}

func (HashImpl) CompareHashAndPassword(hashedPassword, password []byte) error {
	return bcrypt.CompareHashAndPassword(hashedPassword, password)
}
