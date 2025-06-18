package container

import (
	"booking/config"
	"booking/internal/booking"
	"booking/internal/category"
	"booking/internal/facility"
	"booking/internal/space"
	spacefacility "booking/internal/space_facility"
	"booking/internal/user"
	"booking/pkg/database"
	"booking/pkg/logger"
	"booking/pkg/middleware"
	"booking/pkg/redis"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/sarulabs/di/v2"
	"gorm.io/gorm"
)

func NewContainer() (di.Container, error) {
	builder, err := di.NewBuilder()
	if err != nil {
		return di.Container{}, err
	}

	defs := []di.Def{
		{
			Name: ConfigDefName,
			Build: func(ctn di.Container) (interface{}, error) {
				return config.LoadConfig()
			},
		},
		{
			Name: LoggerDefName,
			Build: func(ctn di.Container) (interface{}, error) {
				return logger.NewLogger(), nil
			},
		},
		{
			Name: DBDefName,
			Build: func(ctn di.Container) (interface{}, error) {
				cfg := ctn.Get(ConfigDefName).(config.Config)
				return database.InitDB(cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)
			},
		},
		{
			Name: RedisClientDefName,
			Build: func(ctn di.Container) (interface{}, error) {
				cfg := ctn.Get(ConfigDefName).(config.Config)
				return redis.NewRedisClient(cfg.RedisAddr, cfg.RedisPassword, cfg.RedisDB), nil
			},
		},
		{
			Name: UserServiceDefName,
			Build: func(ctn di.Container) (interface{}, error) {
				cfg := ctn.Get(ConfigDefName).(config.Config)
				db := ctn.Get(DBDefName).(*gorm.DB)
				logger := ctn.Get(LoggerDefName).(logger.Logger)
				redisClient := ctn.Get(RedisClientDefName).(*redis.RedisClient)
				return user.NewUserService(db, cfg.JWTSecret, logger, redisClient), nil
			},
		},
		{
			Name: CategoryServiceDefName,
			Build: func(ctn di.Container) (interface{}, error) {
				db := ctn.Get(DBDefName).(*gorm.DB)
				return category.NewCategoryService(db), nil
			},
		},
		{
			Name: UserHandlerDefName,
			Build: func(ctn di.Container) (interface{}, error) {
				userService := ctn.Get(UserServiceDefName).(user.UserServiceInterface)
				return user.NewUserHandler(userService), nil
			},
		},
		{
			Name: CategoryHandlerDefName,
			Build: func(ctn di.Container) (interface{}, error) {
				categoryService := ctn.Get(CategoryServiceDefName).(category.CategoryServiceInterface)
				return category.NewCategoryHandler(categoryService), nil
			},
		},
		{
			Name: AuthMiddlewareDefName,
			Build: func(ctn di.Container) (interface{}, error) {
				cfg := ctn.Get(ConfigDefName).(config.Config)
				userService := ctn.Get(UserServiceDefName).(user.UserServiceInterface)
				return middleware.AuthMiddleware(userService, cfg.JWTSecret), nil
			},
		},
		{
			Name: AdminAuthMiddlewareDefName,
			Build: func(ctn di.Container) (interface{}, error) {
				return middleware.AdminMiddleware(), nil
			},
		},
		{
			Name: ValidatorDefName,
			Build: func(ctn di.Container) (interface{}, error) {
				return validator.New(), nil
			},
		},
		{
			Name: EchoDefName,
			Build: func(ctn di.Container) (interface{}, error) {
				e := echo.New()
				validate := ctn.Get(ValidatorDefName).(*validator.Validate)
				e.Validator = &CustomValidator{validator: validate}
				return e, nil
			},
		},
		{
			Name: SpaceServiceDefName,
			Build: func(ctn di.Container) (interface{}, error) {
				db := ctn.Get(DBDefName).(*gorm.DB)
				return space.NewSpaceService(db), nil
			},
		},
		{
			Name: SpaceHandlerDefName,
			Build: func(ctn di.Container) (interface{}, error) {
				spaceService := ctn.Get(SpaceServiceDefName).(space.SpaceServiceInterface)
				categoryService := ctn.Get(CategoryServiceDefName).(category.CategoryServiceInterface)
				return space.NewSpaceHandler(spaceService, categoryService), nil
			},
		},
		{
			Name: FacilityServiceDefName,
			Build: func(ctn di.Container) (interface{}, error) {
				db := ctn.Get(DBDefName).(*gorm.DB)
				logger := ctn.Get(LoggerDefName).(logger.Logger)
				return facility.NewFacilityService(db, logger), nil
			},
		},
		{
			Name: FacilityHandlerDefName,
			Build: func(ctn di.Container) (interface{}, error) {
				facilityService := ctn.Get(FacilityServiceDefName).(facility.FacilityServiceInterface)
				return facility.NewFacilityHandler(facilityService), nil
			},
		},
		{
			Name: SpaceFacilityServiceDefName,
			Build: func(ctn di.Container) (interface{}, error) {
				db := ctn.Get(DBDefName).(*gorm.DB)
				logger := ctn.Get(LoggerDefName).(logger.Logger)
				return spacefacility.NewSpaceFacilityService(db, logger), nil
			},
		},
		{
			Name: SpaceFacilityHandlerDefName,
			Build: func(ctn di.Container) (interface{}, error) {
				spaceFacilityService := ctn.Get(SpaceFacilityServiceDefName).(spacefacility.SpaceFacilityServiceInterface)
				return spacefacility.NewSpaceFacilityHandler(spaceFacilityService), nil
			},
		},
		{
			Name: BookingServiceDefName,
			Build: func(ctn di.Container) (interface{}, error) {
				db := ctn.Get(DBDefName).(*gorm.DB)
				logger := ctn.Get(LoggerDefName).(logger.Logger)
				userService := ctn.Get(UserServiceDefName).(user.UserServiceInterface)
				spaceService := ctn.Get(SpaceServiceDefName).(space.SpaceServiceInterface)
				return booking.NewBookingService(db, logger, userService, spaceService), nil
			},
		},
		{
			Name: BookingHandlerDefName,
			Build: func(ctn di.Container) (interface{}, error) {
				bookingService := ctn.Get(BookingServiceDefName).(booking.BookingServiceInterface)
				logger := ctn.Get(LoggerDefName).(logger.Logger)
				return booking.NewBookingHandler(bookingService, logger), nil
			},
		},
	}

	if err := builder.Add(defs...); err != nil {
		return di.Container{}, err
	}

	return builder.Build(), nil
}

// CustomValidator adalah custom validator untuk Echo
type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}
