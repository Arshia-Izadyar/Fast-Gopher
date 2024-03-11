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

// Add godoc
// @Summary Add a device to the whitelist
// @Description Adds a device IP and its identifier to the user's whitelist, ensuring the device is allowed to access the service.
// @Tags Whitelist
// @Accept json
// @Produce json
// @Param Device-Id header string true "Device-Id"
// @Success 201 {object} helper.Response "Successfully whitelisted the device"
// @Failure 500 {object} helper.Response "Internal Server Error"
// @Router /w [get]
// @Security AuthBearer
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
		return c.Status(fiber.StatusInternalServerError).JSON(helper.GenerateResponseWithError(err, false))
	}
	return c.Status(fiber.StatusCreated).JSON(helper.GenerateResponse("whitelisted", true))
}

// Remove godoc
// @Summary remove a device to the whitelist
// @Description removes a device IP and its identifier to the user's whitelist, ensuring the device is not allowed to access the service.
// @Tags Whitelist
// @Accept json
// @Produce json
// @Param DeviceIdKey header string true "Device-Id"
// @Success 201 {object} helper.Response "Successfully whitelisted the device"
// @Failure 500 {object} helper.Response "Internal Server Error"
// @Router /rw [get]
// @Security AuthBearer
func (w *WhiteListHandler) Remove(c *fiber.Ctx) error {
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
	err = w.service.WhiteListRemove(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(helper.GenerateResponseWithError(err, false))
	}
	return c.Status(fiber.StatusCreated).JSON(helper.GenerateResponse("whitelist removed", true))
}
