package common

import (
	"unicode"

	"github.com/Arshia-Izadyar/Fast-Gopher/src/pkg/service_errors"
)

func ValidatePassword(password string) error {
	var (
		hasMinLength = false
		hasUpper     = false
		hasLower     = false
		hasNumber    = false
		hasSpecial   = false
	)

	if len(password) >= 8 {
		hasMinLength = true
	}

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if !hasMinLength {
		return &service_errors.ServiceError{EndUserMessage: "password must be at least 8 characters long"}
	}
	if !hasUpper {
		return &service_errors.ServiceError{EndUserMessage: "password must include at least one uppercase letter"}
	}
	if !hasLower {
		return &service_errors.ServiceError{EndUserMessage: "password must include at least one lowercase letter"}
	}
	if !hasNumber {
		return &service_errors.ServiceError{EndUserMessage: "password must include at least one digit"}
	}
	if !hasSpecial {
		return &service_errors.ServiceError{EndUserMessage: "password must include at least one special character"}
	}

	return nil
}
