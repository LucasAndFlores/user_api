package dto

import (
	"time"

	"github.com/LucasAndFlores/user_api/internal/model"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type UserDTO struct {
	Name        string `json:"name" validate:"required,min=2"`
	Email       string `json:"email" validate:"email,required,min=2"`
	ExternalId  string `json:"id" validate:"uuid,required"`
	DateOfBirth string `json:"date_of_birth" validate:"required"`
}

func (d *UserDTO) ConvertToUserModel() (model.User, error) {
	uuid, err := uuid.Parse(d.ExternalId)

	if err != nil {
		return model.User{}, err
	}

	date, err := time.Parse(time.RFC3339, d.DateOfBirth)

	if err != nil {
		return model.User{}, err
	}

	return model.User{
		Name:        d.Name,
		Email:       d.Email,
		ExternalId:  uuid,
		DateOfBirth: date,
	}, nil
}

var Validator = validator.New()

type RequestBodyError struct {
	Field string
	Tag   string
	Value string
}

func ValidateUserRequestBody(fi *fiber.Ctx) error {
	var errors []*RequestBodyError

	var user UserDTO

	fi.BodyParser(&user)

	_, err := time.Parse(time.RFC3339, user.DateOfBirth)

	if err != nil {
		var el RequestBodyError
		el.Field = "DateOfBirth"
		el.Tag = "required"
		el.Value = "invalid date format"
		errors = append(errors, &el)
	}

	err = Validator.Struct(user)

	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var el RequestBodyError
			el.Field = err.Field()
			el.Tag = err.Tag()
			el.Value = err.Param()
			errors = append(errors, &el)
		}

	}

	if len(errors) != 0 {
		return fi.Status(fiber.ErrUnprocessableEntity.Code).JSON(errors)
	}

	return fi.Next()
}
