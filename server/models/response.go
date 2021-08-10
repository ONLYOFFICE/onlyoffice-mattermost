package models

type UserInfoResponse struct {
	Id          string      `json:"id"`
	Username    string      `json:"username"`
	Permissions Permissions `json:"permissions"`
}
