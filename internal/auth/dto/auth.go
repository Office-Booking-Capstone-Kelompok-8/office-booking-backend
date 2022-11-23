package dto

import "office-booking-backend/pkg/entity"

type SignupRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
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
