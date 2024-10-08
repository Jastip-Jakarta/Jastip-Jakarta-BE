package data

import (
	ad "jastip-jakarta/features/admin/data"
	"jastip-jakarta/features/order"
	ud "jastip-jakarta/features/user/data"

	"time"

	"gorm.io/gorm"
)

type UserOrder struct {
	ID uint `gorm:"primaryKey" json:"id"`
	gorm.Model
	UserID         uint
	ItemName       string
	TrackingNumber string
	OnlineStore    string
	WhatsappNumber int
	RegionCodeID   string
	User           ud.User       `gorm:"foreignKey:UserID"`
	Region         ad.RegionCode `gorm:"foreignKey:RegionCodeID"`
	OrderDetail    OrderDetail
}

type OrderDetail struct {
	gorm.Model
	UserOrderID           uint
	AdminID               *uint `gorm:"default:null"`
	Status                string
	WeightItem            float64
	TrackingNumberJastip  string
	DeliveryBatchID       *string `gorm:"default:null"`
	EstimatedDeliveryTime *time.Time
	Admin                 ad.Admin         `gorm:"foreignKey:AdminID"`
	DeliveryBatch         ad.DeliveryBatch `gorm:"foreignKey:DeliveryBatchID"`
}

type PhotoOrder struct {
	gorm.Model
	DeliveryBatchID string
	PhotoPacked     string
	PhotoReceived   string
	RegionCodeID    string
	UserID          uint
	User            ud.User          `gorm:"foreignKey:UserID"`
	Region          ad.RegionCode    `gorm:"foreignKey:RegionCodeID"`
	DeliveryBatch   ad.DeliveryBatch `gorm:"foreignKey:DeliveryBatchID"`
}

func OrderDetailStatusToModel(updateStatus order.OrderDetail) OrderDetail {
	return OrderDetail{
		Status: updateStatus.Status,
	}
}

func PhotoOrderToModel(input order.PhotoOrder) PhotoOrder {
	return PhotoOrder{
		DeliveryBatchID: input.DeliveryBatchID,
		PhotoPacked:     input.PhotoPacked,
		PhotoReceived:   input.PhotoReceived,
		RegionCodeID:    input.RegionCodeID,
		UserID:          input.UserID,
	}
}

func (o PhotoOrder) ModelToPhotoOrder() order.PhotoOrder {
	return order.PhotoOrder{
		ID:              o.ID,
		UserID:          o.UserID,
		DeliveryBatchID: o.DeliveryBatchID,
		PhotoPacked:     o.PhotoPacked,
		PhotoReceived:   o.PhotoReceived,
		RegionCodeID:    o.RegionCodeID,
	}
}

func UserOrderToModel(input order.UserOrder) UserOrder {
	return UserOrder{
		ID:             input.ID,
		UserID:         input.UserID,
		ItemName:       input.ItemName,
		TrackingNumber: input.TrackingNumber,
		OnlineStore:    input.OnlineStore,
		WhatsappNumber: input.WhatsAppNumber,
		RegionCodeID:   input.RegionCode,
	}
}

func UserOrderUpdateToModel(input order.UpdateOrderByID) (UserOrder, OrderDetail) {
	userOrder := UserOrder{
		ItemName:       input.ItemName,
		TrackingNumber: input.TrackingNumber,
		OnlineStore:    input.OnlineStore,
		WhatsappNumber: input.WhatsAppNumber,
		RegionCodeID:     input.RegionCode,
	}

	orderDetail := OrderDetail{
		WeightItem:           input.WeightItem,
		TrackingNumberJastip: input.TrackingNumberJastip,
		DeliveryBatchID:      &input.DeliveryBatch,
	}

	return userOrder, orderDetail
}

func (uo UserOrder) ModelToUserOrderWait() order.UserOrder {
	return order.UserOrder{
		ID:             uo.ID,
		UserID:         uo.UserID,
		ItemName:       uo.ItemName,
		TrackingNumber: uo.TrackingNumber,
		OnlineStore:    uo.OnlineStore,
		WhatsAppNumber: uo.WhatsappNumber,
		Region:         uo.Region.ModelToRegionCode(),
		User:           uo.User.ModelToUser(),
		OrderDetails:   uo.OrderDetail.ModelToOrderDetail(),
	}
}

func (o UserOrder) ModelToUserOrderWaits() *order.UserOrder {
	if o.OrderDetail.Status != "Menunggu Diterima" {
		return nil
	}
	return &order.UserOrder{
		ID:             o.ID,
		UserID:         o.UserID,
		ItemName:       o.ItemName,
		TrackingNumber: o.TrackingNumber,
		OnlineStore:    o.OnlineStore,
		WhatsAppNumber: o.WhatsappNumber,
		Region:         o.Region.ModelToRegionCode(),
		User:           o.User.ModelToUser(),
		OrderDetails:   o.OrderDetail.ModelToOrderDetail(),
	}
}

func OrderDetailToModel(input order.OrderDetail) OrderDetail {
	return OrderDetail{
		AdminID:               input.AdminID,
		Status:                input.Status,
		WeightItem:            input.WeightItem,
		DeliveryBatchID:       input.DeliveryBatchID,
		TrackingNumberJastip:  input.TrackingNumberJastip,
		EstimatedDeliveryTime: input.EstimatedDeliveryTime,
	}
}

func (o OrderDetail) ModelToOrderDetail() order.OrderDetail {
	return order.OrderDetail{
		ID:                    o.ID,
		UserOrderID:           o.UserOrderID,
		Status:                o.Status,
		WeightItem:            o.WeightItem,
		DeliveryBatchID:       o.DeliveryBatchID,
		EstimatedDeliveryTime: o.EstimatedDeliveryTime,
		TrackingNumberJastip:  o.TrackingNumberJastip,
	}
}
