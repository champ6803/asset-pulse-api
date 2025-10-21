package transformer

import (
	"asset-pulse-api/handler/dto"
	"asset-pulse-api/utils/apperrs"
)

func SuccessResponse(status int, data interface{}) dto.BaseResponse {
	return dto.BaseResponse{
		Msg:  "success",
		Data: data,
	}
}

func ExceptionResponse(status int, err *apperrs.AppError) dto.BaseResponse {
	errResponse := map[string]interface{}{
		"code":    err.Code,
		"message": err.Message,
	}

	if err.Cause != nil {
		errResponse["cause"] = err.Cause.Error()
	}

	return dto.BaseResponse{
		Msg:   "error",
		Error: errResponse,
	}
}

type ExceptionResponseDetail struct {
	Code  string `json:"code"`
	Msg   string `json:"message"`
	Cause string `json:"cause,omitempty"`
}

