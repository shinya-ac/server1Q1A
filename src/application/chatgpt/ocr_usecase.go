package chatgpt

import (
	"context"

	config "github.com/shinya-ac/server1Q1A/configs"
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

func (uc *OcrUseCase) HandleOcr(ctx context.Context, imageURL string) (string, error) {
	OcrContentID := config.Config.MicrocmsOcrContentID
	response, err := uc.chatGptClient.Ocr(ctx, imageURL, OcrContentID)
	if err != nil {
		return "", err
	}

	logging.Logger.Info("レスポンス", "response", response)
	return response, nil
}
