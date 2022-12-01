package validator

import "github.com/stretchr/testify/mock"

type ValidatorMock struct {
	mock.Mock
}

func (m *ValidatorMock) ValidateStruct(s interface{}) *ErrorsResponse {
	args := m.Called(s)
	return args.Get(0).(*ErrorsResponse)
}

func (m *ValidatorMock) ValidateVar(field interface{}, tag string) *ErrorsResponse {
	args := m.Called(field, tag)
	return args.Get(0).(*ErrorsResponse)
}
