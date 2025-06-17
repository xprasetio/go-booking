package model

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Facility struct {
	ID        uuid.UUID      `json:"id" gorm:"type:char(36);primary_key;default:(UUID())"`
	Name      string         `json:"name" gorm:"unique;not null"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

type CreateFacilityInput struct {
	Name string `json:"name" validate:"required"`
}

func NewFacility(input CreateFacilityInput) (*Facility, error) {
	if input.Name == "" {
		return nil, errors.New("Facility name is required")
	}
	facility := &Facility{
		ID:   uuid.New(),
		Name: input.Name,
	}

	return facility, nil
}
