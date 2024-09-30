package chatgpt

import "context"

type ChatGPTRepository interface {
	Ocr(ctx context.Context, imageURL string) (string, error)
}
