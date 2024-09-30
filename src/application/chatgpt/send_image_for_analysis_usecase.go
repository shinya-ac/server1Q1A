package chatgpt

import (
	"context"

	"github.com/shinya-ac/server1Q1A/domain/chatgpt"
	"github.com/shinya-ac/server1Q1A/pkg/logging"
)

type OcrUseCase struct {
	repository chatgpt.ChatGPTRepository
}

func NewChatGPTUseCase(repo chatgpt.ChatGPTRepository) *OcrUseCase {
	return &OcrUseCase{
		repository: repo,
	}
}

func (uc *OcrUseCase) HandleImageAnalysis(ctx context.Context, imageURL string) (string, error) {
	response, err := uc.repository.Ocr(ctx, imageURL)
	if err != nil {
		return "", err
	}

	logging.Logger.Info("レスポンス", "response", response)
	return response, nil
}
