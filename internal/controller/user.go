package controller

import (
	"github.com/LucasAndFlores/user_api/internal/dto"
	"github.com/LucasAndFlores/user_api/internal/service"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type Controller interface {
	HandleCreateUser(*fiber.Ctx) error
	HandleFindUserByExternalId(*fiber.Ctx) error
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

func (c *UserController) HandleFindUserByExternalId(fi *fiber.Ctx) error {
	id := fi.Params("id")

	uuid, err := uuid.Parse(id)

	if err != nil {
		return fi.Status(fiber.StatusBadRequest).JSON(map[string]string{"message": "unable to parse the id"})
	}

	status, body := c.service.FindUserByExternalId(fi.Context(), uuid)

	return fi.Status(status).JSON(body)
}
