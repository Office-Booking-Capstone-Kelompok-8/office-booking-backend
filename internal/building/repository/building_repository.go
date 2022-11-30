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
	CreateBuilding(ctx context.Context, building *entity.Building) error
	UpdateBuildingByID(ctx context.Context, building *entity.Building) error
	UpdateBuildingPictures(ctx context.Context, picture *entity.Picture) error
	CountBuildingPicturesByID(ctx context.Context, buildingId string) (int64, error)
	CheckBuilding(ctx context.Context, buildingId string) (bool, error)
	AddPicture(ctx context.Context, picture *entity.Picture) error
}
