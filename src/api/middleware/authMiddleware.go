package authentication

import (
	"strings"

	"github.com/Arshia-Izadyar/Fast-Gopher/src/api/helper"
	"github.com/Arshia-Izadyar/Fast-Gopher/src/common"
	"github.com/Arshia-Izadyar/Fast-Gopher/src/config"
	"github.com/Arshia-Izadyar/Fast-Gopher/src/constants"
	"github.com/Arshia-Izadyar/Fast-Gopher/src/data/redis"
	"github.com/Arshia-Izadyar/Fast-Gopher/src/pkg/service_errors"
	"github.com/gofiber/fiber/v2"
)

func New(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(helper.GenerateResponseWithError(&service_errors.ServiceError{EndUserMessage: service_errors.TokenNotPresent}, false))
		}
		slicedToken := strings.Split(authHeader, " ")
		if len(slicedToken) != 2 || slicedToken[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(helper.GenerateResponseWithError(&service_errors.ServiceError{EndUserMessage: service_errors.TokenInvalidFormat}, false))
		}
		_, err := redis.Get[bool](slicedToken[1])
		if err == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(helper.GenerateResponseWithError(&service_errors.ServiceError{EndUserMessage: service_errors.TokenInvalid}, false))
		}
		tk, err := common.ValidateToken(slicedToken[1], cfg)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(helper.GenerateResponseWithError(err, false))
		}
		c.Locals(constants.UserIdKey, tk[constants.UserIdKey].(string))
		c.Locals(constants.AuthenticationKey, slicedToken[1])
		return c.Next()
	}
}
