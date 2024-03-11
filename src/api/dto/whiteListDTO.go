package dto

import "github.com/google/uuid"

type WhiteListAddDTO struct {
	UserId       uuid.UUID `json:"user_id"`
	UserDeviceID uuid.UUID `json:"user_device_id"`
	UserIp       string    `json:"user_ip"`
}
