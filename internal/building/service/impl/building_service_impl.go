package impl

import (
	"context"
	"log"
	"office-booking-backend/internal/building/dto"
	"office-booking-backend/internal/building/repository"
	"office-booking-backend/internal/building/service"
	err2 "office-booking-backend/pkg/errors"
	"time"
)

type BuildingServiceImpl struct {
	repo repository.BuildingRepository
}

func NewBuildingServiceImpl(repo repository.BuildingRepository) service.BuildingService {
	return &BuildingServiceImpl{
		repo: repo,
	}
}

func (b *BuildingServiceImpl) GetAllBuildings(ctx context.Context, q string, cityID int, districtID int, startDate time.Time, endDate time.Time, limit int, page int) (*dto.BriefBuildingsResponse, int64, error) {
	//	check if start date is after end date
	if startDate.After(endDate) {
		return nil, 0, err2.ErrStartDateAfterEndDate
	}

	offset := (page - 1) * limit

	//	get all buildings
	buildings, count, err := b.repo.GetAllBuildings(ctx, q, cityID, districtID, startDate, endDate, limit, offset)
	if err != nil {
		log.Println("error when getting all buildings: ", err)
		return nil, 0, err
	}

	buildingsResponse := dto.NewBriefBuildingsResponse(buildings)
	return buildingsResponse, count, nil
}

func (b *BuildingServiceImpl) GetBuildingDetailByID(ctx context.Context, id string) (*dto.FullBuildingResponse, error) {
	building, err := b.repo.GetBuildingDetailByID(ctx, id)
	if err != nil {
		log.Println("error when getting building detail by id: ", err)
		return nil, err
	}

	buildingResponse := dto.NewFullBuildingResponse(building)
	return buildingResponse, nil
}

func (b *BuildingServiceImpl) GetFacilityCategories(ctx context.Context) (*dto.FacilityCategoriesResponse, error) {
	facilityCategories, err := b.repo.GetFacilityCategories(ctx)
	if err != nil {
		log.Println("error when getting facility categories: ", err)
		return nil, err
	}

	facilityCategoriesResponse := dto.NewFacilityCategoriesResponse(facilityCategories)
	return facilityCategoriesResponse, nil

}
