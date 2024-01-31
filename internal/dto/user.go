package dto

import (
	"time"

	"github.com/LucasAndFlores/user_api/internal/model"
	"github.com/google/uuid"
)

type UserDTO struct {
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	ExternalId  uuid.UUID `json:"id"`
	DateOfBirth time.Time `json:"date_of_birth"`
}

func (d UserDTO) ConvertToUserModel() model.User {
	return model.User{
		Name:        d.Name,
		Email:       d.Email,
		ExternalId:  d.ExternalId,
		DateOfBirth: d.DateOfBirth,
	}
}
