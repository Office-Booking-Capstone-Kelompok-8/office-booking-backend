package repository

import (
	"context"
	"office-booking-backend/pkg/entity"
	"time"
)

type BuildingRepository interface {
	GetAllBuildings(ctx context.Context, q string, cityID int, districtID int, startDate time.Time, endDate time.Time, limit int, offset int) (*entity.Buildings, int64, error)
	GetBuildingDetailByID(ctx context.Context, id string) (*entity.Building, error)
	GetFacilityCategories(ctx context.Context) (*entity.Categories, error)
}
