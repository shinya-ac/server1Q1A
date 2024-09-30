package repository

import (
	"context"
	"database/sql"

	questionDomain "github.com/shinya-ac/server1Q1A/domain/question"
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
