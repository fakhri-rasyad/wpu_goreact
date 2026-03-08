package repositories

import (
	"github.com/fakhri-rasyad/wpu_goreact/config"
	"github.com/fakhri-rasyad/wpu_goreact/models"
	"github.com/google/uuid"
)

type ListRepository interface {
	Create(list *models.List) error
	Update(list *models.List) error
	Delete(id uint) error
	UpdatePosition(boardPubId string, position []string) error
	GetCardOrder(listPublicId string) ([]uuid.UUID, error)
	FindByBoardId(boardId string) ([]models.List, error)
	FindByPubId(pubId string) (*models.List, error)
	FindById(id uint) (*models.List, error)
}

type ListRepositoryImpl struct{}

func NewListRepository() ListRepository {
	return &ListRepositoryImpl{}
}

func (r *ListRepositoryImpl) Create(list *models.List) error {
	return config.DB.Create(list).Error
}

func (r *ListRepositoryImpl) Update(list *models.List) error {
	return config.DB.Model(&models.List{}).Where("public_id = ?", list.PublicID).Updates(map[string]interface{}{
		"title" : list.Title,
	}).Error
}

func (r *ListRepositoryImpl) Delete(id uint) error {
	return config.DB.Delete(&models.List{}, id).Error
}

func (r *ListRepositoryImpl) UpdatePosition(boardPubId string, position []string) error {
	return config.DB.Model(&models.ListPosition{}).Where("board_internal_id = (SELECT internal_id FROM boards WHERE public_id = ?)", boardPubId).
	Update("list_order", position).Error
}

func (r *ListRepositoryImpl) GetCardOrder(listPublicId string) ([]uuid.UUID, error){
	var position models.CardPosition

	err := config.DB.Joins("JOIN lists ON list.internal_id = card_positions.list_internal_id").
	Where("list.public_id = ?", listPublicId).Error

	return position.CardOrder, err
}

func (r *ListRepositoryImpl) FindByBoardId(boardId string) ([]models.List, error) {
	var list []models.List
	err := config.DB.Where("board_public_id = ?",boardId).Order("internal_id ASC").Find(&list).Error
	return list, err
}
func (r *ListRepositoryImpl) FindByPubId(pubId string) (*models.List, error){
	var list models.List
	err := config.DB.Where("public_id = ?", pubId).First(&list).Error

	return &list, err
}

func (r *ListRepositoryImpl) FindById(id uint) (*models.List, error) {
	var list models.List

	err := config.DB.First(&list, id).Error

	return &list, err
}