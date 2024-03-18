package common

import (
	"fmt"
	"time"

	"github.com/Arshia-Izadyar/Fast-Gopher/src/api/dto"
	"github.com/Arshia-Izadyar/Fast-Gopher/src/config"
	"github.com/Arshia-Izadyar/Fast-Gopher/src/constants"
	"github.com/Arshia-Izadyar/Fast-Gopher/src/pkg/service_errors"
	"github.com/golang-jwt/jwt/v5"
)

// pre compiled
var (
	signingMethod = jwt.SigningMethodHS256
)

func GenerateJwt(key *dto.KeyDTO, cfg *config.Config) (*dto.KeyAcDTO, error) {
	var err error
	res := &dto.KeyAcDTO{}

	// Use a single time.Now()
	now := time.Now()

	expirationTimeAccessToken := now.Add(cfg.JWT.AccessTokenExpireDuration * time.Minute).Unix()
	expirationTimeRefreshToken := now.Add(cfg.JWT.RefreshTokenExpireDuration * time.Minute).Unix()

	atClaims := jwt.MapClaims{
		constants.ExpKey:       expirationTimeAccessToken,
		constants.Key:          key.Key,
		constants.Premium:      key.Premium,
		constants.SessionIdKey: key.SessionId,
		constants.AccessType:   true,
	}

	// Using pre compiled signing method
	tk := jwt.NewWithClaims(signingMethod, atClaims)
	res.AccessToken, err = tk.SignedString([]byte(cfg.JWT.Secret))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	rfClaims := jwt.MapClaims{
		constants.ExpKey:       expirationTimeRefreshToken,
		constants.Key:          key.Key,
		constants.SessionIdKey: key.SessionId,
		constants.AccessType:   false,
	}

	rt := jwt.NewWithClaims(signingMethod, rfClaims)
	res.RefreshToken, err = rt.SignedString([]byte(cfg.JWT.Secret))
	if err != nil {
		return nil, err
	}
	res.Key = key.Key
	return res, nil
}

func ValidateToken(token string, cfg *config.Config) (jwt.MapClaims, *service_errors.ServiceErrors) {

	tk, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, &service_errors.ServiceErrors{
				EndUserMessage: service_errors.TokenInvalid,
			}
		}
		return []byte(cfg.JWT.Secret), nil
	})
	if err != nil {
		return nil, &service_errors.ServiceErrors{EndUserMessage: err.Error(), Status: 401}
	}
	if claims, ok := tk.Claims.(jwt.MapClaims); ok {
		exp := time.Unix(int64((claims[constants.ExpKey]).(float64)), 0)
		now := time.Now()
		if now.After(exp) {
			return nil, &service_errors.ServiceErrors{EndUserMessage: service_errors.TokenExpired, Status: 401}
		}
		return claims, nil
	}
	return nil, &service_errors.ServiceErrors{EndUserMessage: service_errors.TokenInvalid, Status: 401}
}
