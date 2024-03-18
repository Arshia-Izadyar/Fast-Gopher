package dto

type TokenDTO struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshTokenDTO struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type KeyDTO struct {
	Key       string `json:"key" validate:"required"`
	SessionId string `json:"session_id" validate:"required"`
	Premium   bool   `json:"premium"`
}

type GenerateKeyDTO struct {
	SessionId string `json:"session_id" validate:"required"`
}

type KeyAcDTO struct {
	Key          string `json:"key"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
