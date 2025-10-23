package dto

import (
	"asset-pulse-api/usecase/models"
)

type UserDTO struct {
	ID             int64  `json:"id"`
	Username       string `json:"username"`
	Email          string `json:"email"`
	DisplayName    string `json:"display_name"`
	CompanyCode    string `json:"company_code"`
	DepartmentCode string `json:"department_code"`
	Role           string `json:"role"`
	Status         string `json:"status"`
}

func NewUserDTO(user *models.UserResponse) *UserDTO {
	dto := &UserDTO{
		ID: user.ID,
	}

	if user.Username != nil {
		dto.Username = *user.Username
	}
	if user.Email != nil {
		dto.Email = *user.Email
	}
	if user.DisplayName != nil {
		dto.DisplayName = *user.DisplayName
	}
	if user.CompanyCode != nil {
		dto.CompanyCode = *user.CompanyCode
	}
	if user.DepartmentCode != nil {
		dto.DepartmentCode = *user.DepartmentCode
	}
	if user.Status != nil {
		dto.Status = *user.Status
	}

	// Role field doesn't exist in UserResponse, set default
	dto.Role = "employee"

	return dto
}
