package repository

import (
	"context"
	"database/sql"
	"strings"

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

func (r *MySQLAnswerRepository) BulkCreate(ctx context.Context, answers []*answerDomain.Answer) error {
	if len(answers) == 0 {
		return nil
	}

	query := `INSERT INTO answers (id, user_id, question_id, folder_id, content, created_at, updated_at) VALUES `
	values := []string{}
	args := []interface{}{}

	for _, a := range answers {
		values = append(values, "(?, ?, ?, ?, ?, ?, ?)")
		args = append(args, a.Id, a.UserId, a.QuestionId, a.FolderId, a.Content, a.CreatedAt, a.UpdatedAt)
	}

	query += strings.Join(values, ",")
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		logging.Logger.Error("回答バルク作成用のクエリの準備に失敗しました", "error", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, args...)
	if err != nil {
		logging.Logger.Error("回答のバルク作成に失敗しました", "error", err)
		return err
	}

	return nil
}
