package repositories

import (
	"github.com/fakhri-rasyad/wpu_goreact/config"
	"github.com/fakhri-rasyad/wpu_goreact/models"
)

type BoardRepository interface {
	Create(board *models.Board) error 
	Update(board *models.Board) error
	FindByPublicID(uuid string) (*models.Board, error)
}

type BoardRepositoryImpl struct {
}

func NewBoardRepostory() BoardRepository {
	return &BoardRepositoryImpl{}
}

func (r *BoardRepositoryImpl) Create(board *models.Board) error {
	return config.DB.Create(board).Error
}

func (r *BoardRepositoryImpl) Update(board *models.Board) error {
	return config.DB.Model(&models.Board{}).Where("public_id = ?", board.PublicID).Updates(map[string]interface{}{
		"title" : board.Title,
		"description" : board.Description,
		"duedate" : board.Duedate,
	}).Error
}

func (r *BoardRepositoryImpl) FindByPublicID(uuid string) (*models.Board, error) {
	var board models.Board
	err := config.DB.Model(&models.Board{}).Where("public_id = ?", uuid).First(&board).Error
	return &board, err
}