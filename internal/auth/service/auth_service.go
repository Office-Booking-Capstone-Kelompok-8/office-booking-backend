package service

import (
	"context"
	"office-booking-backend/internal/auth/dto"
)

type AuthService interface {
	RegisterUser(ctx context.Context, user *dto.SignupRequest) (string, error)
	RegisterAdmin(ctx context.Context, user *dto.SignupRequest) (string, error)
	LoginUser(ctx context.Context, user *dto.LoginRequest) (*dto.TokenPair, error)
	LogoutUser(ctx context.Context, uid string) error
	RefreshToken(ctx context.Context, token *dto.RefreshTokenRequest) (*dto.TokenPair, error)
	RequestPasswordResetOTP(ctx context.Context, email string) error
	VerifyPasswordResetOTP(ctx context.Context, otp *dto.ResetPasswordOTPVerifyRequest) (*string, error)
	RequestEmailOTP(ctx context.Context, userID string) error
	VerifyEmailOTP(ctx context.Context, userID string, otp *dto.VerifyEmailOTOPVerifyRequest) error
	ResetPassword(ctx context.Context, password *dto.PasswordResetRequest) error
	ChangePassword(ctx context.Context, uid string, password *dto.ChangePasswordRequest) error
}
