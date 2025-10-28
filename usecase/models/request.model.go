package models

import "time"

type GetPendingRequestsInp struct {
	UserID      *int64
	CompanyCode *string
	Limit       int
}

type PendingRequestItem struct {
	ID                      int64      `json:"id"`
	TicketNo                string     `json:"ticket_no"`
	CompanyCode             *string    `json:"company_code"`
	Type                    string     `json:"type"`
	RequesterUserID         *int64     `json:"requester_user_id"`
	ScopeLevel              *string    `json:"scope_level"`
	ScopeRefID              *int64     `json:"scope_ref_id"`
	PayloadJSON             *string    `json:"payload_json"`
	Status                  *string    `json:"status"`
	CreatedAt               time.Time  `json:"created_at"`
	UpdatedAt               *time.Time `json:"updated_at"`
	RequesterUsername       *string    `json:"requester_username"`
	RequesterEmail          *string    `json:"requester_email"`
	RequesterFullName       *string    `json:"requester_full_name"`
	RequesterDepartmentCode *string    `json:"requester_department_code"`
}

type GetPendingRequestsResp struct {
	Requests []PendingRequestItem `json:"requests"`
	Total    int64                `json:"total"`
}
