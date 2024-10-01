package route

import (
	"fmt"

	"github.com/labstack/echo/v4"

	chatgptApp "github.com/shinya-ac/server1Q1A/application/chatgpt"
	folderApp "github.com/shinya-ac/server1Q1A/application/folder"
	qaApp "github.com/shinya-ac/server1Q1A/application/qa"
	"github.com/shinya-ac/server1Q1A/infrastructure/chatgpt"
	"github.com/shinya-ac/server1Q1A/infrastructure/mysql/db"
	"github.com/shinya-ac/server1Q1A/infrastructure/mysql/repository"
	"github.com/shinya-ac/server1Q1A/middlewares/auth0"
	chatgptPre "github.com/shinya-ac/server1Q1A/presentation/chatgpt"
	folderPre "github.com/shinya-ac/server1Q1A/presentation/folder"
	"github.com/shinya-ac/server1Q1A/presentation/health_handler"
	qaPre "github.com/shinya-ac/server1Q1A/presentation/qa"
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
	chatRoute(protectedV1)
	qaRoute(protectedV1)
	// api.GET("/v1/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

func folderRoute(r *echo.Group) {
	fmt.Printf("folderRoute実行")
	folderRepository := repository.NewFolderRepository(db.GetDB())
	cuc := folderApp.NewCreateFolderUseCase(folderRepository)
	duc := folderApp.NewDeleteFolderUseCase(folderRepository)
	uuc := folderApp.NewUpdateFolderUseCase(folderRepository)
	h := folderPre.NewHandler(cuc, duc, uuc)

	group := r.Group("/folders")
	group.POST("/", h.CreateFolders)
	group.PATCH("/:id", h.UpdateFolder)
	group.DELETE("/:id", h.DeleteFolder)
}

func chatRoute(r *echo.Group) {
	ChatGPTRepository := chatgpt.NewChatGPTAPI()
	chatgptUsecase := chatgptApp.NewChatGPTUseCase(ChatGPTRepository)
	chatHandler := chatgptPre.NewHandler(chatgptUsecase)

	group := r.Group("/gpt")
	group.POST("/image", chatHandler.Ocr)
}

func qaRoute(r *echo.Group) {
	questionRepository := repository.NewMySQLQuestionRepository(db.GetDB())
	answerRepository := repository.NewMySQLAnswerRepository(db.GetDB())
	transactionManager := repository.NewTransactionManager()
	createQaUseCase := qaApp.NewCreateQaPairUseCase(questionRepository, answerRepository, transactionManager)
	qaHandler := qaPre.NewHandler(createQaUseCase)

	group := r.Group("/folders/:folder_id/qa")
	group.POST("/", qaHandler.CreateQaPairs)
}
