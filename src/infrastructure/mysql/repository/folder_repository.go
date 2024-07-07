package repository

import (
	"context"
	"database/sql"

	errDomain "github.com/shinya-ac/server1Q1A/domain/error"
	"github.com/shinya-ac/server1Q1A/domain/folder"
	"github.com/shinya-ac/server1Q1A/pkg/logging"
)

type FolderRepository struct {
	db *sql.DB
}

func NewFolderRepository(db *sql.DB) folder.FolderRepository {
	return &FolderRepository{db: db}
}

func (r *FolderRepository) Create(ctx context.Context, folder *folder.Folder) error {
	if folder == nil {
		logging.Logger.Error("Folderがnil")
		err := errDomain.NewError("Folderがnilです。")
		return err
	}
	logging.Logger.Info("Create実行", "folder:", *folder)
	query := `INSERT INTO folders (id, title) VALUES(?, ?)`

	_, err := r.db.ExecContext(ctx, query, folder.Id, folder.Title)
	if err != nil {
		logging.Logger.Error("SQL実行に失敗", "error", err)
		return err
	}
	return nil
}
