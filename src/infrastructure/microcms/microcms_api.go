package microcms

import (
	"fmt"

	"github.com/microcmsio/microcms-go-sdk"
)

type MicrocmsClient struct {
	client *microcms.Client
}

func NewMicrocmsClient(apiKey string) *MicrocmsClient {
	client := microcms.New("prompt", apiKey)
	return &MicrocmsClient{client: client}
}

type microCMSPromptResponse struct {
	Prompt string `json:"prompt"`
}

func (r *MicrocmsClient) GetPrompt(contentID string) (string, error) {
	var content microCMSPromptResponse

	fmt.Printf("contentId:%v", contentID)
	err := r.client.Get(
		microcms.GetParams{
			Endpoint:  "get",
			ContentID: contentID,
		},
		&content,
	)
	if err != nil {
		return "", fmt.Errorf("failed to fetch prompt from microCMS: %w", err)
	}

	return content.Prompt, nil
}
