package seed

import (
	"log"

	"github.com/fakhri-rasyad/wpu_goreact/config"
	"github.com/fakhri-rasyad/wpu_goreact/models"
	"github.com/fakhri-rasyad/wpu_goreact/utils"
)

func SeedAdmin() {
	pass, _ := utils.HashPassword("admin123")

	admin := models.User{
		Name: "Super admin",
		Email: "admin@example.com",
		Password: pass,
		Role: "admin",
	}

	if err := config.DB.FirstOrCreate(&admin, models.User{Email: admin.Email}).Error ; err != nil {
		log.Println("Failed to seed admin", err)
	} else {
		log.Panicln("Admin user seeded")
	}
}