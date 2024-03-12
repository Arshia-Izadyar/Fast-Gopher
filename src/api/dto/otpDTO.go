package dto

type SendOtpDTO struct {
	Email string `json:"email"`
	Otp   string `json:"otp"`
}

type OtpDTO struct {
	Used bool   `json:"used"`
	Otp  string `json:"otp"`
}
