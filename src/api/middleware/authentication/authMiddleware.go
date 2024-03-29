package authentication

import (
	"strings"

	"github.com/Arshia-Izadyar/Fast-Gopher/src/api/helper"
	"github.com/Arshia-Izadyar/Fast-Gopher/src/common"
	"github.com/Arshia-Izadyar/Fast-Gopher/src/config"
	"github.com/Arshia-Izadyar/Fast-Gopher/src/constants"
	"github.com/Arshia-Izadyar/Fast-Gopher/src/data/cache"
	"github.com/Arshia-Izadyar/Fast-Gopher/src/pkg/service_errors"
	"github.com/gofiber/fiber/v2"
)

func New(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(helper.GenerateResponseWithError(&service_errors.ServiceErrors{EndUserMessage: service_errors.TokenNotPresent}, false))
		}
		slicedToken := strings.Split(authHeader, " ")
		if len(slicedToken) != 2 || slicedToken[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(helper.GenerateResponseWithError(&service_errors.ServiceErrors{EndUserMessage: service_errors.TokenInvalidFormat}, false))
		}
		_, err := cache.Get[bool](slicedToken[1])
		if err == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(helper.GenerateResponseWithError(&service_errors.ServiceErrors{EndUserMessage: service_errors.TokenInvalid}, false))
		}
		tk, e := common.ValidateToken(slicedToken[1], cfg)
		if e != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(helper.GenerateResponseWithError(e, false))
		}
		c.Locals(constants.Key, tk[constants.Key].(string))
		c.Locals(constants.SessionIdKey, tk[constants.SessionIdKey].(string))
		return c.Next()
	}
}
