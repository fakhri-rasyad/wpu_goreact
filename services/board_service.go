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
	AddMembers(boardPubId string, userPubIds []string) error
}

type BoardServiceImpl struct {
	repo repositories.BoardRepository
	userRepo repositories.UserRepository
	boardmemberRepo repositories.BoardMemberRepository
}

func NewBoardService(repo repositories.BoardRepository, userRepo repositories.UserRepository, boardMemberRepo repositories.BoardMemberRepository) BoardService {
	return &BoardServiceImpl{repo: repo, userRepo: userRepo, boardmemberRepo: boardMemberRepo}
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

func (s *BoardServiceImpl) AddMembers(boardPubId string, userPubIds []string) error {
	board, err := s.repo.FindByPublicID(boardPubId)
	if err != nil {
		return errors.New("Board not found")
	}

	var userInteralId []uint
	for _, pubId := range userPubIds {
		user, err := s.userRepo.FindByPublicID(pubId)
		if err != nil {
			return errors.New("User not found" + pubId)
		}

		userInteralId = append(userInteralId, uint(user.InternalID))
	}

	existingMembers, err := s.boardmemberRepo.GetMembers(string(board.PublicID.String()))

	if err != nil {
		return err
	}

	memberMap := make(map[uint]bool)
	for _, member := range existingMembers {
		memberMap[uint(member.InternalID)] = true
	}

	var newMembersId []uint
	for _, userIds := range userInteralId {
		if !memberMap[userIds] {
			newMembersId = append(newMembersId, userIds)
		}
	}

	if len(newMembersId) == 0{
		return nil
	}

	return s.repo.AddMember(uint(board.InternalID), newMembersId)
}