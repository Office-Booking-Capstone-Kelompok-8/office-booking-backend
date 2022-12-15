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
	IsVerified bool   `json:"isVerified"`
}

func NewUserResponse(user *entity.User) *UserResponse {
	picture := user.Detail.Picture.Url
	if picture == "" {
		picture = config.DEFAULT_USER_AVATAR
	}

	return &UserResponse{
		ID:         user.ID,
		Email:      user.Email,
		Name:       user.Detail.Name,
		Phone:      user.Detail.Phone,
		PictureURL: picture,
		Role:       user.Role,
		IsVerified: *user.IsVerified,
	}
}

type BriefUserResponse struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	Picture string `json:"picture"`
}

func NewBriefUserResponse(user *entity.User) *BriefUserResponse {
	picture := user.Detail.Picture.Url
	if picture == "" {
		picture = config.DEFAULT_USER_AVATAR
	}

	return &BriefUserResponse{
		ID:      user.ID,
		Name:    user.Detail.Name,
		Email:   user.Email,
		Phone:   user.Detail.Phone,
		Picture: picture,
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

type RegisteredStatResponse struct {
	Month string `json:"month"`
	Total int64  `json:"total"`
}

func NewRegisteredStatResponse(stat *entity.MonthlyRegisteredStat) *RegisteredStatResponse {
	return &RegisteredStatResponse{
		Month: stat.Month,
		Total: stat.Total,
	}
}

type RegisteredStatResponseList []RegisteredStatResponse

func NewRegisteredStatResponseList(stats *entity.MonthlyRegisteredStatList) *RegisteredStatResponseList {
	var statList RegisteredStatResponseList
	for _, stat := range *stats {
		statList = append(statList, *NewRegisteredStatResponse(&stat))
	}
	return &statList
}

type TotalByTimeFrame struct {
	Today     int64 `json:"today"`
	ThisWeek  int64 `json:"thisWeek"`
	ThisMonth int64 `json:"thisMonth"`
	ThisYear  int64 `json:"thisYear"`
	AllTime   int64 `json:"allTime"`
}

func NewTimeframeStat(stats *entity.TimeframeStat) *TotalByTimeFrame {
	return &TotalByTimeFrame{
		Today:     stats.Day,
		ThisWeek:  stats.Week,
		ThisMonth: stats.Month,
		ThisYear:  stats.Year,
		AllTime:   stats.All,
	}
}
