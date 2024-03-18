package handler

import (
	"log"

	"github.com/Arshia-Izadyar/Fast-Gopher/src/api/dto"
	"github.com/Arshia-Izadyar/Fast-Gopher/src/api/helper"
	"github.com/Arshia-Izadyar/Fast-Gopher/src/config"
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

// Logout godoc
// @Summary Key logout
// @Description Logs out a user by invalidating the user's session.
// @Tags Key
// @Accept json
// @Produce json
// @Param AuthenticationKey header string true "Authentication Token"
// @Param DeviceIdKey header string true "Device-Id"
// @Success 204 {object} map[string]interface{} "message: user logged out"
// @Failure 400 {object} map[string]interface{} "message: error message"
// @Router /logout [get]
// @Security ApiKeyAuth
func (h *KeyHandler) Refresh(c *fiber.Ctx) error {
	req := &dto.RefreshTokenDTO{}
	err := c.BodyParser(&req)
	if err != nil {
		log.Fatal(err)
	}
	tk, sErr := h.service.Refresh(req)
	if sErr != nil {
		return c.Status(sErr.Status).JSON(helper.GenerateResponseWithError(sErr, true))
	}

	return c.Status(fiber.StatusOK).JSON(helper.GenerateResponse(tk, true))
}

// Refresh godoc
// @Summary Key Refresh
// @Description generate a new token from refresh.
// @Tags Key
// @Produce json
// @Param Request body dto.RefreshTokenDTO true "Create a new token"
// @Success 201 {object} dto.UserTokenDTO "message: new rToken and aToken"
// @Failure 400 {object} helper.Response "message: error message"
// @Router /refresh [POST]
// func (h *UserHandler) Refresh(c *fiber.Ctx) error {
// 	req := &dto.RefreshTokenDTO{}
// 	err := c.BodyParser(&req)
// 	if err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(helper.GenerateResponseWithError(&service_errors.ServiceErrors{EndUserMessage: "cant parse body", Err: err}, false))
// 	}
// 	if err := validate.Struct(req); err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(helper.GenerateResponseWithValidationError(err, false))
// 	}

// 	res, err := h.service.Refresh(req)
// 	if err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(helper.GenerateResponseWithError(err, false))
// 	}
// 	return c.Status(fiber.StatusCreated).JSON(helper.GenerateResponse(res, true))
// }

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
