package microcms

type MicrocmsClient interface {
	GetPrompt(contentID string) (string, error)
}
