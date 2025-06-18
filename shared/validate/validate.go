package validate

import (
	"booking/pkg/response"

	"github.com/labstack/echo/v4"
)

func BindAndValidate[T any](c echo.Context, dst *T) error {
	if err := c.Bind(dst); err != nil {
		return response.BadRequest(c, "invalid request payload", err)
	}
	if err := c.Validate(dst); err != nil {
		return response.ValidationError(c, err)
	}
	return nil
}
