package model

import (
	"testing"

	"booking/shared/constants"
	errs "booking/shared/errors"

	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"
)

type UserTestSuite struct {
	suite.Suite
}

func TestUserSuite(t *testing.T) {
	suite.Run(t, new(UserTestSuite))
}

func (s *UserTestSuite) TestNewUser() {
	tests := []struct {
		name        string
		input       RegisterInput
		expectedErr error
	}{
		{
			name: "success create new user",
			input: RegisterInput{
				Name:     "Test User",
				Email:    "test@example.com",
				Password: "password123",
			},
			expectedErr: nil,
		},
		{
			name: "error password too short",
			input: RegisterInput{
				Name:     "Test User",
				Email:    "test@example.com",
				Password: "123",
			},
			expectedErr: errs.ErrShortPassword,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			user, err := NewUser(tt.input)
			if tt.expectedErr != nil {
				s.Error(err)
				s.Equal(tt.expectedErr, err)
				s.Nil(user)
			} else {
				s.NoError(err)
				s.NotNil(user)
				s.Equal(tt.input.Name, user.Name)
				s.Equal(tt.input.Email, user.Email)
				s.Equal(constants.RoleUser, user.Role)
				s.NotEmpty(user.Password)
			}
		})
	}
}

func (s *UserTestSuite) TestSetPassword() {
	user := &User{
		Name:  "Test User",
		Email: "test@example.com",
	}

	tests := []struct {
		name        string
		password    string
		expectedErr error
	}{
		{
			name:        "success set password",
			password:    "password123",
			expectedErr: nil,
		},
		{
			name:        "error empty password",
			password:    "",
			expectedErr: nil, // bcrypt will still hash empty string
		},
		{
			name:        "error password too long",
			password:    string(make([]byte, 73)), // bcrypt has a limit of 72 bytes
			expectedErr: bcrypt.ErrPasswordTooLong,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			err := user.SetPassword(tt.password)
			if tt.expectedErr != nil {
				s.Error(err)
				s.Equal(tt.expectedErr, err)
			} else {
				s.NoError(err)
				s.NotEmpty(user.Password)
			}
		})
	}
}

func (s *UserTestSuite) TestCheckPassword() {
	user := &User{
		Name:  "Test User",
		Email: "test@example.com",
	}
	correctPassword := "password123"
	err := user.SetPassword(correctPassword)
	s.NoError(err)

	tests := []struct {
		name        string
		password    string
		expectedErr error
	}{
		{
			name:        "success check correct password",
			password:    correctPassword,
			expectedErr: nil,
		},
		{
			name:        "error check wrong password",
			password:    "wrongpassword",
			expectedErr: errs.ErrInvalidPassword,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			err := user.CheckPassword(tt.password)
			if tt.expectedErr != nil {
				s.Error(err)
				s.Equal(tt.expectedErr, err)
			} else {
				s.NoError(err)
			}
		})
	}
}

func (s *UserTestSuite) TestIsAdmin() {
	tests := []struct {
		name     string
		user     *User
		expected bool
	}{
		{
			name: "user is admin",
			user: &User{
				Name:  "Admin User",
				Email: "admin@example.com",
				Role:  constants.RoleAdmin,
			},
			expected: true,
		},
		{
			name: "user is not admin",
			user: &User{
				Name:  "Regular User",
				Email: "user@example.com",
				Role:  constants.RoleUser,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			result := tt.user.IsAdmin()
			s.Equal(tt.expected, result)
		})
	}
}
