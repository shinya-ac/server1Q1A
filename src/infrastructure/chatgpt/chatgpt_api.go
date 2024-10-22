package chatgpt

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	appChatGPT "github.com/shinya-ac/server1Q1A/application/chatgpt"
	"github.com/shinya-ac/server1Q1A/application/microcms"
	config "github.com/shinya-ac/server1Q1A/configs"
	"github.com/shinya-ac/server1Q1A/pkg/logging"
)

const (
	apiURL    = "https://api.openai.com/v1/chat/completions"
	modelName = "gpt-4o"
)

type ChatGptClient struct {
	apiKey         string
	httpClient     *http.Client
	microcmsClient microcms.MicrocmsClient
}

func NewChatGptClient(microcmsClient microcms.MicrocmsClient) *ChatGptClient {
	chatgptApiKey := config.Config.ChatGptApiKey
	return &ChatGptClient{
		apiKey:         chatgptApiKey,
		httpClient:     &http.Client{Timeout: 50 * time.Second},
		microcmsClient: microcmsClient,
	}
}

type chatGPTResponse struct {
	Choices []struct {
		Message struct {
			Content interface{} `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func (client *ChatGptClient) sendRequest(ctx context.Context, reqBody map[string]interface{}) (*chatGPTResponse, error) {
	reqBodyJSON, err := json.Marshal(reqBody)
	if err != nil {
		logging.Logger.Error("リクエストボディの変換に失敗", "error", err)
		return nil, fmt.Errorf("リクエストボディの変換に失敗しました: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, apiURL, bytes.NewBuffer(reqBodyJSON))
	if err != nil {
		logging.Logger.Error("リクエストの作成に失敗", "error", err)
		return nil, fmt.Errorf("リクエストの作成に失敗しました: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+client.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.httpClient.Do(req)
	if err != nil {
		logging.Logger.Error("リクエストの送信に失敗", "error", err)
		return nil, fmt.Errorf("リクエストの送信に失敗しました: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errBody, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			logging.Logger.Error("エラーレスポンスの読み取りに失敗", "error", readErr)
		} else {
			logging.Logger.Error("APIからのエラーレスポンス", "body", string(errBody))
		}
		return nil, fmt.Errorf("予期しないステータスコードが返されました: %d", resp.StatusCode)
	}

	var result chatGPTResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		logging.Logger.Error("レスポンスのデコードに失敗", "error", err)
		return nil, fmt.Errorf("レスポンスのデコードに失敗しました: %w", err)
	}

	if len(result.Choices) == 0 {
		logging.Logger.Error("ChatGPTからのレスポンスがありません")
		return nil, fmt.Errorf("ChatGPTからのレスポンスがありません")
	}

	return &result, nil
}

func (client *ChatGptClient) Ocr(ctx context.Context, imageURL string, contentID string) (string, error) {
	logging.Logger.Info("Ocr 実行開始", "imageURL", imageURL)

	prompt, err := client.microcmsClient.GetPrompt(contentID)
	if err != nil {
		return "", fmt.Errorf("プロンプトの取得に失敗しました: %w", err)
	}
	if prompt == "" {
		return "", fmt.Errorf("プロンプトが空です")
	}

	reqBody := map[string]interface{}{
		"model": modelName,
		"messages": []map[string]interface{}{
			{
				"role": "user",
				"content": []map[string]interface{}{
					{
						"type": "text",
						"text": prompt,
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

	result, err := client.sendRequest(ctx, reqBody)
	if err != nil {
		return "", err
	}

	content, ok := result.Choices[0].Message.Content.(string)
	if !ok {
		return "", fmt.Errorf("レスポンスのコンテンツが期待する形式ではありません")
	}

	return content, nil
}

func (client *ChatGptClient) GenerateQas(ctx context.Context, content string) ([]*appChatGPT.Qas, error) {
	logging.Logger.Info("GenerateQas 実行開始")

	reqBody := map[string]interface{}{
		"model": modelName,
		"messages": []map[string]interface{}{
			{
				"role": "user",
				"content": []map[string]interface{}{
					{
						"type": "text",
						"text": "上記の文章から一問一答を５つ作成してください。質問は「Q.」、解答は「A.」で始めてください。",
					},
					{
						"type": "text",
						"text": content,
					},
				},
			},
		},
		"max_tokens": 4000,
	}

	result, err := client.sendRequest(ctx, reqBody)
	if err != nil {
		return nil, err
	}

	contentStr, ok := result.Choices[0].Message.Content.(string)
	if !ok {
		return nil, fmt.Errorf("レスポンスのコンテンツが期待する形式ではありません")
	}

	qas, err := appChatGPT.ParseQas(contentStr)
	if err != nil {
		return nil, fmt.Errorf("QAペアの解析に失敗しました: %w", err)
	}

	return qas, nil
}
