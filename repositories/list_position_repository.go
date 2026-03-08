package repositories

import (
	"github.com/fakhri-rasyad/wpu_goreact/config"
	"github.com/fakhri-rasyad/wpu_goreact/models"
	"github.com/google/uuid"
)

type ListPositionRepository interface{
	GetByBoard(boardPubId string) (*models.ListPosition, error)
	CreateOrUpdate(boardPubId string, position []uuid.UUID) error
	GetListOrder(boardPubId string) ([]uuid.UUID, error)
	UpdateListOrder(position *models.ListPosition) error
}

type ListPositionRepositoryImpl struct {
}

func NewListPositionRepository() ListPositionRepository {
	return &ListPositionRepositoryImpl{}
}

func (r *ListPositionRepositoryImpl) GetByBoard(boardPubId string) (*models.ListPosition, error){
	var positions models.ListPosition

	err:= config.DB.Joins("JOIN boards ON boards.internal_id = list_positions.board_internal_id").
	Where("boards.public_id = ?", boardPubId).First(&positions).Error

	return &positions, err
}

func (r *ListPositionRepositoryImpl) CreateOrUpdate(boardPubId string, listOrder []uuid.UUID) error {
	return config.DB.Exec(`INSERT INTO list_positions (board_internal_id, list_order)
	SELECT internal_id, ? FROM boards Where public_id = ?
	ON CONFLICT (board_internal_id)
	DO UPDATE SET list_order = EXCLUDE.list_order 
	`, listOrder, boardPubId).Error
}

func (r *ListPositionRepositoryImpl)GetListOrder(boardPubId string) ([]uuid.UUID, error){
	position, err := r.GetByBoard(boardPubId)

	if err != nil {
		return nil, err
	}

	return position.ListOrder, nil


}
func (r *ListPositionRepositoryImpl) UpdateListOrder(position *models.ListPosition) error {
	return config.DB.Model(position).
	Where("internal_id = ?", position.InternalID).
	Update("list_order", position.ListOrder).Error
}