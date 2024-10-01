package question

import (
	"context"
)

type QuestionRepository interface {
	Create(ctx context.Context, question *Question) error
	BulkCreate(ctx context.Context, questions []*Question) error
}
