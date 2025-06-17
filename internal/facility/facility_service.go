package facility

import (
	facilityModel "booking/internal/facility/model"
	"booking/pkg/logger"
	"errors"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type FacilityServiceInterface interface {
	Create(input facilityModel.CreateFacilityInput) (*facilityModel.Facility, error)
	GetAll() ([]facilityModel.Facility, error)
	GetByID(id string) (*facilityModel.Facility, error)
	Update(id string, input facilityModel.CreateFacilityInput) (*facilityModel.Facility, error)
	Delete(id string) error
	CheckExists(name string) (bool, error)
}

type FacilityService struct {
	logger logger.Logger
	db     *gorm.DB
}

func NewFacilityService(db *gorm.DB, logger logger.Logger) *FacilityService {
	return &FacilityService{
		db:     db,
		logger: logger,
	}
}

func (s *FacilityService) CheckExists(name string) (bool, error) {
	var count int64
	err := s.db.Model(&facilityModel.Facility{}).Where("name = ?", name).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (s *FacilityService) Create(input facilityModel.CreateFacilityInput) (*facilityModel.Facility, error) {
	if input.Name == "" {
		return nil, errors.New("name is required")
	}

	facility, err := facilityModel.NewFacility(input)
	if err != nil {
		s.logger.WithFields(logrus.Fields{
			"name":  input.Name,
			"error": err.Error(),
		}).Error("Failed to create new facility")
		return nil, err
	}
	if err := s.db.Create(facility).Error; err != nil {
		s.logger.WithFields(logrus.Fields{
			"name":  input.Name,
			"error": err.Error(),
		}).Error("Failed to saving data to database")
		return nil, err
	}
	return facility, nil
}

func (s *FacilityService) GetAll() ([]facilityModel.Facility, error) {
	var facilities []facilityModel.Facility
	if err := s.db.Find(&facilities).Error; err != nil {
		s.logger.WithFields(logrus.Fields{
			"name":  "",
			"error": err.Error(),
		}).Error("Failed to get all facility")
		return nil, err
	}
	return facilities, nil
}

func (s *FacilityService) GetByID(id string) (*facilityModel.Facility, error) {
	var facility facilityModel.Facility
	if err := s.db.First(&facility, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("facility not found")
		}
		return nil, err
	}
	return &facility, nil
}

func (s *FacilityService) Update(id string, input facilityModel.CreateFacilityInput) (*facilityModel.Facility, error) {
	var facility facilityModel.Facility
	if err := s.db.First(&facility, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("facility not found")
		}
		return nil, err
	}

	facility.Name = input.Name
	if err := s.db.Save(&facility).Error; err != nil {
		return nil, err
	}

	return &facility, nil
}

func (s *FacilityService) Delete(id string) error {
	var facility facilityModel.Facility
	if err := s.db.First(&facility, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("facility not found")
		}
		return err
	}

	if err := s.db.Delete(&facility).Error; err != nil {
		return err
	}

	return nil
}
