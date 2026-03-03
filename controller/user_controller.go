package controller

import (
	"math"
	"strconv"

	"github.com/fakhri-rasyad/wpu_goreact/models"
	"github.com/fakhri-rasyad/wpu_goreact/services"
	"github.com/fakhri-rasyad/wpu_goreact/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
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

func (c *UserController) GetUserPagination(ctx *fiber.Ctx) error {
	page, _ := strconv.Atoi(ctx.Query("page", "1"))
	limit, _ := strconv.Atoi(ctx.Query("limit", "10"))
	offset := (page - 1) * limit

	filter := ctx.Query("filter", "")
	sort := ctx.Query("filter", "")

	users, total, err := c.service.FindAllPagination(filter, sort, limit, offset)

	if err  != nil {
		return utils.BadRequest(ctx, "Gagal mengambil data", nil, err.Error())
	}

	var userResp []models.UserResponse
	_ = copier.Copy(&userResp, &users)

	meta := utils.PaginationMeta{
		Page: page,
		Limit: limit,
		Total: int(total),
		TotalPage: int(math.Ceil(float64(total) / float64(limit))),
		Filter: filter,
		Sort: sort,
	}

	if total == 0 {
		return utils.PaginationNotFound(ctx, "Data pengguna tidak ditemukan", userResp, meta)
	}

	return utils.PaginationSuccess(ctx, "Success", userResp, meta)
}

func (c *UserController) UpdateUser(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	public_id, err := uuid.Parse(id)

	if err != nil {
		return utils.BadRequest(ctx, "Invalid id format", nil, err.Error())
	}

	var user models.User
	if err := ctx.BodyParser(&user) ; err != nil {
		return utils.BadRequest(ctx, "Gagal parsing data", nil, err.Error())
	}

	user.PublicID = public_id

	err = c.service.Update(&user)
	if err != nil {
		return utils.BadRequest(ctx, "Gagal mengupdate data", nil, err.Error())
	}
	userUpdated, err := c.service.GetByPublicId(id)
	if err != nil {
		return utils.InternalServerError(ctx, "Gagal mengambil data", nil, err.Error())
	}

	var userPublic models.UserResponse
	err = copier.Copy(&userPublic, &userUpdated)
	if err != nil {
		return utils.InternalServerError(ctx, "Gagal mengambil data", nil, err.Error())
	}
	return utils.Success(ctx, "Success", userPublic)
}
