package model

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type SpaceFacility struct {
	SpaceID    uuid.UUID `json:"space_id" gorm:"type:char(36);primaryKey;not null;foreignKey:SpaceID;references:spaces(id);onDelete:CASCADE"`
	FacilityID uuid.UUID `json:"facility_id" gorm:"type:char(36);primaryKey;not null;foreignKey:FacilityID;references:facilities(id);onDelete:CASCADE"`
	CreatedAt  time.Time `json:"created_at" gorm:"autoCreateTime;not null"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"autoUpdateTime;not null"`
}

type SpaceFacilityInput struct {
	SpaceID    uuid.UUID `json:"space_id" validate:"required"`
	FacilityID uuid.UUID `json:"facility_id" validate:"required"`
}

func NewSpaceFacility(input SpaceFacilityInput) (*SpaceFacility, error) {
	if input.SpaceID == uuid.Nil {
		return nil, errors.New("space id is required")
	}
	if input.FacilityID == uuid.Nil {
		return nil, errors.New("facility id is required")
	}
	spaceFacility := &SpaceFacility{
		SpaceID:    input.SpaceID,
		FacilityID: input.FacilityID,
	}
	return spaceFacility, nil
}
