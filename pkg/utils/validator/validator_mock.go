package validator

import "github.com/stretchr/testify/mock"

type ValidatorMock struct {
	mock.Mock
}

func (m *ValidatorMock) Validate(s interface{}) *ErrorsResponse {
	args := m.Called(s)
	return args.Get(0).(*ErrorsResponse)
}
