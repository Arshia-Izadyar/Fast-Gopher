package handler

import (
	"github.com/Arshia-Izadyar/Fast-Gopher/src/api/dto"
	"github.com/Arshia-Izadyar/Fast-Gopher/src/api/helper"
	"github.com/Arshia-Izadyar/Fast-Gopher/src/config"
	"github.com/Arshia-Izadyar/Fast-Gopher/src/services"
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

	return c.JSON(helper.GenerateResponse("user created", true))
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
		return c.JSON(helper.GenerateResponse(res, true))
	}
}
