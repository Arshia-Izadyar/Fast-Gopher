package handler

import (
	"context"
	"errors"
	"io"
	"net/http"

	"github.com/Arshia-Izadyar/Fast-Gopher/src/api/dto"
	"github.com/Arshia-Izadyar/Fast-Gopher/src/api/helper"
	"github.com/Arshia-Izadyar/Fast-Gopher/src/common"
	"github.com/Arshia-Izadyar/Fast-Gopher/src/config"
	"github.com/Arshia-Izadyar/Fast-Gopher/src/constants"
	"github.com/Arshia-Izadyar/Fast-Gopher/src/pkg/service_errors"
	"github.com/Arshia-Izadyar/Fast-Gopher/src/services"
	"github.com/bytedance/sonic"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
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

// UserRegister godoc
// @Summary Create a user
// @Description create a new user
// @Tags User
// @Accept json
// @produces json
// @Param Request body dto.UserCreateDTO true "Create a user"
// @Success 200 {object} helper.Response "Create a user response"
// @Failure 400 {object} helper.Response "Bad request"
// @Router /register [post]
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

// UserLogin godoc
// @Summary login a user
// @Description login a user
// @Tags User
// @Accept json
// @produces json
// @Param Request body dto.UserDTO true "Create a user"
// @Success 200 {object} helper.Response "Create a user response"
// @Failure 400 {object} helper.Response "Bad request"
// @Router /login [post]
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

// GoogleLogin godoc
// @Summary login a user with google
// @Description login a user with google
// @Tags User
// @Accept json
// @produces json
// @Success 200 {object} helper.Response "Create a user response"
// @Failure 400 {object} helper.Response "Bad request"
// @Router /google [get]
func (uh *UserHandler) GoogleLogin(c *fiber.Ctx) error {

	url := config.AppConfig.GoogleLoginConfig.AuthCodeURL("randomstate")

	c.Status(fiber.StatusSeeOther)
	c.Redirect(url)
	return c.JSON(url)
}

// GoogleCallback godoc
// @Summary login a user
// @Description login a user
// @Tags User
// @Accept json
// @produces json
// @Success 200 {object} helper.Response "Create a user response"
// @Failure 400 {object} helper.Response "Bad request"
// @Router /google/login [get]
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
		return c.Status(fiber.StatusBadRequest).JSON(helper.GenerateResponseWithError(err, false))
	}

	res, err := io.ReadAll(resp.Body)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(helper.GenerateResponseWithError(err, false))
	}

	var data *dto.GoogleUserInfoDTO
	err = sonic.Unmarshal(res, &data)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(helper.GenerateResponseWithError(err, false))

	}

	if res, err := uh.service.GoogleCallback(data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(helper.GenerateResponseWithError(err, false))
	} else {
		return c.Status(fiber.StatusCreated).JSON(helper.GenerateResponse(res, true))
	}
}

// LoginWithGoogleCode godoc
// @Summary login a user with Code from google call back
// @Description login a user
// @Tags User
// @Accept json
// @produces json
// @Success 200 {object} helper.Response "Create a user response"
// @Failure 400 {object} helper.Response "Bad request"
// @Router /auth/callback/google [get]
// @Security None
func (h *UserHandler) LoginWithGoogleCode(c *fiber.Ctx) error {
	code := c.Query("code", "")
	req := &dto.GoogleCodeLoginDTO{
		Code: code,
	}
	if req.Code == "" {
		return c.Status(fiber.StatusBadRequest).JSON(helper.GenerateResponseWithError(&service_errors.ServiceError{EndUserMessage: "please send google code as a string in body"}, false))
	}
	tk, err := h.service.GoogleLoginWithCode(req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(helper.GenerateResponseWithError(err, false))
	}
	return c.Status(fiber.StatusCreated).JSON(helper.GenerateResponse(tk, true))
}

// Logout godoc
// @Summary User logout
// @Description Logs out a user by invalidating the user's session.
// @Tags User
// @Accept json
// @Produce json
// @Param AuthenticationKey header string true "Authentication Token"
// @Param DeviceIdKey header string true "Device-Id"
// @Success 204 {object} map[string]interface{} "message: user logged out"
// @Failure 400 {object} map[string]interface{} "message: error message"
// @Router /logout [get]
// @Security ApiKeyAuth
func (h *UserHandler) Logout(c *fiber.Ctx) error {

	token := c.Locals(constants.AuthenticationKey).(string)
	userId := c.Locals(constants.UserIdKey).(string)
	devId := c.Get(constants.DeviceIdKey)

	parsedUserId, err := uuid.Parse(userId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(helper.GenerateResponseWithError(&service_errors.ServiceError{EndUserMessage: "user uuid parse failed"}, false))

	}
	deviceId, err := uuid.Parse(devId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(helper.GenerateResponseWithError(&service_errors.ServiceError{EndUserMessage: "user device Id parse failed"}, false))

	}
	req := &dto.UserLogout{
		UserId:       parsedUserId,
		UserDeviceID: deviceId,
		UserIp:       c.IP(),
		UserToken:    token,
	}
	err = h.service.Logout(req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(helper.GenerateResponse("user logged out failed", false))
	}
	return c.Status(fiber.StatusNoContent).JSON(helper.GenerateResponse("user logged out", true))
}

// Refresh godoc
// @Summary User Refresh
// @Description generate a new token from refresh.
// @Tags User
// @Produce json
// @Param Request body dto.RefreshTokenDTO true "Create a new token"
// @Success 201 {object} dto.UserTokenDTO "message: new rToken and aToken"
// @Failure 400 {object} helper.Response "message: error message"
// @Router /refresh [POST]
func (h *UserHandler) Refresh(c *fiber.Ctx) error {
	req := &dto.RefreshTokenDTO{}
	err := c.BodyParser(&req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(helper.GenerateResponseWithError(&service_errors.ServiceError{EndUserMessage: "cant parse body", Err: err}, false))
	}
	if err := validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(helper.GenerateResponseWithValidationError(err, false))
	}

	res, err := h.service.Refresh(req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(helper.GenerateResponseWithError(err, false))
	}
	return c.Status(fiber.StatusCreated).JSON(helper.GenerateResponse(res, true))
}

// Reset godoc
// @Summary User Refresh
// @Description Reset User password.
// @Tags User
// @Produce json
// @Param Request body dto.ResetPasswordDTO true "request for password change"
// @Success 202 {object} helper.Response "message: password Changed"
// @Failure 400 {object} helper.Response "message: error message"
// @Router /reset [PUT]
func (h *UserHandler) ResetPassword(c *fiber.Ctx) error {
	req := &dto.ResetPasswordDTO{}
	err := c.BodyParser(&req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(helper.GenerateResponseWithError(&service_errors.ServiceError{EndUserMessage: "cant parse body", Err: err}, false))
	}
	if err := validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(helper.GenerateResponseWithValidationError(err, false))
	}
	err = common.ValidatePassword(req.NewPassword)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(helper.GenerateResponseWithError(err, false))
	}
	if req.NewPassword != req.NewPasswordConfirm {
		return c.Status(fiber.StatusBadRequest).JSON(helper.GenerateResponseWithError(&service_errors.ServiceError{EndUserMessage: service_errors.PasswordsDontMatch}, false))
	}

	req.UserId = c.Locals(constants.UserIdKey).(string)

	if sErr := h.service.ResetPassword(req); sErr != nil {
		return c.Status(sErr.Status).JSON(helper.GenerateResponseWithError(sErr, false))
	}
	return c.Status(fiber.StatusAccepted).JSON(helper.GenerateResponse("password Changed!", true))
}

// forgotPasswordOTP godoc
// @Summary User password rest
// @Description send otp for password forget.
// @Tags User
// @Produce json
// @Param Request body dto.ForgotPasswordOtpDTO true "send a otp for forgot password"
// @Success 200 {object} helper.Response "message: password otp sent"
// @Failure 400 {object} helper.Response "message: error message"
// @Router /forgot/otp [POST]
func (h *UserHandler) ForgotPasswordOtp(c *fiber.Ctx) error {
	req := &dto.ForgotPasswordOtpDTO{}
	c.BodyParser(&req)

	if err := validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(helper.GenerateResponseWithValidationError(err, false))
	}

	if sErr := h.service.ForgotPasswordOtp(req); sErr != nil {
		return c.Status(sErr.Status).JSON(helper.GenerateResponseWithError(sErr, false))
	}
	return c.Status(fiber.StatusOK).JSON(helper.GenerateResponse("password otp sent!", true))
}

// forgotPassword godoc
// @Summary User password forget
// @Description password forget.
// @Tags User
// @Produce json
// @Param Request body dto.ForgotPasswordDTO true "change users password"
// @Success 201 {object} helper.Response "message: password changed"
// @Failure 400 {object} helper.Response "message: error message"
// @Router /forgot [POST]
func (h *UserHandler) ForgotPassword(c *fiber.Ctx) error {
	req := &dto.ForgotPasswordDTO{}
	c.BodyParser(&req)

	err := common.ValidatePassword(req.NewPassword)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(helper.GenerateResponseWithError(err, false))
	}
	if err := validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(helper.GenerateResponseWithValidationError(err, false))
	}

	if req.NewPassword != req.NewPasswordConfirm {
		return c.Status(fiber.StatusBadRequest).JSON(helper.GenerateResponseWithValidationError(errors.New("passwords don't match"), false))
	}

	if sErr := h.service.ForgotPassword(req); sErr != nil {
		return c.Status(sErr.Status).JSON(helper.GenerateResponseWithError(sErr, false))
	}
	return c.Status(fiber.StatusAccepted).JSON(helper.GenerateResponse("password changed", true))
}
