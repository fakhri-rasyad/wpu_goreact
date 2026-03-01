package controller

import (
	"github.com/fakhri-rasyad/wpu_goreact/models"
	"github.com/fakhri-rasyad/wpu_goreact/services"
	"github.com/fakhri-rasyad/wpu_goreact/utils"
	"github.com/gofiber/fiber/v2"
)

type UserController struct {
	service services.UserService
}

func NewUserController(s services.UserService) *UserController{
	return &UserController{service: s}
} 

func (c *UserController) Register(ctx *fiber.Ctx) error {
	user := new(models.User)
	if err := ctx.BodyParser(user); err != nil{
		return utils.BadRequest(ctx, "Gagal parsing data",nil, err.Error())
	}

	if err := c.service.Register(user); err != nil{
		return utils.BadRequest(ctx, "Registrasi gagal", nil, err.Error())
	}

	return utils.Success(ctx, "Register Success", user)
}