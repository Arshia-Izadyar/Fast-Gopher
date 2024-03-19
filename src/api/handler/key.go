package handler

import (
	"log"

	"github.com/Arshia-Izadyar/Fast-Gopher/src/api/dto"
	"github.com/Arshia-Izadyar/Fast-Gopher/src/api/helper"
	"github.com/Arshia-Izadyar/Fast-Gopher/src/config"
	"github.com/Arshia-Izadyar/Fast-Gopher/src/constants"
	"github.com/Arshia-Izadyar/Fast-Gopher/src/pkg/service_errors"
	"github.com/Arshia-Izadyar/Fast-Gopher/src/services"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type KeyHandler struct {
	service *services.KeyService
}

var validate = validator.New()

func NewKeyHandler(cfg *config.Config) *KeyHandler {
	return &KeyHandler{
		service: services.NewKeyService(cfg),
	}
}

// Refresh godoc
// @Summary Refresh JWT with refresh_token
// @Description Refresh JWT with refresh_token and generate new tokens and will blacklist current refresh token.
// @Tags Authentication
// @Accept json
// @Produce json
// @Param Request body dto.RefreshTokenDTO true "Create a new atoken and rtoken"
// @Success 200 {object} dto.KeyAcDTO "message: helper.Response"
// @Failure 400 {object} helper.Response "message: helper.Response"
// @Router /auth/refresh [POST]
// @Security ApiKeyAuth
func (h *KeyHandler) Refresh(c *fiber.Ctx) error {
	req := &dto.RefreshTokenDTO{}
	err := c.BodyParser(&req)
	if err != nil {
		log.Fatal(err)
	}
	if req.RefreshToken == "" {
		return c.Status(fiber.StatusBadRequest).JSON(helper.GenerateResponseWithError(&service_errors.ServiceErrors{EndUserMessage: "please provide 'refresh_token'"}, true))
	}
	tk, sErr := h.service.Refresh(req)
	if sErr != nil {
		return c.Status(sErr.Status).JSON(helper.GenerateResponseWithError(sErr, true))
	}

	return c.Status(fiber.StatusOK).JSON(helper.GenerateResponse(tk, true))
}

// GenerateKey godoc
// @Summary Generate a new key
// @Description generate a new Key when new users install the APP.
// @Tags Key
// @Produce json
// @Param Request body dto.GenerateKeyDTO true "Create a new token"
// @Success 201 {object} dto.KeyAcDTO "message: new rToken and aToken + key"
// @Failure 400 {object} helper.Response "message: error message"
// @Router /key [POST]
func (h *KeyHandler) GenerateKey(c *fiber.Ctx) error {
	req := &dto.GenerateKeyDTO{}
	err := c.BodyParser(&req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(helper.GenerateResponseWithError(&service_errors.ServiceErrors{EndUserMessage: "parsing body failed"}, false))
	}
	if req.SessionId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(helper.GenerateResponseWithError(&service_errors.ServiceErrors{EndUserMessage: "please provide user session_id "}, false))
	}
	key, sErr := h.service.GenerateKey(req)
	if sErr != nil {
		return c.Status(sErr.Status).JSON(helper.GenerateResponseWithError(sErr, false))
	}
	return c.Status(fiber.StatusCreated).JSON(helper.GenerateResponse(key, true))
}

// GenerateTokenFromKey godoc
// @Summary Generate a new Token based on already existing key
// @Description generate a Key when.
// @Tags Key
// @Produce json
// @Param Request body dto.KeyDTO true "Create a new token"
// @Success 201 {object} dto.KeyAcDTO "message: new rToken and aToken"
// @Failure 400 {object} helper.Response "message: error message"
// @Failure 500 {object} helper.Response "message: fuck"
// @Router /key/tk [POST]
func (h *KeyHandler) GenerateTokenFromKey(c *fiber.Ctx) error {
	req := &dto.KeyDTO{}
	err := c.BodyParser(&req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(helper.GenerateResponseWithError(err, false))
	}
	err = validate.Struct(req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(helper.GenerateResponseWithValidationError(err, false))
	}
	res, sErr := h.service.GenerateTokenFromKey(req)
	if sErr != nil {
		return c.Status(sErr.Status).JSON(helper.GenerateResponseWithError(sErr, false))
	}
	return c.Status(fiber.StatusCreated).JSON(helper.GenerateResponse(res, true))

}

// ShowActiveSessions godoc
// @Summary show all of currently active sessions
// @Description show all of active devices.
// @Tags Key
// @Produce json
// @Success 201 {object} []dto.DeviceDTO "message: list of devices"
// @Failure 400 {object} helper.Response "message: error message"
// @Failure 500 {object} helper.Response "message: fuck"
// @Router /show [GET]
// @Security ApiKeyAuth
func (h *KeyHandler) ShowActiveSessions(c *fiber.Ctx) error {
	req := &dto.IKeyDTO{}
	req.Key = c.Locals(constants.Key).(string)
	res, sErr := h.service.ShowAllActiveDevices(req)
	if sErr != nil {
		return c.Status(sErr.Status).JSON(helper.GenerateResponseWithError(sErr, false))
	}
	return c.Status(fiber.StatusOK).JSON(helper.GenerateResponse(res, true))

}

// RemoveDevice godoc
// @Summary Remove a device
// @Description remove a single device from list\nafter removing a device from list you get a 403 error on whitelist end point after that get a new key from /key.
// @Tags Key
// @Produce json
// @Param Request body dto.RemoveDeviceDTO true "info for device"
// @Success 204 {object} helper.Response "message: removed"
// @Failure 400 {object} helper.Response "message: error message"
// @Failure 500 {object} helper.Response "message: fuck"
// @Router /rm [DELETE]
// @Security ApiKeyAuth
func (h *KeyHandler) RemoveDevice(c *fiber.Ctx) error {
	req := &dto.RemoveDeviceDTO{}
	err := c.BodyParser(&req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(helper.GenerateResponseWithError(err, false))
	}
	err = validate.Struct(req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(helper.GenerateResponseWithValidationError(err, false))
	}
	sErr := h.service.DeleteSession(req)
	if sErr != nil {
		return c.Status(sErr.Status).JSON(helper.GenerateResponseWithError(sErr, false))
	}
	return c.SendStatus(fiber.StatusNoContent)
}

// RemoveAllDevices godoc
// @Summary Remove all devices
// @Description remove all devices from list except the current device\nafter removing a device from list you get a 403 error on whitelist end point after that get a new key from /key.
// @Tags Key
// @Produce json
// @Success 204 {object} helper.Response "message: removed"
// @Failure 400 {object} helper.Response "message: error message"
// @Failure 500 {object} helper.Response "message: fuck"
// @Router /rm/all [DELETE]
// @Security ApiKeyAuth
func (h *KeyHandler) RemoveAllDevices(c *fiber.Ctx) error {

	req := &dto.SessionKeyDTO{}
	req.Key = c.Locals(constants.Key).(string)
	req.SessionId = c.Locals(constants.SessionIdKey).(string)
	sErr := h.service.DeleteAllSessions(req)
	if sErr != nil {
		return c.Status(sErr.Status).JSON(helper.GenerateResponseWithError(sErr, false))
	}
	return c.SendStatus(fiber.StatusNoContent)
}
