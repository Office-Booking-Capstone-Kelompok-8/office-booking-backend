package controller

import (
	"errors"
	"office-booking-backend/internal/building/service"
	err2 "office-booking-backend/pkg/errors"
	"office-booking-backend/pkg/response"
	"office-booking-backend/pkg/utils/validator"
	"strconv"
	"time"

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

func (b *BuildingController) GetAllBuildings(c *fiber.Ctx) error {
	q := c.Query("q")
	city := c.Query("city", "0")
	district := c.Query("district", "0")
	startDate := c.Query("startDate", "0001-01-01")
	endDate := c.Query("endDate", "0001-01-01")
	limit := c.Query("limit", "20")
	page := c.Query("page", "1")

	// Parse startDate to time.Time YYYY-MM-DD format
	startDateParsed, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err2.ErrInvalidQueryParams.Error())
	}

	// Parse endDate to time.Time YYYY-MM-DD format
	endDateParsed, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err2.ErrInvalidQueryParams.Error())
	}

	cityInt, err := strconv.Atoi(city)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err2.ErrInvalidQueryParams.Error())
	}

	districtInt, err := strconv.Atoi(district)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err2.ErrInvalidQueryParams.Error())
	}

	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err2.ErrInvalidQueryParams.Error())
	}

	pageInt, err := strconv.Atoi(page)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err2.ErrInvalidQueryParams.Error())
	}

	buildings, total, err := b.buildingService.GetAllBuildings(c.Context(), q, cityInt, districtInt, startDateParsed, endDateParsed, limitInt, pageInt)
	if err != nil {
		switch err {
		case err2.ErrStartDateAfterEndDate:
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		default:
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
	}

	return c.Status(fiber.StatusOK).JSON(response.BaseResponse{
		Message: "buildings fetched successfully",
		Data:    buildings,
		Meta: fiber.Map{
			"limit": limitInt,
			"page":  pageInt,
			"total": total,
		},
	})
}

func (b *BuildingController) GetBuldingDetailByID(c *fiber.Ctx) error {
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
