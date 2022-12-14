package dto

import "office-booking-backend/pkg/entity"

type SignupRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	Name     string `json:"name" validate:"required,min=3"`
	Phone    string `json:"phone" validate:"required,number"`
}

func (s *SignupRequest) ToEntity() *entity.User {
	return &entity.User{
		Email:    s.Email,
		Password: s.Password,
		Detail: entity.UserDetail{
			Name:  s.Name,
			Phone: s.Phone,
		},
	}
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type ResetPasswordOTPRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type ResetPasswordOTPVerifyRequest struct {
	Email string `json:"email" validate:"required,email"`
	Code  string `json:"code" validate:"required"`
}

type VerifyEmailOTOPVerifyRequest struct {
	Code string `json:"code" validate:"required"`
}

type PasswordResetRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	Key      string `json:"key" validate:"required"`
}

func (p *PasswordResetRequest) ToEntity() *entity.User {
	return &entity.User{
		Password: p.Password,
	}
}

type ChangePasswordRequest struct {
	OldPassword string `json:"oldPassword" validate:"required"`
	NewPassword string `json:"newPassword" validate:"required,min=8"`
}

func (c *ChangePasswordRequest) ToEntity() *entity.User {
	return &entity.User{
		Password: c.NewPassword,
	}
}
