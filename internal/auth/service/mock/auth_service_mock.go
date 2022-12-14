package mock

import (
	"context"
	"office-booking-backend/internal/auth/dto"

	"github.com/stretchr/testify/mock"
)

type AuthServiceMock struct {
	mock.Mock
}

func (m *AuthServiceMock) RegisterUser(ctx context.Context, user *dto.SignupRequest) (string, error) {
	args := m.Called(ctx, user)
	return args.String(0), args.Error(1)
}

func (m *AuthServiceMock) RegisterAdmin(ctx context.Context, user *dto.SignupRequest) (string, error) {
	args := m.Called(ctx, user)
	return args.String(0), args.Error(1)
}

func (m *AuthServiceMock) LoginUser(ctx context.Context, user *dto.LoginRequest) (*dto.TokenPair, error) {
	args := m.Called(ctx, user)
	return args.Get(0).(*dto.TokenPair), args.Error(1)
}

func (m *AuthServiceMock) LogoutUser(ctx context.Context, uid string) error {
	args := m.Called(ctx, uid)
	return args.Error(0)
}

func (m *AuthServiceMock) RefreshToken(ctx context.Context, token *dto.RefreshTokenRequest) (*dto.TokenPair, error) {
	args := m.Called(ctx, token)
	return args.Get(0).(*dto.TokenPair), args.Error(1)
}

func (m *AuthServiceMock) RequestPasswordResetOTP(ctx context.Context, email string) error {
	args := m.Called(ctx, email)
	return args.Error(0)
}

func (m *AuthServiceMock) VerifyPasswordResetOTP(ctx context.Context, otp *dto.ResetPasswordOTPVerifyRequest) (*string, error) {
	args := m.Called(ctx, otp)
	return args.Get(0).(*string), args.Error(1)
}

func (m *AuthServiceMock) ResetPassword(ctx context.Context, password *dto.PasswordResetRequest) error {
	args := m.Called(ctx, password)
	return args.Error(0)
}

func (m *AuthServiceMock) ChangePassword(ctx context.Context, uid string, password *dto.ChangePasswordRequest) error {
	args := m.Called(ctx, uid, password)
	return args.Error(0)
}

func (m *AuthServiceMock) RequestEmailOTP(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *AuthServiceMock) VerifyEmailOTP(ctx context.Context, userID string, otp *dto.VerifyEmailOTOPVerifyRequest) error {
	args := m.Called(ctx, userID, otp)
	return args.Error(0)
}
