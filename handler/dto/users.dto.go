package dto

import (
	"asset-pulse-api/usecase/models"

	"github.com/google/uuid"
)

type GetUsersResponse struct {
	Users []UserItem `json:"users"`
	Total int        `json:"total"`
}

type UserItem struct {
	ID             int64      `json:"id"`
	Username       *string    `json:"username"`
	CompanyCode    *string    `json:"company_code"`
	DepartmentCode *string    `json:"department_code"`
	EntraID        *uuid.UUID `json:"entra_id"`
	Email          *string    `json:"email"`
	DisplayName    *string    `json:"display_name"`
	Title          *string    `json:"title"`
	EmployeeID     *string    `json:"employee_id"`
	Status         *string    `json:"status"`
}

func NewGetUsersDTO(resp *models.GetUsersResp) *GetUsersResponse {
	users := make([]UserItem, len(resp.Users))
	for i, user := range resp.Users {
		users[i] = UserItem{
			ID:             user.ID,
			Username:       user.Username,
			CompanyCode:    user.CompanyCode,
			DepartmentCode: user.DepartmentCode,
			EntraID:        user.EntraID,
			Email:          user.Email,
			DisplayName:    user.DisplayName,
			Title:          user.Title,
			EmployeeID:     user.EmployeeID,
			Status:         user.Status,
		}
	}

	return &GetUsersResponse{
		Users: users,
		Total: resp.Total,
	}
}
