package impl

import (
	"context"
	"io"
	"log"
	"office-booking-backend/internal/building/dto"
	"office-booking-backend/internal/building/repository"
	"office-booking-backend/internal/building/service"
	repository2 "office-booking-backend/internal/reservation/repository"
	"office-booking-backend/pkg/entity"
	err2 "office-booking-backend/pkg/errors"
	"office-booking-backend/pkg/utils/imagekit"
	"office-booking-backend/pkg/utils/validator"

	"golang.org/x/sync/errgroup"

	"github.com/google/uuid"
)

type BuildingServiceImpl struct {
	repo            repository.BuildingRepository
	reservationRepo repository2.ReservationRepository
	imgKitService   imagekit.ImgKitService
	validator       validator.Validator
}

func NewBuildingServiceImpl(repo repository.BuildingRepository, reservationRepo repository2.ReservationRepository, imgKitService imagekit.ImgKitService, validator validator.Validator) service.BuildingService {
	return &BuildingServiceImpl{
		repo:            repo,
		reservationRepo: reservationRepo,
		imgKitService:   imgKitService,
		validator:       validator,
	}
}

func (b *BuildingServiceImpl) GetAllPublishedBuildings(ctx context.Context, filter *dto.SearchBuildingQueryParam) (*dto.BriefPublishedBuildingsResponse, int64, error) {
	filter.EndDate = filter.StartDate.ToTime().AddDate(0, filter.Duration, 0)

	count := int64(0)

	filter.Offset = (filter.Page - 1) * filter.Limit
	//	get all buildings
	buildings, count, err := b.repo.GetAllBuildings(ctx, filter, true)
	if err != nil {
		log.Println("error when getting all buildings: ", err)
		return nil, 0, err
	}

	buildingsResponse := dto.NewBriefPublishedBuildingsResponse(buildings)
	return buildingsResponse, count, nil
}

func (b *BuildingServiceImpl) GetAllBuildings(ctx context.Context, filter *dto.SearchBuildingQueryParam) (*dto.BriefBuildingsResponse, int64, error) {
	filter.EndDate = filter.StartDate.ToTime().AddDate(0, filter.Duration, 0)

	count := int64(0)

	filter.Offset = (filter.Page - 1) * filter.Limit
	buildings, count, err := b.repo.GetAllBuildings(ctx, filter, false)
	if err != nil {
		log.Println("error when getting all buildings: ", err)
		return nil, 0, err
	}

	buildingsResponse := dto.NewBriefBuildingsResponse(buildings)
	return buildingsResponse, count, nil
}

func (b *BuildingServiceImpl) GetPublishedBuildingDetailByID(ctx context.Context, id string) (*dto.FullPublishedBuildingResponse, error) {
	building, err := b.repo.GetBuildingDetailByID(ctx, id, true)
	if err != nil {
		log.Println("error when getting building detail by id: ", err)
		return nil, err
	}

	buildingResponse := dto.NewFullPublishedBuildingResponse(building)
	return buildingResponse, nil
}

func (b *BuildingServiceImpl) GetBuildingDetailByID(ctx context.Context, id string) (*dto.FullBuildingResponse, error) {
	building, err := b.repo.GetBuildingDetailByID(ctx, id, false)
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

func (b *BuildingServiceImpl) CreateEmptyBuilding(ctx context.Context, creatorID string) (string, error) {
	building := new(entity.Building)
	building.ID = uuid.New().String()
	building.CreatedByID = creatorID
	err := b.repo.CreateBuilding(ctx, building)
	if err != nil {
		log.Println("error when creating empty building: ", err)
		return "", err
	}

	return building.ID, nil
}

func (b *BuildingServiceImpl) GetCities(ctx context.Context) (*dto.CitiesResponse, error) {
	cities, err := b.repo.GetCities(ctx)
	if err != nil {
		log.Println("error when getting cities: ", err)
		return nil, err
	}

	return dto.NewCitiesResponse(cities), nil
}
func (b *BuildingServiceImpl) GetDistrictsByCityID(ctx context.Context, cityID int) (*dto.DistrictsResponse, error) {
	districts, err := b.repo.GetDistrictsByCityID(ctx, cityID)
	if err != nil {
		log.Println("error when getting districts: ", err)
		return nil, err
	}

	return dto.NewDistrictsResponse(districts), nil
}

func (b *BuildingServiceImpl) GetBuildingStatistics(ctx context.Context) (*dto.BuildingStatResponse, error) {
	statByCities := new(dto.TotalByCities)
	statByTimeFrame := new(dto.TotalByTimeFrame)

	errGroup, c := errgroup.WithContext(ctx)
	errGroup.Go(func() error {
		total, err := b.repo.GetBuildingCountByCity(c)
		if err != nil {
			log.Println("error when getting building count by city: ", err)
			return err
		}

		statByCities = dto.NewTotalByCities(total)
		return nil
	})

	errGroup.Go(func() error {
		total, err := b.repo.GetBuildingCountByTime(c)
		if err != nil {
			log.Println("error when getting building count by time: ", err)
			return err
		}

		statByTimeFrame = dto.NewTimeframeStat(total)
		return nil
	})

	if err := errGroup.Wait(); err != nil {
		return nil, err
	}

	return &dto.BuildingStatResponse{
		ByCity:      statByCities,
		ByTimeFrame: statByTimeFrame,
	}, nil
}

func (b *BuildingServiceImpl) UpdateBuilding(ctx context.Context, building *dto.UpdateBuildingRequest, buildingID string) error {
	buildingEntity := building.ToEntity(buildingID)

	if buildingEntity.DistrictID != 0 || buildingEntity.CityID != 0 {
		//	check if district in the city
		district, err := b.repo.GetDistrictByID(ctx, buildingEntity.DistrictID)
		if err != nil {
			log.Println("error when getting district by id: ", err)
			return err
		}

		if district.CityID != buildingEntity.CityID {
			return err2.ErrDistrictNotInCity
		}
	}

	err := b.repo.UpdateBuildingByID(ctx, buildingEntity)
	if err != nil {
		log.Println("error when updating building: ", err)
		return err
	}

	return nil
}

func (b *BuildingServiceImpl) UpdateBuildingPublishState(ctx context.Context, building *dto.PublishRequest, buildingID string) error {
	buildingEntity := building.ToEntity(buildingID)

	err := b.repo.UpdateBuildingByID(ctx, buildingEntity)
	if err != nil {
		log.Println("error when updating building publish state: ", err)
		return err
	}

	return nil
}

func (b *BuildingServiceImpl) AddBuildingPicture(ctx context.Context, buildingID string, index int, alt string, picture io.Reader) (*dto.AddPictureResponse, error) {
	//	check if building exists
	exists, err := b.repo.IsBuildingExist(ctx, buildingID)
	if err != nil {
		log.Println("error when checking building: ", err)
		return nil, err
	}

	if !exists {
		return nil, err2.ErrBuildingNotFound
	}

	pictureCount, err := b.repo.CountBuildingPicturesByID(ctx, buildingID)
	if err != nil {
		log.Println("error when counting building pictures: ", err)
		return nil, err
	}

	if pictureCount >= 10 {
		return nil, err2.ErrPicureLimitExceeded
	}

	pictureKey := uuid.New().String()
	uploadResult, err := b.imgKitService.UploadFile(ctx, picture, pictureKey, "buildings")
	if err != nil {
		log.Println("error when uploading file: ", err)
		return nil, err2.ErrPictureServiceFailed
	}

	pictureEntity := &entity.Picture{
		ID:           uploadResult.FileId,
		BuildingID:   buildingID,
		Index:        &index,
		Url:          uploadResult.Url,
		ThumbnailUrl: uploadResult.ThumbnailUrl,
		Alt:          alt,
		Key:          pictureKey,
	}

	err = b.repo.AddPicture(ctx, pictureEntity)
	if err != nil {
		log.Println("error when adding building picture: ", err)
		return nil, err
	}

	return dto.NewAddPictureResponse(pictureEntity), nil
}

func (b *BuildingServiceImpl) AddBuildingFacility(ctx context.Context, buildingID string, facilities *dto.AddFacilitiesRequest) error {
	facilitiesEntity := facilities.ToEntity(buildingID)
	err := b.repo.AddFacility(ctx, facilitiesEntity)
	if err != nil {
		log.Println("error when adding building facilities: ", err)
		return err
	}

	return nil
}

func (b *BuildingServiceImpl) ValidateBuilding(ctx context.Context, buildingID string) (*validator.ErrorsResponse, error) {
	//	get building
	building, err := b.repo.GetBuildingDetailByID(ctx, buildingID, false)
	if err != nil {
		log.Println("error when getting building detail by id: ", err)
		return nil, err
	}

	if *building.IsPublished {
		return nil, nil
	}

	buildingDtp := dto.NewFullBuildingResponse(building)
	errs := b.validator.ValidateJSON(buildingDtp)

	indexZero := false
	for _, picture := range building.Pictures {
		if *picture.Index == 0 {
			indexZero = true
			break
		}
	}

	if !indexZero {
		errs.AddError("pictures", "main image is required")
	}

	if errs != nil {
		return errs, err2.ErrNotPublishWorthy
	}

	return nil, nil
}

func (b *BuildingServiceImpl) DeleteBuildingPicture(ctx context.Context, buildingID string, pictureID string) error {
	err := b.repo.DeleteBuildingPicturesByID(ctx, buildingID, pictureID)
	if err != nil {
		log.Println("error when deleting building picture: ", err)
		return err
	}

	err = b.imgKitService.DeleteFile(ctx, pictureID)
	if err != nil {
		log.Println("error when deleting file: ", err)
		return err2.ErrPictureServiceFailed
	}

	return nil
}

func (b *BuildingServiceImpl) DeleteBuildingFacility(ctx context.Context, buildingID string, facilityID int) error {
	err := b.repo.DeleteBuildingFacilityByID(ctx, buildingID, facilityID)
	if err != nil {
		log.Println("error when deleting building facility: ", err)
		return err
	}

	return nil
}

func (b *BuildingServiceImpl) DeleteBuilding(ctx context.Context, buildingID string) error {
	count, err := b.reservationRepo.CountBuildingActiveReservations(ctx, buildingID)
	if err != nil {
		log.Println("error when counting building active reservations: ", err)
		return err
	}

	if count > 0 {
		return err2.ErrBuildingHasReservation
	}

	err = b.repo.DeleteBuildingByID(ctx, buildingID)
	if err != nil {
		log.Println("error when deleting building: ", err)
		return err
	}

	return nil
}

func (b *BuildingServiceImpl) GetBuildingReviews(ctx context.Context, buildingID string, filter *dto.GetBuildingReviewsQueryParam) (*dto.BriefBuildingReviewsResponse, int64, error) {
	var reviews *entity.Reviews
	var total int64

	errGroup, c := errgroup.WithContext(ctx)
	errGroup.Go(func() error {
		filter.Offset = (filter.Page - 1) * filter.Limit
		savedReviews, err := b.repo.GetBuildingReviewsByID(c, buildingID, filter)
		if err != nil {
			log.Println("error when getting reviews: ", err)
			return err
		}

		reviews = savedReviews
		return nil
	})

	errGroup.Go(func() error {
		count, err := b.repo.CountBuildingReviewsByID(c, buildingID)
		if err != nil {
			log.Println("error when counting reviews: ", err)
			return err
		}

		total = count
		return nil
	})

	if err := errGroup.Wait(); err != nil {
		return nil, 0, err
	}

	return dto.NewBriefBuildingReviewsResponse(reviews), total, nil
}
