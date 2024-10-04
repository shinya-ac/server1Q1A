package chatgpt

type generateQasResponse struct {
	Questions []string `json:"questions"`
	Answers   []string `json:"answers"`
}
