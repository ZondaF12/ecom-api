package utils

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

var Validate = validator.New()

func ParseJSON(c echo.Context, payload any) error {
	if c.Request().Body == nil {
		return c.JSON(http.StatusBadRequest, "missing request body")
	}

	return json.NewDecoder(c.Request().Body).Decode(payload)
}

func WriteError(c echo.Context, status int, err error) error {
	return c.JSON(status, map[string]string{"error": err.Error()})
}
