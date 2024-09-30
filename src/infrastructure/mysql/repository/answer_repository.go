package repository

import (
	"context"
	"database/sql"

	answerDomain "github.com/shinya-ac/server1Q1A/domain/answer"
	"github.com/shinya-ac/server1Q1A/pkg/logging"
)

type MySQLAnswerRepository struct {
	db *sql.DB
}

func NewMySQLAnswerRepository(db *sql.DB) *MySQLAnswerRepository {
	return &MySQLAnswerRepository{db: db}
}

func (r *MySQLAnswerRepository) Create(ctx context.Context, a *answerDomain.Answer) error {
	query := `INSERT INTO answers (id, user_id, question_id, folder_id, content, created_at, updated_at) 
	          VALUES (?, ?, ?, ?, ?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query, a.Id, a.UserId, a.QuestionId, a.FolderId, a.Content, a.CreatedAt, a.UpdatedAt)
	if err != nil {
		logging.Logger.Error("Answerテーブルへのインサートに失敗しました", "error", err)
		return err
	}
	return nil
}
