package controller

import (
	"github.com/fakhri-rasyad/wpu_goreact/models"
	"github.com/fakhri-rasyad/wpu_goreact/services"
	"github.com/fakhri-rasyad/wpu_goreact/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ListController struct {
	service services.ListService
}

func NewListController(s services.ListService) *ListController {
	return &ListController{service: s}
}

func (c *ListController) CreateList(ctx *fiber.Ctx) error {
	list := new(models.List)

	if err := ctx.BodyParser(list); err != nil {
		return utils.BadRequest(ctx, "Gagal bind data", nil, err.Error())
	}

	if err := c.service.Create(list); err != nil {
		return utils.BadRequest(ctx, "Gagal membuat list", nil, err.Error())
	}

	return utils.Success(ctx, "List berhasil dibuat", list)
}

func (c *ListController) UpdateList(ctx *fiber.Ctx) error{
	publicID := ctx.Params("id")
	list := new(models.List)

	if err := ctx.BodyParser(list); err != nil {
		return utils.BadRequest(ctx, "Gagal parsing data", nil, err.Error())
	}

	if _, err := uuid.Parse(publicID); err != nil {
		return utils.BadRequest(ctx, "ID tidak valid", nil, err.Error())
	}

	existingList, err := c.service.GetByPublicId(publicID)

	if err != nil {
		return utils.NotFound(ctx, "List tidak ditemukan", nil, err.Error())
	}

	list.InternalID = existingList.InternalID
	list.PublicID = existingList.PublicID
	if err := c.service.Update(list) ; err != nil {
		return utils.BadRequest(ctx, "Gagal update list",nil, err.Error())
	}

	updatedList, err := c.service.GetByPublicId(publicID)

	if err != nil {
		return utils.NotFound(ctx, "List tidak ditemukan", nil, err.Error())
	}

	return utils.Success(ctx, "List update success", updatedList)

}

func (c *ListController) GetListOnBoard(ctx *fiber.Ctx) error {
	boardPubId := ctx.Params("board_id")

	if _, err := uuid.Parse(boardPubId); err != nil {
		return utils.BadRequest(ctx, "Id board tidak valid", nil, err.Error())
	}

	list, err := c.service.GetByBoardId(boardPubId)

	if err != nil {
		return utils.NotFound(ctx, "List tidak ditemukan", nil, err.Error())
	}

	return utils.Success(ctx, "List diperoleh", list)

}

func (c *ListController) DeleteList(ctx *fiber.Ctx)error{
	listPubId := ctx.Params("id")

	parsedUUID, err := uuid.Parse(listPubId)

	if err != nil {
		return utils.BadRequest(ctx, "ID tidak valid", nil, err.Error())
	}

	list, err := c.service.GetByPublicId(parsedUUID.String())

	if err != nil {
		return utils.NotFound(ctx, "List tidak ditemukan", nil, err.Error())
	}

	if err := c.service.Delete(uint(list.InternalID)); err != nil {
		return utils.BadRequest(ctx, "Gagal menghapus list", nil, err.Error())
	}

	return utils.Success(ctx, "List berhasil dihapus", listPubId)
}