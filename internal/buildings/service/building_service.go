package service

import (
	"context"
	"office-booking-backend/internal/buildings/dto"
)

type BuildingService interface {
	GetAllBuildings(ctx context.Context, q string, limit int, offset int) (*dto.GetAllBuildingsResponse, int64, error)
}
