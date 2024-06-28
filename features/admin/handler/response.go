package handler

import (
	"jastip-jakarta/features/admin"
	"jastip-jakarta/utils/time"
)

type AdminResponse struct {
	ID           uint   `json:"admin_id" form:"admin_id"`
	Name         string `json:"name" form:"name"`
	Role         string `json:"role" form:"role"`
	Email        string `json:"email" form:"email"`
	PhoneNumber  int    `json:"phone_number" form:"phone_number"`
	PhotoProfile string `json:"photo_profile" form:"photo_profile"`
	CreatedAt    string `json:"create_account"`
	UpdatedAt    string `json:"last_update"`
}

type RegionCodeResponse struct {
	Code        string `json:"code"`
	Region      string `json:"region"`
	FullAddress string `json:"full_address"`
	PhoneNumber int    `json:"phone_number"`
	AdminID     uint   `json:"admin_id"`
}

type AdminResponseOrder struct {
	Name string `json:"name"`
}

func AdminToResponse(data *admin.Admin) AdminResponse {
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
		FullAddress: data.FullAddress,
		PhoneNumber: data.PhoneNumber,
		AdminID:     data.AdminID,
	}
}
