package dto

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
