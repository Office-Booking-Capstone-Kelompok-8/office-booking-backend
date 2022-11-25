package service

import (
	"context"
	"office-booking-backend/internal/auth/dto"
)

type AuthService interface {
	RegisterUser(ctx context.Context, user *dto.SignupRequest) error
	LoginUser(ctx context.Context, user *dto.LoginRequest) (*dto.TokenPair, error)
	LogoutUser(ctx context.Context, uid string) error
	RefreshToken(ctx context.Context, token *dto.RefreshTokenRequest) (*dto.TokenPair, error)
	RequestOTP(ctx context.Context, email string) error
	VerifyOTP(ctx context.Context, otp *dto.OTPVerifyRequest) (*string, error)
	ResetPassword(ctx context.Context, password *dto.PasswordResetRequest) error
}
