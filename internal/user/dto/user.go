package dto

import "office-booking-backend/pkg/entity"

type UserResponse struct {
	ID         string `json:"id"`
	Email      string `json:"email"`
	Name       string `json:"name"`
	Phone      string `json:"phone"`
	PictureURL string `json:"picture"`
	Role       int    `json:"role"`
}

func NewUserResponse(user *entity.User) *UserResponse {
	return &UserResponse{
		ID:         user.ID,
		Email:      user.Email,
		Name:       user.Detail.Name,
		Phone:      user.Detail.Phone,
		PictureURL: user.Detail.Picture.Url,
		Role:       user.Role,
	}
}
