package chatgpt

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	chatgpt "github.com/shinya-ac/server1Q1A/application/chatgpt"
	config "github.com/shinya-ac/server1Q1A/configs"
	"github.com/shinya-ac/server1Q1A/pkg/logging"
)

const apiUrl = "https://api.openai.com/v1/chat/completions"

type ChatGPTClient struct {
	apiKey string
}

func NewChatGPTAPI() *ChatGPTClient {
	chatgptApiKey := config.Config.ChatGptApiKey
	return &ChatGPTClient{
		apiKey: chatgptApiKey,
	}
}

func (client *ChatGPTClient) Ocr(ctx context.Context, imageURL string) (string, error) {
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
	logging.Logger.Error("ChatGPTからのレスポンスがありません")
	return "", fmt.Errorf("ChatGPTからのレスポンスがありません")
}

func (client *ChatGPTClient) GenerateQas(ctx context.Context, content string) ([]*chatgpt.Qas, error) {
	reqBody := map[string]interface{}{
		"model": "gpt-4",
		"messages": []map[string]interface{}{
			{
				"role": "user",
				"content": []map[string]interface{}{
					{"type": "text", "text": "上記の文章から一問一答を５つ作成してください。質問は「Q.」、解答は「A.」で始めてください。"},
					{"type": "text", "text": content},
				},
			},
		},
		"max_tokens": 4000,
	}

	reqBodyJSON, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("リクエストボディの作成に失敗しました: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(reqBodyJSON))
	if err != nil {
		return nil, fmt.Errorf("リクエストの作成に失敗しました: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+client.apiKey)
	req.Header.Set("Content-Type", "application/json")

	httpClient := &http.Client{Timeout: 30 * time.Second}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("リクエストの送信に失敗しました: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ChatGPT APIからの予期しないレスポンス: %d", resp.StatusCode)
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("レスポンスのデコードに失敗しました: %w", err)
	}

	if len(result.Choices) == 0 {
		return nil, fmt.Errorf("ChatGPTから有効な回答がありません")
	}

	qas, err := parseQas(result.Choices[0].Message.Content)
	if err != nil {
		return nil, fmt.Errorf("QAペアの解析に失敗しました: %w", err)
	}

	return qas, nil
}

func parseQas(content string) ([]*chatgpt.Qas, error) {
	var qas []*chatgpt.Qas
	var questions []string
	var answers []string

	lines := strings.Split(content, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "Q.") {
			questions = append(questions, strings.TrimSpace(line))
		} else if strings.HasPrefix(line, "A.") {
			answers = append(answers, strings.TrimSpace(line))
		}
	}

	if len(questions) != len(answers) {
		return nil, fmt.Errorf("質問と回答の数が一致しません")
	}

	qas = make([]*chatgpt.Qas, len(questions))
	for i := range questions {
		qas[i] = &chatgpt.Qas{
			Question: questions[i],
			Answer:   answers[i],
		}
	}

	return qas, nil
}
