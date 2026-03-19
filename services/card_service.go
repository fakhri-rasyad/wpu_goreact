package services

import (
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/fakhri-rasyad/wpu_goreact/config"
	"github.com/fakhri-rasyad/wpu_goreact/models"
	"github.com/fakhri-rasyad/wpu_goreact/models/types"
	"github.com/fakhri-rasyad/wpu_goreact/repositories"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CardServices interface {
	Create(card *models.Card, listPublicID string) error
	Update(card *models.Card, listPublicID string) error
	Delete(id uint) error

	GetByListID(listPubID string) ([]models.Card, error)
	GetByID(id uint)(*models.Card, error)
	GetByPublicID(cardPubID string)(*models.Card, error)
}

type CardServicesImpl struct {
	cardRepo repositories.CardRepository
	listRepo repositories.ListRepository
	userRepo repositories.UserRepository
}

func NewCardService(
	cardRepo repositories.CardRepository,
	listRepo repositories.ListRepository,
	userRepo repositories.UserRepository,
	) CardServices {
	return &CardServicesImpl{cardRepo: cardRepo, listRepo: listRepo, userRepo: userRepo}
}

func(s *CardServicesImpl) Create(card *models.Card, listPublicID string) error {
	list, err := s.listRepo.FindByPubId(listPublicID)

	if err != nil {
		return fmt.Errorf("List not found: %w", err)
	}

	card.ListID = list.InternalID

	if card.PublicID ==  uuid.Nil {
		card.PublicID = uuid.New()
	}

	card.CreatedAt = time.Now()

	tx := config.DB.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	if err := tx.Create(card).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("Failed to create card: %w", err)
	}

	var position models.CardPosition
	err = tx.Model(&models.CardPosition{}).
	    Where("list_internal_id = ?", list.InternalID).
	    First(&position).Error
	
	if errors.Is(err, gorm.ErrRecordNotFound) {
	    position = models.CardPosition{
	        PublicID:  uuid.New(),
	        ListID:    list.InternalID,
	        CardOrder: types.UUIDArray{card.PublicID},
	    }
	    if err := tx.Create(&position).Error; err != nil {
	        tx.Rollback()
	        return fmt.Errorf("Failed to create card position: %w", err)
	    }
	} else if err != nil {
	    tx.Rollback()
	    return fmt.Errorf("Failed to get card position: %w", err)
	} else {
	    position.CardOrder = append(position.CardOrder, card.PublicID)
	    if err := tx.Model(&models.CardPosition{}).
	        Where("internal_id = ?", position.InternalID).
	        Update("card_order", position.CardOrder).Error; err != nil {
	        tx.Rollback()
	        return fmt.Errorf("Failed to update card position: %w", err)
	    }
	}

	if err:= tx.Commit().Error; err != nil{
		return fmt.Errorf("Transaction commit failed: %w", err)
	}

	return nil
}

func(s *CardServicesImpl) Update(card *models.Card, listPublicID string) error {
	existingCard, err := s.cardRepo.FindByPublicId(card.PublicID.String())
	if err != nil {
		return fmt.Errorf("Card not found: %w", err)
	}

	newList, err := s.listRepo.FindByPubId(listPublicID)
	if err != nil {
		return fmt.Errorf("List not found: %w", err)
	}

	tx := config.DB.Begin()

	defer func(){
		if r := recover(); r != nil{
			tx.Rollback()
			panic(r)
		}
	}()

	// jika pindah list, hapus dari posisi lama

	if existingCard.ListID != newList.InternalID {
		// hapus dari list lama
		var oldPos models.CardPosition
		if err := tx.Where("list_internal_id = ?", existingCard.ListID).First(&oldPos).Error; err !=nil {
			filtered := make(types.UUIDArray, 0, len(oldPos.CardOrder))
			for _, id := range oldPos.CardOrder{
				if id != existingCard.PublicID{
					filtered = append(filtered, id)
				}
			}

			if err := tx.Model(&models.CardPosition{}).Where("internal_id = ?", oldPos.InternalID).
			Update("card_order", types.UUIDArray(filtered)).Error; err != nil {
				tx.Rollback()
				return fmt.Errorf("Failed to update old card position: %w", err)
			}
		} else if !errors.Is(err, gorm.ErrRecordNotFound){
			tx.Rollback()
			return fmt.Errorf("Failed to get old card position: %w", err)
		}
	}

	// tambah ke list baru
	var newPost models.CardPosition

	res := tx.Where("list_internal_id = ?", newList.InternalID).First(&newPost)

	if errors.Is(res.Error, gorm.ErrRecordNotFound){
		newPost = models.CardPosition{
			PublicID: uuid.New(),
			ListID: newList.InternalID,
			CardOrder: types.UUIDArray{existingCard.PublicID},
		}

		if err := tx.Create(&newPost).Error ; err != nil {
			tx.Rollback()
			return fmt.Errorf("Gagal membuat card baru: %w", err)
		}
	} else if res.Error == nil {
		updateOrder := append(newPost.CardOrder, existingCard.PublicID)
		if err := tx.Model(&models.CardPosition{}).Where("internal_id", newPost.InternalID).Update("card_order", types.UUIDArray(updateOrder)).Error ; err != nil {
			tx.Rollback()
			return fmt.Errorf("Failed to update old card position: %w", err)
		}
	} else {
		tx.Rollback()
		return fmt.Errorf("Failed to get new card position: %w", res.Error)
	}

	// Update data card

	card.InternalID = existingCard.InternalID
	card.PublicID = existingCard.PublicID
	card.ListID = existingCard.ListID

	if err := tx.Save(card).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("Failed to update card: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("Failed to commit changes: %w", err)
	}

	return nil
}
func(s *CardServicesImpl) Delete(id uint) error {
	return s.cardRepo.Delete(id)
}

func(s *CardServicesImpl) GetByListID(listPubID string) ([]models.Card, error) {
	list, err := s.listRepo.FindByPubId(listPubID)

	if err != nil {
		return nil, fmt.Errorf("List not found")
	}

	position, err := s.cardRepo.FindCardPositionByListID(list.InternalID)
	if err != nil {
		return nil, fmt.Errorf("Failed to get card position: %w", err)
	}

	cards, err := s.cardRepo.FindByListID(listPubID)
	if err != nil {
		return nil, fmt.Errorf("Failed to get cards using list pub id: %w", err)
	}

	if position != nil && len(position.CardOrder) > 0 {
		cards = sortCardByPosition(cards, position.CardOrder)
	}

	return cards, nil
}

func sortCardByPosition(cards []models.Card, order []uuid.UUID) []models.Card{
	orderMap := make(map[uuid.UUID]int)
	for i, uid := range order {
		orderMap[uid] = i
	}

	defaultIndex := len(order)

	sort.SliceStable(cards, func (i, j int) bool {
		idxI, okI := orderMap[cards[i].PublicID]
		if !okI {
			idxI = defaultIndex
		}

		idxJ, okJ := orderMap[cards[j].PublicID]
		if !okJ {
			idxJ = defaultIndex
		}

		if idxI == idxJ {
			return cards[i].CreatedAt.Before(cards[j].CreatedAt)
		}

		return idxI < idxJ
	})

	return cards
}

func(s *CardServicesImpl) GetByID(id uint)(*models.Card, error) {
	return s.cardRepo.FindByID(id)
}
func(s *CardServicesImpl) GetByPublicID(cardPubID string)(*models.Card, error) {
	return s.cardRepo.FindByPublicId(cardPubID)
}