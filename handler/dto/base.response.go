package dto

const (
	X_USER_ID = "x-user-id"
)

type BaseResponse struct {
	Msg   string      `json:"message"`
	Data  interface{} `json:"data,omitempty"`
	Error interface{} `json:"error,omitempty"`
}

