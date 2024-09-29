package folder

import (
	"context"

	folderDomain "github.com/shinya-ac/server1Q1A/domain/folder"
	"github.com/shinya-ac/server1Q1A/pkg/auth"
	"github.com/shinya-ac/server1Q1A/pkg/logging"
)

type UpdateFolderUseCase struct {
	folderRepo folderDomain.FolderRepository
}

func NewUpdateFolderUseCase(
	folderRepo folderDomain.FolderRepository,
) *UpdateFolderUseCase {
	return &UpdateFolderUseCase{
		folderRepo: folderRepo,
	}
}

type UpdateFolderUseCaseInputDto struct {
	Id    string
	Title string
}

func (uc *UpdateFolderUseCase) Run(
	ctx context.Context,
	dto UpdateFolderUseCaseInputDto,
) error {
	sub, err := auth.GetUserIDFromContext(ctx)
	if err != nil {
		logging.Logger.Error("UserIdの取得に失敗しました", err)
		return err
	}

	folder, err := uc.folderRepo.FindById(ctx, dto.Id)
	if err != nil {
		logging.Logger.Error("Folderが存在しません", "error", err)
		return err
	}

	updatedFolder, err := folderDomain.ReconstructFolder(folder.Id, dto.Title, sub)
	if err != nil {
		logging.Logger.Error("Folderの作成に失敗しました", "error", err)
		return err
	}

	err = uc.folderRepo.Update(ctx, updatedFolder)
	if err != nil {
		logging.Logger.Error("Folderの更新に失敗しました", "error", err)
		return err
	}

	return nil
}
