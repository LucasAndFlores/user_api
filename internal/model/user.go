package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id          int       `gorm:"type:int;primary_key"`
	Name        string    `gorm:"not null"`
	Email       string    `gorm:"unique;not null"`
	ExternalId  uuid.UUID `gorm:"column:external_id;type:uuid;unique;not null"`
	DateOfBirth time.Time `gorm:"column:date_of_birth;type:timestamp with time zone"`
}
