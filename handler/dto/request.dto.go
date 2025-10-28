package dto

import "time"

// GetPendingRequestsQueryParams defines query parameters for fetching pending requests
type GetPendingRequestsQueryParams struct {
	Limit int `form:"limit"`
}

// PendingRequestDTO represents a single pending request
type PendingRequestDTO struct {
	ID                int64     `json:"id"`
	TicketNo          string    `json:"ticket_no"`
	Type              string    `json:"type"`
	Status            *string   `json:"status"`
	ApprovalsCurrStep *int      `json:"approvals_curr_step"`
	ApprovalsStep     *int      `json:"approvals_step"`
	AppID             *int64    `json:"app_id"`
	AppName           *string   `json:"app_name"`
	CreatedAt         time.Time `json:"created_at"`
}

// GetPendingRequestsResponse defines the response structure
type GetPendingRequestsResponse struct {
	Requests []PendingRequestDTO `json:"requests"`
	Total    int64               `json:"total"`
}
