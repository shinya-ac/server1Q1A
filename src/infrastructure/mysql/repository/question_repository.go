package repository

import (
	"context"
	"database/sql"
	"strings"

	questionDomain "github.com/shinya-ac/server1Q1A/domain/question"
	"github.com/shinya-ac/server1Q1A/pkg/logging"
)

type MySQLQuestionRepository struct {
	db *sql.DB
}

func NewMySQLQuestionRepository(db *sql.DB) *MySQLQuestionRepository {
	return &MySQLQuestionRepository{db: db}
}

func (r *MySQLQuestionRepository) Create(ctx context.Context, q *questionDomain.Question) error {
	query := "INSERT INTO questions (id, user_id, folder_id, content, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)"
	_, err := r.db.ExecContext(ctx, query, q.Id, q.UserId, q.FolderId, q.Content, q.CreatedAt, q.UpdatedAt)
	return err
}

func (r *MySQLQuestionRepository) BulkCreate(ctx context.Context, questions []*questionDomain.Question) error {
	if len(questions) == 0 {
		return nil
	}

	query := `INSERT INTO questions (id, user_id, folder_id, content, created_at, updated_at) VALUES `
	values := []string{}
	args := []interface{}{}

	for _, q := range questions {
		values = append(values, "(?, ?, ?, ?, ?, ?)")
		args = append(args, q.Id, q.UserId, q.FolderId, q.Content, q.CreatedAt, q.UpdatedAt)
	}

	query += strings.Join(values, ",")
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		logging.Logger.Error("質問バルク作成用のクエリの準備に失敗しました", "error", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, args...)
	if err != nil {
		logging.Logger.Error("質問のバルク作成に失敗しました", "error", err)
		return err
	}

	return nil
}

func (r *MySQLQuestionRepository) GetQuestionsByFolderId(ctx context.Context, folderId string) ([]*questionDomain.Question, error) {
	var questions []*questionDomain.Question
	query := "SELECT id, user_id, folder_id, content, created_at, updated_at FROM questions WHERE folder_id = ?"

	rows, err := r.db.QueryContext(ctx, query, folderId)
	if err != nil {
		logging.Logger.Error("質問取得中にエラーが発生", "error", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var q questionDomain.Question
		var createdAt, updatedAt []byte

		if err := rows.Scan(&q.Id, &q.UserId, &q.FolderId, &q.Content, &createdAt, &updatedAt); err != nil {
			logging.Logger.Error("質問の読み込みに失敗", "error", err)
			return nil, err
		}

		if err := q.ParseTimeFields(createdAt, updatedAt); err != nil {
			logging.Logger.Error("日時フィールドのパースに失敗", "error", err)
			return nil, err
		}

		questions = append(questions, &q)
	}

	if err = rows.Err(); err != nil {
		logging.Logger.Error("質問の取得完了中にエラーが発生", "error", err)
		return nil, err
	}

	return questions, nil
}
