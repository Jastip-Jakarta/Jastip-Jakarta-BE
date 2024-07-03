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
	Status               string  `json:"status"`
	WeightItem           float64 `json:"weight_item"`
	DeliveryBatch        string  `json:"delivery_batch"`
	TrackingNumberjastip string  `json:"tracking_number_jastip"`
}

type UpdateStatusRequest struct {
	Status string `json:"status"`
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
