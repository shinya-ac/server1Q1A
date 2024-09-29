package folder

import (
	"context"

	folderDomain "github.com/shinya-ac/server1Q1A/domain/folder"
	"github.com/shinya-ac/server1Q1A/pkg/auth"
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
	sub, err := auth.GetUserIDFromContext(ctx)
	if err != nil {
		logging.Logger.Error("Failed to get user ID", "error", err)
		return nil, err
	}

	t, err := folderDomain.NewFolder(sub, dto.Title)
	if err != nil {
		logging.Logger.Error("Failed to create new folder", "error", err)
		return nil, err
	}

	err = uc.folderRepo.Create(ctx, t)
	if err != nil {
		logging.Logger.Error("Failed to create folder in repository", "error", err)
		return nil, err
	}

	return &CreateFolderUseCaseOutputDto{
		Id: t.GetId(),
	}, nil
}
