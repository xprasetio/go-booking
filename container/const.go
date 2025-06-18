package container

const (
	ConfigDefName              string = "config"
	DBDefName                  string = "db"
	LoggerDefName              string = "logger"
	EchoDefName                string = "echo"
	ValidatorDefName           string = "validator"
	RedisClientDefName         string = "redisClient"
	AuthMiddlewareDefName      string = "authMiddleware"
	AdminAuthMiddlewareDefName string = "adminAuthMiddleware"

	//Service
	UserServiceDefName          string = "user.service"
	CategoryServiceDefName      string = "category.service"
	SpaceServiceDefName         string = "space.service"
	SpaceFacilityServiceDefName string = "space_facility.service"
	FacilityServiceDefName      string = "facility.service"

	//Handler
	UserHandlerDefName          string = "user.handler"
	CategoryHandlerDefName      string = "category.handler"
	SpaceHandlerDefName         string = "space.handler"
	SpaceFacilityHandlerDefName string = "space_facility.handler"
	FacilityHandlerDefName      string = "facility.handler"
)
