package dto

type TokenDTO struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshTokenDTO struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type KeyDTO struct {
	Key string `json:"key" validate:"required"`
	// SessionId  string `json:"session_id" validate:"required"`
	DeviceName string `json:"device_name" validate:"required"`
}

type IKeyDTO struct {
	Key string `json:"key" validate:"required"`
}

type GenerateKeyDTO struct {
	// SessionId  string `json:"session_id" validate:"required"`
	DeviceName string `json:"device_name"`
}

type KeyAcDTO struct {
	Key          string `json:"key"`
	SessionId    string `json:"session_id"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type DeviceDTO struct {
	DeviceName string `json:"device_name"`
	SessionId  string `json:"session_id"`
	Ip         string `json:"ip"`
}

type RemoveDeviceDTO struct {
	Key       string `json:"-"`
	SessionId string `json:"session_id" validate:"required"`
}

type SessionKeyDTO struct {
	Key       string `json:"key" validate:"required"`
	SessionId string `json:"session_id" validate:"required"`
}

type RenameDeviceDTO struct {
	NewDeviceName string `json:"new_device_name" validate:"required"`
	SessionId     string `json:"session_id" validate:"required"`
	Key           string `json:"-"`
}
