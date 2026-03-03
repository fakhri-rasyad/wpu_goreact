package utils

import "github.com/gofiber/fiber/v2"

type Response struct {
	Status       string      `json:"status"`
	ResponseCode int         `json:"response_code"`
	Message      string      `json:"message,omitempty"`
	Data         interface{} `json:"data,omitempty"`
	Error        interface{} `json:"error,omitempty"`
}

type PaginationMeta struct {
	Page int `json:"page" example:"1"`
	Limit int `json:"limit" example:"10"`
	Total int `json:"total" example:"10"`
	TotalPage int `json:"total_pages" example:"10"`
	Filter string `json:"filter" example:"nama=triady"`
	Sort string `json:"sort" example:"-id"`
}

type SuccessPagination struct {
	Status       string      `json:"status"`
	ResponseCode int         `json:"response_code"`
	Message      string      `json:"message,omitempty"`
	Data         interface{} `json:"data,omitempty"`
	Error        interface{} `json:"error,omitempty"`
	Meta PaginationMeta `json:"meta"`
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

func InternalServerError(ctx *fiber.Ctx, message string, data interface{}, err string) error {
	return ctx.Status(fiber.StatusInternalServerError).JSON(Response{
		Status: "Internal server error",
		ResponseCode: fiber.StatusInternalServerError,
		Message: message,
		Data: data,
		Error: err,
	})
}

func UnauthorizedRequest(ctx *fiber.Ctx, message string, data interface{}, err string) error {
	return ctx.Status(fiber.StatusUnauthorized).JSON(Response{
		Status: "Error unauthorized",
		ResponseCode: fiber.StatusUnauthorized,
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

func PaginationSuccess(ctx *fiber.Ctx, message string, data interface{}, meta PaginationMeta) error {
	return ctx.Status(fiber.StatusOK).JSON(SuccessPagination{
		Status: "Success",
		ResponseCode: fiber.StatusOK,
		Message: message,
		Data: data,
		Meta: meta,
	})
}

func PaginationNotFound(ctx *fiber.Ctx, message string, data interface{}, meta PaginationMeta) error {
	return ctx.Status(fiber.StatusNotFound).JSON(SuccessPagination{
		Status: "Not found",
		ResponseCode: fiber.StatusNotFound,
		Message: message,
		Data: data,
		Meta: meta,
	})
}