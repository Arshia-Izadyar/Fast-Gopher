package middleware

import (
	"fmt"

	"github.com/Arshia-Izadyar/Fast-Gopher/src/api/helper"
	"github.com/Arshia-Izadyar/Fast-Gopher/src/constants"
	"github.com/Arshia-Izadyar/Fast-Gopher/src/data/postgres"
	"github.com/Arshia-Izadyar/Fast-Gopher/src/pkg/service_errors"
	"github.com/gofiber/fiber/v2"
)

func Premium() fiber.Handler {
	db := postgres.GetDB()
	q := `SELECT premium FROM ac_keys WHERE id = $1;`
	return func(c *fiber.Ctx) error {
		key := c.Locals(constants.Key).(string)
		var premium bool
		err := db.QueryRow(q, key).Scan(&premium) // maybe handle error
		if err != nil {
			fmt.Println(err)
			return c.Status(fiber.StatusInternalServerError).JSON(helper.GenerateResponseWithError(&service_errors.ServiceErrors{EndUserMessage: "An error occurred"}, false))
		}
		if premium {

			return c.Next()
		}
		return c.Status(fiber.StatusForbidden).JSON(helper.GenerateResponseWithError(&service_errors.ServiceErrors{EndUserMessage: "you don't have permission to access to premium servers"}, false))
	}
}
