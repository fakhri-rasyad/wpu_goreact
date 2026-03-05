package controller

import (
	"github.com/fakhri-rasyad/wpu_goreact/models"
	"github.com/fakhri-rasyad/wpu_goreact/services"
	"github.com/fakhri-rasyad/wpu_goreact/utils"
	"github.com/gofiber/fiber/v2"
)

type BoardController struct {
	service services.BoardService
}

func NewBoardController(s services.BoardService) *BoardController{
	return &BoardController{service: s}
}

func (c *BoardController) CreateBoard(ctx *fiber.Ctx) error {
	board := new(models.Board)

	err := ctx.BodyParser(board)

	if err != nil {
		return utils.BadRequest(ctx, "Data yang diberikan kurang", nil, err.Error())
	}

	err = c.service.Create(board)

	if err != nil {
		return utils.BadRequest(ctx, "Masalah saat memasukkan data", nil, err.Error())
	}

	return utils.Success(ctx, "Sukses memasukkan data", board)
}