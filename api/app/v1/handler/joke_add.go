package handler

import (
	"context"

	"github.com/aldy505/jokes-bapak2-api/api/app/v1/models"
	"github.com/gofiber/fiber/v2"
)

func AddNewJoke(c *fiber.Ctx) error {
	var body models.RequestJokePost
	err := c.BodyParser(&body)
	if err != nil {
		return err
	}

	sql, args, err := psql.Insert("jokesbapak2").Columns("link", "key").Values(body.Link, body.Key).ToSql()
	if err != nil {
		return err
	}

	_, err = db.Query(context.Background(), sql, args...)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(models.ResponseJoke{
		Link: body.Link,
	})
}