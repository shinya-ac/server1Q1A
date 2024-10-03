package chatgpt

import (
	"context"

	"github.com/shinya-ac/server1Q1A/pkg/logging"
)

type OcrUseCase struct {
	chatGptClient ChatGptClient
}

func NewChatGPTUseCase(client ChatGptClient) *OcrUseCase {
	return &OcrUseCase{
		chatGptClient: client,
	}
}

func (uc *OcrUseCase) HandleImageAnalysis(ctx context.Context, imageURL string) (string, error) {
	response, err := uc.chatGptClient.Ocr(ctx, imageURL)
	if err != nil {
		return "", err
	}

	logging.Logger.Info("レスポンス", "response", response)
	return response, nil
}
