package chatgpt

import "context"

type Qas struct {
	Question string
	Answer   string
}

type ChatGptClient interface {
	GenerateQas(ctx context.Context, content string) ([]*Qas, error)
	Ocr(ctx context.Context, imageURL string) (string, error)
}
