package facility

import (
	"booking/pkg/response"
	"net/http"

	facilityModel "booking/internal/facility/model"

	"github.com/labstack/echo/v4"
)

type FacilityHandler struct {
	facilityService FacilityServiceInterface
}

func NewFacilityHandler(facilityService FacilityServiceInterface) *FacilityHandler {
	return &FacilityHandler{
		facilityService: facilityService,
	}
}

func (h *FacilityHandler) Create(c echo.Context) error {
	var input facilityModel.CreateFacilityInput
	if err := c.Bind(&input); err != nil {
		return response.BadRequest(c, "invalid request payload", err)
	}

	if err := c.Validate(&input); err != nil {
		return response.BadRequest(c, err.Error(), nil)
	}

	// Check if facility already exists
	exists, err := h.facilityService.CheckExists(input.Name)
	if err != nil {
		return response.InternalServerError(c, "failed to check facility existence", err)
	}
	if exists {
		return response.BadRequest(c, "facility already exists", nil)
	}

	facility, err := h.facilityService.Create(input)
	if err != nil {
		return response.BadRequest(c, "failed to create facility", err)
	}

	return response.Success(c, http.StatusCreated, "Facility created successfully", facility)
}

func (h *FacilityHandler) GetAll(c echo.Context) error {
	facilities, err := h.facilityService.GetAll()
	if err != nil {
		return response.InternalServerError(c, "failed to get facilities", err)
	}

	return response.Success(c, http.StatusOK, "Facilities retrieved successfully", facilities)
}

func (h *FacilityHandler) GetByID(c echo.Context) error {
	id := c.Param("id")
	facility, err := h.facilityService.GetByID(id)
	if err != nil {
		return response.NotFound(c, "facility not found", err)
	}

	return response.Success(c, http.StatusOK, "Facility retrieved successfully", facility)
}

func (h *FacilityHandler) Update(c echo.Context) error {
	id := c.Param("id")
	var input facilityModel.CreateFacilityInput
	if err := c.Bind(&input); err != nil {
		return response.BadRequest(c, "invalid request payload", err)
	}

	if err := c.Validate(&input); err != nil {
		return response.BadRequest(c, err.Error(), nil)
	}

	// Check if facility exists
	exists, err := h.facilityService.CheckExists(input.Name)
	if err != nil {
		return response.InternalServerError(c, "failed to check facility existence", err)
	}
	if exists {
		return response.BadRequest(c, "facility name already exists", nil)
	}

	facility, err := h.facilityService.Update(id, input)
	if err != nil {
		return response.BadRequest(c, "failed to update facility", err)
	}

	return response.Success(c, http.StatusOK, "Facility updated successfully", facility)
}

func (h *FacilityHandler) Delete(c echo.Context) error {
	id := c.Param("id")
	if err := h.facilityService.Delete(id); err != nil {
		return response.BadRequest(c, "failed to delete facility", err)
	}

	return response.Success(c, http.StatusOK, "Facility deleted successfully", nil)
}
