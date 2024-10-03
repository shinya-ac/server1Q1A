package chatgpt

import (
	"context"
)

type GenerateQasUseCase struct {
	chatGptClient ChatGptClient
}

func NewGenerateQasUseCase(client ChatGptClient) *GenerateQasUseCase {
	return &GenerateQasUseCase{
		chatGptClient: client,
	}
}

func (uc *GenerateQasUseCase) Run(ctx context.Context, content string) ([]*Qas, error) {
	return uc.chatGptClient.GenerateQas(ctx, content)
}
