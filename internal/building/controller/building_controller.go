package controller

import (
	"errors"
	"github.com/golang-jwt/jwt/v4"
	dto2 "office-booking-backend/internal/building/dto"
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

func (b *BuildingController) UploadBuildingPicture(c *fiber.Ctx) error {
	form, err := c.MultipartForm()
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err2.ErrInvalidRequestBody.Error())
	}

	// Get buildingID from form
	buildingID := form.Value["buildingID"][0]
	if buildingID == "" {
		return fiber.NewError(fiber.StatusBadRequest, err2.ErrInvalidRequestBody.Error())
	}

	// Get file alt text from form
	altText := form.Value["alt"][0]
	if altText == "" {
		return fiber.NewError(fiber.StatusBadRequest, err2.ErrInvalidRequestBody.Error())
	}

	// Get file from form
	file, err := form.File["picture"][0].Open()
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err2.ErrInvalidRequestBody.Error())
	}

	result, err := b.buildingService.AddBuildingPicture(c.Context(), buildingID, altText, file)
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

	return c.Status(fiber.StatusOK).JSON(response.BaseResponse{
		Message: "building picture uploaded successfully",
		Data:    result,
	})
}

func (b *BuildingController) CreateBuilding(c *fiber.Ctx) error {
	building := new(dto2.CreateBuildingRequest)
	if err := c.BodyParser(building); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err2.ErrInvalidRequestBody.Error())
	}

	if err := b.validator.Validate(building); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err2.ErrInvalidRequestBody.Error())
	}

	if err := b.buildingService.CreateBuilding(c.Context(), building); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(response.BaseResponse{
		Message: "building created successfully",
	})
}
