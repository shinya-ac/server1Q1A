package folder

import (
	"context"
)

type FolderRepository interface {
	Create(ctx context.Context, folder *Folder) error
	FindById(ctx context.Context, id string) (*Folder, error)
	Delete(ctx context.Context, folderId string) error
	Update(ctx context.Context, folder *Folder) error
	GetFoldersByUserId(ctx context.Context, userId string) ([]*Folder, error)
}
