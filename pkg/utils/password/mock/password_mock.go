package mock

import "github.com/stretchr/testify/mock"

type PasswordHashService struct {
	mock.Mock
}

func (m *PasswordHashService) GenerateFromPassword(password []byte, cost int) ([]byte, error) {
	args := m.Called(password, cost)
	return args.Get(0).([]byte), args.Error(1)
}

func (m *PasswordHashService) CompareHashAndPassword(hashedPassword, password []byte) error {
	args := m.Called(hashedPassword, password)
	return args.Error(0)
}
