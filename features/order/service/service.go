package service

import (
	"errors"
	"jastip-jakarta/features/admin"
	"jastip-jakarta/features/order"
	"mime/multipart"
	"time"
)

type orderService struct {
	orderData    order.OrderDataInterface
	adminService admin.AdminServiceInterface
}

func New(repo order.OrderDataInterface, adminService admin.AdminServiceInterface) order.OrderServiceInterface {
	return &orderService{
		orderData:    repo,
		adminService: adminService,
	}
}

// CreateOrder implements order.OrderServiceInterface.
func (o *orderService) CreateUserOrder(userIdLogin int, inputOrder order.UserOrder) error {
	if inputOrder.ItemName == "" {
		return errors.New("nama Barang harus diisi")
	}
	if inputOrder.TrackingNumber == "" {
		return errors.New("nomor resi harus diisi")
	}
	if inputOrder.OnlineStore == "" {
		return errors.New("nama toko online harus diisi")
	}
	if inputOrder.WhatsAppNumber == 0 {
		return errors.New("nomor WhatsApp harus diisi")
	}
	if inputOrder.RegionCode == "" {
		return errors.New("kode wilayah harus diisi")
	}

	_, err := o.adminService.GettByIdRegion(inputOrder.RegionCode)
	if err != nil {
		return err
	}

	return o.orderData.InsertUserOrder(userIdLogin, inputOrder)
}

// UpdateUserOrder implements order.OrderServiceInterface.
func (o *orderService) UpdateUserOrder(userIdLogin int, userOrderId uint, inputOrder order.UserOrder) error {
	// Mengecek status terlebih dahulu
	status, err := o.orderData.CheckOrderStatus(userOrderId)
	if err != nil {
		return err
	}

	// Melakukan update jika status 'Menunggu Diterima'
	if status == "Menunggu Diterima" {
		err = o.orderData.PutUserOrder(userIdLogin, userOrderId, inputOrder)
		if err != nil {
			return err
		}
	} else {
		return errors.New("order tidak dapat diupdate karena status bukan 'Menunggu Diterima'")
	}

	return nil
}

// SelectUserOrderWait implements order.OrderServiceInterface.
func (o *orderService) GetUserOrderWait(userIdLogin int) ([]order.UserOrder, error) {
	userOrders, err := o.orderData.SelectUserOrderWait(userIdLogin)
	if err != nil {
		return nil, err
	}
	return userOrders, nil
}

// GetById implements order.OrderServiceInterface.
func (o *orderService) GetById(IdOrder uint) (*order.UserOrder, error) {
	result, err := o.orderData.SelectById(IdOrder)
	return result, err
}

// CreateOrderDetail implements order.OrderServiceInterface.
func (o *orderService) CreateOrderDetail(adminIdLogin int, userOrderId uint, inputOrder order.OrderDetail) error {
	adminCheck, err := o.adminService.GetById(adminIdLogin)
	if err != nil || adminCheck == nil {
		return errors.New("anda bukan admin")
	}

	orderIdCheck, err := o.orderData.SelectById(userOrderId)
	if err != nil {
		return err
	}

	if orderIdCheck.ID != userOrderId {
		return errors.New("ID Order tidak ditemukan atau salah")
	}

	if inputOrder.Status == "" {
		return errors.New("status Harus Di Isi")
	}

	if inputOrder.WeightItem == 0 {
		return errors.New("berat Tidak Boleh Nol")
	}

	if *inputOrder.DeliveryBatchID == "" {
		return errors.New("batch Pengiriman Tidak Boleh Kosong")
	}

	if inputOrder.TrackingNumberJastip == "" {
		return errors.New("no Resi JASTIP Tidak Boleh Kosong")
	}

	return o.orderData.InsertOrderDetail(adminIdLogin, userOrderId, inputOrder)
}

// GetUserOrderProcess implements order.OrderServiceInterface.
func (o *orderService) GetUserOrderProcess(userIdLogin int) ([]order.UserOrder, error) {
	userOrders, err := o.orderData.SelectUserOrderProcess(userIdLogin)
	if err != nil {
		return nil, err
	}
	return userOrders, nil
}

// SearchUserOrder implements order.OrderServiceInterface.
func (o *orderService) SearchUserOrder(userIdLogin int, itemName string) ([]order.UserOrder, error) {
	userOrders, err := o.orderData.SearchUserOrder(userIdLogin, itemName)
	if err != nil {
		return nil, err
	}
	return userOrders, nil
}

// GetAllUserOrderWait implements order.OrderServiceInterface.
func (o *orderService) GetAllUserOrderWait(adminIdLogin int) ([]order.UserOrder, error) {
	adminCheck, err := o.adminService.GetById(adminIdLogin)
	if err != nil || adminCheck == nil {
		return nil, errors.New("anda bukan admin")
	}

	userOrders, err := o.orderData.SelectAllUserOrderWait()
	if err != nil {
		return nil, err
	}

	return userOrders, nil
}

// GetDeliveryBatchWithRegion implements order.OrderServiceInterface.
func (o *orderService) GetDeliveryBatchWithRegion(adminIdLogin int) ([]order.DeliveryBatchWithRegion, error) {
	adminCheck, err := o.adminService.GetById(adminIdLogin)
	if err != nil || adminCheck == nil {
		return nil, errors.New("anda bukan admin")
	}

	deliveryBatchWithRegion, err := o.orderData.FetchDeliveryBatchWithRegion()
	if err != nil {
		return nil, err
	}
	return deliveryBatchWithRegion, nil
}

// GetNameByUserOrder implements order.OrderServiceInterface.
func (o *orderService) GetNameByUserOrder(adminIdLogin int, code, batch string) ([]order.UserOrder, error) {
	adminCheck, err := o.adminService.GetById(adminIdLogin)
	if err != nil || adminCheck == nil {
		return nil, errors.New("anda bukan admin")
	}

	codeCheck, err := o.adminService.GettByIdRegion(code)
	if err != nil || codeCheck == nil {
		return nil, errors.New("code region tidak ada")
	}

	batchCheck, err := o.adminService.GetDeliveryBatch(batch)
	if err != nil || batchCheck == nil {
		return nil, errors.New("delivery batch tidak ada")
	}

	userOrders, err := o.orderData.SelectNameByUserOrder(code, batch)
	if err != nil {
		return nil, err
	}
	return userOrders, nil
}

// GetOrderByUserOrderNameUser implements order.OrderServiceInterface.
func (o *orderService) GetOrderByUserOrderNameUser(adminIdLogin int, code string, batch string, name string) ([]order.UserOrder, error) {
	adminCheck, err := o.adminService.GetById(adminIdLogin)
	if err != nil || adminCheck == nil {
		return nil, errors.New("anda bukan admin")
	}

	codeCheck, err := o.adminService.GettByIdRegion(code)
	if err != nil || codeCheck == nil {
		return nil, errors.New("code region tidak ada")
	}

	batchCheck, err := o.adminService.GetDeliveryBatch(batch)
	if err != nil || batchCheck == nil {
		return nil, errors.New("delivery batch tidak ada")
	}

	userOrders, err := o.orderData.SelectOrderByUserOrderNameUser(code, batch, name)
	if err != nil {
		return nil, err
	}
	return userOrders, nil
}

// UpdateEstimationForOrders implements order.OrderServiceInterface.
func (o *orderService) UpdateEstimationForOrders(adminIdLogin int, code, batch string, estimation *time.Time) error {
	adminCheck, err := o.adminService.GetById(adminIdLogin)
	if err != nil || adminCheck == nil {
		return errors.New("anda bukan admin")
	}

	codeCheck, err := o.adminService.GettByIdRegion(code)
	if err != nil || codeCheck == nil {
		return errors.New("code region tidak ada")
	}

	batchCheck, err := o.adminService.GetDeliveryBatch(batch)
	if err != nil || batchCheck == nil {
		return errors.New("delivery batch tidak ada")
	}

	err = o.orderData.UpdateEstimationForOrders(code, batch, estimation)
	if err != nil {
		return err
	}

	return nil
}

// UpdateOrderStatus implements order.OrderServiceInterface.
func (o *orderService) UpdateOrderStatus(adminIdLogin int, userOrderId uint, status string) error {
	adminCheck, err := o.adminService.GetById(adminIdLogin)
	if err != nil || adminCheck.Role != "Perwakilan" {
		return errors.New("anda bukan admin perwakilan")
	}

	err = o.orderData.UpdateOrderStatus(userOrderId, status)
	if err != nil {
		return err
	}
	return nil
}

// UploadFotoPacked implements order.OrderServiceInterface.
func (o *orderService) UploadFotoPacked(adminIdLogin int, inputOrder order.PhotoOrder, photoPacked *multipart.FileHeader) error {
	adminCheck, err := o.adminService.GetById(adminIdLogin)
	if err != nil || adminCheck.Role != "Jakarta" {
		return errors.New("anda bukan admin Jakarta")
	}

	codeCheck, err := o.adminService.GettByIdRegion(inputOrder.RegionCodeID)
	if err != nil || codeCheck == nil {
		return errors.New("code region tidak ada")
	}

	batchCheck, err := o.adminService.GetDeliveryBatch(inputOrder.DeliveryBatchID)
	if err != nil || batchCheck == nil {
		return errors.New("delivery batch tidak ada")
	}

	err = o.orderData.UploadFotoPacked(inputOrder, photoPacked)
	if err != nil {
		return err
	}
	return nil
}

// UploadFotoReceived implements order.OrderServiceInterface.
func (o *orderService) UploadFotoReceived(adminIdLogin int, idFoto uint, photoReceived *multipart.FileHeader) error {
	adminCheck, err := o.adminService.GetById(adminIdLogin)
	if err != nil || adminCheck.Role != "Perwakilan" {
		return errors.New("anda bukan admin perwakilan")
	}

	err = o.orderData.UploadFotoReceived(idFoto, photoReceived)
	if err != nil {
		return err
	}
	return nil
}

// GenerateCSVByBatch implements order.OrderServiceInterface.
func (o *orderService) GenerateCSVByBatch(batch, filePath string) error {
	batchCheck, err := o.adminService.GetDeliveryBatch(batch)
	if err != nil || batchCheck == nil {
		return errors.New("batch tidak ada")
	}

	err = o.orderData.GenerateCSVByBatch(batch, filePath)
	if err != nil {
		return err
	}

	return nil
}

// GetFoto implements order.OrderServiceInterface.
func (o *orderService) GetFoto(batch string, code string, userId int) (*order.PhotoOrder, error) {
	fotoOrders, err := o.orderData.GetFoto(batch, code, userId)
	if err != nil {
		return nil, err
	}
	return fotoOrders, nil
}

// SearchOrders implements order.OrderServiceInterface.
func (o *orderService) SearchOrders(adminIdLogin int, searchQuery string) ([]order.UserOrder, error) {
	adminCheck, err := o.adminService.GetById(adminIdLogin)
	if err != nil || adminCheck == nil {
		return nil, errors.New("anda bukan admin")
	}

	if searchQuery == "" {
		return nil, errors.New("anda belum mengetikan sesuatu")
	}

	orderResponse, err := o.orderData.SearchOrders(searchQuery)
	if err != nil {
		return nil, err
	}

	return orderResponse, nil
}

// UpdateOrderByID implements order.OrderServiceInterface.
func (o *orderService) UpdateOrderByID(adminIdLogin int, orderID uint, inputOrder order.UpdateOrderByID) error {
	adminCheck, err := o.adminService.GetById(adminIdLogin)
	if err != nil || adminCheck.Role != "Super" {
		return errors.New("anda bukan admin super")
	}

	_, err = o.orderData.SelectById(orderID)
	if err != nil {
		return errors.New("order tidak ada")
	}

	err = o.orderData.UpdateOrderByID(orderID, inputOrder)
	if err != nil {
		return err
	}

	return nil
}

// FetchRegionStatsByBatch implements order.OrderServiceInterface.
func (o *orderService) FetchRegionStatsByBatch(adminIdLogin int, batch string) ([]order.RegionBatchStats, error) {
	adminCheck, err := o.adminService.GetById(adminIdLogin)
	if err != nil || adminCheck.Role != "Super" {
		return nil, errors.New("anda bukan admin super")
	}

	statsResponse, err := o.orderData.FetchRegionStatsByBatch(batch)
	if err != nil {
		return nil, err
	}

	return statsResponse, nil
}
