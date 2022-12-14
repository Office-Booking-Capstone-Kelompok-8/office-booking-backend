package controller

import (
	"errors"
	"mime/multipart"
	"office-booking-backend/internal/building/dto"
	"office-booking-backend/internal/building/service"
	err2 "office-booking-backend/pkg/errors"
	"office-booking-backend/pkg/response"
	"office-booking-backend/pkg/utils/validator"
	"reflect"
	"strconv"

	"github.com/golang-jwt/jwt/v4"

	"github.com/gofiber/fiber/v2"
)

type BuildingController struct {
	buildingService service.BuildingService
	validator       validator.Validator
}

func NewBuildingController(buildingService service.BuildingService, validator validator.Validator) *BuildingController {
	return &BuildingController{
		buildingService: buildingService,
		validator:       validator,
	}
}

func (b *BuildingController) GetAllPublishedBuildings(c *fiber.Ctx) error {
	filter := new(dto.SearchBuildingQueryParam)
	if err := c.QueryParser(filter); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if errs := b.validator.ValidateQuery(filter); errs != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.BaseResponse{
			Message: err2.ErrInvalidQueryParams.Error(),
			Data:    errs,
		})
	}

	buildings, total, err := b.buildingService.GetAllPublishedBuildings(c.Context(), filter)
	if err != nil {
		if errors.Is(err, err2.ErrInvalidDateRange) {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(response.BaseResponse{
		Message: "buildings fetched successfully",
		Data:    buildings,
		Meta: fiber.Map{
			"limit": filter.Limit,
			"page":  filter.Page,
			"total": total,
		},
	})
}

func (b *BuildingController) GetAllBuildings(c *fiber.Ctx) error {
	filter := new(dto.SearchBuildingQueryParam)
	if err := c.QueryParser(filter); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if errs := b.validator.ValidateQuery(filter); errs != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.BaseResponse{
			Message: err2.ErrInvalidQueryParams.Error(),
			Data:    errs,
		})
	}

	buildings, total, err := b.buildingService.GetAllBuildings(c.Context(), filter)
	if err != nil {
		if errors.Is(err, err2.ErrInvalidDateRange) {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(response.BaseResponse{
		Message: "buildings fetched successfully",
		Data:    buildings,
		Meta: fiber.Map{
			"limit": filter.Limit,
			"page":  filter.Page,
			"total": total,
		},
	})
}

func (b *BuildingController) GetPublishedBuildingDetailByID(c *fiber.Ctx) error {
	id := c.Params("buildingID")

	building, err := b.buildingService.GetPublishedBuildingDetailByID(c.Context(), id)
	if err != nil {
		if errors.Is(err, err2.ErrBuildingNotFound) {
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}

		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(response.BaseResponse{
		Message: "building fetched successfully",
		Data:    building,
	})
}

func (b *BuildingController) GetBuildingDetailByID(c *fiber.Ctx) error {
	id := c.Params("buildingID")

	building, err := b.buildingService.GetBuildingDetailByID(c.Context(), id)
	if err != nil {
		if errors.Is(err, err2.ErrBuildingNotFound) {
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}

		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(response.BaseResponse{
		Message: "building fetched successfully",
		Data:    building,
	})
}

func (b *BuildingController) GetFacilityCategories(c *fiber.Ctx) error {
	categories, err := b.buildingService.GetFacilityCategories(c.Context())
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(response.BaseResponse{
		Message: "facility categories fetched successfully",
		Data:    categories,
	})
}

func (b *BuildingController) GetCities(c *fiber.Ctx) error {
	cities, err := b.buildingService.GetCities(c.Context())
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(response.BaseResponse{
		Message: "cities fetched successfully",
		Data:    cities,
	})
}

func (b *BuildingController) GetDistricts(c *fiber.Ctx) error {
	cityID := c.Query("cityId", "")
	cityIDInt, err := strconv.Atoi(cityID)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err2.ErrInvalidQueryParams.Error())
	}

	districts, err := b.buildingService.GetDistrictsByCityID(c.Context(), cityIDInt)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(response.BaseResponse{
		Message: "districts fetched successfully",
		Data:    districts,
	})
}

func (b *BuildingController) GetBuildingTotal(c *fiber.Ctx) error {
	total, err := b.buildingService.GetBuildingStatistics(c.Context())
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(response.BaseResponse{
		Message: "building total fetched successfully",
		Data:    total,
	})
}

func (b *BuildingController) GetBuildingReviews(c *fiber.Ctx) error {
	buildingID := c.Params("buildingID")

	filter := new(dto.GetBuildingReviewsQueryParam)
	if err := c.QueryParser(filter); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err2.ErrInvalidQueryParams.Error())
	}

	if errs := b.validator.ValidateQuery(filter); errs != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.BaseResponse{
			Message: err2.ErrInvalidQueryParams.Error(),
			Data:    errs,
		})
	}

	reviews, total, err := b.buildingService.GetBuildingReviews(c.Context(), buildingID, filter)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(response.BaseResponse{
		Message: "building reviews fetched successfully",
		Data:    reviews,
		Meta: fiber.Map{
			"limit": filter.Limit,
			"page":  filter.Page,
			"total": total,
		},
	})
}

func (b *BuildingController) RequestNewBuildingID(c *fiber.Ctx) error {
	token := c.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	creatorID := claims["uid"].(string)

	id, err := b.buildingService.CreateEmptyBuilding(c.Context(), creatorID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(response.BaseResponse{
		Message: "building id requested successfully",
		Data: fiber.Map{
			"id": id,
		},
	})
}

func (b *BuildingController) AddBuildingPicture(c *fiber.Ctx) error {
	buildingID := c.Params("buildingID")

	altText := c.FormValue("alt", "")
	index := c.FormValue("index", "-1")
	indexInt, err := strconv.Atoi(index)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err2.ErrInvalidRequestBody.Error())
	}

	fileHeader, err := c.FormFile("picture")
	validatorDto := struct {
		AltText string                `json:"alt" validate:"omitempty,min=3,max=100"`
		Index   int                   `json:"index" validate:"gte=0,lte=9"`
		Picture *multipart.FileHeader `json:"picture" validate:"multipartImage"`
	}{
		AltText: altText,
		Index:   indexInt,
		Picture: fileHeader,
	}

	errs := b.validator.ValidateJSON(validatorDto)
	if errs != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.BaseResponse{
			Message: err2.ErrInvalidRequestBody.Error(),
			Data:    errs,
		})
	}

	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	file, err := fileHeader.Open()
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err2.ErrInvalidRequestBody.Error())
	}
	defer file.Close()

	result, err := b.buildingService.AddBuildingPicture(c.Context(), buildingID, indexInt, altText, file)
	if err != nil {
		switch err {
		case err2.ErrBuildingNotFound:
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		case err2.ErrPicureLimitExceeded:
			return fiber.NewError(fiber.StatusConflict, err.Error())
		default:
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
	}

	return c.Status(fiber.StatusCreated).JSON(response.BaseResponse{
		Message: "building picture uploaded successfully",
		Data:    result,
	})
}

func (b *BuildingController) AddBuildingFacilities(c *fiber.Ctx) error {
	buildingID := c.Params("buildingID")

	facilities := new(dto.AddFacilitiesRequest)
	if err := c.BodyParser(facilities); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err2.ErrInvalidRequestBody.Error())
	}

	if errs := b.validator.ValidateVar(facilities, "required,dive"); errs != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.BaseResponse{
			Message: err2.ErrInvalidRequestBody.Error(),
			Data:    errs,
		})
	}
	if err := b.buildingService.AddBuildingFacility(c.Context(), buildingID, facilities); err != nil {
		switch err {
		case err2.ErrBuildingNotFound:
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		case err2.ErrInvalidCategoryID:
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		default:
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
	}

	return c.Status(fiber.StatusOK).JSON(response.BaseResponse{
		Message: "building facilities added successfully",
	})
}

func (b *BuildingController) UpdateBuilding(c *fiber.Ctx) error {
	buildingID := c.Params("buildingID")

	building := new(dto.UpdateBuildingRequest)
	if err := c.BodyParser(building); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err2.ErrInvalidRequestBody.Error())
	}

	if errs := b.validator.ValidateJSON(building); errs != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.BaseResponse{
			Message: err2.ErrInvalidRequestBody.Error(),
			Data:    errs,
		})
	}

	if reflect.DeepEqual(*building, dto.UpdateBuildingRequest{}) {
		return fiber.NewError(fiber.StatusBadRequest, err2.ErrInvalidRequestBody.Error())
	}

	if err := b.buildingService.UpdateBuilding(c.Context(), building, buildingID); err != nil {
		switch err {
		case err2.ErrBuildingNotFound:
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		case err2.ErrPictureNotFound:
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		case err2.ErrFacilityNotFound:
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		default:
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
	}

	return c.Status(fiber.StatusOK).JSON(response.BaseResponse{
		Message: "building updated successfully",
	})
}

func (b *BuildingController) UpdateBuildingPublishState(c *fiber.Ctx) error {
	buildingID := c.Params("buildingID")

	publishState := new(dto.PublishRequest)
	if err := c.BodyParser(publishState); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err2.ErrInvalidRequestBody.Error())
	}

	if errs := b.validator.ValidateJSON(publishState); errs != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.BaseResponse{
			Message: err2.ErrInvalidRequestBody.Error(),
			Data:    errs,
		})
	}

	if *publishState.IsPublished {
		if errs, err := b.buildingService.ValidateBuilding(c.Context(), buildingID); err != nil {
			return c.Status(fiber.StatusConflict).JSON(response.BaseResponse{
				Message: err.Error(),
				Data:    errs,
			})
		}
	}

	if err := b.buildingService.UpdateBuildingPublishState(c.Context(), publishState, buildingID); err != nil {
		switch err {
		case err2.ErrBuildingNotFound:
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		default:
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
	}

	return c.Status(fiber.StatusOK).JSON(response.BaseResponse{
		Message: "building publish state updated successfully",
	})
}

func (b *BuildingController) DeleteBuildingPicture(c *fiber.Ctx) error {
	buildingID := c.Params("buildingID")
	pictureID := c.Params("pictureID")

	if err := b.buildingService.DeleteBuildingPicture(c.Context(), buildingID, pictureID); err != nil {
		switch err {
		case err2.ErrPictureNotFound:
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		default:
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
	}

	return c.Status(fiber.StatusOK).JSON(response.BaseResponse{
		Message: "building picture deleted successfully",
	})
}

func (b *BuildingController) DeleteBuildingFacility(c *fiber.Ctx) error {
	buildingID := c.Params("buildingID")
	facilityID := c.Params("facilityID")
	facilityIDInt, err := strconv.Atoi(facilityID)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err2.ErrInvalidFacilityID.Error())
	}

	if err := b.buildingService.DeleteBuildingFacility(c.Context(), buildingID, facilityIDInt); err != nil {
		switch err {
		case err2.ErrFacilityNotFound:
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		default:
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
	}

	return c.Status(fiber.StatusOK).JSON(response.BaseResponse{
		Message: "building facility deleted successfully",
	})
}

func (b *BuildingController) DeleteBuilding(c *fiber.Ctx) error {
	buildingID := c.Params("buildingID")

	if err := b.buildingService.DeleteBuilding(c.Context(), buildingID); err != nil {
		switch err {
		case err2.ErrBuildingNotFound:
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		case err2.ErrBuildingHasReservation:
			return fiber.NewError(fiber.StatusConflict, err.Error())
		default:
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
	}

	return c.Status(fiber.StatusOK).JSON(response.BaseResponse{
		Message: "building deleted successfully",
	})
}
