package service

import (
	"context"

	"github.com/LucasAndFlores/user_api/internal/dto"
	"github.com/LucasAndFlores/user_api/internal/repository"
	"github.com/gofiber/fiber/v2"
)

type responseBody = map[string]interface{}

type Service interface {
	Create(context.Context, dto.UserDTO) (int, responseBody)
}

type UserService struct {
	repo repository.Repository
}

func NewUserService(r repository.Repository) Service {
	return &UserService{repo: r}
}

func (s *UserService) Create(ctx context.Context, user dto.UserDTO) (int, responseBody) {
	exists, err := s.repo.CheckIfUserExist(ctx, user)

	if err != nil {
		return fiber.StatusInternalServerError, responseBody{"message": "internal server error"}
	}

	if exists {
		return fiber.StatusConflict, responseBody{"message": "user already exists"}
	}

	userModel, err := user.ConvertToUserModel()

	if err != nil {
		return fiber.StatusInternalServerError, responseBody{"message": "internal server error"}
	}

	err = s.repo.Insert(ctx, userModel)

	if err != nil {
		return fiber.StatusInternalServerError, responseBody{"message": "internal server error"}
	}

	return fiber.StatusCreated, responseBody{"message": "user successfully created"}

}
