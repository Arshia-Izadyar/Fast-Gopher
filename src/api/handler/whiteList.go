package handler

import (
	"github.com/Arshia-Izadyar/Fast-Gopher/src/api/dto"
	"github.com/Arshia-Izadyar/Fast-Gopher/src/api/helper"
	"github.com/Arshia-Izadyar/Fast-Gopher/src/common"
	"github.com/Arshia-Izadyar/Fast-Gopher/src/config"
	"github.com/Arshia-Izadyar/Fast-Gopher/src/constants"
	"github.com/Arshia-Izadyar/Fast-Gopher/src/services"
	"github.com/gofiber/fiber/v2"
)

type WhiteListHandler struct {
	service *services.WhiteListService
	cfg     *config.Config
}

func NewWhiteListHandler(cfg *config.Config) *WhiteListHandler {
	return &WhiteListHandler{
		service: services.NewWhiteListService(cfg),
		cfg:     cfg,
	}
}

// Add godoc
// @Summary Add a device to the whitelist
// @Description Adds a device IP and its identifier to the user's whitelist, ensuring the device is allowed to access the service.
// @Tags Whitelist
// @Accept json
// @Produce json
// @Param        tk    query     string  false  "token"
// @Success 201 {object} helper.Response "Successfully whitelisted the device"
// @Failure 500 {object} helper.Response "Internal Server Error"
// @Router /w [get]
// @Security ApiKeyAuth
func (w *WhiteListHandler) Add(c *fiber.Ctx) error {
	token := c.Query("tk")
	claims, sErr := common.ValidateToken(token, w.cfg)
	if sErr != nil {
		return c.Status(sErr.Status).JSON(helper.GenerateResponseWithError(sErr, false))
	}
	key := claims[constants.Key].(string)
	sessionId := claims[constants.SessionIdKey].(string)

	req := &dto.WhiteListAddDTO{
		Key:       key,
		SessionId: sessionId,
		UserIp:    c.IP(),
	}
	err := w.service.WhiteListRequest(req)
	if err != nil {
		return c.Status(err.Status).JSON(helper.GenerateResponseWithError(err, false))
	}
	return c.Status(fiber.StatusCreated).JSON(helper.GenerateResponse("whitelisted", true))
}

// Add godoc
// @Summary Add a device to the whitelist (free premium)
// @Description Adds a device IP and its identifier to the user's whitelist, ensuring the device is allowed to access the service.
// @Tags Whitelist
// @Accept json
// @Produce json
// @Param        tk    query     string  false  "token"
// @Success 201 {object} helper.Response "Successfully whitelisted the device"
// @Failure 500 {object} helper.Response "Internal Server Error"
// @Router /w/premium [get]
// @Security ApiKeyAuth
func (w *WhiteListHandler) PremiumAdd(c *fiber.Ctx) error {
	token := c.Query("tk")
	claims, sErr := common.ValidateToken(token, w.cfg)
	if sErr != nil {
		return c.Status(sErr.Status).JSON(helper.GenerateResponseWithError(sErr, false))
	}
	key := claims[constants.Key].(string)
	sessionId := claims[constants.SessionIdKey].(string)

	req := &dto.WhiteListAddDTO{
		Key:       key,
		SessionId: sessionId,
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
// @Param        tk    query     string  false  "token"
// @Success 200 {object} helper.Response "Successfully removed whitelisted the device"
// @Failure 500 {object} helper.Response "Internal Server Error"
// @Router /rw [get]
// @Security ApiKeyAuth
func (w *WhiteListHandler) Remove(c *fiber.Ctx) error {
	token := c.Query("tk")
	claims, sErr := common.ValidateToken(token, w.cfg)
	if sErr != nil {
		return c.Status(sErr.Status).JSON(helper.GenerateResponseWithError(sErr, false))
	}
	key := claims[constants.Key].(string)
	sessionId := claims[constants.SessionIdKey].(string)

	req := &dto.WhiteListAddDTO{
		Key:       key,
		SessionId: sessionId,
		UserIp:    c.IP(),
	}
	err := w.service.WhiteListRemove(req)
	if err != nil {
		return c.Status(err.Status).JSON(helper.GenerateResponseWithError(err, false))
	}
	return c.Status(fiber.StatusOK).JSON(helper.GenerateResponse("device removed from whitelist", true))
}
