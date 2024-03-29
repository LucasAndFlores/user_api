package service

import (
	"context"

	"github.com/LucasAndFlores/user_api/internal/dto"
	"github.com/LucasAndFlores/user_api/internal/repository"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type responseBody = map[string]interface{}

const INTERNAL_SERVER_ERROR_MESSAGE = "internal server error"

type Service interface {
	Create(context.Context, dto.UserDTO) (int, responseBody)
	FindUserByExternalId(context.Context, uuid.UUID) (int, responseBody)
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
		return fiber.StatusInternalServerError, responseBody{"message": INTERNAL_SERVER_ERROR_MESSAGE}
	}

	if exists {
		return fiber.StatusConflict, responseBody{"message": "user already exists"}
	}

	userModel, err := user.ConvertToUserModel()

	if err != nil {
		return fiber.StatusInternalServerError, responseBody{"message": INTERNAL_SERVER_ERROR_MESSAGE}
	}

	err = s.repo.Insert(ctx, &userModel)

	if err != nil {
		return fiber.StatusInternalServerError, responseBody{"message": INTERNAL_SERVER_ERROR_MESSAGE}
	}

	return fiber.StatusCreated, responseBody{"message": "user successfully created"}

}

func (s *UserService) FindUserByExternalId(ctx context.Context, externalId uuid.UUID) (int, responseBody) {
	found, err := s.repo.FindByExternalId(ctx, externalId)

	if err != nil {
		return fiber.StatusInternalServerError, responseBody{"message": INTERNAL_SERVER_ERROR_MESSAGE}
	}

	if found.Id == 0 {
		return fiber.StatusNotFound, responseBody{"message": "user not found"}
	}

	var userDTO dto.UserDTO

	userDTO.ConvertToUserDTO(found)

	return fiber.StatusOK, responseBody{"user": userDTO}
}
