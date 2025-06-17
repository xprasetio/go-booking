package response

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

func Success(c echo.Context, statusCode int, message string, data interface{}) error {
	return c.JSON(statusCode, Response{
		Status:  "success",
		Message: message,
		Data:    data,
	})
}

func Error(c echo.Context, code int, message string, err error) error {
	response := Response{
		Status:  "error",
		Message: message,
	}

	if err != nil {
		response.Error = err.Error()
	}

	return c.JSON(code, response)
}

func BadRequest(c echo.Context, message string, err error) error {
	return Error(c, http.StatusBadRequest, message, err)
}

func Created(c echo.Context, message string, data any) error {
	return Success(c, http.StatusCreated, message, data)
}

func Unauthorized(c echo.Context, message string, err error) error {
	return Error(c, http.StatusUnauthorized, message, err)
}

func Forbidden(c echo.Context, message string, err error) error {
	return Error(c, http.StatusForbidden, message, err)
}

func Ok(c echo.Context, message string, data any) error {
	return Success(c, http.StatusOK, message, data)
}

// ValidationError formats validator.ValidationErrors into readable JSON
func ValidationError(c echo.Context, err error) error {
	if ve, ok := err.(validator.ValidationErrors); ok {
		errors := make(map[string]string)
		for _, e := range ve {
			switch e.Tag() {
			case "required":
				errors[e.Field()] = "This field is required"
			case "email":
				errors[e.Field()] = "Invalid email format"
			case "min":
				errors[e.Field()] = "Minimum length is " + e.Param()
			default:
				errors[e.Field()] = "Invalid value"
			}
		}
		return c.JSON(http.StatusBadRequest, Response{
			Status:  "error",
			Message: "Validation failed",
			Error:   errors,
		})
	}

	// fallback
	return Error(c, http.StatusBadRequest, "Bad request", err)
}

func InternalServerError(c echo.Context, message string, err error) error {
	return Error(c, http.StatusInternalServerError, message, err)
}

func NotFound(c echo.Context, message string, err error) error {
	return c.JSON(http.StatusNotFound, Response{
		Status:  "error",
		Message: message,
		Error:   err.Error(),
	})
}
