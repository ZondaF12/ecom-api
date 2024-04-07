package user

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/zondaf12/ecom-api/config"
	"github.com/zondaf12/ecom-api/service/auth"
	"github.com/zondaf12/ecom-api/types"
	"github.com/zondaf12/ecom-api/utils"
)

type Handler struct {
	store types.UserStore
}

func NewHandler(store types.UserStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *echo.Group) {
	router.POST("/login", h.HandleLogin)
	router.POST("/register", h.HandleRegister)
}

func (h *Handler) HandleLogin(c echo.Context) error {
	// Parse payload
	var payload types.LoginUserPayload
	if err := utils.ParseJSON(c, &payload); err != nil {
		return utils.WriteError(c, http.StatusBadRequest, err)
	}

	// Validate payload
	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		return utils.WriteError(c, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
	}

	u, err := h.store.GetUserByEmail(payload.Email)
	if err != nil {
		return utils.WriteError(c, http.StatusBadRequest, fmt.Errorf("not found, invalid email or password"))
	}

	if !auth.ComparePassword(u.Password, []byte(payload.Password)) {
		return utils.WriteError(c, http.StatusBadRequest, fmt.Errorf("not found, invalid email or password"))
	}

	secret := []byte(config.Envs.JWTSecret)
	token, err := auth.CreateJWT([]byte(secret), u.ID)
	if err != nil {
		return utils.WriteError(c, http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, map[string]string{"token": token})
}

func (h *Handler) HandleRegister(c echo.Context) error {
	// Parse payload
	var payload types.RegisterUserPayload
	if err := utils.ParseJSON(c, &payload); err != nil {
		return utils.WriteError(c, http.StatusBadRequest, err)
	}

	// Validate payload
	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		return utils.WriteError(c, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
	}

	// Check if user exists
	_, err := h.store.GetUserByEmail(payload.Email)
	if err == nil {
		return utils.WriteError(c, http.StatusBadRequest, fmt.Errorf("user with email %s already exists", payload.Email))
	}

	hashedPassword, err := auth.HashPassword(payload.Password)
	if err != nil {
		return utils.WriteError(c, http.StatusInternalServerError, err)
	}

	// Create user
	err = h.store.CreateUser(types.User{
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		Email:     payload.Email,
		Password:  hashedPassword,
	})

	if err != nil {
		utils.WriteError(c, http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusCreated, "User Created")
}
