package dto

type WhiteListAddDTO struct {
	Key       string `json:"key"`
	SessionId string `json:"session_id"`
	UserIp    string `json:"user_ip"`
}
