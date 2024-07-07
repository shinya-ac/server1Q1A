package folder

import (
	"context"

	folderDomain "github.com/shinya-ac/server1Q1A/domain/folder"
	"github.com/shinya-ac/server1Q1A/pkg/logging"
)

type CreateFolderUseCase struct {
	folderRepo folderDomain.FolderRepository
}

func NewCreateFolderUseCase(
	folderRepo folderDomain.FolderRepository,
) *CreateFolderUseCase {
	return &CreateFolderUseCase{
		folderRepo: folderRepo,
	}
}

type CreateFolderUseCaseInputDto struct {
	Title string
}

type CreateFolderUseCaseOutputDto struct {
	Id string
}

func (uc *CreateFolderUseCase) Run(
	ctx context.Context,
	dto CreateFolderUseCaseInputDto,
) (*CreateFolderUseCaseOutputDto, error) {
	t, err := folderDomain.NewFolder(dto.Title)
	if err != nil {
		logging.Logger.Error("サーバーエラー", "error", err)
		return nil, err
	}
	err = uc.folderRepo.Create(ctx, t)
	if err != nil {
		logging.Logger.Error("サーバーエラー", "error", err)
		return nil, err
	}
	return &CreateFolderUseCaseOutputDto{
		Id: t.GetId(),
	}, nil
}
