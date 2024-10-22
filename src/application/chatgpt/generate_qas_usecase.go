package chatgpt

import (
	"context"

	config "github.com/shinya-ac/server1Q1A/configs"
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
	generateQAsContentID := config.Config.MicrocmsGenerateQAsContentID
	response, err := uc.chatGptClient.GenerateQas(ctx, content, generateQAsContentID)
	if err != nil {
		return nil, err
	}
	return response, nil
}
