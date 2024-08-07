package handler

import (
	"jastip-jakarta/features/order"
	"math/rand"
	"time"
)

type UserOrderRequest struct {
	ID             uint
	ItemName       string `json:"item_name"`
	TrackingNumber string `json:"tracking_number"`
	OnlineStore    string `json:"online_store"`
	WhatsAppNumber int    `json:"whatsapp_number"`
	Code           string `json:"code"`
}

type OrderDetailRequest struct {
	Status               string `json:"status"`
	WeightItem           int    `json:"weight_item"`
	DeliveryBatch        string `json:"delivery_batch"`
	TrackingNumberjastip string `json:"tracking_number_jastip"`
}

type UploadFotoRequest struct {
	Batch  string `form:"batch"`
	Code   string `form:"code"`
	UserID uint   `form:"user_id"`
}

type UpdateStatusRequest struct {
	Status string `json:"status"`
}

type UpdateEstimationRequest struct {
	Estimation string `json:"estimation"`
}

func RequestToUserOrder(input UserOrderRequest) order.UserOrder {
	return order.UserOrder{
		ID:             generateID(),
		ItemName:       input.ItemName,
		TrackingNumber: input.TrackingNumber,
		OnlineStore:    input.OnlineStore,
		WhatsAppNumber: input.WhatsAppNumber,
		RegionCode:     input.Code,
	}
}

func RequestToPhotoOrder(input UploadFotoRequest) order.PhotoOrder {
	return order.PhotoOrder{
		UserID:          input.UserID,
		DeliveryBatchID: input.Batch,
		RegionCodeID:    input.Code,
	}
}

func RequestUpdateToUserOrder(input UserOrderRequest) order.UserOrder {
	return order.UserOrder{
		ItemName:       input.ItemName,
		TrackingNumber: input.TrackingNumber,
		OnlineStore:    input.OnlineStore,
		WhatsAppNumber: input.WhatsAppNumber,
		RegionCode:     input.Code,
	}
}

func RequestToOrderDetail(input OrderDetailRequest) order.OrderDetail {
	deliveryBatch := input.DeliveryBatch
	return order.OrderDetail{
		Status:               input.Status,
		WeightItem:           input.WeightItem,
		TrackingNumberJastip: input.TrackingNumberjastip,
		DeliveryBatchID:      &deliveryBatch,
	}
}

func generateID() uint {
	rand.Seed(time.Now().UnixNano())
	randomNumber := rand.Int63n(99999-10000) + 10000
	return uint(randomNumber)
}

func ParseEstimationDate(estimation string) (*time.Time, error) {
	// Format tanggal dd/mm/yy
	layout := "02/01/06"
	t, err := time.Parse(layout, estimation)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// func RequestUpdateEstimasi(input UpdateEstimationRequest) (*time.Time, error) {
//     return ParseEstimationDate(input.Estimation)
// }
