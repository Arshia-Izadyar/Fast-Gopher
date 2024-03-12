package common

import (
	"fmt"
	"time"

	"github.com/Arshia-Izadyar/Fast-Gopher/src/api/dto"
	"github.com/Arshia-Izadyar/Fast-Gopher/src/config"
	"github.com/Arshia-Izadyar/Fast-Gopher/src/constants"
	"github.com/Arshia-Izadyar/Fast-Gopher/src/pkg/service_errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	signingMethod = jwt.SigningMethodHS256
)

func GenerateJwt(userId uuid.UUID, cfg *config.Config) (*dto.UserTokenDTO, error) {
	res := &dto.UserTokenDTO{}

	// Use a single time.Now()
	now := time.Now()

	expirationTimeAccessToken := now.Add(cfg.JWT.AccessTokenExpireDuration * time.Minute).Unix()
	expirationTimeRefreshToken := now.Add(cfg.JWT.RefreshTokenExpireDuration * time.Minute).Unix()

	atClaims := jwt.MapClaims{
		constants.ExpKey:     expirationTimeAccessToken,
		constants.UserIdKey:  userId,
		constants.AccessType: true,
	}

	// Using pre compiled signing method
	tk := jwt.NewWithClaims(signingMethod, atClaims)
	var err error
	res.AccessToken, err = tk.SignedString([]byte(cfg.JWT.Secret))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	rfClaims := jwt.MapClaims{
		constants.ExpKey:     expirationTimeRefreshToken,
		constants.UserIdKey:  userId,
		constants.AccessType: false,
	}

	rt := jwt.NewWithClaims(signingMethod, rfClaims)
	res.RefreshToken, err = rt.SignedString([]byte(cfg.JWT.Secret))
	if err != nil {
		return nil, err
	}

	return res, nil
}

func ValidateToken(token string, cfg *config.Config) (jwt.MapClaims, error) {

	tk, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, &service_errors.ServiceError{
				EndUserMessage: service_errors.TokenInvalid,
			}
		}
		return []byte(cfg.JWT.Secret), nil
	})
	if err != nil {
		return nil, &service_errors.ServiceError{EndUserMessage: service_errors.TokenInvalid, Err: err}
	}
	if claims, ok := tk.Claims.(jwt.MapClaims); ok {
		exp := time.Unix(int64((claims[constants.ExpKey]).(float64)), 0)
		now := time.Now()
		if now.After(exp) {
			return nil, &service_errors.ServiceError{EndUserMessage: service_errors.TokenExpired}
		}
		return claims, nil
	}
	return nil, nil
}
