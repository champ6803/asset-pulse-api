package transformer

import (
	"asset-pulse-api/handler/dto"
)

func SuccessResponse(status int, data interface{}) dto.BaseResponse {
	return dto.BaseResponse{
		Msg:  "success",
		Data: data,
	}
}

func ExceptionResponse(status int, err error) dto.BaseResponse {
	errResponse := map[string]interface{}{
		"code":    "INTERNAL_ERROR",
		"message": err.Error(),
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
