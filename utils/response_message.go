package utils

import "github.com/gofiber/fiber/v2"

type ResponseStatus string

const (
	StatusOK    ResponseStatus = "ok"
	StatusError ResponseStatus = "error"
)

type ResponseMessageStruct struct {
	Status  ResponseStatus `json:"status"`
	Message string         `json:"message"`
}

func ResponseMessageSuccess(c *fiber.Ctx, status int, msg string) error {
	return c.Status(status).JSON(ResponseMessageStruct{
		Status:  StatusOK,
		Message: msg,
	})
}

func ResponseMessageError(c *fiber.Ctx, status int, msg string) error {
	return c.Status(status).JSON(ResponseMessageStruct{
		Status:  StatusError,
		Message: msg,
	})
}
