package dto

import "office-booking-backend/pkg/entity"

type SignupRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	Name     string `json:"name" validate:"required"`
	Phone    string `json:"phone" validate:"required"`
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

type OTPRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type OTPVerifyRequest struct {
	Email string `json:"email" validate:"required,email"`
	Code  string `json:"code" validate:"required"`
}

type PasswordResetRequest struct {
	Token    string `json:"token" validate:"required"`
	Password string `json:"password" validate:"required,min=8"`
}

func (p *PasswordResetRequest) ToEntity() *entity.User {
	return &entity.User{
		Password: p.Password,
	}
}
