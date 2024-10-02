package question

import (
	"context"
)

type QuestionRepository interface {
	Create(ctx context.Context, question *Question) error
	BulkCreate(ctx context.Context, questions []*Question) error
	GetQuestionsByFolderId(ctx context.Context, folderId string) ([]*Question, error)
}
