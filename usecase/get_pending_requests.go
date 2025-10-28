package usecase

import (
	"asset-pulse-api/usecase/models"
	"context"
	"time"
)

func (u *useCase) GetPendingRequests(ctx context.Context, in *models.GetPendingRequestsInp) (*models.GetPendingRequestsResp, error) {
	// Fetch pending requests from repository
	data, total, err := u.dbRepo.GetPendingRequests(ctx, in.UserID, in.CompanyCode, in.Limit)
	if err != nil {
		return nil, err
	}

	// Transform repository data to usecase model
	requests := make([]models.PendingRequestItem, 0, len(data))
	for _, item := range data {
		request := models.PendingRequestItem{}

		// ID
		if id, ok := item["id"].(int64); ok {
			request.ID = id
		}

		// Ticket number (generated in SQL)
		if ticketNo, ok := item["ticket_no"].(string); ok {
			request.TicketNo = ticketNo
		}

		// Company code
		if companyCode, ok := item["company_code"].(string); ok {
			request.CompanyCode = &companyCode
		}

		// Type
		if reqType, ok := item["type"].(string); ok {
			request.Type = reqType
		}

		// Requester user ID
		if requesterUserID, ok := item["requester_user_id"].(int64); ok {
			request.RequesterUserID = &requesterUserID
		}

		// Scope level
		if scopeLevel, ok := item["scope_level"].(string); ok {
			request.ScopeLevel = &scopeLevel
		}

		// Scope ref ID
		if scopeRefID, ok := item["scope_ref_id"].(int64); ok {
			request.ScopeRefID = &scopeRefID
		}

		// Payload JSON
		if payloadJSON, ok := item["payload_json"].(string); ok {
			request.PayloadJSON = &payloadJSON
		}

		// Status
		if status, ok := item["status"].(string); ok {
			request.Status = &status
		}

		// Created at
		if createdAt, ok := item["created_at"].(time.Time); ok {
			request.CreatedAt = createdAt
		}

		// Updated at
		if updatedAt, ok := item["updated_at"].(time.Time); ok {
			request.UpdatedAt = &updatedAt
		}

		// Requester username
		if requesterUsername, ok := item["requester_username"].(string); ok {
			request.RequesterUsername = &requesterUsername
		}

		// Requester email
		if requesterEmail, ok := item["requester_email"].(string); ok {
			request.RequesterEmail = &requesterEmail
		}

		// Requester full name
		if requesterFullName, ok := item["requester_full_name"].(string); ok {
			request.RequesterFullName = &requesterFullName
		}

		// Requester department code
		if requesterDepartmentCode, ok := item["requester_department_code"].(string); ok {
			request.RequesterDepartmentCode = &requesterDepartmentCode
		}

		requests = append(requests, request)
	}

	return &models.GetPendingRequestsResp{
		Requests: requests,
		Total:    total,
	}, nil
}
