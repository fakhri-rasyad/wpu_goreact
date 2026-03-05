package controller

import (
	"github.com/fakhri-rasyad/wpu_goreact/models"
	"github.com/fakhri-rasyad/wpu_goreact/services"
	"github.com/fakhri-rasyad/wpu_goreact/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type BoardController struct {
	service services.BoardService
}

func NewBoardController(s services.BoardService) *BoardController{
	return &BoardController{service: s}
}

func (c *BoardController) CreateBoard(ctx *fiber.Ctx) error {
	board := new(models.Board)
	user := ctx.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	err := ctx.BodyParser(board)

	if err != nil {
		return utils.BadRequest(ctx, "Data yang diberikan kurang", nil, err.Error())
	}

	userId, err := uuid.Parse(claims["pub_id"].(string))

	if err != nil {
		return utils.BadRequest(ctx, "Data yang diberikan kurang", nil, err.Error())
	}

	board.OwnerPublicID = userId

	err = c.service.Create(board)

	if err != nil {
		return utils.BadRequest(ctx, "Masalah saat memasukkan data", nil, err.Error())
	}

	return utils.Success(ctx, "Sukses memasukkan data", board)
}