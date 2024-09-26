package folder

import (
	"context"
	"errors"

	"github.com/form3tech-oss/jwt-go"
	folderDomain "github.com/shinya-ac/server1Q1A/domain/folder"
	"github.com/shinya-ac/server1Q1A/middlewares/auth0"
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
	token := auth0.GetJWT(ctx)
	if token == nil {
		logging.Logger.Error("JWT token not found in context")
		return nil, errors.New("token not found")
	}

	claims := token.Claims.(jwt.MapClaims)

	sub, ok := claims["sub"].(string)
	if !ok {
		logging.Logger.Error("sub claim not found in token claims")
		return nil, errors.New("sub not found")
	}

	t, err := folderDomain.NewFolder(dto.Title, sub)
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
