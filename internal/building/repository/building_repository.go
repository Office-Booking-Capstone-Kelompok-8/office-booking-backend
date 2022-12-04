package repository

import (
	"context"
	"office-booking-backend/pkg/entity"
	"time"
)

type BuildingRepository interface {
	GetAllBuildings(ctx context.Context, q string, cityID int, districtID int, startDate time.Time, endDate time.Time, limit int, offset int, isPublishedOnly bool) (*entity.Buildings, int64, error)
	GetBuildingDetailByID(ctx context.Context, id string, isPublishedOnly bool) (*entity.Building, error)
	GetFacilityCategories(ctx context.Context) (*entity.Categories, error)
	AddPicture(ctx context.Context, picture *entity.Picture) error
	AddFacility(ctx context.Context, facility *entity.Facilities) error
	CreateBuilding(ctx context.Context, building *entity.Building) error
	UpdateBuildingByID(ctx context.Context, building *entity.Building) error
	CountBuildingPicturesByID(ctx context.Context, buildingID string) (int64, error)
	IsBuildingExist(ctx context.Context, buildingID string) (bool, error)
	IsBuildingPublished(ctx context.Context, buildingID string) (bool, error)
	DeleteBuildingPicturesByID(ctx context.Context, buildingID string, pictureID string) error
	DeleteBuildingFacilityByID(ctx context.Context, buildingID string, facilityID int) error
	DeleteBuildingByID(ctx context.Context, buildingID string) error
}
