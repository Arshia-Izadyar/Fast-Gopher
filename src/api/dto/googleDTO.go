package dto

type GoogleUserInfoDTO struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
}

type GoogleCodeLoginDTO struct {
	Code string `json:"code"`
}
