package handler

import (
	"github.com/Arshia-Izadyar/Fast-Gopher/src/api/dto"
	"github.com/Arshia-Izadyar/Fast-Gopher/src/config"
	"github.com/Arshia-Izadyar/Fast-Gopher/src/services"
	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	service *services.UserService
}

func NewUserHandler(cfg *config.Config) *UserHandler {
	return &UserHandler{
		service: services.NewUserService(cfg),
	}
}

func (uh *UserHandler) TestHandler(c *fiber.Ctx) error {
	req := &dto.UserDTO{}
	err := c.BodyParser(&req)

	if err != nil {
		return err
	}
	err = uh.service.Init(req)
	if err != nil {
		return err
	}
	return nil
}
