package controller

import (
	"errors"
	"time"

	"github.com/fakhri-rasyad/wpu_goreact/models"
	"github.com/fakhri-rasyad/wpu_goreact/services"
	"github.com/fakhri-rasyad/wpu_goreact/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type CardController struct {
	service services.CardServices
}

func NewCardController(s services.CardServices) *CardController{
	return &CardController{
		service: s,
	}
}

func (c *CardController) CreateCard(ctx *fiber.Ctx) error {
	type CreateCardRequest struct {
		ListPubId string `json:"list_id"`
		Title string `json:"title"`
		Description string `json:"description"`
		DueDate time.Time `json:"due_date"`
		Position int `json:"position"`

	}

	var req CreateCardRequest
	if err := ctx.BodyParser(&req) ; err != nil {
		return utils.BadRequest(ctx, "Gagal mengambil data", nil, err.Error())
	}

	card := &models.Card{
		Title: req.Title,
		Description: req.Description,
		DueDate: &req.DueDate,
		Position: req.Position,
	}

	if err := c.service.Create(card, req.ListPubId) ; err != nil {
		return utils.InternalServerError(ctx, "Gagal membuat card", nil, err.Error())
	}

	return utils.Success(ctx, "Sukses menambahkan card", card)
}

func (c *CardController) UpdateCard(ctx *fiber.Ctx) error {
	publicId := ctx.Params("id")

	type updateCardRequets struct{
		ListPubId string `json:"list_id"`
		Title string `json:"title"`
		Description string `json:"description"`
		DueDate *time.Time `json:"due_date"`
		Position int `json:"position"`
	}

	var req updateCardRequets

	if err := ctx.BodyParser(&req); err != nil {
		return utils.BadRequest(ctx, "Gagal parsing data", nil ,err.Error())
	}

	if _, err := uuid.Parse(publicId); err != nil {
		return utils.InternalServerError(ctx, "ID tidak valid", nil, err.Error())
	}

	card := &models.Card{
		Title: req.Title,
		Description: req.Description,
		Position: req.Position,
		DueDate: req.DueDate,
		PublicID: uuid.MustParse(publicId),
	}

	if err := c.service.Update(card, req.ListPubId); err != nil {
		return utils.BadRequest(ctx, "Failed to update card", nil, err.Error())
	}

	return utils.Success(ctx, "Card has been updated", card)
}

func (c *CardController) DeleteCard(ctx *fiber.Ctx) error {
	cardPubId := ctx.Params("id")

	if _, err := uuid.Parse(cardPubId); err != nil {
		return utils.BadRequest(ctx, "Gagal parsing uid", nil, err.Error())
	}

	card, err := c.service.GetByPublicID(cardPubId)
	if err != nil {
		return utils.NotFound(ctx, "Card not found", nil, err.Error())
	}

	if err := c.service.Delete(uint(card.InternalID)) ; err != nil {
		return utils.InternalServerError(ctx, "Failed to delete card", nil, err.Error())
	}

	return utils.Success(ctx, "Card has been deleted", cardPubId)
}

func (c *CardController) GetCardDetail(ctx *fiber.Ctx) error {
	cardPubId := ctx.Params("id")

	if _, err := uuid.Parse(cardPubId); err != nil {
		return utils.BadRequest(ctx, "Gagal parsing uid", nil, err.Error())
	}

	card, err := c.service.GetByPublicID(cardPubId)
	if err != nil {
		return utils.BadRequest(ctx, "Failed to get card", nil, err.Error())
	}

	if card == nil {
		return utils.NotFound(ctx, "Card not found", nil, errors.New("Card was not found").Error())
	}

	return utils.Success(ctx, "Success", card)
}