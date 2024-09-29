package folder

import (
	"context"
)

type FolderRepository interface {
	Create(ctx context.Context, folder *Folder) error
	FindById(ctx context.Context, id string) (*Folder, error)
	Delete(ctx context.Context, folderId string) error
}
