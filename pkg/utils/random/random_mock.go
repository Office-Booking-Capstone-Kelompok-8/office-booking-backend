package random

import "github.com/stretchr/testify/mock"

type GeneratorMock struct {
	mock.Mock
}

func (m *GeneratorMock) GenerateRandomIntString(digit int) (string, error) {
	args := m.Called(digit)
	return args.String(0), args.Error(1)
}
