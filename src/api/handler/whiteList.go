package handler

import (
	"github.com/Arshia-Izadyar/Fast-Gopher/src/api/dto"
	"github.com/Arshia-Izadyar/Fast-Gopher/src/api/helper"
	"github.com/Arshia-Izadyar/Fast-Gopher/src/config"
	"github.com/Arshia-Izadyar/Fast-Gopher/src/constants"
	"github.com/Arshia-Izadyar/Fast-Gopher/src/pkg/service_errors"
	"github.com/Arshia-Izadyar/Fast-Gopher/src/services"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type WhiteListHandler struct {
	service *services.WhiteListService
}

func NewWhiteListHandler(cfg *config.Config) *WhiteListHandler {
	return &WhiteListHandler{
		service: services.NewWhiteListService(cfg),
	}
}

func (w *WhiteListHandler) Add(c *fiber.Ctx) error {
	v := c.Locals(constants.UserIdKey).(string)
	devId := c.Get(constants.DeviceIdKey)
	uid, err := uuid.Parse(v)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(helper.GenerateResponseWithError(err, false))
	}

	deviceId, err := uuid.Parse(devId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(helper.GenerateResponseWithError(&service_errors.ServiceError{EndUserMessage: "cant parse device uuid"}, false))
	}

	req := &dto.WhiteListAddDTO{
		UserId:       uid,
		UserDeviceID: deviceId,
		UserIp:       c.IP(),
	}
	err = w.service.WhiteListRequest(req)
	if err != nil {
		return c.JSON(fiber.Map{"err": err})
	}
	return c.Send([]byte("ok"))
}
