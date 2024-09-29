package route

import (
	"fmt"

	"github.com/labstack/echo/v4"

	folderApp "github.com/shinya-ac/server1Q1A/application/folder"
	"github.com/shinya-ac/server1Q1A/infrastructure/mysql/db"
	"github.com/shinya-ac/server1Q1A/infrastructure/mysql/repository"
	"github.com/shinya-ac/server1Q1A/middlewares/auth0"
	folderPre "github.com/shinya-ac/server1Q1A/presentation/folder"
	"github.com/shinya-ac/server1Q1A/presentation/health_handler"
	"github.com/shinya-ac/server1Q1A/presentation/settings"
	// swaggerFiles "github.com/swaggo/files"
	// ginSwagger "github.com/swaggo/gin-swagger"
)

func InitRoute(api *echo.Echo) {
	api.Use(settings.ErrorHandler)

	api.GET("/v1/health", health_handler.HealthCheck)

	protectedV1 := api.Group("/v1")
	protectedV1.Use(echo.WrapMiddleware(auth0.UseJWT))
	folderRoute(protectedV1)
	// api.GET("/v1/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

func folderRoute(r *echo.Group) {
	fmt.Printf("folderRoute実行")
	folderRepository := repository.NewFolderRepository(db.GetDB())
	cuc := folderApp.NewCreateFolderUseCase(folderRepository)
	duc := folderApp.NewDeleteFolderUseCase(folderRepository)
	h := folderPre.NewHandler(cuc, duc)

	group := r.Group("/folders")
	group.POST("/", h.CreateFolders)
	group.DELETE("/:id", h.DeleteFolder)
}
