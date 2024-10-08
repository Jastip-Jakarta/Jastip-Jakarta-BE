package data

import (
	"fmt"
	"jastip-jakarta/features/order"
	"jastip-jakarta/utils/cloudinary"
	"jastip-jakarta/utils/csv"
	"log"
	"mime/multipart"
	"time"

	"gorm.io/gorm"
)

type orderQuery struct {
	db  *gorm.DB
	cld cloudinary.CloudinaryUploaderInterface
	csv csv.CSVGeneratorInterface
}

func New(db *gorm.DB, cloudinaryUploader cloudinary.CloudinaryUploaderInterface, csvGenerator csv.CSVGeneratorInterface) order.OrderDataInterface {
	return &orderQuery{
		db:  db,
		cld: cloudinaryUploader,
		csv: csvGenerator,
	}
}

// InsertUserOrder implements order.OrderDataInterface.
func (o *orderQuery) InsertUserOrder(userIdLogin int, inputOrder order.UserOrder) error {
	newOrder := UserOrderToModel(inputOrder)
	newOrder.UserID = uint(userIdLogin)

	result := o.db.Create(&newOrder)
	if result.Error != nil {
		return result.Error
	}

	detailOrder := OrderDetail{
		UserOrderID: newOrder.ID,
		Status:      "Menunggu Diterima",
		AdminID:     nil,
	}

	result = o.db.Create(&detailOrder)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// PutUserOrder implements order.OrderDataInterface.
func (o *orderQuery) PutUserOrder(userIdLogin int, userOrderId uint, inputOrder order.UserOrder) error {
	putOrder := UserOrderToModel(inputOrder)
	result := o.db.Model(&UserOrder{}).Where("id = ? AND user_id = ?", userOrderId, userIdLogin).Updates(putOrder)
	return result.Error
}

// CheckOrderStatus implements order.OrderDataInterface.
func (o *orderQuery) CheckOrderStatus(userOrderId uint) (string, error) {
	var adminOrder OrderDetail
	result := o.db.Select("status").Where("user_order_id = ?", userOrderId).First(&adminOrder)
	if result.Error != nil {
		return "", result.Error
	}
	return adminOrder.Status, nil
}

// SelectUserOrderWait implements order.OrderDataInterface.
func (o *orderQuery) SelectUserOrderWait(userIdLogin int) ([]order.UserOrder, error) {
	var userOrders []UserOrder

	err := o.db.Preload("User").Preload("Region").Preload("OrderDetail").
		Joins("JOIN order_details ON order_details.user_order_id = user_orders.id AND order_details.status = ?", "Menunggu Diterima").
		Where("user_orders.user_id = ?", userIdLogin).
		Find(&userOrders).Error

	if err != nil {
		return nil, err
	}

	var responseOrders []order.UserOrder
	for _, uo := range userOrders {
		if order := uo.ModelToUserOrderWaits(); order != nil {
			responseOrders = append(responseOrders, *order)
		}
	}

	return responseOrders, nil
}

// SelectById implements order.OrderDataInterface.
func (o *orderQuery) SelectById(IdOrder uint) (*order.UserOrder, error) {
	var userOrderData UserOrder
	err := o.db.Preload("User").
		Preload("Region").
		Preload("OrderDetail").
		First(&userOrderData, IdOrder).Error
	if err != nil {
		log.Printf("Error finding order with ID %d: %v", IdOrder, err)
		return nil, err
	}
	result := userOrderData.ModelToUserOrderWait()
	return &result, nil
}

func (o *orderQuery) InsertOrderDetail(adminIdLogin int, userOrderId uint, inputOrder order.OrderDetail) error {
	// Cari data UserOrder berdasarkan userOrderId
	var userOrder UserOrder
	if err := o.db.First(&userOrder, userOrderId).Error; err != nil {
		return err
	}

	// Konversi OrderDetail ke model yang sesuai dengan struktur database
	newOrder := OrderDetailToModel(inputOrder)

	// Set AdminID dari input adminIdLogin
	adminID := uint(adminIdLogin)
	newOrder.AdminID = &adminID

	// Assign UserOrderID dari userOrder yang sudah ditemukan
	newOrder.UserOrderID = userOrder.ID

	// Lakukan operasi Create pada database
	result := o.db.Create(&newOrder)
	if result.Error != nil {
		return result.Error
	}

	// Lakukan operasi Update status
	updateStatus := OrderDetailStatusToModel(inputOrder)

	result = o.db.Model(&newOrder).Updates(&updateStatus)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// SelectUserOrderProcess implements order.OrderDataInterface.
func (o *orderQuery) SelectUserOrderProcess(userIdLogin int) ([]order.UserOrder, error) {
	var userOrders []UserOrder

	err := o.db.Preload("User").Preload("OrderDetail").Preload("Region").
		Joins("JOIN order_details ON order_details.user_order_id = user_orders.id").
		Where("user_orders.user_id = ?", userIdLogin).
		Where("order_details.status <> ?", "Menunggu Diterima").
		Find(&userOrders).Error

	if err != nil {
		return nil, err
	}

	var responseOrders []order.UserOrder
	for _, uo := range userOrders {
		responseOrders = append(responseOrders, uo.ModelToUserOrderWait())
	}

	return responseOrders, nil
}

// SearchUserOrder implements order.OrderDataInterface.
func (o *orderQuery) SearchUserOrder(userIdLogin int, itemName string) ([]order.UserOrder, error) {
	var userOrders []UserOrder

	err := o.db.Preload("User").Preload("OrderDetail").Preload("Region").
		Joins("JOIN order_details ON order_details.user_order_id = user_orders.id").
		Where("user_orders.user_id = ? AND user_orders.item_name LIKE ?", userIdLogin, "%"+itemName+"%").
		Find(&userOrders).Error

	if err != nil {
		return nil, err
	}

	var responseOrders []order.UserOrder
	for _, uo := range userOrders {
		responseOrders = append(responseOrders, uo.ModelToUserOrderWait())
	}

	return responseOrders, nil
}

// SelectAllUserOrderWait implements order.OrderDataInterface.
func (o *orderQuery) SelectAllUserOrderWait() ([]order.UserOrder, error) {
	var userOrders []UserOrder

	err := o.db.Preload("User").Preload("Region").Preload("OrderDetail").
		Joins("JOIN order_details ON order_details.user_order_id = user_orders.id").
		Where("order_details.status = ?", "Menunggu Diterima").
		Find(&userOrders).Error

	if err != nil {
		return nil, err
	}

	var responseOrders []order.UserOrder
	for _, uo := range userOrders {
		if order := uo.ModelToUserOrderWaits(); order != nil {
			responseOrders = append(responseOrders, *order)
		}
	}

	return responseOrders, nil
}

// FetchDeliveryBatchWithRegion implements order.OrderDataInterface.
func (o *orderQuery) FetchDeliveryBatchWithRegion() ([]order.DeliveryBatchWithRegion, error) {
	var result []order.DeliveryBatchWithRegion

	err := o.db.Model(&UserOrder{}).
		Select("order_details.delivery_batch_id, user_orders.region_code_id as region_code, region_codes.region, order_details.user_order_id").
		Joins("JOIN order_details ON user_orders.id = order_details.user_order_id").
		Joins("JOIN region_codes ON user_orders.region_code_id = region_codes.id").
		Scan(&result).Error
	if err != nil {
		return nil, err
	}

	return result, nil
}

// SelectNameByUserOrder implements order.OrderDataInterface.
func (o *orderQuery) SelectNameByUserOrder(code string, batch string) ([]order.UserOrder, error) {
	var userOrders []UserOrder

	err := o.db.Preload("User").
		Preload("Region").
		Preload("OrderDetail").
		Joins("JOIN order_details ON order_details.user_order_id = user_orders.id").
		Where("user_orders.region_code_id = ? AND order_details.delivery_batch_id = ?", code, batch).
		Find(&userOrders).Error

	if err != nil {
		return nil, err
	}

	var responseOrders []order.UserOrder
	for _, uo := range userOrders {
		responseOrders = append(responseOrders, uo.ModelToUserOrderWait())
	}

	return responseOrders, nil
}

// SelectOrderByUserOrderNameUser implements order.OrderDataInterface.
func (o *orderQuery) SelectOrderByUserOrderNameUser(code string, batch string, name string) ([]order.UserOrder, error) {
	var userOrders []UserOrder

	err := o.db.Preload("User").
		Preload("Region").
		Preload("OrderDetail").
		Joins("JOIN order_details ON order_details.user_order_id = user_orders.id").
		Joins("JOIN users ON users.id = user_orders.user_id").
		Where("user_orders.region_code_id = ? AND order_details.delivery_batch_id = ? AND users.name = ?", code, batch, name).
		Where("order_details.status <> ?", "Menunggu Diterima").
		Find(&userOrders).Error

	if err != nil {
		return nil, err
	}

	var responseOrders []order.UserOrder
	for _, uo := range userOrders {
		responseOrders = append(responseOrders, uo.ModelToUserOrderWait())
	}

	return responseOrders, nil
}

// UpdateEstimationForOrders implements order.OrderDataInterface.
func (o *orderQuery) UpdateEstimationForOrders(code, batch string, estimation *time.Time) error {
	subQuery := o.db.Model(&UserOrder{}).
		Select("id").
		Where("region_code_id = ?", code)

	return o.db.Model(&OrderDetail{}).
		Where("user_order_id IN (?)", subQuery).
		Where("delivery_batch_id = ?", batch).
		Update("estimated_delivery_time", estimation).Error
}

// UpdateOrderStatus implements order.OrderDataInterface.
func (o *orderQuery) UpdateOrderStatus(userOrderId uint, status string) error {
	return o.db.Model(&OrderDetail{}).
		Where("user_order_id = ?", userOrderId).
		Update("status", status).Error
}

// UploadFotoPacked implements order.OrderDataInterface.
func (o *orderQuery) UploadFotoPacked(inputOrder order.PhotoOrder, photoPacked *multipart.FileHeader) error {
	imageURL, err := o.cld.UploadImage(photoPacked)
	if err != nil {
		return err
	}

	dataGorm := PhotoOrderToModel(inputOrder)
	dataGorm.PhotoPacked = imageURL

	result := o.db.Create(&dataGorm)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// UploadFotoReceived implements order.OrderDataInterface.
func (o *orderQuery) UploadFotoReceived(idFoto uint, photoReceived *multipart.FileHeader) error {
	imageURL, err := o.cld.UploadImage(photoReceived)
	if err != nil {
		return err
	}

	tx := o.db.Model(&PhotoOrder{}).Where("id = ?", idFoto).Update("photo_received", imageURL)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

// FetchOrdersByBatch implements order.OrderDataInterface.
func (o *orderQuery) FetchOrdersByBatch(batch string) ([]order.UserOrder, error) {
	var userOrders []UserOrder

	err := o.db.Preload("User").
		Preload("Region").
		Preload("OrderDetail").
		Joins("JOIN order_details ON order_details.user_order_id = user_orders.id").
		Where("order_details.delivery_batch_id = ?", batch).
		Find(&userOrders).Error

	if err != nil {
		return nil, err
	}

	var responseOrders []order.UserOrder
	for _, uo := range userOrders {
		responseOrders = append(responseOrders, uo.ModelToUserOrderWait())
	}
	return responseOrders, nil
}

// GenerateCSVByBatch implements order.OrderDataInterface.
func (o *orderQuery) GenerateCSVByBatch(batch string, filePath string) error {
	data, err := o.FetchOrdersByBatch(batch)
	if err != nil {
		return err
	}

	var csvData []csv.UserOrderCSV
	for _, order := range data {
		csvData = append(csvData, csv.UserOrderCSV{
			NamaUser:             order.User.Name,
			NomorTeleponWhatsapp: fmt.Sprintf("%d", order.WhatsAppNumber),
			NomorResiJastip:      order.OrderDetails.TrackingNumberJastip,
			NomorResi:            order.TrackingNumber,
			NomorOrder:           fmt.Sprintf("%d", order.ID),
			KodeWilayah:          fmt.Sprintf("%s - %s", order.Region.ID, order.Region.Region),
			HargaPerKodeWilayah:  fmt.Sprintf("%d", order.Region.Price),
			Berat:                fmt.Sprintf("%2f", float64(order.OrderDetails.WeightItem)),
			NamaBarang:           order.ItemName,
			BatchPengiriman:      batch,
		})
	}

	err = o.csv.GenerateCSV(filePath, csvData)
	if err != nil {
		return err
	}
	return nil
}

// GetFoto implements order.OrderDataInterface.
func (o *orderQuery) GetFoto(batch string, code string, userId int) (*order.PhotoOrder, error) {
	var photoOrders PhotoOrder

	err := o.db.Preload("User").
		Preload("Region").
		Preload("DeliveryBatch").
		Where("delivery_batch_id = ? AND region_code_id = ? AND user_id = ?", batch, code, userId).
		First(&photoOrders).Error

	if err != nil {
		return nil, err
	}

	result := photoOrders.ModelToPhotoOrder()

	return &result, nil
}

// SearchOrders implements order.OrderDataInterface.
func (o *orderQuery) SearchOrders(searchQuery string) ([]order.UserOrder, error) {
	var userOrders []UserOrder

	searchPattern := "%" + searchQuery + "%"

	err := o.db.Preload("User").
		Preload("Region").
		Preload("OrderDetail").
		Joins("JOIN order_details ON order_details.user_order_id = user_orders.id").
		Joins("JOIN users ON users.id = user_orders.user_id").
		Joins("JOIN region_codes ON region_codes.id = user_orders.region_code_id").
		Where("user_orders.item_name LIKE ? OR "+
			"user_orders.tracking_number LIKE ? OR "+
			"user_orders.online_store LIKE ? OR "+
			"user_orders.region_code_id LIKE ? OR "+
			"user_orders.whatsapp_number LIKE ? OR "+
			"region_codes.region LIKE ? OR "+
			"order_details.status LIKE ? OR "+
			"CAST(order_details.weight_item AS CHAR) LIKE ? OR "+
			"order_details.tracking_number_jastip LIKE ?",
			searchPattern, searchPattern, searchPattern, searchPattern, searchPattern, searchPattern, searchPattern, searchPattern, searchPattern).
		Find(&userOrders).Error

	if err != nil {
		return nil, err
	}

	var responseOrders []order.UserOrder
	for _, uo := range userOrders {
		responseOrders = append(responseOrders, uo.ModelToUserOrderWait())
	}

	return responseOrders, nil
}

// UpdateOrderByID updates an order's details based on the given order ID.
func (o *orderQuery) UpdateOrderByID(orderID uint, inputOrder order.UpdateOrderByID) error {
	userOrder, orderDetail := UserOrderUpdateToModel(inputOrder)

	// Update UserOrder
	result := o.db.Model(&UserOrder{}).Where("id = ?", orderID).Updates(userOrder)
	if result.Error != nil {
		return result.Error
	}

	// Update OrderDetail
	result = o.db.Model(&OrderDetail{}).Where("user_order_id = ?", orderID).Updates(orderDetail)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// FetchRegionStatsByBatch implements order.OrderDataInterface.
func (o *orderQuery) FetchRegionStatsByBatch(batch string) ([]order.RegionBatchStats, error) {
	var results []order.RegionBatchStats

	err := o.db.Table("user_orders").
		Select(`
			user_orders.region_code_id AS region_code,
			COALESCE(region_codes.price, 0) AS price_per_code,
			SUM(order_details.weight_item) AS total_weight,
			COUNT(DISTINCT user_orders.user_id) AS total_users,
			COUNT(user_orders.id) AS total_orders,
			(COALESCE(region_codes.price, 0) * COUNT(user_orders.id)) AS total_price`).
		Joins(`JOIN order_details ON user_orders.id = order_details.user_order_id`).
		Joins(`JOIN region_codes ON user_orders.region_code_id = region_codes.id`).
		Where(`order_details.delivery_batch_id = ?`, batch).
		Group(`user_orders.region_code_id, region_codes.price`).
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	fmt.Println("Response:")
	for _, result := range results {
		fmt.Printf("Kode Wilayah: %s\n", result.RegionCode)
		fmt.Printf("Total per kode wilayah: %d\n", result.PricePerCode)
		fmt.Printf("Total Berat: %.2f\n", result.TotalWeight)
		fmt.Printf("Total Users: %d\n", result.TotalUsers)
		fmt.Printf("Total Orders: %d\n", result.TotalOrders)
		fmt.Printf("Total Harga: %d\n", result.TotalPrice)
		fmt.Println("--------------------")
	}

	return results, nil
}
