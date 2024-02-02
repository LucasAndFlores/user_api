package routes

import (
	"github.com/LucasAndFlores/user_api/internal/controller"
	"github.com/LucasAndFlores/user_api/internal/dto"
	"github.com/LucasAndFlores/user_api/internal/repository"
	"github.com/LucasAndFlores/user_api/internal/service"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupUserRoutes(api fiber.Router, db *gorm.DB) {
	repo := repository.NewUserRepository(db)
	service := service.NewUserService(repo)

	userController := controller.NewUserController(service)

	api.Post("/save", dto.ValidateUserRequestBody, userController.HandleCreateUser)
	api.Get("/:id", userController.HandleFindUserByExternalId)
}
