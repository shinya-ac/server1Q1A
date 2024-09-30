package chatgpt

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	config "github.com/shinya-ac/server1Q1A/configs"
	chatgpt "github.com/shinya-ac/server1Q1A/domain/chatgpt"
	"github.com/shinya-ac/server1Q1A/pkg/logging"
)

const apiUrl = "https://api.openai.com/v1/chat/completions"

type ChatGPTRepository struct {
	apiKey string
}

func NewChatGPTAPI() chatgpt.ChatGPTRepository {
	chatgptApiKey := config.Config.ChatGptApiKey
	return &ChatGPTRepository{
		apiKey: chatgptApiKey,
	}
}

func (client *ChatGPTRepository) Ocr(ctx context.Context, imageURL string) (string, error) {
	logging.Logger.Info("Ocr 実行開始", "imageURL", imageURL)
	reqBody := map[string]interface{}{
		"model": "gpt-4o",
		"messages": []map[string]interface{}{
			{
				"role": "user",
				"content": []map[string]interface{}{
					{
						"type": "text",
						"text": "この画像に書かれている文章を抽出してください。出力はただその文章の文言を出力するだけで大丈夫です。",
					},
					{
						"type": "image_url",
						"image_url": map[string]string{
							"url": imageURL,
						},
					},
				},
			},
		},
		"max_tokens": 4000,
	}
	reqBodyJSON, err := json.Marshal(reqBody)
	if err != nil {
		logging.Logger.Error("リクエストボディの変換に失敗", "error", err)
		return "", fmt.Errorf("bodyの変換に失敗しました %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(reqBodyJSON))
	if err != nil {
		logging.Logger.Error("リクエストの作成に失敗", "error", err)
		return "", fmt.Errorf("リクエストの作成に失敗しました: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+client.apiKey)
	req.Header.Set("Content-Type", "application/json")

	clientHTTP := &http.Client{Timeout: 10 * time.Second}
	resp, err := clientHTTP.Do(req)
	if err != nil {
		logging.Logger.Error("リクエストの送信に失敗", "error", err)
		return "", fmt.Errorf("リクエストの送信に失敗しました: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResponse map[string]interface{}
		errBody, err := io.ReadAll(resp.Body)
		if err != nil {
			logging.Logger.Error("エラーレスポンスの読み取りに失敗", "error", err)
		} else {
			logging.Logger.Error("APIからのエラーレスポンス", "body", string(errBody))
			if err := json.Unmarshal(errBody, &errResponse); err == nil {
				logging.Logger.Error("エラーの詳細", "エラーメッセージ", errResponse)
			} else {
				logging.Logger.Error("エラーレスポンスのデコードに失敗", "error", err)
			}
		}
		return "", fmt.Errorf("予期しないコードが返却されました: %d", resp.StatusCode)
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		logging.Logger.Error("レスポンスのデコードに失敗", "error", err)
		return "", fmt.Errorf("レスポンスのデコードに失敗しました: %w", err)
	}

	if len(result.Choices) > 0 {
		return result.Choices[0].Message.Content, nil
	}
	logging.Logger.Error("OpenAIからのレスポンスがありません")
	return "", fmt.Errorf("OpenAIからのレスポンスがありません")
}
