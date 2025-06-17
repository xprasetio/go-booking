package user

import (
	"context"
	"strings"
	"time"

	"booking/internal/user/model"
	"booking/pkg/jwt"
	"booking/pkg/logger"
	"booking/pkg/redis"
	"booking/shared/constants"
	userErr "booking/shared/errors"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// UserServiceInterface mendefinisikan kontrak untuk UserService
type UserServiceInterface interface {
	Register(ctx context.Context, input model.RegisterInput) (*model.User, error)
	Login(ctx context.Context, input model.LoginInput) (string, error)
	GetUserByID(ctx context.Context, userID string) (*model.User, error)
	GetStoredToken(ctx context.Context, userID string) (string, error)
	Logout(ctx context.Context, userID string) error
	UpdateProfile(ctx context.Context, userID string, input model.UpdateProfileInput) (*model.User, error)
	DeleteAccount(ctx context.Context, userID string) error
	GetAllUsers(ctx context.Context) ([]model.User, error)
	UpdateUserRole(ctx context.Context, adminID string, targetUserID string, newRole constants.Role) (*model.User, error)
}

type UserService struct {
	db          *gorm.DB
	jwtSecret   string
	logger      logger.Logger
	redisClient *redis.RedisClient
}

func NewUserService(db *gorm.DB, jwtSecret string, logger logger.Logger, redisClient *redis.RedisClient) *UserService {
	if db == nil {
		panic("database connection is required")
	}
	if logger == nil {
		panic("logger is required")
	}
	if redisClient == nil {
		panic("redis client is required")
	}
	if jwtSecret == "" {
		panic("jwt secret is required")
	}

	return &UserService{
		db:          db,
		jwtSecret:   jwtSecret,
		logger:      logger,
		redisClient: redisClient,
	}
}

func (s *UserService) Register(ctx context.Context, input model.RegisterInput) (*model.User, error) {
	var existingUser model.User
	if err := s.db.WithContext(ctx).Where("email = ?", input.Email).First(&existingUser).Error; err == nil {
		s.logger.WithFields(logrus.Fields{
			"email": input.Email,
			"error": userErr.ErrEmailAlreadyRegistered.Error(),
		}).Error("Email sudah terdaftar")
		return nil, userErr.ErrEmailAlreadyRegistered
	}

	user, err := model.NewUser(input)
	if err != nil {
		s.logger.WithFields(logrus.Fields{
			"email": input.Email,
			"error": err.Error(),
		}).Error("Gagal membuat user baru")
		return nil, err
	}

	if err := s.db.WithContext(ctx).Create(user).Error; err != nil {
		s.logger.WithFields(logrus.Fields{
			"email": input.Email,
			"error": err.Error(),
		}).Error("Gagal menyimpan user ke database")
		return nil, err
	}

	return user, nil
}

func (s *UserService) Login(ctx context.Context, input model.LoginInput) (string, error) {
	var user model.User
	if err := s.db.WithContext(ctx).Where("email = ?", input.Email).First(&user).Error; err != nil {
		s.logger.WithFields(logrus.Fields{
			"email": input.Email,
			"error": userErr.ErrInvalidCredentials.Error(),
		}).Error("Kredensial login tidak valid")
		return "", userErr.ErrInvalidCredentials
	}

	if err := user.CheckPassword(input.Password); err != nil {
		s.logger.WithFields(logrus.Fields{
			"email": input.Email,
			"error": userErr.ErrInvalidCredentials.Error(),
		}).Error("Password tidak valid")
		return "", userErr.ErrInvalidCredentials
	}

	token, err := jwt.GenerateToken(user.ID.String(), s.jwtSecret)
	if err != nil {
		s.logger.WithFields(logrus.Fields{
			"user_id": user.ID,
			"error":   err.Error(),
		}).Error("Gagal generate token")
		return "", err
	}
	if err := s.redisClient.SetToken(ctx, user.ID.String(), token, 24*time.Hour); err != nil {
		s.logger.WithFields(logrus.Fields{
			"user_id": user.ID,
			"error":   err.Error(),
		}).Error("Gagal menyimpan token di Redis")
		return "", err
	}

	return token, nil
}

func (s *UserService) GetUserByID(ctx context.Context, userID string) (*model.User, error) {
	var user model.User
	if err := s.db.WithContext(ctx).First(&user, "id = ?", userID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserService) GetStoredToken(ctx context.Context, userID string) (string, error) {
	if s.redisClient == nil {
		return "", userErr.ErrInvalidCredentials
	}
	return s.redisClient.GetToken(ctx, userID)
}

func (s *UserService) Logout(ctx context.Context, userID string) error {
	if err := s.redisClient.DeleteToken(ctx, userID); err != nil {
		s.logger.WithFields(logrus.Fields{
			"user_id": userID,
			"error":   err.Error(),
		}).Error("Gagal menghapus token dari Redis")
		return err
	}
	return nil
}

func (s *UserService) UpdateProfile(ctx context.Context, userID string, input model.UpdateProfileInput) (*model.User, error) {
	var user model.User
	if err := s.db.WithContext(ctx).First(&user, "id = ?", userID).Error; err != nil {
		return nil, err
	}

	// Update fields
	user.Name = strings.TrimSpace(input.Name)
	if input.Email != "" {
		user.Email = strings.TrimSpace(input.Email)
	}
	if input.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(strings.TrimSpace(input.Password)), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		user.Password = string(hashedPassword)
	}

	if err := s.db.WithContext(ctx).Save(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *UserService) DeleteAccount(ctx context.Context, userID string) error {
	return s.db.WithContext(ctx).Delete(&model.User{}, "id = ?", userID).Error
}

// GetAllUsers mengambil semua user tanpa password
func (s *UserService) GetAllUsers(ctx context.Context) ([]model.User, error) {
	var users []model.User
	if err := s.db.WithContext(ctx).Find(&users).Error; err != nil {
		s.logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Gagal mengambil daftar user")
		return nil, err
	}

	// Hapus password dari setiap user
	for i := range users {
		users[i].Password = ""
	}

	return users, nil
}

// UpdateUserRole mengubah role user oleh superadmin
func (s *UserService) UpdateUserRole(ctx context.Context, adminID string, targetUserID string, newRole constants.Role) (*model.User, error) {
	// Cek apakah admin adalah superadmin
	var admin model.User
	if err := s.db.WithContext(ctx).Where("id = ?", adminID).First(&admin).Error; err != nil {
		return nil, err
	}

	if admin.Role != constants.RoleSuperAdmin {
		s.logger.WithFields(logrus.Fields{
			"admin_id": adminID,
			"error":    "Unauthorized: Only superadmin can change user roles",
		}).Error("Akses ditolak untuk mengubah role")
		return nil, userErr.ErrUnauthorized
	}

	// Cek target user
	var targetUser model.User
	if err := s.db.WithContext(ctx).Where("id = ?", targetUserID).First(&targetUser).Error; err != nil {
		return nil, err
	}

	// Update role
	targetUser.Role = newRole
	if err := s.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", targetUserID).Update("role", newRole).Error; err != nil {
		s.logger.WithFields(logrus.Fields{
			"target_user_id": targetUserID,
			"new_role":       newRole,
			"error":          err.Error(),
		}).Error("Gagal mengubah role user")
		return nil, err
	}

	return &targetUser, nil
}
