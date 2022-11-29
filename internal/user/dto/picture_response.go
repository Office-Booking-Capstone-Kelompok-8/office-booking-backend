package dto

import (
	"github.com/imagekit-developer/imagekit-go/api/uploader"
	"office-booking-backend/pkg/entity"
)

func NewPictureEntity(picture *uploader.UploadResult) *entity.ProfilePicture {
	return &entity.ProfilePicture{
		ID:  picture.FileId,
		Key: picture.Name,
		Url: picture.Url,
	}
}
