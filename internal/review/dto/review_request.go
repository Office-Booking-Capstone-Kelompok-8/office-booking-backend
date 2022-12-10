package dto

import "office-booking-backend/pkg/entity"

type AddReviewRequest struct {
	Rating  int    `json:"rating" validate:"required"`
	Message string `json:"comment" validate:"required"`
}

func (a *AddReviewRequest) ToEntity() *entity.Review {
	return &entity.Review{
		Rating:  a.Rating,
		Message: a.Message,
	}
}
