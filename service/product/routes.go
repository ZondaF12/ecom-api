package product

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/zondaf12/ecom-api/types"
	"github.com/zondaf12/ecom-api/utils"
)

type Handler struct {
	store types.ProductStore
}

func NewHandler(store types.ProductStore) *Handler {
	return &Handler{store}
}

func (h *Handler) RegisterRoutes(router *echo.Group) {
	router.GET("/products", h.HandleGetProducts)
	router.POST("/products", h.HandlerCreateProduct)
}

func (h *Handler) HandleGetProducts(c echo.Context) error {
	ps, err := h.store.GetProducts()
	if err != nil {
		return utils.WriteError(c, http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, ps)
}

func (h *Handler) HandlerCreateProduct(c echo.Context) error {
	// Parse payload
	var product types.CreateProductPayload
	if err := utils.ParseJSON(c, &product); err != nil {
		return utils.WriteError(c, http.StatusBadRequest, err)
	}

	// Validate payload
	if err := utils.Validate.Struct(product); err != nil {
		errors := err.(validator.ValidationErrors)
		return utils.WriteError(c, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
	}

	err := h.store.CreateNewProduct(product)
	if err != nil {
		return utils.WriteError(c, http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusCreated, product)
}
