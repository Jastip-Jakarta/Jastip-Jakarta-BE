package handler

import (
	"jastip-jakarta/features/order"
	"jastip-jakarta/utils/middlewares"
	"jastip-jakarta/utils/responses"
	"net/http"

	"github.com/labstack/echo/v4"
)

type OrderHandler struct {
	orderService order.OrderServiceInterface
}

func New(os order.OrderServiceInterface) *OrderHandler {
	return &OrderHandler{
		orderService: os,
	}
}

func (handler *OrderHandler) CreateUserOrder(c echo.Context) error {
	userIdLogin := middlewares.ExtractTokenUserId(c)
	if userIdLogin == 0 {
		return c.JSON(http.StatusUnauthorized, responses.WebResponse("Silahkan login terlebih dahulu", nil))
	}

	newOrder := UserOrderRequest{}
	errBind := c.Bind(&newOrder)
	if errBind != nil {
		return c.JSON(http.StatusBadRequest, responses.WebResponse("error bind data. data order not valid", nil))
	}

	orderCore := RequestToUserOrder(newOrder)
	errInsert := handler.orderService.CreateOrder(userIdLogin, orderCore)
	if errInsert != nil {
		return c.JSON(http.StatusInternalServerError, responses.WebResponse(errInsert.Error(), nil))
	}

	return c.JSON(http.StatusOK, responses.WebResponse("Berhasil Membuat Orderan Jastip", nil))
}