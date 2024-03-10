package services

import (
	"database/sql"
	"fmt"

	"github.com/Arshia-Izadyar/Fast-Gopher/src/api/dto"
	"github.com/Arshia-Izadyar/Fast-Gopher/src/common"
	"github.com/Arshia-Izadyar/Fast-Gopher/src/config"
	"github.com/Arshia-Izadyar/Fast-Gopher/src/data/models"
	"github.com/Arshia-Izadyar/Fast-Gopher/src/data/postgres"
	"github.com/Arshia-Izadyar/Fast-Gopher/src/pkg/service_errors"
	"github.com/gofiber/fiber/v2"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	db  *sql.DB
	cfg *config.Config
}

func NewUserService(cfg *config.Config) *UserService {
	db := postgres.GetDB()
	return &UserService{
		db:  db,
		cfg: cfg,
	}
}

func HasPassword(password string) (string, error) {
	bs, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(bs), nil

}

func (us *UserService) CreateUser(req *dto.UserCreateDTO) (err error) {

	if req.UserPassword != req.UserPasswordConfirm {
		return &service_errors.ServiceError{EndUserMessage: service_errors.PasswordsDontMatch}
	}

	req.UserPassword, err = HasPassword(req.UserPassword)
	if err != nil {
		return &service_errors.ServiceError{EndUserMessage: service_errors.InternalError, Err: err}
	}

	q := `
	INSERT INTO public.users (email, user_password)
	VALUES ($1, $2) returning id;
	`
	// "returning" is a psql feature

	var usrId uuid.UUID
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
		return nil, err
	}
	return res, nil
}

func (us *UserService) GoogleLogin(c *fiber.Ctx) error {

	url := config.AppConfig.GoogleLoginConfig.AuthCodeURL("randomstate")

	c.Status(fiber.StatusSeeOther)
	c.Redirect(url)
	return c.JSON(url)
}

func (us *UserService) GoogleCallback(req *dto.GoogleUserInfoDTO) (*dto.UserTokenDTO, error) {

	q := `
	SELECT id FROM 
	users WHERE email = $1;
	`

	var userId uuid.UUID
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
