package folder

import (
	"context"
	"errors"

	folderDomain "github.com/shinya-ac/server1Q1A/domain/folder"
	"github.com/shinya-ac/server1Q1A/pkg/auth"
	"github.com/shinya-ac/server1Q1A/pkg/logging"
)

type DeleteFolderUseCase struct {
	folderRepo folderDomain.FolderRepository
}

func NewDeleteFolderUseCase(folderRepo folderDomain.FolderRepository) *DeleteFolderUseCase {
	return &DeleteFolderUseCase{
		folderRepo: folderRepo,
	}
}

func (uc *DeleteFolderUseCase) Run(ctx context.Context, folderId string) error {
	sub, err := auth.GetUserIDFromContext(ctx)
	if err != nil {
		logging.Logger.Error("Failed to get user ID", "error", err)
		return err
	}

	folder, err := uc.folderRepo.FindById(ctx, folderId)
	if err != nil {
		return errors.New("folderが見つかりません")
	}

	if folder.UserId != sub {
		return errors.New("folderを削除する権限がありません")
	}

	err = uc.folderRepo.Delete(ctx, folderId)
	if err != nil {
		logging.Logger.Error("Folderの削除に失敗しました", "error", err)
		return err
	}

	return nil
}
