package model

import (
	"booking/shared/constants"
	"errors"
	"time"

	"github.com/google/uuid"
)

type Booking struct {
	ID         uuid.UUID `json:"id" gorm:"type:char(36);primary_key"`
	UserID     uuid.UUID `json:"user_id" gorm:"type:char(36);not null"`
	SpaceID    uuid.UUID `json:"space_id" gorm:"type:char(36);not null"`
	StartDate  time.Time `json:"start_date" gorm:"not null"`
	EndDate    time.Time `json:"end_date" gorm:"not null"`
	TotalPrice float64   `json:"total_price" gorm:"not null"`
	Status     string    `json:"status" gorm:"type:varchar(20);not null;default:'pending'"`
	CreatedAt  time.Time `json:"created_at" gorm:"not null"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"not null"`
}

type CreateBookingInput struct {
	UserID    uuid.UUID `json:"user_id"`
	SpaceID   uuid.UUID `json:"space_id"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
}

type BookingResponse struct {
	ID         uuid.UUID `json:"id"`
	UserID     uuid.UUID `json:"user_id"`
	SpaceID    uuid.UUID `json:"space_id"`
	StartDate  string    `json:"start_date"`
	EndDate    string    `json:"end_date"`
	TotalPrice float64   `json:"total_price"`
	Status     string    `json:"status"`
	CreatedAt  string    `json:"created_at"`
	UpdatedAt  string    `json:"updated_at"`
}

func NewBooking(input CreateBookingInput, totalPrice float64) (*Booking, error) {
	now := time.Now()
	return &Booking{
		ID:         uuid.New(),
		UserID:     input.UserID,
		SpaceID:    input.SpaceID,
		StartDate:  input.StartDate,
		EndDate:    input.EndDate,
		TotalPrice: totalPrice,
		Status:     "pending",
		CreatedAt:  now,
		UpdatedAt:  now,
	}, nil
}

func (b *Booking) ToResponse() BookingResponse {
	return BookingResponse{
		ID:         b.ID,
		UserID:     b.UserID,
		SpaceID:    b.SpaceID,
		StartDate:  b.StartDate.Format("2006-01-02 15:04:05"),
		EndDate:    b.EndDate.Format("2006-01-02 15:04:05"),
		TotalPrice: b.TotalPrice,
		Status:     b.Status,
		CreatedAt:  b.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:  b.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

func (b *Booking) UpdateStatus(status constants.BookingStatus) error {
	switch status {
	case constants.BookingStatusPending, constants.BookingStatusPaid, constants.BookingStatusCancelled:
		b.Status = string(status)
		return nil
	default:
		return errors.New("invalid booking status")
	}
}
