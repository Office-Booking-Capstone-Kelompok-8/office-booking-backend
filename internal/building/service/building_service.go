package service

import (
	"context"
	"io"
	"office-booking-backend/internal/building/dto"
	"office-booking-backend/pkg/utils/validator"
	"time"
)

type BuildingService interface {
	GetAllPublishedBuildings(ctx context.Context, q string, cityID int, districtID int, startDate time.Time, endDate time.Time, limit int, page int) (*dto.BriefPublishedBuildingsResponse, int64, error)
	GetAllBuildings(ctx context.Context, q string, cityID int, districtID int, startDate time.Time, endDate time.Time, limit int, page int) (*dto.BriefBuildingsResponse, int64, error)
	GetPublishedBuildingDetailByID(ctx context.Context, id string) (*dto.FullPublishedBuildingResponse, error)
	GetBuildingDetailByID(ctx context.Context, id string) (*dto.FullBuildingResponse, error)
	GetFacilityCategories(ctx context.Context) (*dto.FacilityCategoriesResponse, error)
	CreateEmptyBuilding(ctx context.Context, creatorID string) (string, error)
	UpdateBuilding(ctx context.Context, building *dto.UpdateBuildingRequest, buildingID string) error
	AddBuildingPicture(ctx context.Context, buildingID string, index int, alt string, picture io.Reader) (*dto.AddPictureResponse, error)
	AddBuildingFacility(ctx context.Context, buildingID string, facilities *dto.AddFacilitiesRequest) error
	ValidateBuilding(ctx context.Context, buildingID string) (*validator.ErrorsResponse, error)
	DeleteBuildingPicture(ctx context.Context, buildingID string, pictureID string) error
}
