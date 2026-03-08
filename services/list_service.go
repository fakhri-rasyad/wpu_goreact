package services

import (
	"errors"
	"fmt"

	"github.com/fakhri-rasyad/wpu_goreact/config"
	"github.com/fakhri-rasyad/wpu_goreact/models"
	"github.com/fakhri-rasyad/wpu_goreact/models/types"
	"github.com/fakhri-rasyad/wpu_goreact/repositories"
	"github.com/fakhri-rasyad/wpu_goreact/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ListWithPosition struct {
	Positions []uuid.UUID
	Lists []models.List
}

type ListService interface{
	GetByBoardId(boardPublicId string) (*ListWithPosition, error)
	GetById(id uint) (*models.List, error)
	GetByPublicId(publicId string)(*models.List, error)
	Create(list *models.List) error
	Update(list *models.List) error
	Delete(id uint) error
	UpdatePosition(boardPubId string, posiitions []uuid.UUID) error
}

type ListServiceImpl struct {
	listRepository repositories.ListRepository
	boardRepository repositories.BoardRepository
	listPositionRepository repositories.ListPositionRepository
}



func NewListService(listRepo repositories.ListRepository, boardRepo repositories.BoardRepository, listPositionRepo repositories.ListPositionRepository) ListService {
	return &ListServiceImpl{
		listRepository: listRepo,
		boardRepository: boardRepo,
		listPositionRepository: listPositionRepo,
	}
}

func (s *ListServiceImpl) GetByBoardId(boardPubId string) (*ListWithPosition, error) {
	_, err := s.boardRepository.FindByPublicID(boardPubId)

	if err != nil {
		return nil, errors.New("Board not found")
	}

	positions,err :=  s.listPositionRepository.GetListOrder(boardPubId)
	if err != nil {
		return nil, errors.New("Failed to get list order" + err.Error())
	}

	if len(positions) == 0 {
		return nil, errors.New("List position not found")
	}

	lists, err := s.listRepository.FindByBoardId(boardPubId)
	if err != nil {
		return nil, errors.New("Failed to get list" + err.Error())
	}

	orderedList := utils.SortListByPosition(lists, positions)

	return &ListWithPosition{
		Positions: positions,
		Lists: orderedList,
	}, nil
}

func (s *ListServiceImpl) GetById(id uint) (*models.List, error) {
	return s.listRepository.FindById(id)
}

func (s *ListServiceImpl) GetByPublicId(publicId string)(*models.List, error) {
	return s.listRepository.FindByPubId(publicId)
}

func (s *ListServiceImpl) Create(list *models.List) error {
	board, err := s.boardRepository.FindByPublicID(list.BoardPublicID.String())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("Board not found")
		}
		return fmt.Errorf("Failed to get board: %w", err)
	}

	list.BoardInternalID = board.InternalID
	if list.PublicID == uuid.Nil {
		list.PublicID = uuid.New()
	}

	tx := config.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Create(list).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("Failed to create list: %w", err)
	}

	var position models.ListPosition
	res := tx.Where("board_internal_id = ?", board.InternalID).First(&position)

	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		position = models.ListPosition{
			PublicID:  uuid.New(),
			BoardID:   board.InternalID,
			ListOrder: types.UUIDArray{list.PublicID},
		}

		if err := tx.Create(&position).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("Failed to create list position: %w", err)
		}
	} else if res.Error != nil {
        // ✅ Was incorrectly calling tx.Create — should rollback and return the error
		tx.Rollback()
		return fmt.Errorf("Failed to get list position: %w", res.Error)
	} else {
		position.ListOrder = append(position.ListOrder, list.PublicID)

		if err := tx.Model(&position).Update("list_order", position.ListOrder).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("Failed to update list position: %w", err)
		}
	}

    // ✅ Commit once at the end, covering all branches
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("Transaction commit failed: %w", err)
	}

	return nil
}
func (s *ListServiceImpl)Update(list *models.List) error {
	return s.listRepository.Update(list)
}
func (s *ListServiceImpl)Delete(id uint) error {
	return s.listRepository.Delete(id)
}

func (s *ListServiceImpl) UpdatePosition(boardPubId string, posiitions []uuid.UUID) error {
	board, err := s.boardRepository.FindByPublicID(boardPubId)

	if err != nil {
		return errors.New("Board not found")
	}

	position, err := s.listPositionRepository.GetByBoard(board.PublicID.String())

	if err != nil {
		return errors.New("List position not found")
	}

	position.ListOrder = posiitions
	return s.listPositionRepository.UpdateListOrder(position)
}