package repositories

import (
	"time"

	"github.com/fakhri-rasyad/wpu_goreact/config"
	"github.com/fakhri-rasyad/wpu_goreact/models"
)

type BoardRepository interface {
	Create(board *models.Board) error 
	Update(board *models.Board) error
	FindByPublicID(uuid string) (*models.Board, error)
	AddMember(boardId uint, userIDs []uint ) error
	RemoveMembers(boardId uint, userIDs []uint) error
	FindAllByUserPaginate(userPubId, filter, sort string, limit, offset int) ([]models.Board, int64, error)
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

func (r *BoardRepositoryImpl) AddMember(boardId uint, userIDs []uint ) error{
	if len(userIDs) == 0 {
		return nil
	}

	now := time.Now()
	var members []models.BoardMember
	for _, userId := range userIDs {
		members = append(members, models.BoardMember{
			BoardID: int64(boardId),
			UserID: int64(userId),
			JoinedAt: now,
		})
	}
	return config.DB.Create(&members).Error
}

func (r *BoardRepositoryImpl) RemoveMembers(boardId uint, userIDs []uint) error{
	if len(userIDs) == 0 {
		return nil
	}

	return config.DB.Where("board_internal_id = ? AND user_internal_id IN (?)", boardId, userIDs).Delete(&models.BoardMember{}).Error
}

func (r *BoardRepositoryImpl) FindAllByUserPaginate(userPubId, filter, sort string, limit, offset int) ([]models.Board, int64, error) {
	var board []models.Board
	var total int64

	query := config.DB.Model(&models.Board{}).Where("board_public_id = ? OR user_internal_id = ? IN (" + 
	"SELECT board_members.board_internal_id FROM board_members" +
	"SELECT board_members.user_internal_id FROM board_members" +
	"WHERE users.public_id = ?", userPubId, userPubId)

	if filter != "" {
		query = query.Where("title ILIKE ?", "%" + filter + "%")
	}


	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	
	if sort != "" {
		query = query.Order(sort)
	} else {
		query = query.Order("created_at DESC")
	}

	if err := query.Limit(limit).Offset(offset).Find(&board).Error; err != nil {
		return nil, 0, err
	}

	return board, total, nil
}