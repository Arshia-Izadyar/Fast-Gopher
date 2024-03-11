package dto

import "github.com/google/uuid"

type UserCreateDTO struct {
	Email               string `json:"email" validate:"required,email"`
	UserPassword        string `json:"password" validate:"required"`
	UserPasswordConfirm string `json:"password_confirm" validate:"required"`
}

type UserDTO struct {
	Email        string `json:"email" validate:"required,email"`
	UserPassword string `json:"password" validate:"required"`
}

type UserTokenDTO struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type UserLogout struct {
	UserId       uuid.UUID `json:"user_id"`
	UserDeviceID uuid.UUID `json:"user_device_id"`
	UserIp       string    `json:"user_ip"`
	UserToken    string    `json:"user_token"`
}

type RefreshTokenDTO struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}
