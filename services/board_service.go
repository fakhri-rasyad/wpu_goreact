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
	RemoveMembers(boardPubId string, memberIDs []string) error
	GetAllUserPaginatreBy (userPubId, filter, sort string, limit, offset int) ([]models.Board, int64, error)
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

func (s *BoardServiceImpl) RemoveMembers(boardPubId string, memberIDs []string) error {
	board, err := s.repo.FindByPublicID(boardPubId)

	if err != nil {
		return errors.New("Board tidak ditemukan")
	}

	var memberIntIDs []uint
	for _, memberPubId := range memberIDs {
		user, err := s.userRepo.FindByPublicID(memberPubId)
		if err != nil {
			return errors.New("User tidak ditemukan" + memberPubId)
		}

		memberIntIDs = append(memberIntIDs, uint(user.InternalID))
	}

	existingMembers , err := s.boardmemberRepo.GetMembers(boardPubId)
	if err != nil {
		return errors.New("Board tidak memiliki id")
	}

	memberMap := make(map[uint]bool)
	for _, member := range existingMembers{
		memberMap[uint(member.InternalID)] = true
	}

	var memberToRemove []uint
	for _, memIntId := range memberIntIDs{
		if memberMap[memIntId] {
			memberToRemove = append(memberToRemove, memIntId)
		}
	}

	return s.repo.RemoveMembers(uint(board.InternalID), memberToRemove)

}

func (s *BoardServiceImpl) GetAllUserPaginatreBy (userPubId, filter, sort string, limit, offset int) ([]models.Board, int64, error) {
	return s.repo.FindAllByUserPaginate(userPubId, filter, sort, limit, offset)
}