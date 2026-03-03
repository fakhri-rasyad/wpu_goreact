package routes

import (
	"log"

	"github.com/fakhri-rasyad/wpu_goreact/config"
	"github.com/fakhri-rasyad/wpu_goreact/controller"
	"github.com/fakhri-rasyad/wpu_goreact/utils"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/joho/godotenv"
)

func Setup(app *fiber.App, uc *controller.UserController){
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	app.Post("v1/auth/register", uc.Register)
	app.Post("v1/auth/login", uc.Login)

	// JWT Protected routes
	api := app.Group("api/v1", jwtware.New(jwtware.Config{
		SigningKey: []byte(config.APPConfig.JWTSecret),
		ContextKey: "user",
		ErrorHandler: func (c *fiber.Ctx, err error) error {
			return utils.UnauthorizedRequest(c, "Unauthorized", nil, err.Error())
		},
	}))

	userGroup := api.Group("/users")
	userGroup.Get("/page", uc.GetUserPagination)
	userGroup.Get("/:id", uc.GetUser)
	userGroup.Put("/:id", uc.UpdateUser)
	userGroup.Delete("/:id", uc.DeleteUser)
}