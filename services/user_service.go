package services

import (
	"errors"

	"github.com/fakhri-rasyad/wpu_goreact/models"
	"github.com/fakhri-rasyad/wpu_goreact/repositories"
	"github.com/fakhri-rasyad/wpu_goreact/utils"
	"github.com/google/uuid"
)

type UserService interface {
	Register(user *models.User) error
}

type userServiceImpl struct {
	repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) UserService {
	return &userServiceImpl{repo}
}

func (s *userServiceImpl) Register(user *models.User) error {
	//Ngecek email yang sudah terdaftar
	// Hashing password
	// Simpan user
	existingUser, _ := s.repo.FindByEmail(user.Email)
	if existingUser.InternalID != 0 {
		return errors.New("Email already registred")
	}

	hashed , _ := utils.HashPassword(user.Password)

	user.Password = hashed
	user.Role = "user"
	user.PublicID = uuid.New()

	return s.repo.Create(user)
}