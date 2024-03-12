package services

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/Arshia-Izadyar/Fast-Gopher/src/api/dto"
	"github.com/Arshia-Izadyar/Fast-Gopher/src/config"
	"github.com/Arshia-Izadyar/Fast-Gopher/src/constants"
	"github.com/Arshia-Izadyar/Fast-Gopher/src/data/cache"
	"github.com/Arshia-Izadyar/Fast-Gopher/src/data/postgres"
	"github.com/Arshia-Izadyar/Fast-Gopher/src/pkg/service_errors"
)

type OtpService struct {
	cfg *config.Config
	db  *sql.DB
	exp time.Duration
}

func NewOtpService(cfg *config.Config) *OtpService {
	db := postgres.GetDB()
	t := cfg.Otp.ExpireTime * time.Minute
	return &OtpService{
		cfg: cfg,
		db:  db,
		exp: t,
	}
}

func (o *OtpService) SendOtp(req *dto.SendOtpDTO) *service_errors.ServiceErrors {

	var emailExists int

	q := `
	SELECT COUNT(*) FROM users WHERE email = $1;
	`
	err := o.db.QueryRow(q, req.Email).Scan(&emailExists)
	if err != nil {
		return &service_errors.ServiceErrors{EndUserMessage: service_errors.InternalError, Status: 500}

	}
	fmt.Println(emailExists)
	if emailExists == 0 {
		return nil
	}
	key := fmt.Sprintf("%s_%s", constants.DefaultRedisKey, req.Email)
	storedOtp, err := cache.Get[dto.OtpDTO](key)
	fmt.Println(storedOtp)
	fmt.Println(err)
	if err == nil && !storedOtp.Used {
		return &service_errors.ServiceErrors{EndUserMessage: service_errors.OtpExists, Status: 429}
	}

	otpRedisStore := &dto.OtpDTO{
		Used: false,
		Otp:  req.Otp,
	}
	err = cache.Set[dto.OtpDTO](key, *otpRedisStore, o.exp)
	if err != nil {
		return &service_errors.ServiceErrors{EndUserMessage: service_errors.InternalError, Status: 500}
	}
	// TODO: Send email
	return nil
}

func (o *OtpService) ValidateOtp(email, otp string) error {
	key := fmt.Sprintf("%s_%s", constants.DefaultRedisKey, email)
	v, err := cache.Get[dto.OtpDTO](key)
	if err != nil {
		return &service_errors.ServiceError{EndUserMessage: service_errors.OtpInvalid}
	}
	if v.Used {
		return &service_errors.ServiceError{EndUserMessage: service_errors.OtpUsed}
	}
	if !v.Used && v.Otp != otp {
		return &service_errors.ServiceError{EndUserMessage: service_errors.OtpInvalid}
	}
	v.Used = true
	err = cache.Set[dto.OtpDTO](key, *v, o.exp)
	if err != nil {
		return &service_errors.ServiceError{EndUserMessage: service_errors.InternalError, Err: err}
	}
	return nil
}
