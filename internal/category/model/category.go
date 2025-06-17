package model

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Category struct {
	ID          uuid.UUID `json:"id" gorm:"type:char(36);primary_key;default:(UUID())"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description string    `json:"description"`
	CreatedBy   uuid.UUID `json:"created_by" gorm:"type:char(36)"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateCategoryInput struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

func generateSlug(name string) string {
	// Mengubah nama menjadi huruf kecil
	slug := strings.ToLower(name)
	// Mengganti spasi dengan tanda strip
	slug = strings.ReplaceAll(slug, " ", "-")
	return slug
}

func NewCategory(input CreateCategoryInput, userID uuid.UUID) (*Category, error) {
	if input.Name == "" {
		return nil, errors.New("category name is required")
	}

	return &Category{
		ID:          uuid.New(),
		Name:        input.Name,
		Slug:        generateSlug(input.Name),
		Description: input.Description,
		CreatedBy:   userID,
	}, nil
}

func (c *Category) Update(input CreateCategoryInput) error {
	if input.Name == "" {
		return errors.New("category name is required")
	}

	c.Name = input.Name
	c.Slug = generateSlug(input.Name)
	c.Description = input.Description
	return nil
}
