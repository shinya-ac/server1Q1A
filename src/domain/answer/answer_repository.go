package answer

import (
	"context"
)

type AnswerRepository interface {
	Create(ctx context.Context, answer *Answer) error
}
