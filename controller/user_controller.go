package controller

import (
	"github.com/fakhri-rasyad/wpu_goreact/models"
	"github.com/fakhri-rasyad/wpu_goreact/services"
	"github.com/fakhri-rasyad/wpu_goreact/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/copier"
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

	var userResponse models.UserResponse
	_ = copier.Copy(&userResponse, &user)
	return utils.Success(ctx, "Register Success", userResponse)
}

func (c *UserController) Login(ctx *fiber.Ctx) error {
	var body struct{
		Email string `json:"email"`
		Password string `json:"password"`
	}

	if err := ctx.BodyParser(&body); err != nil {
		return utils.BadRequest(ctx, "Gagal parsing data", nil, err.Error())
	}

	user , err := c.service.Login(body.Email, body.Password)

	if err != nil {
		return utils.UnauthorizedRequest(ctx, "Login gagal", nil, err.Error())
	}

	token, err := utils.GenerateJWTToken(user.InternalID, user.Role, user.Email, user.PublicID)
	refreshToken, err := utils.RefreshJWTToken(user.InternalID)
	var userResponse models.UserResponse
	_ = copier.Copy(&userResponse, &user)

	return utils.Success(ctx, "Login Success", fiber.Map{
		"access_token" : token,
		"refresh_token" : refreshToken,
		"user" : userResponse,
	})
}

func (c *UserController) GetUser(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	user, err := c.service.GetByPublicId(id)
	if err != nil {
		return utils.NotFound(ctx, "Data not found", nil, err.Error())
	}

	var userResp models.UserResponse
	err = copier.Copy(&userResp, &user)

	if err != nil {
		return utils.BadRequest(ctx, "Internal Server Error", nil, err.Error())
	}

	return utils.Success(ctx, "Request Success", userResp)
}
