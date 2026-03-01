package seed

import (
	"log"

	"github.com/fakhri-rasyad/wpu_goreact/config"
	"github.com/fakhri-rasyad/wpu_goreact/models"
	"github.com/fakhri-rasyad/wpu_goreact/utils"
	"github.com/google/uuid"
)

func SeedAdmin() {
	pass, _ := utils.HashPassword("admin123")

	admin := models.User{
		Name: "Super admin",
		Email: "admin@example.com",
		Password: pass,
		Role: "admin",
		PublicID: uuid.New(),
	}

	if err := config.DB.FirstOrCreate(&admin, models.User{Email: admin.Email}).Error ; err != nil {
		log.Println("Failed to seed admin", err)
	} else {
		log.Println("Admin user seeded")
	}
}