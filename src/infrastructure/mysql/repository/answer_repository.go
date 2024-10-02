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

func (r *MySQLAnswerRepository) GetAnswersByQuestionIds(ctx context.Context, questionIds []string) ([]*answerDomain.Answer, error) {
	var answers []*answerDomain.Answer

	placeholders := make([]string, len(questionIds))
	args := make([]interface{}, len(questionIds))
	for i, id := range questionIds {
		placeholders[i] = "?"
		args[i] = id
	}

	query := "SELECT id, user_id, question_id, folder_id, content, created_at, updated_at FROM answers WHERE question_id IN (" + strings.Join(placeholders, ",") + ")"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		logging.Logger.Error("解答取得中にエラーが発生", "error", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var a answerDomain.Answer
		var createdAt, updatedAt []byte

		if err := rows.Scan(&a.Id, &a.UserId, &a.QuestionId, &a.FolderId, &a.Content, &createdAt, &updatedAt); err != nil {
			logging.Logger.Error("解答の読み込みに失敗", "error", err)
			return nil, err
		}

		if err := a.ParseTimeFields(createdAt, updatedAt); err != nil {
			logging.Logger.Error("日時フィールドのパースに失敗", "error", err)
			return nil, err
		}

		answers = append(answers, &a)
	}

	if err = rows.Err(); err != nil {
		logging.Logger.Error("解答の取得完了中にエラーが発生", "error", err)
		return nil, err
	}

	return answers, nil
}
