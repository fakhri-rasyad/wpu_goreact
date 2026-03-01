package utils

import "github.com/gofiber/fiber/v2"

type Response struct {
	Status       string      `json:"status"`
	ResponseCode int         `json:"response_code"`
	Message      string      `json:"message,omitempty"`
	Data         interface{} `json:"data,omitempty"`
	Error        interface{} `json:"error,omitempty"`
}

func Success(ctx *fiber.Ctx, message string, data interface{}) error {
	return ctx.Status(fiber.StatusOK).JSON(Response{
		Status: "Success",
		ResponseCode: fiber.StatusOK,
		Message: message,
		Data: data,
	})
}

func Created(ctx *fiber.Ctx, message string, data interface{}) error {
	return ctx.Status(fiber.StatusCreated).JSON(Response{
		Status: "Created",
		ResponseCode: fiber.StatusCreated,
		Message: message,
		Data: data,
	})
}

func BadRequest(ctx *fiber.Ctx, message string, data interface{}, err string) error {
	return ctx.Status(fiber.StatusBadRequest).JSON(Response{
		Status: "Error bad request",
		ResponseCode: fiber.StatusBadRequest,
		Message: message,
		Data: data,
		Error: err,
	})
}

func NotFound(ctx *fiber.Ctx, message string, data interface{}, err string) error {
	return ctx.Status(fiber.StatusNotFound).JSON(Response{
		Status: "Error not found",
		ResponseCode: fiber.StatusNotFound,
		Message: message,
		Data: data,
		Error: err,
	})
}