package mail

import (
	"context"
	"github.com/stretchr/testify/mock"
)

type ClientMock struct {
	mock.Mock
}

func (m *ClientMock) SendMail(ctx context.Context, mail *Mail) error {
	args := m.Called(ctx, mail)
	return args.Error(0)
}
