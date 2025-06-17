package user

import (
	"net/http"

	"booking/internal/user/model"
	"booking/pkg/response"
	"booking/shared/constants"

	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	userService UserServiceInterface
}

func NewUserHandler(userService UserServiceInterface) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// helper: parse & validate payload
func bindAndValidate[T any](c echo.Context, dst *T) error {
	if err := c.Bind(dst); err != nil {
		return response.BadRequest(c, "invalid request payload", err)
	}
	if err := c.Validate(dst); err != nil {
		return response.ValidationError(c, err)
	}
	return nil
}

// helper: get user_id safely
func getUserID(c echo.Context) (string, error) {
	userID, ok := c.Get("user_id").(string)
	if !ok || userID == "" {
		return "", echo.NewHTTPError(http.StatusUnauthorized, "user ID not found in context")
	}
	return userID, nil
}

func (h *UserHandler) Register(c echo.Context) error {
	var input model.RegisterInput
	if err := bindAndValidate(c, &input); err != nil {
		return err
	}

	user, err := h.userService.Register(c.Request().Context(), input)
	if err != nil {
		return response.BadRequest(c, "registration failed", err)
	}

	return response.Success(c, http.StatusCreated, "User registered successfully", user)
}

func (h *UserHandler) Login(c echo.Context) error {
	var input model.LoginInput
	if err := bindAndValidate(c, &input); err != nil {
		return err
	}

	token, err := h.userService.Login(c.Request().Context(), input)
	if err != nil {
		return response.Unauthorized(c, "invalid credentials", err)
	}

	return response.Success(c, http.StatusOK, "Login successful", map[string]string{
		"token": token,
	})
}

func (h *UserHandler) Logout(c echo.Context) error {
	userID, err := getUserID(c)
	if err != nil {
		return response.Unauthorized(c, "unauthorized", err)
	}

	if err := h.userService.Logout(c.Request().Context(), userID); err != nil {
		return response.InternalServerError(c, "logout failed", err)
	}

	return response.Success(c, http.StatusOK, "Logout successful", nil)
}

func (h *UserHandler) GetMe(c echo.Context) error {
	userID, err := getUserID(c)
	if err != nil {
		return response.Unauthorized(c, "unauthorized", err)
	}

	user, err := h.userService.GetUserByID(c.Request().Context(), userID)
	if err != nil {
		return response.NotFound(c, "user not found", err)
	}

	return response.Success(c, http.StatusOK, "User profile retrieved successfully", user)
}

func (h *UserHandler) UpdateProfile(c echo.Context) error {
	userID, err := getUserID(c)
	if err != nil {
		return response.Unauthorized(c, "unauthorized", err)
	}

	var input model.UpdateProfileInput
	if err := bindAndValidate(c, &input); err != nil {
		return err
	}

	user, err := h.userService.UpdateProfile(c.Request().Context(), userID, input)
	if err != nil {
		return response.BadRequest(c, "failed to update profile", err)
	}

	return response.Success(c, http.StatusOK, "Profile updated successfully", user)
}

func (h *UserHandler) DeleteAccount(c echo.Context) error {
	userID, err := getUserID(c)
	if err != nil {
		return response.Unauthorized(c, "unauthorized", err)
	}

	if err := h.userService.DeleteAccount(c.Request().Context(), userID); err != nil {
		return response.InternalServerError(c, "failed to delete account", err)
	}

	return response.Success(c, http.StatusOK, "Account deleted successfully", nil)
}
func (h *UserHandler) GetAllUsers(c echo.Context) error {
	users, err := h.userService.GetAllUsers(c.Request().Context())
	if err != nil {
		return response.InternalServerError(c, "failed to get all users", err)
	}

	return response.Success(c, http.StatusOK, "Users retrieved successfully", users)
}

// UpdateUserRole mengubah role user oleh superadmin
func (h *UserHandler) UpdateUserRole(c echo.Context) error {
	adminID, err := getUserID(c)
	if err != nil {
		return response.Unauthorized(c, "unauthorized", err)
	}

	targetUserID := c.Param("id")
	var input struct {
		Role constants.Role `json:"role" validate:"required"`
	}

	if err := bindAndValidate(c, &input); err != nil {
		return err
	}

	user, err := h.userService.UpdateUserRole(c.Request().Context(), adminID, targetUserID, input.Role)
	if err != nil {
		return response.BadRequest(c, "failed to update user role", err)
	}

	return response.Success(c, http.StatusOK, "User role updated successfully", user)
}
