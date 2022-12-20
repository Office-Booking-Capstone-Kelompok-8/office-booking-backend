package service

import (
	"context"
	"io"
	"office-booking-backend/internal/building/dto"
	"office-booking-backend/pkg/utils/validator"
)

type BuildingService interface {
	GetAllPublishedBuildings(ctx context.Context, filter *dto.SearchBuildingQueryParam) (*dto.BriefPublishedBuildingsResponse, int64, error)
	GetAllBuildings(ctx context.Context, filter *dto.SearchBuildingQueryParam) (*dto.BriefBuildingsResponse, int64, error)
	GetPublishedBuildingDetailByID(ctx context.Context, id string) (*dto.FullPublishedBuildingResponse, error)
	GetBuildingDetailByID(ctx context.Context, id string) (*dto.FullBuildingResponse, error)
	GetFacilityCategories(ctx context.Context) (*dto.FacilityCategoriesResponse, error)
	GetCities(ctx context.Context) (*dto.CitiesResponse, error)
	GetDistrictsByCityID(ctx context.Context, cityID int) (*dto.DistrictsResponse, error)
	GetBuildingReviews(ctx context.Context, buildingID string, filter *dto.GetBuildingReviewsQueryParam) (*dto.BriefBuildingReviewsResponse, int64, error)
	GetBuildingStatistics(ctx context.Context) (*dto.BuildingStatResponse, error)
	CreateEmptyBuilding(ctx context.Context, creatorID string) (string, error)
	UpdateBuilding(ctx context.Context, building *dto.UpdateBuildingRequest, buildingID string) error
	UpdateBuildingPublishState(ctx context.Context, building *dto.PublishRequest, buildingID string) error
	AddBuildingPicture(ctx context.Context, buildingID string, index int, alt string, picture io.Reader) (*dto.AddPictureResponse, error)
	AddBuildingFacility(ctx context.Context, buildingID string, facilities *dto.AddFacilitiesRequest) error
	ValidateBuilding(ctx context.Context, buildingID string) (*validator.ErrorsResponse, error)
	DeleteBuildingPicture(ctx context.Context, buildingID string, pictureID string) error
	DeleteBuildingFacility(ctx context.Context, buildingID string, facilityID int) error
	DeleteBuilding(ctx context.Context, buildingID string) error
}
