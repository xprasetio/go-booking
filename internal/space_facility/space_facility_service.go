package spacefacility

import (
	spaceFacilityModel "booking/internal/space_facility/model"
	"booking/pkg/logger"
	"context"
	"errors"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type SpaceFacilityServiceInterface interface {
	Create(ctx context.Context, input spaceFacilityModel.SpaceFacilityInput) (*spaceFacilityModel.SpaceFacility, error)
}
type SpaceFacilityService struct {
	db     *gorm.DB
	logger logger.Logger
}

func NewSpaceFacilityService(db *gorm.DB, logger logger.Logger) *SpaceFacilityService {
	return &SpaceFacilityService{
		db:     db,
		logger: logger,
	}
}

func (s *SpaceFacilityService) Create(ctx context.Context, input spaceFacilityModel.SpaceFacilityInput) (*spaceFacilityModel.SpaceFacility, error) {
	// cek apakah space ada
	var count int64
	if err := s.db.WithContext(ctx).Table("spaces").Where("id = ?", input.SpaceID).Count(&count).Error; err != nil {
		s.logger.Error(ctx, "failed to check space existence", err)
		return nil, err
	}
	if count == 0 {
		return nil, errors.New("space not found")
	}

	// Cek apakah facility ada
	if err := s.db.WithContext(ctx).Table("facilities").Where("id = ?", input.FacilityID).Count(&count).Error; err != nil {
		s.logger.Error(ctx, "failed to check facility existence", err)
		return nil, err
	}
	if count == 0 {
		return nil, errors.New("facility not found")
	}

	sf, err := spaceFacilityModel.NewSpaceFacility(input)
	if err != nil {
		s.logger.WithFields(logrus.Fields{
			"method":      "Create",
			"space_id":    input.SpaceID,
			"facility_id": input.FacilityID,
		}).Error(ctx, "failed to create space facility", err)
		return nil, err
	}

	// Simpan ke database
	if err := s.db.WithContext(ctx).Create(sf).Error; err != nil {
		s.logger.WithFields(logrus.Fields{
			"method":      "Create",
			"space_id":    input.SpaceID,
			"facility_id": input.FacilityID,
		}).Error(ctx, "failed to save space facility to database", err)
		return nil, err
	}

	return sf, nil
}
