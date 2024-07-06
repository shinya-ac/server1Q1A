package route

import (
	"github.com/labstack/echo/v4"

	"github.com/shinya-ac/server1Q1A/presentation/health_handler"
	"github.com/shinya-ac/server1Q1A/presentation/settings"
	// swaggerFiles "github.com/swaggo/files"
	// ginSwagger "github.com/swaggo/gin-swagger"
)

func InitRoute(api *echo.Echo) {
	api.Use(settings.ErrorHandler)

	api.GET("/v1/health", health_handler.HealthCheck)
	// api.GET("/v1/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
