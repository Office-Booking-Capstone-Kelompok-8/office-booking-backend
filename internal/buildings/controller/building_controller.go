package controller

import (
	"office-booking-backend/internal/buildings/service"
	err2 "office-booking-backend/pkg/errors"
	"office-booking-backend/pkg/response"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type BuildingController struct {
	buildingService service.BuildingService
}

func NewBuildingController(buildingService service.BuildingService) *BuildingController {
	return &BuildingController{
		buildingService: buildingService,
	}
}

func (b *BuildingController) GetAllBuildings(c *fiber.Ctx) error {
	q := c.Query("q")
	limit := c.Query("limit", "20")
	page := c.Query("page", "1")

	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err2.ErrInvalidQueryParams.Error())
	}

	pageInt, err := strconv.Atoi(page)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err2.ErrInvalidQueryParams.Error())
	}

	buildings, total, err := b.buildingService.GetAllBuildings(c.Context(), q, limitInt, pageInt)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(response.BaseResponse{
		Message: "user fetched successfully",
		Data:    buildings,
		Meta: fiber.Map{
			"limit": limitInt,
			"page":  pageInt,
			"total": total,
		},
	})
}
