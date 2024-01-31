package repository

import (
	"context"
	"errors"

	"github.com/LucasAndFlores/user_api/internal/dto"
	"github.com/LucasAndFlores/user_api/internal/model"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

type Repository interface {
	Insert(context.Context, model.User) error
	CheckIfUserExist(context.Context, dto.UserDTO) (bool, error)
}

func NewUserRepository(d *gorm.DB) Repository {
	return &UserRepository{
		db: d,
	}
}

func (r *UserRepository) Insert(ctx context.Context, user model.User) error {
	result := r.db.Create(&user)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *UserRepository) CheckIfUserExist(ctx context.Context, user dto.UserDTO) (bool, error) {
	var foundUser model.User

	err := r.db.Select("email", "external_id").Where("email = ?", user.Email).Or("external_id = ?", user.ExternalId).Take(&foundUser).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}

	if err != nil {
		return true, err
	}

	return true, nil
}
