package services

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Arshia-Izadyar/Fast-Gopher/src/api/dto"
	"github.com/Arshia-Izadyar/Fast-Gopher/src/common"
	"github.com/Arshia-Izadyar/Fast-Gopher/src/config"
	"github.com/Arshia-Izadyar/Fast-Gopher/src/constants"
	"github.com/Arshia-Izadyar/Fast-Gopher/src/data/models"
	"github.com/Arshia-Izadyar/Fast-Gopher/src/data/postgres"
	"github.com/Arshia-Izadyar/Fast-Gopher/src/data/redis"
	"github.com/Arshia-Izadyar/Fast-Gopher/src/pkg/service_errors"
	"github.com/bytedance/sonic"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	db               *sql.DB
	cfg              *config.Config
	whiteListService *WhiteListService
}

func NewUserService(cfg *config.Config) *UserService {
	db := postgres.GetDB()
	wl := NewWhiteListService(cfg)
	return &UserService{
		db:               db,
		cfg:              cfg,
		whiteListService: wl,
	}
}

func HashPassword(password string) (string, error) {
	bs, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(bs), nil
}

func (us *UserService) CreateUser(req *dto.UserCreateDTO) (err error) {
	var usrId uuid.UUID

	if req.UserPassword != req.UserPasswordConfirm {
		return &service_errors.ServiceError{EndUserMessage: service_errors.PasswordsDontMatch}
	}
	err = common.ValidatePassword(req.UserPassword)
	if err != nil {
		return err
	}

	req.UserPassword, err = HashPassword(req.UserPassword)
	if err != nil {
		return &service_errors.ServiceError{EndUserMessage: "hashing password gone wrong", Err: err}
	}

	q := `
	INSERT INTO public.users (email, user_password)
	VALUES ($1, $2) returning id;
	`
	// "returning" is a psql feature

	err = us.db.QueryRow(q, req.Email, req.UserPassword).Scan(&usrId)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return &service_errors.ServiceError{
				EndUserMessage: "user with this email already exists",
			}
		}
		return &service_errors.ServiceError{EndUserMessage: service_errors.BadRequest}
	}
	return nil
}

func (us *UserService) LoginUser(req *dto.UserDTO) (res *dto.UserTokenDTO, err error) {

	q := `
	SELECT id, user_password
	FROM users where email = $1;
	`

	var user models.User

	err = us.db.QueryRow(q, req.Email).Scan(&user.ID, &user.UserPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, &service_errors.ServiceError{EndUserMessage: "user not found"}
		}
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.UserPassword), []byte(req.UserPassword))
	if err != nil {
		return nil, &service_errors.ServiceError{EndUserMessage: service_errors.WrongPassword}
	}
	res, err = common.GenerateJwt(user.ID, us.cfg)
	if err != nil {
		return nil, &service_errors.ServiceError{EndUserMessage: "can't create JWT", Err: err}
	}
	return res, nil
}

func (us *UserService) GoogleCallback(req *dto.GoogleUserInfoDTO) (*dto.UserTokenDTO, error) {

	var userId uuid.UUID
	q := `
	SELECT id FROM 
	users WHERE email = $1;
	`

	err := us.db.QueryRow(q, req.Email).Scan(&userId)

	if err != nil {
		if err == sql.ErrNoRows {
			q := `
			INSERT INTO public.users (email, user_password)
			VALUES ($1, $2) returning id;
			`
			err = us.db.QueryRow(q, req.Email, nil).Scan(&userId)
			if err != nil {
				fmt.Println(err)
				return nil, &service_errors.ServiceError{EndUserMessage: service_errors.InternalError}
			}
		} else {
			return nil, &service_errors.ServiceError{EndUserMessage: service_errors.InternalError}
		}
	}

	tk, err := common.GenerateJwt(userId, us.cfg)
	if err != nil {
		fmt.Println(err)

		return nil, &service_errors.ServiceError{EndUserMessage: service_errors.InternalError, Err: err}
	}
	return tk, nil
}

func (us *UserService) GoogleLoginWithCode(req *dto.GoogleCodeLoginDTO) (*dto.UserTokenDTO, error) {
	googlecon := config.AppConfig.GoogleLoginConfig

	token, err := googlecon.Exchange(context.Background(), req.Code)

	if err != nil {
		return nil, &service_errors.ServiceError{EndUserMessage: "Code-Token Exchange Failed"}
	}

	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		return nil, &service_errors.ServiceError{EndUserMessage: "User Data Fetch Failed"}
	}

	res, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, &service_errors.ServiceError{EndUserMessage: "User Data READ Failed"}
	}

	var data *dto.GoogleUserInfoDTO
	err = sonic.Unmarshal(res, &data)
	if err != nil {
		return nil, &service_errors.ServiceError{EndUserMessage: "User Data Json unmarshal Failed"}
	}
	return us.GoogleCallback(data)

}

func (us *UserService) Logout(req *dto.UserLogout) error {

	err := redis.Set[bool](req.UserToken, true, us.cfg.JWT.AccessTokenExpireDuration*time.Minute)
	if err != nil {
		return err
	}

	err = us.whiteListService.WhiteListRemove(&dto.WhiteListAddDTO{
		UserId:       req.UserId,
		UserDeviceID: req.UserDeviceID,
		UserIp:       req.UserIp,
	})
	if err != nil {
		return err
	}
	// TODO: send command to remove from whitelist
	return nil
}

// refresh
func (us *UserService) Refresh(req *dto.RefreshTokenDTO) (*dto.UserTokenDTO, error) {

	// 1. check if refresh is used
	// 2. check if its is a refresh
	// 3. blacklist refresh
	// 4. issue new jwt

	claims, err := common.ValidateToken(req.RefreshToken, us.cfg)
	if err != nil {
		return nil, err
	}
	if claims[constants.AccessType] == true {
		return nil, &service_errors.ServiceError{EndUserMessage: service_errors.NotRefreshToken}
	}

	_, err = redis.Get[bool](req.RefreshToken)
	if err == nil {
		return nil, &service_errors.ServiceError{EndUserMessage: service_errors.TokenInvalid}
	}

	go func() {
		redis.Set[bool](req.RefreshToken, true, time.Minute*us.cfg.JWT.RefreshTokenExpireDuration)
	}()
	userUUid, err := uuid.Parse(claims[constants.UserIdKey].(string))
	if err != nil {
		return nil, &service_errors.ServiceError{EndUserMessage: service_errors.InternalError, Err: err}
	}
	res, err := common.GenerateJwt(userUUid, us.cfg)
	if err != nil {
		return nil, &service_errors.ServiceError{EndUserMessage: "JWT generation gone wrong", Err: err}
	}
	return res, nil
}
