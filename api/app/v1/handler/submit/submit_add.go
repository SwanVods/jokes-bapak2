package submit

import (
	"context"
	"jokes-bapak2-api/app/v1/core"
	"jokes-bapak2-api/app/v1/handler"
	"jokes-bapak2-api/app/v1/models"
	"strings"
	"time"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/gofiber/fiber/v2"
)

func SubmitJoke(c *fiber.Ctx) error {
	var body models.Submission
	err := c.BodyParser(&body)
	if err != nil {
		return err
	}

	// Image and/or Link should not be empty
	if body.Image == "" && body.Link == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.Error{
			Error: "a link or an image should be supplied in a form of multipart/form-data",
		})
	}

	// Author should be supplied
	if body.Author == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.Error{
			Error: "an author key consisting on the format \"yourname <youremail@mail>\" must be supplied",
		})
	} else {
		// Validate format
		valid := core.ValidateAuthor(body.Author)
		if !valid {
			return c.Status(fiber.StatusBadRequest).JSON(models.Error{
				Error: "please stick to the format of \"yourname <youremail@mail>\" and within 200 characters",
			})
		}
	}

	var url string

	// Check link validity if link was provided
	if body.Link != "" {
		valid, err := core.CheckImageValidity(handler.Client, body.Link)
		if err != nil {
			return err
		}
		if !valid {
			return c.Status(fiber.StatusBadRequest).JSON(models.Error{
				Error: "URL provided is not a valid image",
			})
		}

		url = body.Link
	}

	// If image was provided
	if body.Image != "" {
		image := strings.NewReader(body.Image)

		url, err = core.UploadImage(handler.Client, image)
		if err != nil {
			return err
		}
	}

	now := time.Now().UTC().Format(time.RFC3339)

	sql, args, err := handler.Psql.
		Insert("submission").
		Columns("link", "created_at", "author").
		Values(url, now, body.Author).
		Suffix("RETURNING id,created_at,link,author,status").
		ToSql()
	if err != nil {
		return err
	}

	var submission []models.Submission
	result, err := handler.Db.Query(context.Background(), sql, args...)
	if err != nil {
		return err
	}
	defer result.Close()

	err = pgxscan.ScanAll(&submission, result)
	if err != nil {
		return err
	}

	return c.
		Status(fiber.StatusCreated).
		JSON(models.ResponseSubmission{
			Message: "Joke submitted. Please wait for a few days for admin to approve your submission.",
			Data:    submission[0],
		})
}
