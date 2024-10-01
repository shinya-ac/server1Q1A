package folder

import (
	"context"

	folderDomain "github.com/shinya-ac/server1Q1A/domain/folder"
	"github.com/shinya-ac/server1Q1A/pkg/auth"
	"github.com/shinya-ac/server1Q1A/pkg/logging"
)

type ReadFoldersUseCase struct {
	folderRepo folderDomain.FolderRepository
}

func NewReadFoldersUseCase(
	folderRepo folderDomain.FolderRepository,
) *ReadFoldersUseCase {
	return &ReadFoldersUseCase{
		folderRepo: folderRepo,
	}
}

func (uc *ReadFoldersUseCase) Run(ctx context.Context) ([]*folderDomain.Folder, error) {
	sub, err := auth.GetUserIDFromContext(ctx)
	if err != nil {
		logging.Logger.Error("ユーザーIDの取得に失敗しました", "error", err)
		return nil, err
	}

	folders, err := uc.folderRepo.GetFoldersByUserId(ctx, sub)
	if err != nil {
		logging.Logger.Error("フォルダー一覧の取得に失敗", "error", err)
		return nil, err
	}

	return folders, nil
}
