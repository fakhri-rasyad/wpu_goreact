package repositories

import (
	"github.com/fakhri-rasyad/wpu_goreact/config"
	"github.com/fakhri-rasyad/wpu_goreact/models"
)

type BoardRepository interface {
	Create(board *models.Board) error 
}

type BoardRepositoryImpl struct {
}

func NewBoardRepostory() BoardRepository {
	return &BoardRepositoryImpl{}
}

func (r *BoardRepositoryImpl) Create(board *models.Board) error {
	return config.DB.Create(board).Error
}