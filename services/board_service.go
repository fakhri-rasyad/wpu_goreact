package services

import (
	"errors"

	"github.com/fakhri-rasyad/wpu_goreact/models"
	"github.com/fakhri-rasyad/wpu_goreact/repositories"
	"github.com/google/uuid"
)

type BoardService interface{
	Create(board *models.Board) error
	Update(board *models.Board) error
	GetByID(uuid string) (*models.Board, error)
}

type BoardServiceImpl struct {
	repo repositories.BoardRepository
	userRepo repositories.UserRepository
}

func NewBoardService(repo repositories.BoardRepository, userRepo repositories.UserRepository) BoardService {
	return &BoardServiceImpl{repo: repo, userRepo: userRepo}
}

func (s *BoardServiceImpl) Create(board *models.Board) error {
	user, err := s.userRepo.FindByPublicID(board.OwnerPublicID.String())
	if err != nil {
		return errors.New("Owner not found")
	}
	board.PublicID = uuid.New()
	board.OwnerID = user.InternalID
	return s.repo.Create(board)
}

func (s *BoardServiceImpl) Update(board *models.Board) error{
	return s.repo.Update(board)
}

func (s *BoardServiceImpl) GetByID(uuid string) (*models.Board, error){
	return s.repo.FindByPublicID(uuid)
}