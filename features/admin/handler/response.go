package handler

import (
	"jastip-jakarta/features/admin"
	"jastip-jakarta/utils/time"
	uh "jastip-jakarta/features/user"
)

type AdminResponse struct {
	ID           uint   `json:"admin_id"`
	Name         string `json:"name"`
	Role         string `json:"role"`
	Email        string `json:"email"`
	PhoneNumber  int    `json:"phone_number"`
	PhotoProfile string `json:"photo_profile"`
	CreatedAt    string `json:"create_account"`
	UpdatedAt    string `json:"last_update"`
}

type RegionCodeResponse struct {
	Code        string `json:"code"`
	Region      string `json:"region"`
	Price       int    `json:"price"`
	FullAddress string `json:"full_address"`
	PhoneNumber int    `json:"phone_number"`
	AdminID     uint   `json:"admin_id"`
}

type AdminResponseOrder struct {
	Name string `json:"name"`
}

type DeliveryBatchResponse struct {
	DeliveryBatch string `json:"delivery_batch"`
	Batch         int    `json:"batch"`
	Year          int    `json:"year"`
	Month         int    `json:"month"`
}

type UserResponse struct {
    ID           uint   `json:"id"`
    Name         string `json:"name"`
    Email        string `json:"email"`
    PhoneNumber  int    `json:"phone_number"`
    PhotoProfile string `json:"photo_profile"`
}

func UserToResponse(user uh.User) UserResponse {
    return UserResponse{
        ID:           user.ID,
        Name:         user.Name,
        Email:        user.Email,
        PhoneNumber:  user.PhoneNumber,
        PhotoProfile: user.PhotoProfile,
    }
}

func AdminToResponse(data admin.Admin) AdminResponse {
	return AdminResponse{
		ID:           data.ID,
		Name:         data.Name,
		Email:        data.Email,
		PhoneNumber:  data.PhoneNumber,
		PhotoProfile: data.PhotoProfile,
		Role:         data.Role,
		CreatedAt:    time.FormatDateToIndonesian(data.CreatedAt),
		UpdatedAt:    time.FormatDateToIndonesian(data.UpdatedAt),
	}
}

func CoreToResponseRegionCode(data admin.RegionCode) RegionCodeResponse {
	return RegionCodeResponse{
		Code:        data.ID,
		Region:      data.Region,
		Price:       data.Price,
		FullAddress: data.FullAddress,
		PhoneNumber: data.PhoneNumber,
		AdminID:     data.AdminID,
	}
}

func CoreToResponseDeliveryBatch(data admin.DeliveryBatch) DeliveryBatchResponse {
	return DeliveryBatchResponse{
		DeliveryBatch: data.ID,
		Batch:         data.Batch,
		Year:          data.Year,
		Month:         data.Month,
	}
}
