package server

import (
	"context"

	"strings"
	"test/test_app/app/api/cors"
	"test/test_app/app/api/middleware/auth"
	"test/test_app/app/api/middleware/jwt"
	"test/test_app/app/constants"
	authC "test/test_app/app/controller/auth"
	"test/test_app/app/controller/file"
	"test/test_app/app/service/aws"
	"test/test_app/app/service/logger"
	"test/test_app/app/service/supabase"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/spf13/viper"
)

func Init() {
	if strings.EqualFold(viper.GetString("Environment"), "production") {
		gin.SetMode(gin.ReleaseMode)
	}
	r := NewRouter(context.Background())
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	err := r.Run(viper.GetString("server.port"))
	if err != nil {
		logger.SugarLogger.Error("Server not able to startup with error: ", err)
	}
}

func NewRouter(ctx context.Context) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(cors.CORSMiddleware())
	// health := new(controller.HealthController)
	// router.GET("test/user/health", health.Status)

	router.Use(uuidInjectionMiddleware())
	jwt := jwt.NewJwtService()

	log := logger.Logger(ctx)


	//aws service
	supabaseClient := supabase.NewSupabaseService()
	awsInit, err := aws.InitAwsStr(ctx)
	if err != nil {
		log.Fatalf("Aws connection failed with error: %v", err)
	}

	userFileController := file.NewFileController(awsInit, supabaseClient)
	authController := authC.NewAuthController(awsInit, supabaseClient)
	// jwt := jwt.NewJwtService()
	v1 := router.Group("/v1")
	{
		v1.POST(SIGNUP, authController.SignUp())
		v1.POST(LOGIN, authController.Login())
	}

	authToken := v1.Group("")
	{
		authToken.Use(auth.Authentication(jwt))
		authToken.POST(FILE, userFileController.UploadFile())
		authToken.GET(FILE, userFileController.GetFile())
	}

	return router
}

func uuidInjectionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		correlationId := c.GetHeader(constants.CorrelationId)
		if len(correlationId) == 0 {
			correlationID, _ := uuid.NewUUID()
			correlationId = correlationID.String()
			c.Request.Header.Set(constants.CorrelationId, correlationId)
		}
		c.Writer.Header().Set(constants.CorrelationId, correlationId)

		c.Next()
	}
}
