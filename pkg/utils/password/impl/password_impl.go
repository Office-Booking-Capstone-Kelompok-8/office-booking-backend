package impl

import (
	"office-booking-backend/pkg/utils/password"

	"golang.org/x/crypto/bcrypt"
)

type PasswordServiceImpl struct {
}

func NewPasswordFuncImpl() password.PasswordService {
	return &PasswordServiceImpl{}
}

func (PasswordServiceImpl) GenerateFromPassword(password []byte, cost int) ([]byte, error) {
	return bcrypt.GenerateFromPassword(password, cost)
}

func (PasswordServiceImpl) CompareHashAndPassword(hashedPassword, password []byte) error {
	return bcrypt.CompareHashAndPassword(hashedPassword, password)
}
