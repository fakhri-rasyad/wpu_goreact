package repositories

import (
	"github.com/fakhri-rasyad/wpu_goreact/config"
	"github.com/fakhri-rasyad/wpu_goreact/models"
)

type BoardMemberRepository interface {
	GetMembers(boardPublicId string) ([]models.User, error)
}

type BoardMemberRepostioryImpl struct {
	
}

func NewBoardMemberRepository() BoardMemberRepository {
	return &BoardMemberRepostioryImpl{}
}

func (repo *BoardMemberRepostioryImpl) GetMembers(boardPublicId string) ([]models.User, error) {
	var users []models.User 
	err := config.DB.Joins("JOIN board_members ON board_members.user_internal_id = users.internal_id").
	Joins("JOIN board_members ON board_members.board_internal_id = boards.internal_id").
	Where("boards.public_id = ?", boardPublicId).
	Find(&users).Error

	return users, err
}