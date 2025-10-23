package models

type AuthenticateUserInp struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthenticateUserResp struct {
	User  UserResponse `json:"user"`
	Token string       `json:"token"`
}

type LoginResponse struct {
	User  UserResponse `json:"user"`
	Token string       `json:"token"`
}
