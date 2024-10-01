package answer

import (
	"context"
)

type AnswerRepository interface {
	Create(ctx context.Context, answer *Answer) error
	BulkCreate(ctx context.Context, answers []*Answer) error
}
