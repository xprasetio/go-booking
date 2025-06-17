package category

import (
	"errors"
	"strings"

	categoryModel "booking/internal/category/model"
	userModel "booking/internal/user/model"

	"gorm.io/gorm"
)

// CategoryServiceInterface mendefinisikan kontrak untuk CategoryService
type CategoryServiceInterface interface {
	Create(input categoryModel.CreateCategoryInput, user *userModel.User) (*categoryModel.Category, error)
	GetAll() ([]categoryModel.Category, error)
	GetByID(id string) (*categoryModel.Category, error)
	Update(id string, input categoryModel.CreateCategoryInput, user *userModel.User) (*categoryModel.Category, error)
	Delete(id string, user *userModel.User) error
}

type CategoryService struct {
	db *gorm.DB
}

func NewCategoryService(db *gorm.DB) *CategoryService {
	return &CategoryService{
		db: db,
	}
}

func (s *CategoryService) Create(input categoryModel.CreateCategoryInput, user *userModel.User) (*categoryModel.Category, error) {
	// Membersihkan spasi di awal dan akhir
	input.Name = strings.TrimSpace(input.Name)
	input.Description = strings.TrimSpace(input.Description)

	// Menghapus spasi berlebih di tengah
	input.Name = strings.Join(strings.Fields(input.Name), " ")
	input.Description = strings.Join(strings.Fields(input.Description), " ")

	category, err := categoryModel.NewCategory(input, user.ID)
	if err != nil {
		return nil, err
	}

	if err := s.db.Create(category).Error; err != nil {
		return nil, err
	}

	return category, nil
}

func (s *CategoryService) GetAll() ([]categoryModel.Category, error) {
	var categories []categoryModel.Category
	if err := s.db.Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

func (s *CategoryService) GetByID(id string) (*categoryModel.Category, error) {
	var category categoryModel.Category
	if err := s.db.First(&category, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

func (s *CategoryService) Update(id string, input categoryModel.CreateCategoryInput, user *userModel.User) (*categoryModel.Category, error) {
	if !user.IsAdmin() {
		return nil, errors.New("only admin can update category")
	}

	category, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	if err := category.Update(input); err != nil {
		return nil, err
	}

	if err := s.db.Save(category).Error; err != nil {
		return nil, err
	}

	return category, nil
}

func (s *CategoryService) Delete(id string, user *userModel.User) error {
	if !user.IsAdmin() {
		return errors.New("only admin can delete category")
	}

	return s.db.Delete(&categoryModel.Category{}, "id = ?", id).Error
}
