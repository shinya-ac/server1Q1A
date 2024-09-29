package folder

import (
	"github.com/labstack/echo/v4"

	"github.com/shinya-ac/server1Q1A/application/folder"
	"github.com/shinya-ac/server1Q1A/pkg/logging"
	validator "github.com/shinya-ac/server1Q1A/pkg/validator"
	"github.com/shinya-ac/server1Q1A/presentation/settings"
)

type handler struct {
	createFolderUseCase *folder.CreateFolderUseCase
	deleteFolderUseCase *folder.DeleteFolderUseCase
}

func NewHandler(
	createFolderUseCase *folder.CreateFolderUseCase,
	deleteFolderUseCase *folder.DeleteFolderUseCase,
) handler {
	return handler{
		createFolderUseCase: createFolderUseCase,
		deleteFolderUseCase: deleteFolderUseCase,
	}
}

// CreateFolders godoc
// @Summary Folderを登録する
// @Description パラメーターから新規Todoを作成する
// @Tags Folder
// @Accept json
// @Produce json
// @Param request body CreateFolderParams true "Folder登録"
// @Success 201 {object} createFolderResponse
// @Router /v1/folders [post]

func (h handler) CreateFolders(ctx echo.Context) error {
	logging.Logger.Info("CreateFolder実行開始")

	var params CreateFolderParams
	err := ctx.Bind(&params)
	if err != nil {
		logging.Logger.Error("paramsの形式が不正", "error", err)
		settings.ReturnBadRequest(ctx, err)
		return err
	}
	validate := validator.GetValidator()
	err = validate.Struct(params)
	if err != nil {
		logging.Logger.Error("paramsの内容が不正", "error", err)
		settings.ReturnStatusBadRequest(ctx, err)
		return err
	}

	input := folder.CreateFolderUseCaseInputDto{
		Title: params.Title,
	}

	dto, err := h.createFolderUseCase.Run(ctx.Request().Context(), input)
	if err != nil {
		logging.Logger.Error("usecaseの実行に失敗", "error", err)
		settings.ReturnError(ctx, err)
		return err
	}

	response := createFolderResponse{
		FolderId: dto.Id,
	}
	return settings.ReturnStatusCreated(ctx, response)
}

// CreateFolders godoc
// @Summary Folderを削除する
// @Description idを指定してFolderを削除する
// @Tags Folder
// @Accept json
// @Produce json
// @Param id path string true "フォルダーID"
// @Success 204 {object} nil "フォルダーが正常に削除されました"
// @Router /v1/folders/{id} [delete]

func (h handler) DeleteFolder(ctx echo.Context) error {
	logging.Logger.Info("DeleteFolder実行開始")

	folderId := ctx.Param("id")

	err := h.deleteFolderUseCase.Run(ctx.Request().Context(), folderId)
	if err != nil {
		if err.Error() == "folderが見つからない" {
			return settings.ReturnNotFound(ctx, err)
		}
		if err.Error() == "folderを削除する権限がない" {
			return settings.ReturnForbidden(ctx, err)
		}
		return settings.ReturnStatusInternalServerError(ctx, err)
	}

	return settings.ReturnStatusNoContent(ctx)
}
