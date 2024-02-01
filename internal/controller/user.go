package controller

import (
	"github.com/LucasAndFlores/user_api/internal/dto"
	"github.com/LucasAndFlores/user_api/internal/service"
	"github.com/gofiber/fiber/v2"
)

type Controller interface {
	HandleCreateUser(*fiber.Ctx) error
}

type UserController struct {
	service service.Service
}

func NewUserController(s service.Service) Controller {
	return &UserController{service: s}
}

func (c *UserController) HandleCreateUser(fi *fiber.Ctx) error {
	var userDTO dto.UserDTO

	err := fi.BodyParser(&userDTO)

	if err != nil {
		return fi.Status(fiber.StatusInternalServerError).JSON(map[string]string{"message": "internal server error"})
	}

	status, body := c.service.Create(fi.Context(), userDTO)

	return fi.Status(status).JSON(body)
}
