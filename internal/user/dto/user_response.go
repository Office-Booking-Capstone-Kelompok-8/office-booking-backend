package dto

import (
	"office-booking-backend/pkg/config"
	"office-booking-backend/pkg/entity"
)

type UserResponse struct {
	ID         string `json:"id"`
	Email      string `json:"email"`
	Name       string `json:"name"`
	Phone      string `json:"phone"`
	PictureURL string `json:"picture"`
	Role       int    `json:"role"`
}

func NewUserResponse(user *entity.User) *UserResponse {
	response := &UserResponse{
		ID:         user.ID,
		Email:      user.Email,
		Name:       user.Detail.Name,
		Phone:      user.Detail.Phone,
		PictureURL: user.Detail.Picture.Url,
		Role:       user.Role,
	}

	if user.Detail.Picture.Url == "" {
		response.PictureURL = config.DEFAULT_USER_AVATAR
	}

	return response
}

type BriefUserResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

func NewBriefUserResponse(user *entity.User) *BriefUserResponse {
	return &BriefUserResponse{
		ID:    user.ID,
		Name:  user.Detail.Name,
		Email: user.Email,
		Phone: user.Detail.Phone,
	}
}

type BriefUsersResponse []BriefUserResponse

func NewBriefUsersResponse(users *entity.Users) *BriefUsersResponse {
	var briefUsers BriefUsersResponse
	for _, user := range *users {
		briefUsers = append(briefUsers, *NewBriefUserResponse(&user))
	}
	return &briefUsers
}
