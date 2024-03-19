package handler

import (
	"github.com/Arshia-Izadyar/Fast-Gopher/src/api/dto"
	"github.com/Arshia-Izadyar/Fast-Gopher/src/api/helper"
	"github.com/Arshia-Izadyar/Fast-Gopher/src/config"
	"github.com/Arshia-Izadyar/Fast-Gopher/src/constants"
	"github.com/Arshia-Izadyar/Fast-Gopher/src/services"
	"github.com/gofiber/fiber/v2"
)

type WhiteListHandler struct {
	service *services.WhiteListService
}

func NewWhiteListHandler(cfg *config.Config) *WhiteListHandler {
	return &WhiteListHandler{
		service: services.NewWhiteListService(cfg),
	}
}

// Add godoc
// @Summary Add a device to the whitelist
// @Description Adds a device IP and its identifier to the user's whitelist, ensuring the device is allowed to access the service.
// @Tags Whitelist
// @Accept json
// @Produce json
// @Success 201 {object} helper.Response "Successfully whitelisted the device"
// @Failure 500 {object} helper.Response "Internal Server Error"
// @Router /w [get]
// @Security ApiKeyAuth
func (w *WhiteListHandler) Add(c *fiber.Ctx) error {
	key := c.Locals(constants.Key).(string)
	SessionId := c.Locals(constants.SessionIdKey).(string)
	// sessionId := c.Get(constants.SessionIdKey)

	req := &dto.WhiteListAddDTO{
		Key:       key,
		SessionId: SessionId,
		UserIp:    c.IP(),
	}
	err := w.service.WhiteListRequest(req)
	if err != nil {
		return c.Status(err.Status).JSON(helper.GenerateResponseWithError(err, false))
	}
	return c.Status(fiber.StatusCreated).JSON(helper.GenerateResponse("whitelisted", true))
}

// Remove godoc
// @Summary remove a device to the whitelist
// @Description removes a device IP and its identifier to the user's whitelist, ensuring the device is not allowed to access the service.
// @Tags Whitelist
// @Accept json
// @Produce json
// @Success 201 {object} helper.Response "Successfully whitelisted the device"
// @Failure 500 {object} helper.Response "Internal Server Error"
// @Router /rw [get]
// @Security ApiKeyAuth
func (w *WhiteListHandler) Remove(c *fiber.Ctx) error {
	key := c.Locals(constants.Key).(string)
	sessionId := c.Locals(constants.SessionIdKey).(string)
	req := &dto.WhiteListAddDTO{
		Key:       key,
		SessionId: sessionId,
		UserIp:    "",
	}
	err := w.service.WhiteListRemove(req)
	if err != nil {
		return c.Status(err.Status).JSON(helper.GenerateResponseWithError(err, false))
	}
	return c.Status(fiber.StatusCreated).JSON(helper.GenerateResponse("device removed from whitelist", true))
}
