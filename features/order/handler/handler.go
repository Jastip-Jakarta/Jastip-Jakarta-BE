package handler

import (
	"fmt"
	"jastip-jakarta/features/order"
	"jastip-jakarta/utils/middlewares"
	"jastip-jakarta/utils/responses"
	"net/http"
	"strconv"

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
	errInsert := handler.orderService.CreateUserOrder(userIdLogin, orderCore)
	if errInsert != nil {
		return c.JSON(http.StatusInternalServerError, responses.WebResponse(errInsert.Error(), nil))
	}

	return c.JSON(http.StatusOK, responses.WebResponse("Berhasil Membuat Orderan Jastip", nil))
}

func (handler *OrderHandler) UpdateUserOrder(c echo.Context) error {
	userIdLogin := middlewares.ExtractTokenUserId(c)
	if userIdLogin == 0 {
		return c.JSON(http.StatusUnauthorized, responses.WebResponse("Silahkan login terlebih dahulu", nil))
	}

	userOrderId, errParse := strconv.ParseUint(c.Param("order_id"), 10, 32)
	if errParse != nil {
		return c.JSON(http.StatusBadRequest, responses.WebResponse("ID order tidak valid", nil))
	}

	updateOrder := UserOrderRequest{}
	errBind := c.Bind(&updateOrder)
	if errBind != nil {
		return c.JSON(http.StatusBadRequest, responses.WebResponse("error bind data. data order not valid", nil))
	}

	orderCore := RequestUpdateToUserOrder(updateOrder)
	errUpdate := handler.orderService.UpdateUserOrder(userIdLogin, uint(userOrderId), orderCore)
	if errUpdate != nil {
		return c.JSON(http.StatusInternalServerError, responses.WebResponse(errUpdate.Error(), nil))
	}

	return c.JSON(http.StatusOK, responses.WebResponse("Order berhasil diperbarui", nil))
}

func (handler *OrderHandler) GetUserOrderWait(c echo.Context) error {
	userIdLogin := middlewares.ExtractTokenUserId(c)
	if userIdLogin == 0 {
		return c.JSON(http.StatusUnauthorized, responses.WebResponse("Silahkan login terlebih dahulu", nil))
	}

	userOrders, err := handler.orderService.GetUserOrderWait(userIdLogin)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.WebResponse(err.Error(), nil))
	}

	var userOrderWaitResponses []UserOrderWaitResponse
	for _, userOrder := range userOrders {
		userOrderWaitResponses = append(userOrderWaitResponses, CoreToResponseUserOrderWait(userOrder))
	}

	return c.JSON(http.StatusOK, responses.WebResponse("Berhasil mendapatkan orderan yang menunggu dikirim", userOrderWaitResponses))
}

func (handler *OrderHandler) GetOrderById(c echo.Context) error {
	orderIdStr := c.Param("order_id")
	orderId, err := strconv.ParseUint(orderIdStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, responses.WebResponse("ID order tidak valid", nil))
	}

	result, errSelect := handler.orderService.GetById(uint(orderId))
	if errSelect != nil {
		return c.JSON(http.StatusNotFound, responses.WebResponse("Order tidak ditemukan", nil))
	}

	var orderResult = CoreToResponseUserOrderById(*result)
	return c.JSON(http.StatusOK, responses.WebResponse("success read data.", orderResult))
}

func (handler *OrderHandler) CreateOrderDetail(c echo.Context) error {
	adminIdLogin := middlewares.ExtractTokenUserId(c)
	if adminIdLogin == 0 {
		return c.JSON(http.StatusUnauthorized, responses.WebResponse("Silahkan login terlebih dahulu", nil))
	}

	orderId, err := strconv.Atoi(c.Param("order_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, responses.WebResponse("ID order tidak valid", nil))
	}

	newOrder := OrderDetailRequest{}
	errBind := c.Bind(&newOrder)
	if errBind != nil {
		return c.JSON(http.StatusBadRequest, responses.WebResponse("error bind data. data order not valid", nil))
	}

	// Fetch UserOrder based on orderId
	userOrder, errFetch := handler.orderService.GetById(uint(orderId))
	if errFetch != nil {
		return c.JSON(http.StatusInternalServerError, responses.WebResponse("Gagal mendapatkan detail order", nil))
	}

	orderCore := RequestToOrderDetail(newOrder, *userOrder)
	errInsert := handler.orderService.CreateOrderDetail(adminIdLogin, uint(orderId), orderCore)
	if errInsert != nil {
		return c.JSON(http.StatusInternalServerError, responses.WebResponse(errInsert.Error(), nil))
	}

	return c.JSON(http.StatusOK, responses.WebResponse("Berhasil Membuat Orderan Jastip", nil))
}

func (handler *OrderHandler) GetUserOrderProcess(c echo.Context) error {
	userIdLogin := middlewares.ExtractTokenUserId(c)
	if userIdLogin == 0 {
		return c.JSON(http.StatusUnauthorized, responses.WebResponse("Silahkan login terlebih dahulu", nil))
	}

	userOrders, err := handler.orderService.GetUserOrderProcess(userIdLogin)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.WebResponse(err.Error(), nil))
	}

	groupedResponses := CoreToGroupedOrderResponse(userOrders, handler.orderService.GetFoto)

	return c.JSON(http.StatusOK, responses.WebResponse("Berhasil mendapatkan orderan yang diproses", groupedResponses))
}

func (handler *OrderHandler) SearchUserOrder(c echo.Context) error {
	userIdLogin := middlewares.ExtractTokenUserId(c)
	if userIdLogin == 0 {
		return c.JSON(http.StatusUnauthorized, responses.WebResponse("Silahkan login terlebih dahulu", nil))
	}

	itemName := c.QueryParam("item_name")
	if itemName == "" {
		return c.JSON(http.StatusBadRequest, responses.WebResponse("Nama barang tidak boleh kosong", nil))
	}

	userOrders, err := handler.orderService.SearchUserOrder(userIdLogin, itemName)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.WebResponse(err.Error(), nil))
	}

	var userOrderResponses []UserOrderWaitResponse
	for _, userOrder := range userOrders {
		userOrderResponses = append(userOrderResponses, CoreToResponseUserOrderWait(userOrder))
	}

	return c.JSON(http.StatusOK, responses.WebResponse("Berhasil mencari orderan", userOrderResponses))
}

func (handler *OrderHandler) GetAllUserOrderWait(c echo.Context) error {
	adminIdLogin := middlewares.ExtractTokenUserId(c)
	if adminIdLogin == 0 {
		return c.JSON(http.StatusUnauthorized, responses.WebResponse("Silahkan login terlebih dahulu", nil))
	}

	userOrders, err := handler.orderService.GetAllUserOrderWait(adminIdLogin)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.WebResponse(err.Error(), nil))
	}

	var userOrderResponses []UserOrderWaitResponse
	for _, userOrder := range userOrders {
		userOrderResponses = append(userOrderResponses, CoreToResponseUserOrderWait(userOrder))
	}

	return c.JSON(http.StatusOK, responses.WebResponse("Berhasil mendapatkan semua orderan yang menunggu diterima", userOrderResponses))
}

func (handler *OrderHandler) GetDeliveryBatchWithRegion(c echo.Context) error {
	adminIdLogin := middlewares.ExtractTokenUserId(c)
	if adminIdLogin == 0 {
		return c.JSON(http.StatusUnauthorized, responses.WebResponse("Silahkan login terlebih dahulu", nil))
	}

	deliveryBatchWithRegion, err := handler.orderService.GetDeliveryBatchWithRegion(adminIdLogin)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.WebResponse(err.Error(), nil))
	}

	groupedResult := CoreToResponseDeliveryBatches(deliveryBatchWithRegion)

	return c.JSON(http.StatusOK, responses.WebResponse("Berhasil mendapatkan batch pengiriman dengan kode wilayah", groupedResult))
}

func (handler *OrderHandler) GetUserOrderNames(c echo.Context) error {
	adminIdLogin := middlewares.ExtractTokenUserId(c)
	if adminIdLogin == 0 {
		return c.JSON(http.StatusUnauthorized, responses.WebResponse("Silahkan login terlebih dahulu", nil))
	}

	code := c.QueryParam("code")
	batch := c.QueryParam("batch")
	if code == "" || batch == "" {
		return c.JSON(http.StatusBadRequest, responses.WebResponse("Code dan batch tidak boleh kosong", nil))
	}

	userOrders, err := handler.orderService.GetNameByUserOrder(adminIdLogin, code, batch)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.WebResponse(err.Error(), nil))
	}

	response := CoreToGetCustomerResponse(userOrders, batch, code)

	return c.JSON(http.StatusOK, responses.WebResponse("Berhasil mendapatkan nama-nama user order", response))
}

func (handler *OrderHandler) GetOrderByNameUser(c echo.Context) error {
	adminIdLogin := middlewares.ExtractTokenUserId(c)
	if adminIdLogin == 0 {
		return c.JSON(http.StatusUnauthorized, responses.WebResponse("Silahkan login terlebih dahulu", nil))
	}

	code := c.QueryParam("code")
	batch := c.QueryParam("batch")
	if code == "" || batch == "" {
		return c.JSON(http.StatusBadRequest, responses.WebResponse("Code dan batch tidak boleh kosong", nil))
	}

	name := c.QueryParam("name")
	if name == "" {
		return c.JSON(http.StatusBadRequest, responses.WebResponse("Name tidak boleh kosong", nil))
	}

	userOrders, err := handler.orderService.GetOrderByUserOrderNameUser(adminIdLogin, code, batch, name)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.WebResponse(err.Error(), nil))
	}

	response := CoreToGroupedAdminOrderResponse(userOrders, batch, code, handler.orderService.GetFoto)

	return c.JSON(http.StatusOK, responses.WebResponse("Berhasil mendapatkan user order", response))
}

func (handler *OrderHandler) UpdateEstimationForOrders(c echo.Context) error {
    adminIdLogin := middlewares.ExtractTokenUserId(c)
    if adminIdLogin == 0 {
        return c.JSON(http.StatusUnauthorized, responses.WebResponse("Silahkan login terlebih dahulu", nil))
    }

    code := c.QueryParam("code")
    batch := c.QueryParam("batch")
    if code == "" || batch == "" {
        return c.JSON(http.StatusBadRequest, responses.WebResponse("Code dan batch tidak boleh kosong", nil))
    }

    var req UpdateEstimationRequest
    errBind := c.Bind(&req)
    if errBind != nil {
        return c.JSON(http.StatusBadRequest, responses.WebResponse("Error bind data. Data order tidak valid", nil))
    }

    estimasiTime, err := ParseEstimationDate(req.Estimation)
    if err != nil {
        return c.JSON(http.StatusBadRequest, responses.WebResponse("Format tanggal tidak valid. Gunakan format dd/mm/yyyy", nil))
    }

    err = handler.orderService.UpdateEstimationForOrders(adminIdLogin, code, batch, estimasiTime)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, responses.WebResponse(err.Error(), nil))
    }

    return c.JSON(http.StatusOK, responses.WebResponse(fmt.Sprintf("Estimasi berhasil diperbarui menjadi %s untuk semua pesanan", estimasiTime.Format("02/01/2006")), nil))
}

func (handler *OrderHandler) UpdateOrderStatus(c echo.Context) error {
    adminIdLogin := middlewares.ExtractTokenUserId(c)
    if adminIdLogin == 0 {
        return c.JSON(http.StatusUnauthorized, responses.WebResponse("Silahkan login terlebih dahulu", nil))
    }

    userOrderIdStr := c.Param("order_id")
    userOrderId, err := strconv.ParseUint(userOrderIdStr, 10, 64)
    if err != nil {
        return c.JSON(http.StatusBadRequest, responses.WebResponse("Invalid userOrderId", nil))
    }

    var req UpdateStatusRequest
    err = c.Bind(&req)
    if err != nil {
        return c.JSON(http.StatusBadRequest, responses.WebResponse("error bind data. data order not valid", nil))
    }

    err = handler.orderService.UpdateOrderStatus(adminIdLogin, uint(userOrderId), req.Status)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, responses.WebResponse(err.Error(), nil))
    }

    return c.JSON(http.StatusOK, responses.WebResponse("Status berhasil diperbarui", nil))
}

func (handler *OrderHandler) UploadFotoPacked(c echo.Context) error {
	adminIdLogin := middlewares.ExtractTokenUserId(c)
	if adminIdLogin == 0 {
		return c.JSON(http.StatusUnauthorized, responses.WebResponse("Silahkan login terlebih dahulu", nil))
	}
	
	uploadRequest := UploadFotoRequest{}
	errBind := c.Bind(&uploadRequest)
	if errBind != nil {
		return c.JSON(http.StatusBadRequest, responses.WebResponse("Error bind data. Data upload tidak valid", nil))
	}
	
	fileHeaderPacked, err := c.FormFile("photo_packed")
	if err != nil && err != http.ErrMissingFile {
		return c.JSON(http.StatusBadRequest, responses.WebResponse("error retrieving the file", nil))
	}

	photoCore := RequestToPhotoOrder(uploadRequest)
	errUpload := handler.orderService.UploadFotoPacked(adminIdLogin, photoCore, fileHeaderPacked)
	if errUpload != nil {
		return c.JSON(http.StatusInternalServerError, responses.WebResponse(errUpload.Error(), nil))
	}

	return c.JSON(http.StatusOK, responses.WebResponse("Berhasil mengunggah foto", nil))
}

func (handler *OrderHandler) UploadFotoReceived(c echo.Context) error {
	adminIdLogin := middlewares.ExtractTokenUserId(c)
	if adminIdLogin == 0 {
		return c.JSON(http.StatusUnauthorized, responses.WebResponse("Silahkan login terlebih dahulu", nil))
	}

	idFoto, err := strconv.Atoi(c.Param("id_foto"))
	if idFoto == 0 {
		return c.JSON(http.StatusBadRequest, responses.WebResponse("id Foto query kosong", nil))
	}
	if err != nil {
		return c.JSON(http.StatusBadRequest, responses.WebResponse("error parsing foto id", nil))
	}
	
	fileHeaderReceived, err := c.FormFile("photo_received")
	if err != nil && err != http.ErrMissingFile {
		return c.JSON(http.StatusBadRequest, responses.WebResponse("error retrieving the file", nil))
	}

	errUpload := handler.orderService.UploadFotoReceived(adminIdLogin, uint(idFoto), fileHeaderReceived)
	if errUpload != nil {
		return c.JSON(http.StatusInternalServerError, responses.WebResponse(errUpload.Error(), nil))
	}

	return c.JSON(http.StatusOK, responses.WebResponse("Berhasil mengunggah foto", nil))
}

func (handler *OrderHandler) GenerateCSVByBatch(c echo.Context) error {
	batch := c.QueryParam("batch")
	if batch == "" {
		return c.JSON(http.StatusBadRequest, responses.WebResponse("Batch pengiriman tidak boleh kosong", nil))
	}

	filePath := "batch_pengiriman_" + batch + ".csv"
	
	err := handler.orderService.GenerateCSVByBatch(batch, filePath)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.WebResponse(err.Error(), nil))
	}

	c.Response().Header().Set(echo.HeaderContentDisposition, fmt.Sprintf("attachment; filename=%s", filePath))
	c.Response().Header().Set(echo.HeaderContentType, "text/csv")
	return c.File(filePath)
}

func (handler *OrderHandler) SearchOrder(c echo.Context) error {
	adminIdLogin := middlewares.ExtractTokenUserId(c)
	if adminIdLogin == 0 {
		return c.JSON(http.StatusUnauthorized, responses.WebResponse("Silahkan login terlebih dahulu", nil))
	}

	searchQuery := c.QueryParam("jastip")

	userOrders, err := handler.orderService.SearchOrders(adminIdLogin, searchQuery)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.WebResponse(err.Error(), nil))
	}

	var userOrderResponses []UserOrderWaitResponse
	for _, userOrder := range userOrders {
		userOrderResponses = append(userOrderResponses, CoreToResponseUserOrderWait(userOrder))
	}

	return c.JSON(http.StatusOK, responses.WebResponse("Berhasil mencari orderan", userOrderResponses))
}

func (handler *OrderHandler) UpdateOrderById(c echo.Context) error {
	adminIdLogin := middlewares.ExtractTokenUserId(c)
	if adminIdLogin == 0 {
		return c.JSON(http.StatusUnauthorized, responses.WebResponse("Silahkan login terlebih dahulu", nil))
	}

	userOrderId, errParse := strconv.ParseUint(c.Param("order_id"), 10, 32)
	if errParse != nil {
		return c.JSON(http.StatusBadRequest, responses.WebResponse("ID order tidak valid", nil))
	}

	updateOrder := UpdateOrderByID{}
	errBind := c.Bind(&updateOrder)
	if errBind != nil {
		return c.JSON(http.StatusBadRequest, responses.WebResponse("error bind data. data order not valid", nil))
	}

	orderCore := RequestToUserOrderUpdate(updateOrder)
	errUpdate := handler.orderService.UpdateOrderByID(adminIdLogin, uint(userOrderId), orderCore)
	if errUpdate != nil {
		return c.JSON(http.StatusInternalServerError, responses.WebResponse(errUpdate.Error(), nil))
	}

	return c.JSON(http.StatusOK, responses.WebResponse("Order berhasil diperbarui", nil))
}

func (handler *OrderHandler) GetOrderSStats(c echo.Context) error {
	adminIdLogin := middlewares.ExtractTokenUserId(c)
	if adminIdLogin == 0 {
		return c.JSON(http.StatusUnauthorized, responses.WebResponse("Silahkan login terlebih dahulu", nil))
	}

	batch := c.Param("batch")
	if batch == "" {
		return c.JSON(http.StatusBadRequest, responses.WebResponse("batch tidak boleh kosong", nil))
	}

	orderStats, err := handler.orderService.FetchRegionStatsByBatch(adminIdLogin, batch)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.WebResponse(err.Error(), nil))
	}

	var orderStatsResponses []RegionBatchStats
	for _, orderStat := range orderStats {
		orderStatsResponses = append(orderStatsResponses, CoreToResponseRegionBatchStats(orderStat))
	}

	return c.JSON(http.StatusOK, responses.WebResponse("Berhasil mendapatkan statistik", orderStatsResponses))
}