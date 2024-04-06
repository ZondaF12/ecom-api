package user

import (
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/zondaf12/ecom-api/types"
)

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) RegisterRoutes(router *echo.Group) {
	router.POST("/login", h.HandleLogin)
	router.POST("/register", h.HandleRegister)
}

func (h *Handler) HandleLogin(c echo.Context) error {
	return c.JSON(http.StatusOK, "Login")
}

func (h *Handler) HandleRegister(c echo.Context) error {
	if c.Request().Body == nil {
		return c.JSON(http.StatusBadRequest, "Invalid payload")
	}

	var payload types.RegisrerUserPayload
	err := json.NewDecoder(c.Request().Body).Decode(payload)

	return c.JSON(http.StatusOK, "Register")
}
