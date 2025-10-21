package models

import "github.com/google/uuid"

type GetUsersInp struct {
	CompanyCode *string
	Status      *string
	Page        int
	PageSize    int
}

type GetUsersResp struct {
	Users []UserResponse `json:"users"`
	Total int            `json:"total"`
}

type UserResponse struct {
	ID             int64      `json:"id"`
	CompanyCode    *string    `json:"company_code"`
	DepartmentCode *string    `json:"department_code"`
	EntraID        *uuid.UUID `json:"entra_id"`
	Email          *string    `json:"email"`
	DisplayName    *string    `json:"display_name"`
	Title          *string    `json:"title"`
	EmployeeID     *string    `json:"employee_id"`
	Status         *string    `json:"status"`
}

