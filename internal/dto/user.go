package dto

import (
	"time"

	"github.com/LucasAndFlores/user_api/internal/model"
	"github.com/google/uuid"
)

type UserDTO struct {
	Name        string `json:"name" validate:"required,min=2"`
	Email       string `json:"email" validate:"email,required,min=2"`
	ExternalId  string `json:"id" validate:"uuid,required"`
	DateOfBirth string `json:"date_of_birth" validate:"required"`
}

func (d *UserDTO) ConvertToUserDTO(u *model.User) {
	d.Name = u.Name
	d.Email = u.Email
	d.ExternalId = u.ExternalId.String()
	d.DateOfBirth = u.DateOfBirth.String()
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
