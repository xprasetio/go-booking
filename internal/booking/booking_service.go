package booking

import (
	"booking/internal/booking/model"
	"booking/internal/space"
	"booking/internal/user"
	"booking/pkg/logger"
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type BookingServiceInterface interface {
	Create(ctx context.Context, input model.CreateBookingInput) (*model.Booking, error)
	GetByID(ctx context.Context, id string) (*model.Booking, error)
	GetAll(ctx context.Context, userID uuid.UUID) ([]model.Booking, error)
	Cancel(ctx context.Context, bookingID string, userID uuid.UUID) error
}

type BookingService struct {
	db           *gorm.DB
	logger       logger.Logger
	userService  user.UserServiceInterface
	spaceService space.SpaceServiceInterface
}

func NewBookingService(db *gorm.DB, logger logger.Logger, userService user.UserServiceInterface, spaceService space.SpaceServiceInterface) *BookingService {
	return &BookingService{
		db:           db,
		logger:       logger,
		userService:  userService,
		spaceService: spaceService,
	}
}

func (s *BookingService) Create(ctx context.Context, input model.CreateBookingInput) (*model.Booking, error) {
	// Validasi user exists
	if _, err := s.userService.GetUserByID(ctx, input.UserID.String()); err != nil {
		s.logger.WithFields(logrus.Fields{
			"user_id": input.UserID,
			"error":   err.Error(),
		}).Error(ctx, "failed to get user")
		return nil, errors.New("user not found")
	}

	// Validasi space exists
	space, err := s.spaceService.GetByID(input.SpaceID.String())
	if err != nil {
		s.logger.WithFields(logrus.Fields{
			"space_id": input.SpaceID,
			"error":    err.Error(),
		}).Error(ctx, "failed to get space")
		return nil, errors.New("space not found")
	}

	// Validasi space is active
	if !space.IsActive {
		return nil, errors.New("space is not active")
	}

	// Validasi tanggal tidak boleh lebih kecil dari hari ini
	today := time.Now().Truncate(24 * time.Hour)
	if input.StartDate.Before(today) {
		return nil, errors.New("start date must be today or later")
	}

	// Hitung durasi booking dalam hari
	startDate := time.Date(input.StartDate.Year(), input.StartDate.Month(), input.StartDate.Day(), 0, 0, 0, 0, time.Local)
	endDate := time.Date(input.EndDate.Year(), input.EndDate.Month(), input.EndDate.Day(), 0, 0, 0, 0, time.Local)
	duration := endDate.Sub(startDate).Hours() / 24

	if duration < 1 {
		return nil, errors.New("minimum booking duration is 1 day")
	}

	// Hitung total harga
	totalPrice := space.PricePerNight * float64(duration)

	// Cek apakah ada booking yang overlap
	var count int64
	err = s.db.Model(&model.Booking{}).
		Where("space_id = ? AND status != ? AND ((start_date <= ? AND end_date > ?) OR (start_date < ? AND end_date >= ?) OR (start_date >= ? AND start_date < ?))",
			input.SpaceID,
			"cancelled",
			input.EndDate, input.StartDate,
			input.EndDate, input.StartDate,
			input.StartDate, input.EndDate).
		Count(&count).Error

	if err != nil {
		s.logger.WithFields(logrus.Fields{
			"space_id": input.SpaceID,
			"error":    err.Error(),
		}).Error(ctx, "failed to check booking overlap")
		return nil, errors.New("failed to check booking availability")
	}

	if count > 0 {
		return nil, errors.New("space is already booked for the selected dates")
	}

	// Buat booking baru
	booking, err := model.NewBooking(input, totalPrice)
	if err != nil {
		s.logger.WithFields(logrus.Fields{
			"user_id":  input.UserID,
			"space_id": input.SpaceID,
			"error":    err.Error(),
		}).Error(ctx, "failed to create booking")
		return nil, err
	}

	// Simpan ke database
	if err := s.db.WithContext(ctx).Create(booking).Error; err != nil {
		s.logger.WithFields(logrus.Fields{
			"user_id":  input.UserID,
			"space_id": input.SpaceID,
			"error":    err.Error(),
		}).Error(ctx, "failed to save booking to database")
		return nil, errors.New("failed to create booking")
	}

	return booking, nil
}

func (s *BookingService) GetByID(ctx context.Context, id string) (*model.Booking, error) {
	var booking model.Booking
	if err := s.db.WithContext(ctx).First(&booking, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("booking not found")
		}
		return nil, err
	}
	return &booking, nil
}

func (s *BookingService) GetAll(ctx context.Context, userID uuid.UUID) ([]model.Booking, error) {
	var bookings []model.Booking
	if err := s.db.WithContext(ctx).Where("user_id = ?", userID).Find(&bookings).Error; err != nil {
		return nil, err
	}
	return bookings, nil
}

func (s *BookingService) Cancel(ctx context.Context, bookingID string, userID uuid.UUID) error {
	var booking model.Booking
	if err := s.db.WithContext(ctx).First(&booking, "id = ?", bookingID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("booking not found")
		}
		return err
	}

	// Check if user is authorized to cancel this booking
	if booking.UserID != userID {
		return errors.New("you are not authorized to cancel this booking")
	}

	// Check if booking can be cancelled
	if booking.Status == "cancelled" {
		return errors.New("booking is already cancelled")
	}

	// Update booking status
	booking.Status = "cancelled"
	if err := s.db.WithContext(ctx).Save(&booking).Error; err != nil {
		return err
	}

	return nil
}
