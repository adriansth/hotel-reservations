package middleware

import (
	"fmt"

	"github.com/adriansth/go-hotel-reservations/types"
	"github.com/gofiber/fiber/v2"
)

func AdminAuth(c *fiber.Ctx) error {
	user, ok := c.Context().UserValue("user").(*types.User)
	if !ok {
		return fmt.Errorf("Not authorized")
	}
	if !user.IsAdmin {
		return fmt.Errorf("Not authorized")
	}
	return c.Next()
}
