package dto

import "office-booking-backend/pkg/entity"

type UserUpdateRequest struct {
	Email string `json:"email" validate:"omitempty,email"`
	Name  string `json:"name" validate:"omitempty,min=3,max=50"`
	Phone string `json:"phone" validate:"omitempty,number,min=3,max=50"`
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
