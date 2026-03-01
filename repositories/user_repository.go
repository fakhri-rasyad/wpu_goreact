package repositories

import (
	"github.com/fakhri-rasyad/wpu_goreact/config"
	"github.com/fakhri-rasyad/wpu_goreact/models"
)

type UserRepository interface {
	Create(user *models.User) error
	FindByEmail(email string) (*models.User, error)
}

type userRepositoryImpl struct {}

func NewUserRepostiry() UserRepository {
	return &userRepositoryImpl{}
}

func (r *userRepositoryImpl) Create (user *models.User) error {
	return config.DB.Create(user).Error
}

func (r *userRepositoryImpl) FindByEmail(email string) (*models.User, error){
	var user models.User
	err := config.DB.Where("email = ?", email).First(&user).Error
	return &user, err
}