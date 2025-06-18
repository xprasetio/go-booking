package spacefacility

import (
	"booking/internal/space_facility/model"
	"booking/pkg/response"
	"booking/shared/validate"
	"net/http"

	"github.com/labstack/echo/v4"
)

type SpaceFacilityHandler struct {
	sFService SpaceFacilityServiceInterface
}

func NewSpaceFacilityHandler(sFService SpaceFacilityServiceInterface) *SpaceFacilityHandler {
	return &SpaceFacilityHandler{
		sFService: sFService,
	}
}

func (h *SpaceFacilityHandler) Create(c echo.Context) error {
	var input model.SpaceFacilityInput
	if err := validate.BindAndValidate(c, &input); err != nil {
		return err
	}
	sf, err := h.sFService.Create(c.Request().Context(), input)
	if err != nil {
		return response.BadRequest(c, "failed to create space facility", err)
	}
	return response.Success(c, http.StatusCreated, "space facility created", sf)
}
