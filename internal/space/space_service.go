package space

import (
	categoryModel "booking/internal/category/model"
	spaceModel "booking/internal/space/model"

	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SpaceServiceInterface interface {
	Create(input spaceModel.CreateSpaceInput, category *categoryModel.Category) (*spaceModel.Space, error)
	GetAll() ([]spaceModel.Space, error)
	GetByID(id string) (*spaceModel.Space, error)
	Update(id string, input spaceModel.CreateSpaceInput) (*spaceModel.Space, error)
	Delete(id string) error
}

type SpaceService struct {
	db *gorm.DB
}

func NewSpaceService(db *gorm.DB) *SpaceService {
	return &SpaceService{
		db: db,
	}
}

func (s *SpaceService) Create(input spaceModel.CreateSpaceInput, category *categoryModel.Category) (*spaceModel.Space, error) {
	if category == nil {
		return nil, errors.New("category is required")
	}

	if category.ID == uuid.Nil {
		return nil, errors.New("invalid category ID")
	}

	space, err := spaceModel.NewSpace(input, category.ID)
	if err != nil {
		return nil, err
	}

	if err := s.db.Create(space).Error; err != nil {
		return nil, err
	}

	return space, nil
}

func (s *SpaceService) GetAll() ([]spaceModel.Space, error) {
	var spaces []spaceModel.Space
	if err := s.db.Find(&spaces).Error; err != nil {
		return nil, err
	}
	return spaces, nil
}

func (s *SpaceService) GetByID(id string) (*spaceModel.Space, error) {
	spaceID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid space ID")
	}

	var space spaceModel.Space
	if err := s.db.First(&space, "id = ?", spaceID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("space not found")
		}
		return nil, err
	}

	return &space, nil
}

func (s *SpaceService) Update(id string, input spaceModel.CreateSpaceInput) (*spaceModel.Space, error) {
	spaceID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid space ID")
	}

	var space spaceModel.Space
	if err := s.db.First(&space, "id = ?", spaceID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("space not found")
		}
		return nil, err
	}

	// Update fields
	space.Name = input.Name
	space.Description = input.Description
	space.PricePerNight = input.PricePerNight
	space.CategoryID = input.CategoryID

	if err := s.db.Save(&space).Error; err != nil {
		return nil, err
	}

	return &space, nil
}

func (s *SpaceService) Delete(id string) error {
	spaceID, err := uuid.Parse(id)
	if err != nil {
		return errors.New("invalid space ID")
	}

	result := s.db.Delete(&spaceModel.Space{}, "id = ?", spaceID)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("space not found")
	}

	return nil
}
