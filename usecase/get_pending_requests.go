package usecase

import (
	"asset-pulse-api/usecase/models"
	"context"
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
		request := models.PendingRequestItem{
			ID:                      item.ID,
			TicketNo:                item.TicketNo,
			CompanyCode:             item.CompanyCode,
			Type:                    item.Type,
			RequesterUserID:         item.RequesterUserID,
			ScopeLevel:              item.ScopeLevel,
			ScopeRefID:              item.ScopeRefID,
			PayloadJSON:             item.PayloadJSON,
			Status:                  item.Status,
			CreatedAt:               item.CreatedAt,
			UpdatedAt:               item.UpdatedAt,
			RequesterUsername:       item.RequesterUsername,
			RequesterEmail:          item.RequesterEmail,
			RequesterFullName:       item.RequesterFullName,
			RequesterDepartmentCode: item.RequesterDepartmentCode,
			AppID:                   item.AppID,
			AppName:                 item.AppName,
		}

		// Convert current_step_no from *int64 to *int
		if item.CurrentStepNo != nil {
			stepInt := int(*item.CurrentStepNo)
			request.ApprovalsCurrStep = &stepInt
		}

		// Convert total_steps from *int64 to *int
		if item.TotalSteps != nil {
			totalInt := int(*item.TotalSteps)
			request.ApprovalsStep = &totalInt
		}

		requests = append(requests, request)
	}

	return &models.GetPendingRequestsResp{
		Requests: requests,
		Total:    total,
	}, nil
}
