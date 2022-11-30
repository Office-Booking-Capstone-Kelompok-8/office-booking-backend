package dto

import "office-booking-backend/pkg/entity"

type UserUpdateRequest struct {
	Email string `json:"email" validate:"omitempty,email"`
	Name  string `json:"name"`
	Phone string `json:"phone"`
}

func (u *UserUpdateRequest) ToEntity() *entity.User {
	return &entity.User{
		Email: u.Email,
		Detail: entity.UserDetail{
			Name:  u.Name,
			Phone: u.Phone,
		},
	}
}
