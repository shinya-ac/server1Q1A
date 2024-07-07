package folder

import (
	"context"
)

type FolderRepository interface {
	Create(ctx context.Context, folder *Folder) error
}
