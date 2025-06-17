package space

import (
	"net/http"

	categoryService "booking/internal/category"
	spaceModel "booking/internal/space/model"
	"booking/pkg/response"

	"github.com/labstack/echo/v4"
)

type SpaceHandler struct {
	spaceService    SpaceServiceInterface
	categoryService categoryService.CategoryServiceInterface
}

func NewSpaceHandler(spaceService SpaceServiceInterface, categoryService categoryService.CategoryServiceInterface) *SpaceHandler {
	return &SpaceHandler{
		spaceService:    spaceService,
		categoryService: categoryService,
	}
}

func (h *SpaceHandler) Create(c echo.Context) error {
	var input spaceModel.CreateSpaceInput
	if err := c.Bind(&input); err != nil {
		return response.BadRequest(c, "invalid request payload", err)
	}

	if err := c.Validate(&input); err != nil {
		return response.BadRequest(c, err.Error(), nil)
	}

	// Get category from database using CategoryService
	category, err := h.categoryService.GetByID(input.CategoryID.String())
	if err != nil {
		return response.NotFound(c, "category not found", err)
	}

	space, err := h.spaceService.Create(input, category)
	if err != nil {
		return response.BadRequest(c, "failed to create space", err)
	}

	return response.Success(c, http.StatusCreated, "Space created successfully", space)
}

func (h *SpaceHandler) GetAll(c echo.Context) error {
	spaces, err := h.spaceService.GetAll()
	if err != nil {
		return response.InternalServerError(c, "failed to get spaces", err)
	}

	return response.Success(c, http.StatusOK, "Spaces retrieved successfully", spaces)
}

func (h *SpaceHandler) GetByID(c echo.Context) error {
	id := c.Param("id")
	space, err := h.spaceService.GetByID(id)
	if err != nil {
		return response.NotFound(c, "space not found", err)
	}

	return response.Success(c, http.StatusOK, "Space retrieved successfully", space)
}

func (h *SpaceHandler) Update(c echo.Context) error {
	id := c.Param("id")
	var input spaceModel.CreateSpaceInput
	if err := c.Bind(&input); err != nil {
		return response.BadRequest(c, "invalid request payload", err)
	}

	if err := c.Validate(&input); err != nil {
		return response.BadRequest(c, err.Error(), nil)
	}

	// Verify category exists
	if _, err := h.categoryService.GetByID(input.CategoryID.String()); err != nil {
		return response.NotFound(c, "category not found", err)
	}

	space, err := h.spaceService.Update(id, input)
	if err != nil {
		return response.BadRequest(c, "failed to update space", err)
	}

	return response.Success(c, http.StatusOK, "Space updated successfully", space)
}

func (h *SpaceHandler) Delete(c echo.Context) error {
	id := c.Param("id")
	if err := h.spaceService.Delete(id); err != nil {
		return response.BadRequest(c, "failed to delete space", err)
	}

	return response.Success(c, http.StatusOK, "Space deleted successfully", nil)
}
