package service

import (
	"context"
	"io"
	"office-booking-backend/internal/building/dto"
	"time"
)

type BuildingService interface {
	GetAllBuildings(ctx context.Context, q string, cityID int, districtID int, startDate time.Time, endDate time.Time, limit int, page int) (*dto.BriefBuildingsResponse, int64, error)
	GetBuildingDetailByID(ctx context.Context, id string) (*dto.FullBuildingResponse, error)
	GetFacilityCategories(ctx context.Context) (*dto.FacilityCategoriesResponse, error)
	CreateEmptyBuilding(ctx context.Context, creatorID string) (string, error)
	CreateBuilding(ctx context.Context, building *dto.CreateBuildingRequest) error
	AddBuildingPicture(ctx context.Context, buildingID string, alt string, picture io.Reader) (*dto.AddPictureResponse, error)
}
