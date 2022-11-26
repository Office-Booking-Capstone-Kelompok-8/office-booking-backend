package repository

import (
	"context"
	"office-booking-backend/pkg/entity"
)

type BuildingRepository interface {
	GetAllBuildings(ctx context.Context, q string, limit int, offset int) (*entity.Building, int64, error)
}
