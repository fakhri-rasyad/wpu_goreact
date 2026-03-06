package repositories

import (
	"github.com/fakhri-rasyad/wpu_goreact/config"
	"github.com/fakhri-rasyad/wpu_goreact/models"
)

type ListRepository interface {
}

type ListRepositoryImpl struct{}

func NewListRepository() ListRepository {
	return &ListRepositoryImpl{}
}

func (r *ListRepositoryImpl) Create(list *models.List) error {
	return config.DB.Create(list).Error
}