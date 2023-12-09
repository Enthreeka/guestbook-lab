package controller

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"lab-3/repository"
)

func AuthMiddleware(token repository.SessionRepository, store *session.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		sess, err := store.Get(c)
		if err != nil {
			return c.SendStatus(fiber.StatusForbidden)
		}

		sessionToken := sess.Get("session_token")

		if sessionToken == nil {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		t, err := token.GetToken(context.Background(), sessionToken.(string))
		if err != nil {
			if err == repository.ErrNoRows {
				return c.SendStatus(fiber.StatusUnauthorized)
			}
			return c.Status(fiber.StatusUnauthorized).JSON(err)
		}

		c.Locals("token", t)

		return c.Next()
	}
}
