package repository

import (
	"context"
	"database/sql"
	"errors"

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
	query := `INSERT INTO folders (id, user_id, title, created_at, updated_at) VALUES(?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`

	_, err := r.db.ExecContext(ctx, query, folder.Id, folder.UserId, folder.Title)
	if err != nil {
		logging.Logger.Error("SQL実行に失敗", "error", err)
		return err
	}
	return nil
}

func (r *FolderRepository) FindById(ctx context.Context, folderId string) (*folder.Folder, error) {
	query := `SELECT id, user_id, title FROM folders WHERE id = ?`

	row := r.db.QueryRowContext(ctx, query, folderId)
	var f folder.Folder
	err := row.Scan(&f.Id, &f.UserId, &f.Title)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("folderが見つかりません")
		}
		logging.Logger.Error("folderを取得できませんでした", "error", err)
		return nil, err
	}
	return &f, nil
}

func (r *FolderRepository) Delete(ctx context.Context, folderId string) error {
	query := `DELETE FROM folders WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, folderId)
	if err != nil {
		logging.Logger.Error("folderの削除に失敗しました", "error", err)
		return err
	}
	return nil
}

func (r *FolderRepository) Update(ctx context.Context, folder *folder.Folder) error {
	if folder == nil {
		logging.Logger.Error("Folderがnil")
		err := errDomain.NewError("Folderがnilです。")
		return err
	}
	query := `UPDATE folders SET title = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, folder.Title, folder.Id)
	if err != nil {
		logging.Logger.Error("Folderの更新に失敗しました", "error", err)
		return err
	}
	return nil
}

func (r *FolderRepository) GetFoldersByUserId(ctx context.Context, userId string) ([]*folder.Folder, error) {
	var folders []*folder.Folder
	query := "SELECT id, title, user_id FROM folders WHERE user_id = ?"

	rows, err := r.db.QueryContext(ctx, query, userId)
	if err != nil {
		logging.Logger.Error("フォルダー一覧取得中にエラーが発生", "error", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var f folder.Folder
		if err := rows.Scan(&f.Id, &f.Title, &f.UserId); err != nil {
			logging.Logger.Error("フォルダーの読み込みに失敗", "error", err)
			return nil, err
		}
		folders = append(folders, &f)
	}

	if err = rows.Err(); err != nil {
		logging.Logger.Error("フォルダー一覧の取得完了中にエラーが発生", "error", err)
		return nil, err
	}

	return folders, nil
}
