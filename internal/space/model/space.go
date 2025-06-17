package model

import (
	"errors"

	"github.com/google/uuid"
)

type Space struct {
	ID            uuid.UUID `json:"id" gorm:"type:char(36);primary_key;default:(UUID())"`
	CategoryID    uuid.UUID `json:"category_id" gorm:"type:char(36)"`
	Name          string    `json:"name" gorm:"size:150"`
	Description   string    `json:"description" gorm:"type:text"`
	PricePerNight float64   `json:"price_per_night" gorm:"type:decimal(12,2)"`
	IsActive      bool      `json:"is_active" gorm:"default:true"`
}

type CreateSpaceInput struct {
	CategoryID    uuid.UUID `json:"category_id" binding:"required"`
	Name          string    `json:"name" binding:"required"`
	Description   string    `json:"description" binding:"required"`
	PricePerNight float64   `json:"price_per_night" binding:"required"`
}

func NewSpace(input CreateSpaceInput, categoryID uuid.UUID) (*Space, error) {
	if input.CategoryID == uuid.Nil {
		return nil, errors.New("category ID is required")
	}
	if input.Name == "" {
		return nil, errors.New("space name is required")
	}
	if input.Description == "" {
		return nil, errors.New("description is required")
	}
	if input.PricePerNight <= 0 {
		return nil, errors.New("price must be greater than zero")
	}
	space := &Space{
		ID:            uuid.New(),
		CategoryID:    categoryID,
		Name:          input.Name,
		Description:   input.Description,
		PricePerNight: input.PricePerNight,
		IsActive:      true,
	}

	return space, nil
}
