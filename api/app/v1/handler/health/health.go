package health

import (
	"context"
	"jokes-bapak2-api/app/v1/handler"
	"jokes-bapak2-api/app/v1/models"

	"github.com/gofiber/fiber/v2"
)

func Health(c *fiber.Ctx) error {
	// Ping REDIS database
	err := handler.Redis.Ping(context.Background()).Err()
	if err != nil {
		return c.
			Status(fiber.StatusServiceUnavailable).
			JSON(models.Error{
				Error: "REDIS: " + err.Error(),
			})
	}

	_, err = handler.Db.Query(context.Background(), "SELECT \"id\" FROM \"jokesbapak2\" LIMIT 1")
	if err != nil {
		return c.
			Status(fiber.StatusServiceUnavailable).
			JSON(models.Error{
				Error: "POSTGRESQL: " + err.Error(),
			})
	}
	return c.SendStatus(fiber.StatusOK)
}
