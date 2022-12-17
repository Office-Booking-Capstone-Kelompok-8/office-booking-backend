package repository

import (
	"context"
	"office-booking-backend/internal/building/dto"
	"office-booking-backend/pkg/entity"
)

type BuildingRepository interface {
	GetAllBuildings(ctx context.Context, filter *dto.SearchBuildingQueryParam, isPublishedOnly bool) (*entity.Buildings, int64, error)
	GetBuildingDetailByID(ctx context.Context, id string, isPublishedOnly bool) (*entity.Building, error)
	GetFacilityCategories(ctx context.Context) (*entity.Categories, error)
	GetCities(ctx context.Context) (*entity.Cities, error)
	GetDistrictsByCityID(ctx context.Context, cityID int) (*entity.Districts, error)
	GetDistrictByID(ctx context.Context, districtID int) (*entity.District, error)
	GetBuildingReviewsByID(ctx context.Context, buildingID string, filter *dto.GetBuildingReviewsQueryParam) (*entity.Reviews, error)
	GetBuildingCountByCity(ctx context.Context) (*entity.CitiesStat, error)
	GetBuildingCountByTime(ctx context.Context) (*entity.TimeframeStat, error)
	AddPicture(ctx context.Context, picture *entity.Picture) error
	AddFacility(ctx context.Context, facility *entity.Facilities) error
	CreateBuilding(ctx context.Context, building *entity.Building) error
	UpdateBuildingByID(ctx context.Context, building *entity.Building) error
	CountBuildingPicturesByID(ctx context.Context, buildingID string) (int64, error)
	CountBuildingReviewsByID(ctx context.Context, buildingID string) (int64, error)
	IsBuildingExist(ctx context.Context, buildingID string) (bool, error)
	DeleteBuildingPicturesByID(ctx context.Context, buildingID string, pictureID string) error
	DeleteBuildingFacilityByID(ctx context.Context, buildingID string, facilityID int) error
	DeleteBuildingByID(ctx context.Context, buildingID string) error
}
