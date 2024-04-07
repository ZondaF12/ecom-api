package cart

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/zondaf12/ecom-api/service/auth"
	"github.com/zondaf12/ecom-api/types"
	"github.com/zondaf12/ecom-api/utils"
)

type Handler struct {
	store        types.OrderStore
	productStore types.ProductStore
	userStore    types.UserStore
}

func NewHandler(store types.OrderStore, productStore types.ProductStore, userStore types.UserStore) *Handler {
	return &Handler{store, productStore, userStore}
}

func (h *Handler) RegisterRoutes(router *echo.Group) {
	router.POST("/cart/checkout", auth.WithJWTAuth(h.HandleCheckout, h.userStore))
}

func (h *Handler) HandleCheckout(c echo.Context) error {
	userId := auth.GetUserIDFromContext(c.Request().Context())

	var cart types.CartCheckoutPayload
	if err := utils.ParseJSON(c, &cart); err != nil {
		return utils.WriteError(c, http.StatusBadRequest, err)
	}

	if err := utils.Validate.Struct(cart); err != nil {
		errors := err.(validator.ValidationErrors)
		return utils.WriteError(c, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
	}

	// Get Products
	productIDs, err := getCartItemIDs(cart.Items)
	if err != nil {
		return utils.WriteError(c, http.StatusBadRequest, err)
	}

	ps, err := h.productStore.GetProductsByID(productIDs)
	if err != nil {
		return utils.WriteError(c, http.StatusInternalServerError, err)
	}

	orderId, totalPrice, err := h.createOrder(ps, cart.Items, userId)
	if err != nil {
		return utils.WriteError(c, http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, map[string]any{
		"order_id":    orderId,
		"total_price": totalPrice,
	})
}
