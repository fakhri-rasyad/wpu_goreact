package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserResponse struct {
	PublicID  uuid.UUID      `json:"public_id" `
	Name      string         `json:"name" `
	Email     string         `json:"email" `
	Role      string         `json:"role" `
	CreatedAt time.Time      `json:"created_at" `
	UpdatedAt time.Time      `json:"updated_at" `
	DeletedAt gorm.DeletedAt `json:"-"`
}