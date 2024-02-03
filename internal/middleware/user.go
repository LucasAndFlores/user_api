package middleware

import (
	"time"

	"github.com/LucasAndFlores/user_api/internal/dto"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var Validator = validator.New()

type RequestBodyError struct {
	Field string
	Tag   string
	Value string
}

func ValidateUserRequestBody(fi *fiber.Ctx) error {
	var errors []*RequestBodyError

	var user dto.UserDTO

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

