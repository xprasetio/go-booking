package model

import (
	"time"

	"booking/shared/constants"
	errs "booking/shared/errors"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        uuid.UUID      `json:"id" gorm:"type:char(36);primary_key;default:(UUID())"`
	Name      string         `json:"name" gorm:"size:50" validate:"required,min=2,max=50"`
	Email     string         `json:"email" gorm:"unique" validate:"required,email"`
	Password  string         `json:"-"`
	Role      constants.Role `json:"role"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

// DTO: Register input
type RegisterInput struct {
	Name     string `json:"name" validate:"required,min=2,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

// DTO: Login input
type LoginInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// DTO: Update profile input
type UpdateProfileInput struct {
	Name     string `json:"name" validate:"required,min=2,max=50"`
	Email    string `json:"email" validate:"omitempty,email"`
	Role     string `json:"role" validate:"omitempty`
	Password string `json:"password,omitempty" validate:"omitempty,min=6"`
}

// Factory: Create new user from register input
func NewUser(input RegisterInput) (*User, error) {
	if len(input.Password) < constants.PasswordMinLength {
		return nil, errs.ErrShortPassword
	}

	user := &User{
		ID:    uuid.New(),
		Name:  input.Name,
		Email: input.Email,
		Role:  constants.RoleUser,
	}

	if err := user.SetPassword(input.Password); err != nil {
		return nil, errs.ErrHashingPassword
	}

	return user, nil
}

// Password setter with hashing
func (u *User) SetPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), constants.BcryptCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// Password checker
func (u *User) CheckPassword(password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return errs.ErrInvalidPassword
	}
	return nil
}

// Role checker
func (u *User) IsAdmin() bool {
	return u.Role == constants.RoleAdmin
}
func (u *User) IsSuperAdmin() bool {
	return u.Role == constants.RoleSuperAdmin
}
