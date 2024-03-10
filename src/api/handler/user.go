package handler

import (
	"context"
	"io"
	"net/http"

	"github.com/Arshia-Izadyar/Fast-Gopher/src/api/dto"
	"github.com/Arshia-Izadyar/Fast-Gopher/src/api/helper"
	"github.com/Arshia-Izadyar/Fast-Gopher/src/config"
	"github.com/Arshia-Izadyar/Fast-Gopher/src/pkg/service_errors"
	"github.com/Arshia-Izadyar/Fast-Gopher/src/services"
	"github.com/bytedance/sonic"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	service *services.UserService
}

var validate = validator.New()

func NewUserHandler(cfg *config.Config) *UserHandler {
	return &UserHandler{
		service: services.NewUserService(cfg),
	}
}

func (uh *UserHandler) TestHandler(c *fiber.Ctx) error {
	req := &dto.UserCreateDTO{}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(helper.GenerateResponseWithError(err, false))
	}

	if err := validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(helper.GenerateResponseWithValidationError(err, false))
	}

	if err := uh.service.CreateUser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(helper.GenerateResponseWithError(err, false))
	}

	return c.Status(fiber.StatusCreated).JSON(helper.GenerateResponse("user created", true))
}

func (uh *UserHandler) LoginHandler(c *fiber.Ctx) error {
	req := &dto.UserDTO{}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(helper.GenerateResponseWithError(err, false))
	}

	if err := validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(helper.GenerateResponseWithValidationError(err, false))
	}

	if res, err := uh.service.LoginUser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(helper.GenerateResponseWithError(err, false))
	} else {
		return c.Status(fiber.StatusCreated).JSON(helper.GenerateResponse(res, true))
	}
}

func (uh *UserHandler) GoogleLogin(c *fiber.Ctx) error {

	url := config.AppConfig.GoogleLoginConfig.AuthCodeURL("randomstate")

	c.Status(fiber.StatusSeeOther)
	c.Redirect(url)
	return c.JSON(url)
}

func (uh *UserHandler) GoogleCallback(c *fiber.Ctx) error {
	state := c.Query("state")
	if state != "randomstate" {
		return &service_errors.ServiceError{EndUserMessage: "States don't Match!!"}
	}

	code := c.Query("code")

	googlecon := config.AppConfig.GoogleLoginConfig

	token, err := googlecon.Exchange(context.Background(), code)

	if err != nil {
		return &service_errors.ServiceError{EndUserMessage: "Code-Token Exchange Failed"}
	}

	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(helper.GenerateResponseWithError(&service_errors.ServiceError{EndUserMessage: "User Data Fetch Failed"}, false))

	}

	res, err := io.ReadAll(resp.Body)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(helper.GenerateResponseWithError(&service_errors.ServiceError{EndUserMessage: "User Data Fetch Failed"}, false))
	}

	var data *dto.GoogleUserInfoDTO
	err = sonic.Unmarshal(res, &data)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(helper.GenerateResponseWithError(&service_errors.ServiceError{EndUserMessage: "User Data Fetch Failed"}, false))

	}

	if res, err := uh.service.GoogleCallback(data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(helper.GenerateResponseWithError(err, false))
	} else {
		return c.Status(fiber.StatusCreated).JSON(helper.GenerateResponse(res, true))
	}
}
