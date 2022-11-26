package dto

import "office-booking-backend/pkg/entity"

type UserUpdateRequest struct {
	Email string `json:"email"`
	Name  string `json:"name"`
	Phone string `json:"phone"`
	Role  int    `json:"role"`
}

func (u *UserUpdateRequest) ToEntity() *entity.User {
	return &entity.User{
		Email: u.Email,
		Detail: entity.UserDetail{
			Name:  u.Name,
			Phone: u.Phone,
		},
		Role: u.Role,
	}
}
